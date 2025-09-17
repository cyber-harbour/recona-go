package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// MockAccountClient implements the internal.Client interface for testing
type MockAccountClient struct {
	// MakeRequestFunc allows customization of the MakeRequest behavior
	MakeRequestFunc func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error)

	// Track calls for verification
	calls []MockCall
}

type MockCall struct {
	Method   string
	Endpoint string
	Body     interface{}
	Context  context.Context
}

func (m *MockAccountClient) MakeRequest(
	ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	// Record the call
	m.calls = append(m.calls, MockCall{
		Method:   method,
		Endpoint: endpoint,
		Body:     body,
		Context:  ctx,
	})

	// Use custom function if provided, otherwise return default success
	if m.MakeRequestFunc != nil {
		return m.MakeRequestFunc(ctx, method, endpoint, body)
	}

	// Default behavior - return empty response
	return createAccountMockResponse(`{}`), nil
}

func (m *MockAccountClient) GetCalls() []MockCall {
	return m.calls
}

func (m *MockAccountClient) Reset() {
	m.calls = nil
	m.MakeRequestFunc = nil
}

// Helper function to create mock HTTP responses
func createAccountMockResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// Helper function to load test data from file
func loadTestData(filename string) ([]byte, error) {
	// Try multiple possible paths for the data file
	possiblePaths := []string{
		filepath.Join("data", filename),
		filepath.Join("..", "data", filename),
		filepath.Join("..", "..", "data", filename),
		filename, // Direct path
	}

	for _, path := range possiblePaths {
		// #nosec G304 - This is a test file read, not user input
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("could not find test data file: %s", filename)
}

func TestNewAccountService(t *testing.T) {
	mockAccountClient := &MockAccountClient{}
	service := NewAccountService(mockAccountClient)

	if service == nil {
		t.Fatal("NewAccountService returned nil")
	}

	if service.client != mockAccountClient {
		t.Error("NewAccountService did not set the client correctly")
	}
}

func TestAccountService_GetDetails_Success(t *testing.T) { // nolint: funlen
	// Load test data from account.json
	testData, err := loadTestData("account.json")
	if err != nil {
		// If we can't load the file, create sample data
		t.Logf("Could not load account.json: %v, using sample data", err)
		testData = []byte(`{
		  "data": {
			"id": 1,
			"login": "back",
			"status": 1,
			"nickname": "Backend",
			"subscription_id": 2,
			"subscription_name": "Enterprise",
			"group_id": 1,
			"group_title": "default",
			"role_id": 1,
			"subscription_start_at": "2025-02-06T17:43:29.847Z",
			"subscription_expires_at": "2026-02-06T17:43:32.216Z",
			"organization_id": 1,
			"organization_title": "Recona Team",
			"created_at": "2025-02-06T15:43:48.156109Z",
			"updated_at": "2025-08-11T08:37:25.633165Z",
			"last_seen": "2025-09-15T13:15:58.906281Z",
			"total_request_count": 389,
			"daily_request_count": 126,
			"week_request_count": 129,
			"enabled_two_fa": true,
			"products_permission": {
			  "recona": true
			},
			"permissions": {
			  "ui_rows_limit": 1000,
			  "api_rows_limit": 10000,
			  "request_limit_per_day": 100000,
			  "filter_limit": 15,
			  "request_rate_limit": 15
			},
			"request_count": 126,
			"request_limit_per_day": 100000,
			"start_at": "2025-09-15T00:00:00Z",
			"EndAt": "2025-09-16T00:00:00Z"
		  }
		}
		`)
	}

	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			return createAccountMockResponse(string(testData)), nil
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx := context.Background()

	profile, err := service.GetDetails(ctx)

	// Verify no error occurred
	if err != nil {
		t.Fatalf("GetDetails returned error: %v", err)
	}

	// Verify profile is not nil
	if profile == nil {
		t.Fatal("GetDetails returned nil profile")
	}

	// Verify the correct endpoint was called
	calls := mockAccountClient.GetCalls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 call, got %d", len(calls))
	}

	call := calls[0]
	if call.Method != "GET" {
		t.Errorf("Expected GET method, got %s", call.Method)
	}
	if call.Endpoint != accountEndpoint {
		t.Errorf("Expected endpoint %s, got %s", accountEndpoint, call.Endpoint)
	}
	if call.Body != nil {
		t.Errorf("Expected nil body, got %v", call.Body)
	}

	// Verify profile data (this will depend on your actual models.Profile structure)
	// Add specific field validations based on the actual JSON structure
	t.Logf("Profile data: %+v", profile)
}

func TestAccountService_GetDetails_WithRealData(t *testing.T) {
	// This test specifically uses the actual account.json data
	testData, err := loadTestData("account.json")
	if err != nil {
		t.Skipf("Skipping real data test: %v", err)
		return
	}

	// First, let's validate that the JSON is valid
	var jsonCheck interface{}
	if err := json.Unmarshal(testData, &jsonCheck); err != nil {
		t.Fatalf("Invalid JSON in account.json: %v", err)
	}

	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			return createAccountMockResponse(string(testData)), nil
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx := context.Background()

	profile, err := service.GetDetails(ctx)

	if err != nil {
		t.Fatalf("GetDetails with real data failed: %v", err)
	}

	if profile == nil {
		t.Fatal("GetDetails returned nil profile with real data")
	}

	// Log the actual structure for debugging
	t.Logf("Real profile data structure: %+v", profile)
}

func TestAccountService_GetDetails_HTTPError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantError  string
	}{
		{
			name:       "400 Bad Request",
			statusCode: 400,
			wantError:  "failed to make request",
		},
		{
			name:       "401 Unauthorized",
			statusCode: 401,
			wantError:  "failed to make request",
		},
		{
			name:       "403 Forbidden",
			statusCode: 403,
			wantError:  "failed to make request",
		},
		{
			name:       "404 Not Found",
			statusCode: 404,
			wantError:  "failed to make request",
		},
		{
			name:       "500 Internal Server Error",
			statusCode: 500,
			wantError:  "failed to make request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccountClient := &MockAccountClient{
				MakeRequestFunc: func(
					ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("HTTP %d error", tt.statusCode)
				},
			}

			service := NewAccountService(mockAccountClient)
			ctx := context.Background()

			profile, err := service.GetDetails(ctx)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}

			if profile != nil {
				t.Errorf("Expected nil profile, got %+v", profile)
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.wantError, err.Error())
			}
		})
	}
}

func TestAccountService_GetDetails_JSONDecodeError(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		wantError    string
	}{
		{
			name:         "Invalid JSON",
			responseBody: `{invalid json}`,
			wantError:    "failed to decode response body",
		},
		{
			name:         "Empty response",
			responseBody: ``,
			wantError:    "failed to decode response body",
		},
		{
			name:         "Malformed JSON object",
			responseBody: `{"name": "test", "incomplete":}`,
			wantError:    "failed to decode response body",
		},
		{
			name:         "Wrong JSON type",
			responseBody: `"string instead of object"`,
			wantError:    "failed to decode response body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccountClient := &MockAccountClient{
				MakeRequestFunc: func(
					ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
					return createAccountMockResponse(tt.responseBody), nil
				},
			}

			service := NewAccountService(mockAccountClient)
			ctx := context.Background()

			profile, err := service.GetDetails(ctx)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}

			if profile != nil {
				t.Errorf("Expected nil profile, got %+v", profile)
			}

			if !strings.Contains(err.Error(), tt.wantError) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.wantError, err.Error())
			}
		})
	}
}

func TestAccountService_GetDetails_ContextCancellation(t *testing.T) {
	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			// Simulate context cancellation
			return nil, context.Canceled
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	profile, err := service.GetDetails(ctx)

	if err == nil {
		t.Error("Expected error due to context cancellation")
	}

	if profile != nil {
		t.Errorf("Expected nil profile, got %+v", profile)
	}

	if !strings.Contains(err.Error(), "failed to make request") {
		t.Errorf("Expected error about failed request, got: %v", err)
	}
}

func TestAccountService_GetDetails_ContextTimeout(t *testing.T) {
	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			// Simulate timeout
			return nil, context.DeadlineExceeded
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	profile, err := service.GetDetails(ctx)

	if err == nil {
		t.Error("Expected error due to timeout")
	}

	if profile != nil {
		t.Errorf("Expected nil profile, got %+v", profile)
	}
}

func TestAccountService_GetDetails_NetworkError(t *testing.T) {
	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			return nil, errors.New("network connection failed")
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx := context.Background()

	profile, err := service.GetDetails(ctx)

	if err == nil {
		t.Error("Expected network error")
	}

	if profile != nil {
		t.Errorf("Expected nil profile, got %+v", profile)
	}

	if !strings.Contains(err.Error(), "failed to make request") {
		t.Errorf("Expected error about failed request, got: %v", err)
	}
}

// Test with various valid JSON responses to ensure robustness
func TestAccountService_GetDetails_VariousValidResponses(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		description  string
	}{
		{
			name:         "Minimal profile",
			responseBody: `{"id": 123}`,
			description:  "Profile with only ID",
		},
		{
			name:         "Full profile",
			responseBody: `{"id": 123, "login": "user@test.com", "nickname": "Test User", "status": 1}`,
			description:  "Complete profile information",
		},
		{
			name:         "Profile with nested objects",
			responseBody: `{"id": 123, "permissions": {"ui_rows_limit": 1000, "api_rows_limit": 10000}}`,
			description:  "Profile with nested settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAccountClient := &MockAccountClient{
				MakeRequestFunc: func(
					ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
					return createAccountMockResponse(tt.responseBody), nil
				},
			}

			service := NewAccountService(mockAccountClient)
			ctx := context.Background()

			profile, err := service.GetDetails(ctx)

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}

			if profile == nil {
				t.Errorf("Expected profile for %s, got nil", tt.description)
			}

			t.Logf("%s result: %+v", tt.description, profile)
		})
	}
}

// ThreadSafeMockAccountClient is a thread-safe version of MockAccountClient for concurrent testing
type ThreadSafeMockAccountClient struct {
	// MakeRequestFunc allows customization of the MakeRequest behavior
	MakeRequestFunc func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error)

	// Track calls for verification with mutex protection
	mu    sync.Mutex
	calls []MockCall
}

func (m *ThreadSafeMockAccountClient) GetCalls() []MockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Return a copy to avoid race conditions
	calls := make([]MockCall, len(m.calls))
	copy(calls, m.calls)
	return calls
}

func (m *ThreadSafeMockAccountClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = nil
	m.MakeRequestFunc = nil
}

// ThreadSafeMockAccountClient methods
func (m *ThreadSafeMockAccountClient) MakeRequest(
	ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	// Record the call with mutex protection
	m.mu.Lock()
	m.calls = append(m.calls, MockCall{
		Method:   method,
		Endpoint: endpoint,
		Body:     body,
		Context:  ctx,
	})
	m.mu.Unlock()

	// Use custom function if provided, otherwise return default success
	if m.MakeRequestFunc != nil {
		return m.MakeRequestFunc(ctx, method, endpoint, body)
	}

	// Default behavior - return empty response
	return createAccountMockResponse(`{}`), nil
}

// Test concurrent requests
func TestAccountService_GetDetails_Concurrent(t *testing.T) {
	// Thread-safe mock client
	mockAccountClient := &ThreadSafeMockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			// Simulate some processing time
			time.Sleep(10 * time.Millisecond)
			return createAccountMockResponse(`{"id": 123}`), nil
		},
	}

	service := NewAccountService(mockAccountClient)

	const numGoroutines = 10
	results := make(chan error, numGoroutines)

	// Launch multiple goroutines
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			ctx := context.Background()
			profile, err := service.GetDetails(ctx)

			if err != nil {
				results <- fmt.Errorf("goroutine %d failed: %w", id, err)
				return
			}

			if profile == nil {
				results <- fmt.Errorf("goroutine %d got nil profile", id)
				return
			}

			results <- nil
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		if err := <-results; err != nil {
			t.Errorf("Concurrent test failed: %v", err)
		}
	}

	// Verify all calls were made
	calls := mockAccountClient.GetCalls()
	if len(calls) != numGoroutines {
		t.Errorf("Expected %d calls, got %d", numGoroutines, len(calls))
	}
}

// Benchmark test
func BenchmarkAccountService_GetDetails(b *testing.B) {
	testData := `{"id": 123, "email": "bench@test.com", "nickname": "Benchmark User"}`

	mockAccountClient := &MockAccountClient{
		MakeRequestFunc: func(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
			return createAccountMockResponse(testData), nil
		},
	}

	service := NewAccountService(mockAccountClient)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetDetails(ctx)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// Test that verifies the exact endpoint constant is used
func TestAccountService_EndpointConstant(t *testing.T) {
	if accountEndpoint != "/customers/account" {
		t.Errorf("accountEndpoint constant changed: expected '/customers/account', got '%s'", accountEndpoint)
	}

	mockAccountClient := &MockAccountClient{}
	service := NewAccountService(mockAccountClient)
	ctx := context.Background()

	// This will fail due to nil response, but we're testing the endpoint
	_, _ = service.GetDetails(ctx)

	calls := mockAccountClient.GetCalls()
	if len(calls) == 1 && calls[0].Endpoint == "/customers/account" {
		return
	}

	t.Errorf("Service did not use the correct endpoint constant")
}

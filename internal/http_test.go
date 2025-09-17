package internal

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Test constants
func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"BaseURL", BaseURL, "https://api.recona.io"},
		{"authorizationHeaderName", authorizationHeaderName, "Authorization"},
		{"authorizationType", authorizationType, "Bearer "},
		{"contentTypeHeaderName", contentTypeHeaderName, "Content-Type"},
		{"acceptHeaderName", acceptHeaderName, "Accept"},
		{"defaultContentType", defaultContentType, "application/json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("constant %s = %q, want %q", tt.name, tt.actual, tt.expected)
			}
		})
	}
}

func TestMakeAuthenticatedRequest(t *testing.T) { // nolint: funlen
	tests := []struct {
		name           string
		method         string
		token          string
		body           interface{}
		serverResponse string
		statusCode     int
		wantError      bool
		errorContains  string
		validateReq    func(t *testing.T, req *http.Request)
	}{
		{
			name:           "successful GET request",
			method:         "GET",
			token:          "test-token",
			body:           nil,
			serverResponse: `{"success": true}`,
			statusCode:     200,
			wantError:      false,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Method != "GET" {
					t.Errorf("expected GET method, got %s", req.Method)
				}
				if auth := req.Header.Get("Authorization"); auth != "Bearer test-token" {
					t.Errorf("expected Authorization header 'Bearer test-token', got '%s'", auth)
				}
				if ct := req.Header.Get("Content-Type"); ct != "application/json" {
					t.Errorf("expected Content-Type header 'application/json', got '%s'", ct)
				}
				if accept := req.Header.Get("Accept"); accept != "application/json" {
					t.Errorf("expected Accept header 'application/json', got '%s'", accept)
				}
			},
		},
		{
			name:           "successful POST request with JSON body",
			method:         "POST",
			token:          "another-token",
			body:           map[string]string{"key": "value"},
			serverResponse: `{"id": 123}`,
			statusCode:     201,
			wantError:      false,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Method != "POST" {
					t.Errorf("expected POST method, got %s", req.Method)
				}
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				expected := `{"key":"value"}`
				if string(bodyBytes) != expected {
					t.Errorf("expected request body %s, got %s", expected, string(bodyBytes))
				}
			},
		},
		{
			name:           "successful PUT request with struct body",
			method:         "PUT",
			token:          "put-token",
			body:           struct{ Name string }{Name: "test"},
			serverResponse: `{"updated": true}`,
			statusCode:     200,
			wantError:      false,
			validateReq: func(t *testing.T, req *http.Request) {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				expected := `{"Name":"test"}`
				if string(bodyBytes) != expected {
					t.Errorf("expected request body %s, got %s", expected, string(bodyBytes))
				}
			},
		},
		{
			name:          "400 error response",
			method:        "POST",
			token:         "token",
			body:          nil,
			statusCode:    400,
			wantError:     true,
			errorContains: "API error 400",
		},
		{
			name:          "401 error response",
			method:        "GET",
			token:         "invalid-token",
			body:          nil,
			statusCode:    401,
			wantError:     true,
			errorContains: "API error 401",
		},
		{
			name:          "404 error response",
			method:        "GET",
			token:         "token",
			body:          nil,
			statusCode:    404,
			wantError:     true,
			errorContains: "API error 404",
		},
		{
			name:          "500 error response",
			method:        "POST",
			token:         "token",
			body:          nil,
			statusCode:    500,
			wantError:     true,
			errorContains: "API error 500",
		},
		{
			name:           "error with response body",
			method:         "POST",
			token:          "token",
			body:           nil,
			serverResponse: `{"error": "validation failed"}`,
			statusCode:     400,
			wantError:      true,
			errorContains:  "validation failed",
		},
		{
			name:          "JSON marshal error",
			method:        "POST",
			token:         "token",
			body:          make(chan int), // channels cannot be marshaled to JSON
			statusCode:    200,
			wantError:     true,
			errorContains: "failed to marshal request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.validateReq != nil {
					tt.validateReq(t, r)
				}
				w.WriteHeader(tt.statusCode)
				if tt.serverResponse != "" {
					_, err := w.Write([]byte(tt.serverResponse))
					if err != nil {
						t.Errorf("failed to write response: %v", err)
					}
				}
			}))
			defer server.Close()

			client := server.Client()
			ctx := context.Background()

			resp, err := MakeAuthenticatedRequest(ctx, client, tt.method, server.URL, tt.token, tt.body)

			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if resp == nil {
				t.Error("expected non-nil response")
				return
			}

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.serverResponse != "" {
				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("failed to read response body: %v", err)
				} else if string(bodyBytes) != tt.serverResponse {
					t.Errorf("expected response body '%s', got '%s'", tt.serverResponse, string(bodyBytes))
				}
			}

			_ = resp.Body.Close()
		})
	}
}

func TestMakeAuthenticatedRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := server.Client()
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel the context immediately
	cancel()

	_, err := MakeAuthenticatedRequest(ctx, client, "GET", server.URL, "token", nil)

	if err == nil {
		t.Error("expected error due to context cancellation")
	}

	if !strings.Contains(err.Error(), "context canceled") && !strings.Contains(err.Error(), "request failed") {
		t.Errorf("expected context cancellation error, got: %v", err)
	}
}

func TestMakeAuthenticatedRequest_ContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := server.Client()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := MakeAuthenticatedRequest(ctx, client, "GET", server.URL, "token", nil)

	if err == nil {
		t.Error("expected timeout error")
	}

	if !strings.Contains(err.Error(), "deadline exceeded") && !strings.Contains(err.Error(), "request failed") {
		t.Errorf("expected timeout error, got: %v", err)
	}
}

func TestMakeAuthenticatedRequest_InvalidURL(t *testing.T) {
	client := &http.Client{}
	ctx := context.Background()

	_, err := MakeAuthenticatedRequest(ctx, client, "GET", "://invalid-url", "token", nil)

	if err == nil {
		t.Error("expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "failed to create request") {
		t.Errorf("expected 'failed to create request' error, got: %v", err)
	}
}

func TestMakeAuthenticatedRequest_NilClient(t *testing.T) {
	ctx := context.Background()

	_, err := MakeAuthenticatedRequest(ctx, nil, "GET", "http://example.com", "token", nil)

	if err == nil {
		t.Error("expected error for nil client")
	}

	if !strings.Contains(err.Error(), "request failed") {
		t.Errorf("expected 'request failed' error, got: %v", err)
	}
}

func TestMakeAuthenticatedRequest_EmptyToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer" {
			t.Errorf("expected Authorization header 'Bearer', got '%s'", auth)
		}
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := server.Client()
	ctx := context.Background()

	_, err := MakeAuthenticatedRequest(ctx, client, "GET", server.URL, "", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDecodeJSON(t *testing.T) { // nolint: funlen
	tests := []struct {
		name      string
		jsonInput string
		target    interface{}
		wantError bool
		validate  func(t *testing.T, target interface{})
	}{
		{
			name:      "decode simple object",
			jsonInput: `{"name": "test", "value": 123}`,
			target:    &map[string]interface{}{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				m := target.(*map[string]interface{})
				if (*m)["name"] != "test" {
					t.Errorf("expected name=test, got %v", (*m)["name"])
				}
				if (*m)["value"] != float64(123) { // JSON numbers become float64
					t.Errorf("expected value=123, got %v", (*m)["value"])
				}
			},
		},
		{
			name:      "decode struct",
			jsonInput: `{"Name": "Alice", "Age": 30}`,
			target: &struct {
				Name string
				Age  int
			}{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				s := target.(*struct {
					Name string
					Age  int
				})
				if s.Name != "Alice" {
					t.Errorf("expected Name=Alice, got %s", s.Name)
				}
				if s.Age != 30 {
					t.Errorf("expected Age=30, got %d", s.Age)
				}
			},
		},
		{
			name:      "decode array",
			jsonInput: `[1, 2, 3]`,
			target:    &[]int{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				arr := target.(*[]int)
				expected := []int{1, 2, 3}
				if len(*arr) != len(expected) {
					t.Errorf("expected length %d, got %d", len(expected), len(*arr))
				}
				for i, v := range expected {
					if (*arr)[i] != v {
						t.Errorf("expected arr[%d]=%d, got %d", i, v, (*arr)[i])
					}
				}
			},
		},
		{
			name:      "decode null",
			jsonInput: `null`,
			target:    &map[string]interface{}{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				m := target.(*map[string]interface{})
				if *m != nil {
					t.Errorf("expected nil map, got %v", *m)
				}
			},
		},
		{
			name:      "invalid JSON",
			jsonInput: `{invalid json}`,
			target:    &map[string]interface{}{},
			wantError: true,
		},
		{
			name:      "empty JSON object",
			jsonInput: `{}`,
			target:    &map[string]interface{}{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				m := target.(*map[string]interface{})
				if len(*m) != 0 {
					t.Errorf("expected empty map, got %v", *m)
				}
			},
		},
		{
			name:      "nested object",
			jsonInput: `{"user": {"name": "Bob", "details": {"age": 25}}}`,
			target:    &map[string]interface{}{},
			wantError: false,
			validate: func(t *testing.T, target interface{}) {
				m := target.(*map[string]interface{})
				user := (*m)["user"].(map[string]interface{})
				if user["name"] != "Bob" {
					t.Errorf("expected user.name=Bob, got %v", user["name"])
				}
				details := user["details"].(map[string]interface{})
				if details["age"] != float64(25) {
					t.Errorf("expected user.details.age=25, got %v", details["age"])
				}
			},
		},
		{
			name:      "type mismatch",
			jsonInput: `{"name": "test"}`,
			target:    &[]int{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.jsonInput)
			err := DecodeJSON(reader, tt.target)

			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, tt.target)
			}
		})
	}
}

func TestDecodeJSON_EmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	target := &map[string]interface{}{}

	err := DecodeJSON(reader, target)
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestDecodeJSON_NilTarget(t *testing.T) {
	reader := strings.NewReader(`{"test": true}`)

	err := DecodeJSON(reader, nil)
	if err == nil {
		t.Error("expected error for nil target")
	}
}

// Benchmark tests
func BenchmarkMakeAuthenticatedRequest(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte(`{"success": true}`))
		if err != nil {
			b.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := server.Client()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, err := MakeAuthenticatedRequest(ctx, client, "GET", server.URL, "token", nil)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		err = resp.Body.Close()
		if err != nil {
			b.Fatalf("failed to close response body: %v", err)
		}
	}
}

func BenchmarkDecodeJSON(b *testing.B) {
	jsonData := `{"name": "test", "value": 123, "active": true, "items": [1, 2, 3, 4, 5]}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(jsonData)
		target := &map[string]interface{}{}
		err := DecodeJSON(reader, target)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cyber-harbour/recona-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockClient is a mock implementation of the internal.Client interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) MakeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	args := m.Called(ctx, method, path, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// Helper function to create a mock HTTP response
func createMockResponse(body interface{}) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(bodyBytes)
	} else {
		bodyReader = strings.NewReader("")
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bodyReader),
		Header:     make(http.Header),
	}
}

// Helper function to create a mock response with raw string body
func createMockResponseWithString(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestNewDomainService(t *testing.T) {
	t.Run("should create new domain service with client", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)

		assert.NotNil(t, service)
		assert.Equal(t, mockClient, service.client)
	})

	t.Run("should create service with nil client", func(t *testing.T) {
		service := NewDomainService(nil)
		assert.NotNil(t, service)
		assert.Nil(t, service.client)
	})
}

func TestDomainService_GetDetails(t *testing.T) { // nolint: funlen
	t.Run("should successfully get domain details", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		domainID := "example.com"

		expectedDomain := &models.Domain{
			Name: domainID,
		}

		mockResponse := createMockResponse(expectedDomain)
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
			Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, domainID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedDomain.Name, result.Name)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle client request error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		domainID := "test-domain"
		expectedError := errors.New("network error")

		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
			Return(nil, expectedError)

		// Act
		result, err := service.GetDetails(ctx, domainID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get domain details for ID test-domain")
		assert.Contains(t, err.Error(), "network error")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle JSON decode error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		domainID := "test-domain-id"

		// Create response with invalid JSON
		mockResponse := createMockResponseWithString(200, "invalid json")
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
			Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, domainID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode domain details response")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty domain ID", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		domainID := ""

		expectedDomain := &models.Domain{}
		mockResponse := createMockResponse(expectedDomain)
		mockClient.On("MakeRequest", ctx, "GET", "/domains/", mock.Anything).Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, domainID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately
		domainID := "test-domain-id"

		expectedError := context.Canceled
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
			Return(nil, expectedError)

		// Act
		result, err := service.GetDetails(ctx, domainID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
		mockClient.AssertExpectations(t)
	})
}

func TestDomainService_Search(t *testing.T) { // nolint: funlen
	t.Run("should successfully search domains", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Search: models.Search{
				Query: "example.com",
			},
			Pagination: models.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		expectedResponse := &models.DomainsResponse{
			Domains: []*models.Domain{
				{Name: "example.com"},
				{Name: "test.com"},
			},
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    2,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Domains, 2)
		assert.Equal(t, int64(2), result.TotalItems.Value)
		assert.Equal(t, "example.com", result.Domains[0].Name)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty search results", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Search: models.Search{
				Query: "nonexistent.com",
			},
		}

		expectedResponse := &models.DomainsResponse{
			Domains: []*models.Domain{},
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    0,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Domains, 0)
		assert.Equal(t, int64(0), result.TotalItems.Value)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle client request error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		searchParams := models.SearchRequest{}
		expectedError := errors.New("request failed")

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchParams).
			Return(nil, expectedError)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search domain records")
		assert.Contains(t, err.Error(), "request failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle JSON decode error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()
		searchParams := models.SearchRequest{}

		mockResponse := createMockResponseWithString(200, "invalid json")
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode domain search response")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle large pagination", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Pagination: models.Pagination{
				Limit:  1000,
				Offset: 5000,
			},
		}

		expectedResponse := &models.DomainsResponse{
			Domains: make([]*models.Domain, 1000),
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    1000,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Domains, 1000)
		mockClient.AssertExpectations(t)
	})
}

func TestDomainService_SearchAll(t *testing.T) { // nolint: funlen
	t.Run("should successfully search all domains with single page", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		expectedDomains := []*models.Domain{
			{Name: "example1.com"},
			{Name: "example2.com"},
		}

		expectedResponse := &models.DomainsResponse{
			Domains: expectedDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    2,
					Relation: "equal",
				},
			},
		}

		expectedSearchRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", expectedSearchRequest).
			Return(mockResponse, nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "example1.com", result[0].Name)
		assert.Equal(t, "example2.com", result[1].Name)
		mockClient.AssertExpectations(t)
	})

	t.Run("should successfully search all domains with multiple pages", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		// First page - 100 results
		firstPageDomains := make([]*models.Domain, 100)
		for i := 0; i < 100; i++ {
			firstPageDomains[i] = &models.Domain{
				Name: fmt.Sprintf("example%d.com", i+1),
			}
		}

		firstPageResponse := &models.DomainsResponse{
			Domains: firstPageDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    150,
					Relation: "equal",
				},
			},
		}

		// Second page - 50 results
		secondPageDomains := make([]*models.Domain, 50)
		for i := 0; i < 50; i++ {
			secondPageDomains[i] = &models.Domain{
				Name: fmt.Sprintf("example%d.com", i+101),
			}
		}

		secondPageResponse := &models.DomainsResponse{
			Domains: secondPageDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    150,
					Relation: "equal",
				},
			},
		}

		firstRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		secondRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 100,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", firstRequest).
			Return(createMockResponse(firstPageResponse), nil)

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", secondRequest).
			Return(createMockResponse(secondPageResponse), nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 150)
		assert.Equal(t, "example1.com", result[0].Name)
		assert.Equal(t, "example150.com", result[149].Name)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle max results limit", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		// Create responses for max results scenario (10000 limit)
		// We'll simulate 100 pages of 100 results each
		for page := 0; page < 100; page++ {
			pageRequest := models.SearchRequest{
				Search: baseParams,
				Pagination: models.Pagination{
					Limit:  100,
					Offset: page * 100,
				},
			}

			pageDomains := make([]*models.Domain, 100)
			for i := 0; i < 100; i++ {
				pageDomains[i] = &models.Domain{
					Name: fmt.Sprintf("example%d.com", page*100+i+1),
				}
			}

			pageResponse := &models.DomainsResponse{
				Domains: pageDomains,
				PaginationResponse: models.PaginationResponse{
					TotalItems: models.TotalItems{
						Value:    15000,
						Relation: "equal",
					},
				},
			}

			mockClient.On("MakeRequest", ctx, "POST", "/domains/search", pageRequest).
				Return(createMockResponse(pageResponse), nil)
		}

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 10000) // Should stop at maxResults
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle partial last page", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		// Simulate scenario where last page has fewer than 100 results
		// but total remaining is more than last page size
		for page := 0; page < 99; page++ {
			pageRequest := models.SearchRequest{
				Search: baseParams,
				Pagination: models.Pagination{
					Limit:  100,
					Offset: page * 100,
				},
			}

			pageDomains := make([]*models.Domain, 100)
			for i := 0; i < 100; i++ {
				pageDomains[i] = &models.Domain{
					Name: fmt.Sprintf("example%d.com", page*100+i+1),
				}
			}

			pageResponse := &models.DomainsResponse{
				Domains: pageDomains,
				PaginationResponse: models.PaginationResponse{
					TotalItems: models.TotalItems{
						Value:    9950,
						Relation: "equal",
					},
				},
			}

			mockClient.On("MakeRequest", ctx, "POST", "/domains/search", pageRequest).
				Return(createMockResponse(pageResponse), nil)
		}

		// Last page with remaining 50 results
		lastPageRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100, // Will be limited to remaining 1 result
				Offset: 9900,
			},
		}

		lastPageDomains := make([]*models.Domain, 1)
		lastPageDomains[0] = &models.Domain{
			Name: "example9901.com",
		}

		lastPageResponse := &models.DomainsResponse{
			Domains: lastPageDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    9901,
					Relation: "equal",
				},
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", lastPageRequest).
			Return(createMockResponse(lastPageResponse), nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 9901)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty results", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "nonexistent.com",
		}

		expectedResponse := &models.DomainsResponse{
			Domains: []*models.Domain{},
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    0,
					Relation: "equal",
				},
			},
		}

		expectedRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", expectedRequest).Return(mockResponse, nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle search error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		expectedError := errors.New("search failed")
		expectedRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", expectedRequest).
			Return(nil, expectedError)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search domain records at offset 0")
		assert.Contains(t, err.Error(), "search failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle error on second page", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "example.com",
		}

		// First page succeeds
		firstPageDomains := make([]*models.Domain, 100)
		for i := 0; i < 100; i++ {
			firstPageDomains[i] = &models.Domain{
				Name: fmt.Sprintf("example%d.com", i+1),
			}
		}

		firstPageResponse := &models.DomainsResponse{
			Domains: firstPageDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    200,
					Relation: "equal",
				},
			},
		}

		firstRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		secondRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 100,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", firstRequest).
			Return(createMockResponse(firstPageResponse), nil)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", secondRequest).
			Return(nil, errors.New("second page failed"))

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search domain records at offset 100")
		assert.Contains(t, err.Error(), "second page failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		baseParams := models.Search{
			Query: "example.com",
		}

		expectedRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", expectedRequest).
			Return(nil, context.Canceled)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
		mockClient.AssertExpectations(t)
	})
}

// Benchmark tests
func BenchmarkDomainService_GetDetails(b *testing.B) {
	mockClient := &MockClient{}
	service := NewDomainService(mockClient)
	ctx := context.Background()

	expectedDomain := &models.Domain{
		Name: "example.com",
	}

	mockResponse := createMockResponse(expectedDomain)

	for i := 0; i < b.N; i++ {
		mockClient.On("MakeRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockResponse, nil).Once()
		_, _ = service.GetDetails(ctx, "test-id")
	}
}

func BenchmarkDomainService_Search(b *testing.B) {
	mockClient := &MockClient{}
	service := NewDomainService(mockClient)
	ctx := context.Background()

	searchParams := models.SearchRequest{
		Search: models.Search{Query: "example.com"},
	}

	expectedResponse := &models.DomainsResponse{
		Domains: []*models.Domain{
			{Name: "example.com"},
		},
		PaginationResponse: models.PaginationResponse{
			TotalItems: models.TotalItems{
				Value:    1,
				Relation: "equal",
			},
		},
	}

	mockResponse := createMockResponse(expectedResponse)

	for i := 0; i < b.N; i++ {
		mockClient.On("MakeRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockResponse, nil).Once()
		_, _ = service.Search(ctx, searchParams)
	}
}

// Integration-style tests (testing with real HTTP client behavior simulation)
func TestDomainService_Integration(t *testing.T) {
	t.Run("should handle response body close properly", func(t *testing.T) {
		// This test verifies that response bodies are properly closed
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		// Create a response that tracks if Body.Close() was called
		closeCalled := false
		mockBody := &mockReadCloser{
			data: strings.NewReader(`{"id":"test","name":"example.com"}`),
			closeFunc: func() error {
				closeCalled = true
				return nil
			},
		}

		response := &http.Response{
			StatusCode: 200,
			Body:       mockBody,
			Header:     make(http.Header),
		}

		mockClient.On("MakeRequest", ctx, "GET", "/domains/test-id", mock.Anything).Return(response, nil)

		// Act
		_, err := service.GetDetails(ctx, "test-id")

		// Assert
		assert.NoError(t, err)
		assert.True(t, closeCalled, "Response body should have been closed")
		mockClient.AssertExpectations(t)
	})
}

// mockReadCloser helps test that response bodies are properly closed
type mockReadCloser struct {
	data      io.Reader
	closeFunc func() error
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	return m.data.Read(p)
}

func (m *mockReadCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

// Test SearchAll with edge cases around the 10000 limit
func TestDomainService_SearchAll_EdgeCases(t *testing.T) { // nolint: funlen
	testCases := []struct {
		name            string
		setupMocks      func(*MockClient, context.Context, models.Search)
		expectedResults int
		expectedError   string
	}{
		{
			name: "exactly at max results limit",
			setupMocks: func(mc *MockClient, ctx context.Context, baseParams models.Search) {
				// Setup exactly 100 pages of 100 results each (10000 total)
				for page := 0; page < 100; page++ {
					pageRequest := models.SearchRequest{
						Search: baseParams,
						Pagination: models.Pagination{
							Limit:  100,
							Offset: page * 100,
						},
					}

					pageDomains := make([]*models.Domain, 100)
					for i := 0; i < 100; i++ {
						pageDomains[i] = &models.Domain{
							Name: fmt.Sprintf("example%d.com", page*100+i+1),
						}
					}

					pageResponse := &models.DomainsResponse{
						Domains: pageDomains,
						PaginationResponse: models.PaginationResponse{
							TotalItems: models.TotalItems{
								Value:    10000,
								Relation: "equal",
							},
						},
					}

					mc.On("MakeRequest", ctx, "POST", "/domains/search", pageRequest).
						Return(createMockResponse(pageResponse), nil)
				}
			},
			expectedResults: 10000,
			expectedError:   "",
		},
		{
			name: "partial final page due to max limit",
			setupMocks: func(mc *MockClient, ctx context.Context, baseParams models.Search) {
				// Setup 99 full pages, then a partial page that hits the limit
				for page := 0; page < 99; page++ {
					pageRequest := models.SearchRequest{
						Search: baseParams,
						Pagination: models.Pagination{
							Limit:  100,
							Offset: page * 100,
						},
					}

					pageDomains := make([]*models.Domain, 100)
					for i := 0; i < 100; i++ {
						pageDomains[i] = &models.Domain{
							Name: fmt.Sprintf("example%d.com", page*100+i+1),
						}
					}

					pageResponse := &models.DomainsResponse{
						Domains: pageDomains,
						PaginationResponse: models.PaginationResponse{
							TotalItems: models.TotalItems{
								Value:    15000,
								Relation: "equal",
							},
						},
					}

					mc.On("MakeRequest", ctx, "POST", "/domains/search", pageRequest).
						Return(createMockResponse(pageResponse), nil)
				}

				// Final page with only 100 results (9900 + 100 = 10000)
				finalPageRequest := models.SearchRequest{
					Search: baseParams,
					Pagination: models.Pagination{
						Limit:  100,
						Offset: 9900,
					},
				}

				finalPageDomains := make([]*models.Domain, 100)
				for i := 0; i < 100; i++ {
					finalPageDomains[i] = &models.Domain{
						Name: fmt.Sprintf("example%d.com", 9900+i+1),
					}
				}

				finalPageResponse := &models.DomainsResponse{
					Domains: finalPageDomains,
					PaginationResponse: models.PaginationResponse{
						TotalItems: models.TotalItems{
							Value:    15000,
							Relation: "equal",
						},
					},
				}

				mc.On("MakeRequest", ctx, "POST", "/domains/search", finalPageRequest).
					Return(createMockResponse(finalPageResponse), nil)
			},
			expectedResults: 10000,
			expectedError:   "",
		},
		{
			name: "single result under limit",
			setupMocks: func(mc *MockClient, ctx context.Context, baseParams models.Search) {
				pageRequest := models.SearchRequest{
					Search: baseParams,
					Pagination: models.Pagination{
						Limit:  100,
						Offset: 0,
					},
				}

				pageResponse := &models.DomainsResponse{
					Domains: []*models.Domain{
						{Name: "single.com"},
					},
					PaginationResponse: models.PaginationResponse{
						TotalItems: models.TotalItems{
							Value:    1,
							Relation: "equal",
						},
					},
				}

				mc.On("MakeRequest", ctx, "POST", "/domains/search", pageRequest).
					Return(createMockResponse(pageResponse), nil)
			},
			expectedResults: 1,
			expectedError:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockClient{}
			service := NewDomainService(mockClient)
			ctx := context.Background()

			baseParams := models.Search{
				Query: "example.com",
			}

			tc.setupMocks(mockClient, ctx, baseParams)

			result, err := service.SearchAll(ctx, baseParams)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result, tc.expectedResults)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// Test concurrent access (race condition testing)
func TestDomainService_Concurrency(t *testing.T) {
	t.Run("concurrent GetDetails calls", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)

		// Setup mock to handle multiple concurrent calls
		for i := 0; i < 10; i++ {
			domainID := fmt.Sprintf("example%d.com", i)
			expectedDomain := &models.Domain{
				Name: domainID,
			}
			mockResponse := createMockResponse(expectedDomain)
			mockClient.
				On("MakeRequest", mock.Anything, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
				Return(mockResponse, nil).
				Once() // ensure each expectation is called once
		}

		// Prepare slices for results
		results := make([]*models.Domain, 10)
		mockErrors := make([]error, 10)

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				ctx := context.Background()
				res, err := service.GetDetails(ctx, fmt.Sprintf("example%d.com", index))
				results[index] = res
				mockErrors[index] = err
			}(i)
		}

		wg.Wait() // wait for all goroutines to complete

		// Verify all calls succeeded
		for i := 0; i < 10; i++ {
			assert.NoError(t, mockErrors[i])
			assert.NotNil(t, results[i])
			if results[i] != nil {
				assert.Equal(t, fmt.Sprintf("example%d.com", i), results[i].Name)
			}
		}

		mockClient.AssertExpectations(t)
	})
}

// Test memory usage and resource cleanup
func TestDomainService_ResourceManagement(t *testing.T) {
	t.Run("should handle large responses without memory issues", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		// Create a large response (simulate 1000 domains)
		largeDomains := make([]*models.Domain, 1000)
		for i := 0; i < 1000; i++ {
			largeDomains[i] = &models.Domain{
				Name: fmt.Sprintf("large-example-%d.com", i),
			}
		}

		searchRequest := models.SearchRequest{
			Search:     models.Search{Query: "large"},
			Pagination: models.Pagination{Limit: 1000, Offset: 0},
		}

		largeResponse := &models.DomainsResponse{
			Domains: largeDomains,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    1000,
					Relation: "equal",
				},
			},
		}

		mockResponse := createMockResponse(largeResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/domains/search", searchRequest).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchRequest)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Domains, 1000)

		// Verify data integrity
		assert.Equal(t, "large-example-0.com", result.Domains[0].Name)
		assert.Equal(t, "large-example-999.com", result.Domains[999].Name)

		mockClient.AssertExpectations(t)
	})
}

// Test error handling for different HTTP status codes
func TestDomainService_HTTPStatusCodes(t *testing.T) { // nolint: funlen
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedError  string
		shouldHaveData bool
	}{
		{
			name:           "200 OK with valid data",
			statusCode:     200,
			responseBody:   `{"name":"example.com"}`,
			expectedError:  "",
			shouldHaveData: true,
		},
		{
			name:           "404 Not Found",
			statusCode:     404,
			responseBody:   `{"error":"Domain not found"}`,
			expectedError:  "",
			shouldHaveData: true, // The service doesn't check HTTP status, it just decodes JSON
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     500,
			responseBody:   `{"error":"Internal server error"}`,
			expectedError:  "",
			shouldHaveData: true, // The service doesn't check HTTP status, it just decodes JSON
		},
		{
			name:           "Empty response body",
			statusCode:     200,
			responseBody:   "",
			expectedError:  "failed to decode domain details response",
			shouldHaveData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockClient{}
			service := NewDomainService(mockClient)
			ctx := context.Background()
			domainID := "example.com"

			mockResponse := createMockResponseWithString(tc.statusCode, tc.responseBody)
			mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
				Return(mockResponse, nil)

			result, err := service.GetDetails(ctx, domainID)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tc.shouldHaveData {
					assert.NotNil(t, result)
				}
			}

			mockClient.AssertExpectations(t)
		})
	}
}

// Test nil pointer handling
func TestDomainService_NilPointerHandling(t *testing.T) {
	t.Run("should handle nil client gracefully", func(t *testing.T) {
		service := NewDomainService(nil)
		ctx := context.Background()

		// This should panic or return an error when trying to use nil client
		assert.Panics(t, func() {
			_, _ = service.GetDetails(ctx, "test-id")
		})
	})

	t.Run("should handle nil context", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)

		// The underlying client should handle nil context appropriately
		expectedError := errors.New("nil context")
		mockClient.On("MakeRequest", mock.Anything, "GET", "/domains/test-id", mock.Anything).
			Return(nil, expectedError)

		result, err := service.GetDetails(context.Background(), "test-id")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockClient.AssertExpectations(t)
	})
}

// Performance tests
func TestDomainService_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("should handle rapid sequential calls", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewDomainService(mockClient)
		ctx := context.Background()

		// Setup many mock responses
		for i := 0; i < 1000; i++ {
			domainID := fmt.Sprintf("perf-example%d.com", i)
			expectedDomain := &models.Domain{
				Name: domainID,
			}
			mockResponse := createMockResponse(expectedDomain)
			mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/domains/%s", domainID), mock.Anything).
				Return(mockResponse, nil)
		}

		// Measure time for 1000 sequential calls
		start := time.Now()
		for i := 0; i < 1000; i++ {
			_, err := service.GetDetails(ctx, fmt.Sprintf("perf-example%d.com", i))
			require.NoError(t, err)
		}
		duration := time.Since(start)

		t.Logf("1000 sequential GetDetails calls took: %v", duration)
		assert.Less(t, duration.Milliseconds(), int64(5000), "Should complete 1000 calls in less than 5 seconds")

		mockClient.AssertExpectations(t)
	})
}

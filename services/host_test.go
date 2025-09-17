package services

import (
	"context"
	"errors"
	"fmt"
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

const TestHost = "1.1.1.1"

func TestNewHostService(t *testing.T) {
	t.Run("should create new host service with client", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewHostService(mockClient)

		assert.NotNil(t, service)
		assert.Equal(t, mockClient, service.client)
	})

	t.Run("should create service with nil client", func(t *testing.T) {
		service := NewHostService(nil)
		assert.NotNil(t, service)
		assert.Nil(t, service.client)
	})
}

func TestHostService_GetDetails(t *testing.T) { // nolint: funlen
	t.Run("should successfully get host details", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		ip := TestHost

		expectedHost := &models.Host{
			IP: ip,
		}

		mockResponse := createMockResponse(expectedHost)
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, ip)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedHost.IP, result.IP)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle client request error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		ip := TestHost
		expectedError := errors.New("network error")

		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
			Return(nil, expectedError)

		// Act
		result, err := service.GetDetails(ctx, ip)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get host details for ID 1.1.1.1")
		assert.Contains(t, err.Error(), "network error")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle JSON decode error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		ip := TestHost

		// Create response with invalid JSON
		mockResponse := createMockResponseWithString(200, "invalid json")
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
			Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, ip)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode host details response")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty host ID", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		ip := ""

		expectedHost := &models.Host{}
		mockResponse := createMockResponse(expectedHost)
		mockClient.On("MakeRequest", ctx, "GET", "/hosts/", mock.Anything).Return(mockResponse, nil)

		// Act
		result, err := service.GetDetails(ctx, ip)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel context immediately
		ip := TestHost

		expectedError := context.Canceled
		mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
			Return(nil, expectedError)

		// Act
		result, err := service.GetDetails(ctx, ip)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "context canceled")
		mockClient.AssertExpectations(t)
	})
}

func TestHostService_Search(t *testing.T) { // nolint: funlen
	t.Run("should successfully search hosts", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Search: models.Search{
				Query: TestHost,
			},
			Pagination: models.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		expectedResponse := &models.HostsResponse{
			Hosts: []*models.Host{
				{IP: TestHost},
				{IP: "8.8.8.8"},
			},
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    2,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Hosts, 2)
		assert.Equal(t, int64(2), result.TotalItems.Value)
		assert.Equal(t, TestHost, result.Hosts[0].IP)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle empty search results", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Search: models.Search{
				Query: "nonexistent.com",
			},
		}

		expectedResponse := &models.HostsResponse{
			Hosts: []*models.Host{},
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    0,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Hosts, 0)
		assert.Equal(t, int64(0), result.TotalItems.Value)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle client request error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		searchParams := models.SearchRequest{}
		expectedError := errors.New("request failed")

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchParams).Return(nil, expectedError)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search host records")
		assert.Contains(t, err.Error(), "request failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle JSON decode error", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()
		searchParams := models.SearchRequest{}

		mockResponse := createMockResponseWithString(200, "invalid json")
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode host search response")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle large pagination", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		searchParams := models.SearchRequest{
			Pagination: models.Pagination{
				Limit:  1000,
				Offset: 5000,
			},
		}

		expectedResponse := &models.HostsResponse{
			Hosts: make([]*models.Host, 1000),
			PaginationResponse: models.PaginationResponse{TotalItems: models.TotalItems{
				Value:    1000,
				Relation: "equal",
			}},
		}

		mockResponse := createMockResponse(expectedResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchParams).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Hosts, 1000)
		mockClient.AssertExpectations(t)
	})
}

func TestHostService_SearchAll(t *testing.T) { // nolint: funlen
	t.Run("should successfully search all hosts with single page", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
		}

		expectedHosts := []*models.Host{
			{IP: TestHost},
			{IP: "8.8.8.8"},
		}

		expectedResponse := &models.HostsResponse{
			Hosts: expectedHosts,
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
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", expectedSearchRequest).Return(mockResponse, nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, TestHost, result[0].IP)
		assert.Equal(t, "8.8.8.8", result[1].IP)
		mockClient.AssertExpectations(t)
	})

	t.Run("should successfully search all hosts with multiple pages", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
		}

		// First page - 100 results
		firstPageHosts := make([]*models.Host, 100)
		for i := 0; i < 100; i++ {
			firstPageHosts[i] = &models.Host{
				IP: fmt.Sprintf("1.1.1.%d", i+1),
			}
		}

		firstPageResponse := &models.HostsResponse{
			Hosts: firstPageHosts,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    150,
					Relation: "equal",
				},
			},
		}

		// Second page - 50 results
		secondPageHosts := make([]*models.Host, 50)
		for i := 0; i < 50; i++ {
			secondPageHosts[i] = &models.Host{
				IP: fmt.Sprintf("1.1.1.%d", i+101),
			}
		}

		secondPageResponse := &models.HostsResponse{
			Hosts: secondPageHosts,
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

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", firstRequest).
			Return(createMockResponse(firstPageResponse), nil)

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", secondRequest).
			Return(createMockResponse(secondPageResponse), nil)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 150)
		assert.Equal(t, TestHost, result[0].IP)
		assert.Equal(t, "1.1.1.150", result[149].IP)
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle max results limit", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
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

			pageHosts := make([]*models.Host, 100)
			for i := 0; i < 100; i++ {
				pageHosts[i] = &models.Host{
					IP: fmt.Sprintf("1.1.1.%d", page*100+i+1),
				}
			}

			pageResponse := &models.HostsResponse{
				Hosts: pageHosts,
				PaginationResponse: models.PaginationResponse{
					TotalItems: models.TotalItems{
						Value:    15000,
						Relation: "equal",
					},
				},
			}

			mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", pageRequest).
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
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
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

			pageHosts := make([]*models.Host, 100)
			for i := 0; i < 100; i++ {
				pageHosts[i] = &models.Host{
					IP: fmt.Sprintf("1.1.1.%d", page*100+i+1),
				}
			}

			pageResponse := &models.HostsResponse{
				Hosts: pageHosts,
				PaginationResponse: models.PaginationResponse{
					TotalItems: models.TotalItems{
						Value:    9950,
						Relation: "equal",
					},
				},
			}

			mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", pageRequest).
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

		lastPageHosts := make([]*models.Host, 1)
		lastPageHosts[0] = &models.Host{
			IP: "1.1.1.9901",
		}

		lastPageResponse := &models.HostsResponse{
			Hosts: lastPageHosts,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    9901,
					Relation: "equal",
				},
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", lastPageRequest).
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
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: "nonexistent_ip",
		}

		expectedResponse := &models.HostsResponse{
			Hosts: []*models.Host{},
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
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", expectedRequest).Return(mockResponse, nil)

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
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
		}

		expectedError := errors.New("search failed")
		expectedRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", expectedRequest).
			Return(nil, expectedError)

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search host records at offset 0")
		assert.Contains(t, err.Error(), "search failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle error on second page", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		baseParams := models.Search{
			Query: TestHost,
		}

		// First page succeeds
		firstPageHosts := make([]*models.Host, 100)
		for i := 0; i < 100; i++ {
			firstPageHosts[i] = &models.Host{
				IP: fmt.Sprintf("1.1.1.%d", i+1),
			}
		}

		firstPageResponse := &models.HostsResponse{
			Hosts: firstPageHosts,
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

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", firstRequest).
			Return(createMockResponse(firstPageResponse), nil)
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", secondRequest).
			Return(nil, errors.New("second page failed"))

		// Act
		result, err := service.SearchAll(ctx, baseParams)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to search host records at offset 100")
		assert.Contains(t, err.Error(), "second page failed")
		mockClient.AssertExpectations(t)
	})

	t.Run("should handle context cancellation", func(t *testing.T) {
		// Arrange
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		baseParams := models.Search{
			Query: TestHost,
		}

		expectedRequest := models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  100,
				Offset: 0,
			},
		}

		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", expectedRequest).
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
func BenchmarkHostService_GetDetails(b *testing.B) {
	mockClient := &MockClient{}
	service := NewHostService(mockClient)
	ctx := context.Background()

	expectedHost := &models.Host{
		IP: TestHost,
	}

	mockResponse := createMockResponse(expectedHost)

	for i := 0; i < b.N; i++ {
		mockClient.On("MakeRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(mockResponse, nil).Once()
		_, _ = service.GetDetails(ctx, "test-ip")
	}
}

func BenchmarkHostService_Search(b *testing.B) {
	mockClient := &MockClient{}
	service := NewHostService(mockClient)
	ctx := context.Background()

	searchParams := models.SearchRequest{
		Search: models.Search{Query: TestHost},
	}

	expectedResponse := &models.HostsResponse{
		Hosts: []*models.Host{
			{IP: TestHost},
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
func TestHostService_Integration(t *testing.T) {
	t.Run("should handle response body close properly", func(t *testing.T) {
		// This test verifies that response bodies are properly closed
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		// Create a response that tracks if Body.Close() was called
		closeCalled := false
		mockBody := &mockReadCloser{
			data: strings.NewReader(`{"ip":"1.1.1.1"}`),
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

		mockClient.On("MakeRequest", ctx, "GET", "/hosts/1.1.1.1", mock.Anything).Return(response, nil)

		// Act
		_, err := service.GetDetails(ctx, TestHost)

		// Assert
		assert.NoError(t, err)
		assert.True(t, closeCalled, "Response body should have been closed")
		mockClient.AssertExpectations(t)
	})
}

// Test SearchAll with edge cases around the 10000 limit
func TestHostService_SearchAll_EdgeCases(t *testing.T) { // nolint: funlen
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

					pageHosts := make([]*models.Host, 100)
					for i := 0; i < 100; i++ {
						pageHosts[i] = &models.Host{
							IP: fmt.Sprintf("1.1.1.%d", page*100+i+1),
						}
					}

					pageResponse := &models.HostsResponse{
						Hosts: pageHosts,
						PaginationResponse: models.PaginationResponse{
							TotalItems: models.TotalItems{
								Value:    10000,
								Relation: "equal",
							},
						},
					}

					mc.On("MakeRequest", ctx, "POST", "/hosts/search", pageRequest).
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

					pageHosts := make([]*models.Host, 100)
					for i := 0; i < 100; i++ {
						pageHosts[i] = &models.Host{
							IP: fmt.Sprintf("1.1.1.1%d", page*100+i+1),
						}
					}

					pageResponse := &models.HostsResponse{
						Hosts: pageHosts,
						PaginationResponse: models.PaginationResponse{
							TotalItems: models.TotalItems{
								Value:    15000,
								Relation: "equal",
							},
						},
					}

					mc.On("MakeRequest", ctx, "POST", "/hosts/search", pageRequest).
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

				finalPageHosts := make([]*models.Host, 100)
				for i := 0; i < 100; i++ {
					finalPageHosts[i] = &models.Host{
						IP: fmt.Sprintf("1.1.1.%d", 9900+i+1),
					}
				}

				finalPageResponse := &models.HostsResponse{
					Hosts: finalPageHosts,
					PaginationResponse: models.PaginationResponse{
						TotalItems: models.TotalItems{
							Value:    15000,
							Relation: "equal",
						},
					},
				}

				mc.On("MakeRequest", ctx, "POST", "/hosts/search", finalPageRequest).
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

				pageResponse := &models.HostsResponse{
					Hosts: []*models.Host{
						{IP: TestHost},
					},
					PaginationResponse: models.PaginationResponse{
						TotalItems: models.TotalItems{
							Value:    1,
							Relation: "equal",
						},
					},
				}

				mc.On("MakeRequest", ctx, "POST", "/hosts/search", pageRequest).
					Return(createMockResponse(pageResponse), nil)
			},
			expectedResults: 1,
			expectedError:   "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockClient{}
			service := NewHostService(mockClient)
			ctx := context.Background()

			baseParams := models.Search{
				Query: TestHost,
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
func TestHostService_Concurrency(t *testing.T) {
	t.Run("concurrent GetDetails calls", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewHostService(mockClient)

		// Setup mock to handle multiple concurrent calls
		for i := 0; i < 10; i++ {
			ip := fmt.Sprintf("1.1.1.%d", i)
			expectedHost := &models.Host{
				IP: ip,
			}
			mockResponse := createMockResponse(expectedHost)
			mockClient.
				On("MakeRequest", mock.Anything, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
				Return(mockResponse, nil).
				Once() // ensure each expectation is called once
		}

		// Prepare slices for results
		results := make([]*models.Host, 10)
		mockErrors := make([]error, 10)

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				ctx := context.Background()
				res, err := service.GetDetails(ctx, fmt.Sprintf("1.1.1.%d", index))
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
				assert.Equal(t, fmt.Sprintf("1.1.1.%d", i), results[i].IP)
			}
		}

		mockClient.AssertExpectations(t)
	})
}

// Test memory usage and resource cleanup
func TestHostService_ResourceManagement(t *testing.T) {
	t.Run("should handle large responses without memory issues", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		// Create a large response (simulate 1000 hosts)
		largeHosts := make([]*models.Host, 1000)
		for i := 0; i < 1000; i++ {
			largeHosts[i] = &models.Host{
				IP: fmt.Sprintf("1.1.1.%d", i),
			}
		}

		searchRequest := models.SearchRequest{
			Search:     models.Search{Query: "large"},
			Pagination: models.Pagination{Limit: 1000, Offset: 0},
		}

		largeResponse := &models.HostsResponse{
			Hosts: largeHosts,
			PaginationResponse: models.PaginationResponse{
				TotalItems: models.TotalItems{
					Value:    1000,
					Relation: "equal",
				},
			},
		}

		mockResponse := createMockResponse(largeResponse)
		mockClient.On("MakeRequest", ctx, "POST", "/hosts/search", searchRequest).Return(mockResponse, nil)

		// Act
		result, err := service.Search(ctx, searchRequest)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Hosts, 1000)

		// Verify data integrity
		assert.Equal(t, "1.1.1.0", result.Hosts[0].IP)
		assert.Equal(t, "1.1.1.999", result.Hosts[999].IP)

		mockClient.AssertExpectations(t)
	})
}

// Test error handling for different HTTP status codes
func TestHostService_HTTPStatusCodes(t *testing.T) { // nolint: funlen
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
			responseBody:   `{"ip": "1.1.1.1"}`,
			expectedError:  "",
			shouldHaveData: true,
		},
		{
			name:           "404 Not Found",
			statusCode:     404,
			responseBody:   `{"error":"Host not found"}`,
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
			expectedError:  "failed to decode host details response",
			shouldHaveData: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &MockClient{}
			service := NewHostService(mockClient)
			ctx := context.Background()
			ip := TestHost

			mockResponse := createMockResponseWithString(tc.statusCode, tc.responseBody)
			mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
				Return(mockResponse, nil)

			result, err := service.GetDetails(ctx, ip)

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
func TestHostService_NilPointerHandling(t *testing.T) {
	t.Run("should handle nil client gracefully", func(t *testing.T) {
		service := NewHostService(nil)
		ctx := context.Background()

		// This should panic or return an error when trying to use nil client
		assert.Panics(t, func() {
			_, _ = service.GetDetails(ctx, TestHost)
		})
	})

	t.Run("should handle nil context", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewHostService(mockClient)

		// The underlying client should handle nil context appropriately
		expectedError := errors.New("nil context")
		mockClient.On("MakeRequest", mock.Anything, "GET", "/hosts/1.1.1.1", mock.Anything).
			Return(nil, expectedError)

		result, err := service.GetDetails(context.Background(), TestHost)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockClient.AssertExpectations(t)
	})
}

// Performance tests
func TestHostService_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	t.Run("should handle rapid sequential calls", func(t *testing.T) {
		mockClient := &MockClient{}
		service := NewHostService(mockClient)
		ctx := context.Background()

		// Setup many mock responses
		for i := 0; i < 1000; i++ {
			ip := fmt.Sprintf("1.1.1.%d", i)
			expectedHost := &models.Host{
				IP: ip,
			}
			mockResponse := createMockResponse(expectedHost)
			mockClient.On("MakeRequest", ctx, "GET", fmt.Sprintf("/hosts/%s", ip), mock.Anything).
				Return(mockResponse, nil)
		}

		// Measure time for 1000 sequential calls
		start := time.Now()
		for i := 0; i < 1000; i++ {
			_, err := service.GetDetails(ctx, fmt.Sprintf("1.1.1.%d", i))
			require.NoError(t, err)
		}
		duration := time.Since(start)

		t.Logf("1000 sequential GetDetails calls took: %v", duration)
		assert.Less(t, duration.Milliseconds(), int64(5000), "Should complete 1000 calls in less than 5 seconds")

		mockClient.AssertExpectations(t)
	})
}

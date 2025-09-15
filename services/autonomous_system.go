package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

// ASService handles Autonomous System (AS) operations for the Recona API.
// It provides methods to retrieve AS details, search for AS records, and perform bulk searches.
type ASService struct {
	client internal.Client
}

// NewASService creates a new instance of ASService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewASService(client internal.Client) *ASService {
	return &ASService{client: client}
}

// GetDetails retrieves detailed information for a specific Autonomous System by its number.
// Note: This method returns a *models.Host, which may be incorrect - consider if it should return *models.AS instead.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - number: The AS number as a string (e.g., "12345")
//
// Returns:
//   - *models.Host: The AS details (type may need to be reviewed)
//   - error: Any error that occurred during the request or response parsing
func (s *ASService) GetDetails(ctx context.Context, number string) (*models.Host, error) {
	// Make GET request to retrieve AS details by number
	resp, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/autonomous-system/%s", number), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get AS details for number %s: %w", number, err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize host variable to hold the decoded response
	var host *models.Host

	// Decode the JSON response into the host struct
	if err = internal.DecodeJSON(resp.Body, &host); err != nil {
		return nil, fmt.Errorf("failed to decode AS details response: %w", err)
	}

	return host, nil
}

// Search performs a search for Autonomous Systems based on the provided search parameters.
// It returns paginated results according to the pagination settings in the request.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: Search request containing search criteria and pagination settings
//
// Returns:
//   - *models.ASResponse: The search results with matching AS records
//   - error: Any error that occurred during the request or response parsing
func (s *ASService) Search(ctx context.Context, params models.SearchRequest) (*models.ASResponse, error) {
	// Make POST request to search for AS records
	resp, err := s.client.MakeRequest(ctx, "POST", "/autonomous-system/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search AS records: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.ASResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode AS search response: %w", err)
	}

	return result, nil
}

// SearchAll performs a comprehensive search that retrieves all matching AS records by paginating through results.
// It automatically handles pagination to collect up to maxResults records, making multiple API calls as needed.
// This method is useful when you need to retrieve all matching records without manual pagination handling.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - baseParams: Base search criteria to apply across all paginated requests
//
// Returns:
//   - []*models.AS: A slice containing all matching AS records from all pages
//   - error: Any error that occurred during the search process
func (s *ASService) SearchAll(ctx context.Context, baseParams models.Search) ([]*models.AS, error) {
	const (
		pageSize   = 100   // Number of records to fetch per API call
		maxResults = 10000 // Maximum total records to retrieve (safety limit)
	)

	offset := 0            // Current offset for pagination
	var allAS []*models.AS // Accumulator for all AS records
	limit := pageSize      // Current page size limit

	// Continue fetching until we reach maxResults or no more data is available
	for offset < maxResults {
		// Calculate remaining slots to avoid exceeding maxResults
		remaining := maxResults - offset
		if remaining < pageSize {
			limit = remaining
		}

		// Perform search with current pagination settings
		resp, err := s.Search(ctx, models.SearchRequest{
			Search: baseParams,
			Pagination: models.Pagination{
				Limit:  limit,
				Offset: offset,
			},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search AS records at offset %d: %w", offset, err)
		}

		// Break if no results returned (end of data)
		if len(resp.AutonomousSystems) == 0 {
			break
		}

		// Append current page results to our collection
		allAS = append(allAS, resp.AutonomousSystems...)

		// Update offset for next iteration
		offset += len(resp.AutonomousSystems)

		// Break if we received fewer results than requested (likely last page)
		if len(resp.AutonomousSystems) < limit {
			break
		}
	}

	return allAS, nil
}

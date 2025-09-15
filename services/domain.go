package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

// DomainService handles domain-related operations for the Recona API.
// It provides methods to retrieve domain details, search for domains, and perform bulk searches.
type DomainService struct {
	client internal.Client
}

// NewDomainService creates a new instance of DomainService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewDomainService(client internal.Client) *DomainService {
	return &DomainService{client: client}
}

// GetDetails retrieves detailed information for a specific domain by its ID.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - id: The domain ID as a string
//
// Returns:
//   - *models.Domain: The domain details
//   - error: Any error that occurred during the request or response parsing
func (s *DomainService) GetDetails(ctx context.Context, id string) (*models.Domain, error) {
	// Make GET request to retrieve domain details by ID
	resp, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/domains/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain details for ID %s: %w", id, err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize domain variable to hold the decoded response
	var domain *models.Domain

	// Decode the JSON response into the domain struct
	if err = internal.DecodeJSON(resp.Body, &domain); err != nil {
		return nil, fmt.Errorf("failed to decode domain details response: %w", err)
	}

	return domain, nil
}

// Search performs a search for domains based on the provided search parameters.
// It returns paginated results according to the pagination settings in the request.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: Search request containing search criteria and pagination settings
//
// Returns:
//   - *models.DomainsResponse: The search results with matching domain records
//   - error: Any error that occurred during the request or response parsing
// All possible search parameters can be found here: https://reconatest.io/docs/domain-filters
func (s *DomainService) Search(ctx context.Context, params models.SearchRequest) (*models.DomainsResponse, error) {
	// Make POST request to search for domain records
	resp, err := s.client.MakeRequest(ctx, "POST", "/domains/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search domain records: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.DomainsResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode domain search response: %w", err)
	}

	return result, nil
}

// SearchAll performs a comprehensive search that retrieves all matching domain records by paginating through results.
// It automatically handles pagination to collect up to maxResults records, making multiple API calls as needed.
// This method is useful when you need to retrieve all matching domains without manual pagination handling.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - baseParams: Base search criteria to apply across all paginated requests
//
// Returns:
//   - []*models.Domain: A slice containing all matching domain records from all pages
//   - error: Any error that occurred during the search process
// All possible search parameters can be found here: https://reconatest.io/docs/domain-filters
func (s *DomainService) SearchAll(ctx context.Context, baseParams models.Search) ([]*models.Domain, error) {
	const (
		pageSize   = 100   // Number of records to fetch per API call
		maxResults = 10000 // Maximum total records to retrieve (safety limit)
	)

	offset := 0                     // Current offset for pagination
	var allDomains []*models.Domain // Accumulator for all domain records
	limit := pageSize               // Current page size limit

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
			return nil, fmt.Errorf("failed to search domain records at offset %d: %w", offset, err)
		}

		// Break if no results returned (end of data)
		if len(resp.Domains) == 0 {
			break
		}

		// Append current page results to our collection
		allDomains = append(allDomains, resp.Domains...)

		// Update offset for next iteration
		offset += len(resp.Domains)

		// Break if we received fewer results than requested (likely last page)
		if len(resp.Domains) < limit {
			break
		}
	}

	return allDomains, nil
}

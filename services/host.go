package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

// HostService handles host-related operations for the Recona API.
// It provides methods to retrieve host details, search for hosts, and perform bulk searches.
// Hosts typically represent network endpoints, servers, or devices with associated IP addresses.
type HostService struct {
	client internal.Client
}

// NewHostService creates a new instance of HostService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewHostService(client internal.Client) *HostService {
	return &HostService{client: client}
}

// GetDetails retrieves detailed information for a specific host by its ID.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - id: The host ID as a string
//
// Returns:
//   - *models.Host: The host details including IP addresses, services, and metadata
//   - error: Any error that occurred during the request or response parsing
func (s *HostService) GetDetails(ctx context.Context, id string) (*models.Host, error) {
	// Make GET request to retrieve host details by ID
	resp, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/hosts/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get host details for ID %s: %w", id, err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize host variable to hold the decoded response
	var host *models.Host

	// Decode the JSON response into the host struct
	if err = internal.DecodeJSON(resp.Body, &host); err != nil {
		return nil, fmt.Errorf("failed to decode host details response: %w", err)
	}

	return host, nil
}

// Search performs a search for hosts based on the provided search parameters.
// It returns paginated results according to the pagination settings in the request.
// This is useful for finding hosts that match specific criteria like IP ranges, services, or metadata.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: Search request containing search criteria and pagination settings
//
// Returns:
//   - *models.HostsResponse: The search results with matching host records
//   - error: Any error that occurred during the request or response parsing
// All possible search parameters can be found here: https://reconatest.io/docs/ip-filters
func (s *HostService) Search(ctx context.Context, params models.SearchRequest) (*models.HostsResponse, error) {
	// Make POST request to search for host records
	resp, err := s.client.MakeRequest(ctx, "POST", "/hosts/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search host records: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.HostsResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode host search response: %w", err)
	}

	return result, nil
}

// SearchAll performs a comprehensive search that retrieves all matching host records by paginating through results.
// It automatically handles pagination to collect up to maxResults records, making multiple API calls as needed.
// This method is useful when you need to retrieve all matching hosts without manual pagination handling.
//
// Warning: Use with caution as this method can potentially retrieve large amounts of data.
// Consider using Search() with manual pagination for better control over resource usage.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - baseParams: Base search criteria to apply across all paginated requests
//
// Returns:
//   - []*models.Host: A slice containing all matching host records from all pages
//   - error: Any error that occurred during the search process
// All possible search parameters can be found here: https://reconatest.io/docs/ip-filters
func (s *HostService) SearchAll(ctx context.Context, baseParams models.Search) ([]*models.Host, error) {
	const (
		pageSize   = 100   // Number of records to fetch per API call
		maxResults = 10000 // Maximum total records to retrieve (safety limit)
	)

	offset := 0                 // Current offset for pagination
	var allHosts []*models.Host // Accumulator for all host records
	limit := pageSize           // Current page size limit

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
			return nil, fmt.Errorf("failed to search host records at offset %d: %w", offset, err)
		}

		// Break if no results returned (end of data)
		if len(resp.Hosts) == 0 {
			break
		}

		// Append current page results to our collection
		allHosts = append(allHosts, resp.Hosts...)

		// Update offset for next iteration
		offset += len(resp.Hosts)

		// Break if we received fewer results than requested (likely last page)
		if len(resp.Hosts) < limit {
			break
		}
	}

	return allHosts, nil
}

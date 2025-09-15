package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

// CVEService handles Common Vulnerabilities and Exposures (CVE) operations for the Recona API.
// It provides methods to retrieve CVE details, search for vulnerabilities, and access CWE (Common Weakness
// Enumeration) data.
// CVEs are standardized identifiers for publicly disclosed cybersecurity vulnerabilities.
type CVEService struct {
	client internal.Client
}

// NewCVEService creates a new instance of CVEService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewCVEService(client internal.Client) *CVEService {
	return &CVEService{client: client}
}

// GetDetails retrieves detailed information for a specific CVE by its ID.
// CVE IDs typically follow the format "CVE-YYYY-NNNNN" (e.g., "CVE-2021-44228").
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - id: The CVE ID as a string (e.g., "CVE-2021-44228")
//
// Returns:
//   - *models.CVE: The CVE details including description, severity, affected products, etc.
//   - error: Any error that occurred during the request or response parsing
func (s *CVEService) GetDetails(ctx context.Context, id string) (*models.CVE, error) {
	// Make GET request to retrieve CVE details by ID
	resp, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/cve/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get CVE details for ID %s: %w", id, err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize cve variable to hold the decoded response
	var cve *models.CVE

	// Decode the JSON response into the CVE struct
	if err = internal.DecodeJSON(resp.Body, &cve); err != nil {
		return nil, fmt.Errorf("failed to decode CVE details response: %w", err)
	}

	return cve, nil
}

// Search performs a search for CVE records based on the provided search parameters.
// It returns paginated results according to the pagination settings in the request.
// This is useful for finding vulnerabilities that match specific criteria like severity, date range,
// or affected products.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: Search request containing search criteria and pagination settings
//
// Returns:
//   - *models.CVEResponse: The search results with matching CVE records
//   - error: Any error that occurred during the request or response parsing
// All possible search parameters can be found here: https://recona.io/docs/cve-filters
func (s *CVEService) Search(ctx context.Context, params models.SearchRequest) (*models.CVEResponse, error) {
	// Make POST request to search for CVE records
	resp, err := s.client.MakeRequest(ctx, "POST", "/cve/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search CVE records: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.CVEResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode CVE search response: %w", err)
	}

	return result, nil
}

// GetCWE retrieves Common Weakness Enumeration (CWE) data based on the provided parameters.
// CWE is a category system for hardware and software weaknesses and vulnerabilities.
// This method is useful for understanding the underlying weakness types associated with vulnerabilities.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: CWE parameters specifying which weakness data to retrieve
//
// Returns:
//   - *models.CWEResponse: The CWE data including weakness descriptions and classifications
//   - error: Any error that occurred during the request or response parsing
func (s *CVEService) GetCWE(ctx context.Context, params models.CWEParams) (*models.CWEResponse, error) {
	// Make POST request to retrieve CWE data
	resp, err := s.client.MakeRequest(ctx, "POST", "/cwe", params)
	if err != nil {
		return nil, fmt.Errorf("failed to get CWE data: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.CWEResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode CWE response: %w", err)
	}

	return result, nil
}

// SearchAll performs a comprehensive search that retrieves all matching CVE records by paginating through results.
// It automatically handles pagination to collect up to maxResults records, making multiple API calls as needed.
// This method is useful when you need to retrieve all matching vulnerabilities without manual pagination handling.
//
// Warning: CVE databases can be very large. Use with caution and consider filtering your search criteria
// to avoid retrieving excessive amounts of data. Consider using Search() with manual pagination for better control.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - baseParams: Base search criteria to apply across all paginated requests
//
// Returns:
//   - []*models.NistCVEData: A slice containing all matching NIST CVE data records from all pages
//   - error: Any error that occurred during the search process
// All possible search parameters can be found here: https://recona.io/docs/cve-filters
func (s *CVEService) SearchAll(ctx context.Context, baseParams models.Search) ([]*models.NistCVEData, error) {
	const (
		pageSize   = 100   // Number of records to fetch per API call
		maxResults = 10000 // Maximum total records to retrieve (safety limit)
	)

	offset := 0                      // Current offset for pagination
	var allCVE []*models.NistCVEData // Accumulator for all CVE records
	limit := pageSize                // Current page size limit

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
			return nil, fmt.Errorf("failed to search CVE records at offset %d: %w", offset, err)
		}

		// Break if no results returned (end of data)
		if len(resp.CVEList) == 0 {
			break
		}

		// Append current page results to our collection
		allCVE = append(allCVE, resp.CVEList...)

		// Update offset for next iteration
		offset += len(resp.CVEList)

		// Break if we received fewer results than requested (likely last page)
		if len(resp.CVEList) < limit {
			break
		}
	}

	return allCVE, nil
}

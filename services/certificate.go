package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

// CertificateService handles SSL/TLS certificate operations for the Recona API.
// It provides methods to retrieve certificate details, search for certificates, and perform bulk searches.
type CertificateService struct {
	client internal.Client
}

// NewCertificateService creates a new instance of CertificateService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewCertificateService(client internal.Client) *CertificateService {
	return &CertificateService{client: client}
}

// GetDetails retrieves detailed information for a specific certificate by its ID.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - id: The certificate ID as a string
//
// Returns:
//   - *models.Certificate: The certificate details
//   - error: Any error that occurred during the request or response parsing
func (s *CertificateService) GetDetails(ctx context.Context, id string) (*models.Certificate, error) {
	// Make GET request to retrieve certificate details by ID
	resp, err := s.client.MakeRequest(ctx, "GET", fmt.Sprintf("/certificates/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate details for ID %s: %w", id, err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize certificate variable to hold the decoded response
	var certificate *models.Certificate

	// Decode the JSON response into the certificate struct
	if err = internal.DecodeJSON(resp.Body, &certificate); err != nil {
		return nil, fmt.Errorf("failed to decode certificate details response: %w", err)
	}

	return certificate, nil
}

// Search performs a search for certificates based on the provided search parameters.
// It returns paginated results according to the pagination settings in the request.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - params: Search request containing search criteria and pagination settings
//
// Returns:
//   - *models.CertificatesResponse: The search results with matching certificate records
//   - error: Any error that occurred during the request or response parsing
// All possible search parameters can be found here: https://reconatest.io/docs/certificate-filters
func (s *CertificateService) Search(
	ctx context.Context, params models.SearchRequest) (*models.CertificatesResponse, error) {

	// Make POST request to search for certificate records
	resp, err := s.client.MakeRequest(ctx, "POST", "/certificates/search", params)
	if err != nil {
		return nil, fmt.Errorf("failed to search certificate records: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize result variable to hold the decoded response
	var result *models.CertificatesResponse

	// Decode the JSON response into the result struct
	if err = internal.DecodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode certificate search response: %w", err)
	}

	return result, nil
}

// SearchAll performs a comprehensive search that retrieves all matching certificate records by paginating
// through results.
// It automatically handles pagination to collect up to maxResults records, making multiple API calls as needed.
// This method is useful when you need to retrieve all matching certificates without manual pagination handling.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//   - baseParams: Base search criteria to apply across all paginated requests
//
// Returns:
//   - []*models.Certificate: A slice containing all matching certificate records from all pages
//   - error: Any error that occurred during the search process
// All possible search parameters can be found here: https://reconatest.io/docs/certificate-filters
func (s *CertificateService) SearchAll(ctx context.Context, baseParams models.Search) ([]*models.Certificate, error) {
	const (
		pageSize   = 100   // Number of records to fetch per API call
		maxResults = 10000 // Maximum total records to retrieve (safety limit)
	)

	offset := 0                               // Current offset for pagination
	var allCertificates []*models.Certificate // Accumulator for all certificate records
	limit := pageSize                         // Current page size limit

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
			return nil, fmt.Errorf("failed to search certificate records at offset %d: %w", offset, err)
		}

		// Break if no results returned (end of data)
		if len(resp.Certificates) == 0 {
			break
		}

		// Append current page results to our collection
		allCertificates = append(allCertificates, resp.Certificates...)

		// Update offset for next iteration
		offset += len(resp.Certificates)

		// Break if we received fewer results than requested (likely last page)
		if len(resp.Certificates) < limit {
			break
		}
	}

	return allCertificates, nil
}

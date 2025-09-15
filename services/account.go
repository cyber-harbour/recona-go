package services

import (
	"context"
	"fmt"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/models"
)

const (
	// accountEndpoint defines the API endpoint for account-related operations
	accountEndpoint = "/customers/account"
)

// AccountService handles account operations for the Recona API.
// It provides methods to interact with user account data and profile information.
type AccountService struct {
	client internal.Client
}

// NewAccountService creates a new instance of AccountService with the provided client.
// The client parameter should implement the internal.Client interface for making HTTP requests.
func NewAccountService(c internal.Client) *AccountService {
	return &AccountService{
		client: c,
	}
}

// GetDetails retrieves the account profile details for the authenticated user.
// It returns the user's profile information or an error if the request fails.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout control
//
// Returns:
//   - *models.Profile: The user's profile information
//   - error: Any error that occurred during the request or response parsing
func (s *AccountService) GetDetails(ctx context.Context) (*models.Profile, error) {
	// Make GET request to the account endpoint
	resp, err := s.client.MakeRequest(ctx, "GET", accountEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to account endpoint: %w", err)
	}

	// Ensure response body is always closed to prevent resource leaks
	defer func() {
		_ = resp.Body.Close()
	}()

	// Initialize profile variable to hold the decoded response
	var profile *models.Profile

	// Decode the JSON response into the profile struct
	if err = internal.DecodeJSON(resp.Body, &profile); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return profile, nil
}

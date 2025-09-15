package reconago

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cyber-harbour/recona-go/internal"
	"github.com/cyber-harbour/recona-go/services"

	"golang.org/x/time/rate"
)

// Client provides a centralized HTTP client for interacting with the API.
// It manages authentication, rate limiting, and service-specific endpoints.
type Client struct {
	// Core HTTP configuration
	baseURL    string       // Base URL for all API endpoints
	token      string       // Authentication token for API requests
	httpClient *http.Client // Underlying HTTP client with timeout configuration

	// Rate limiting
	rateLimiter *rate.Limiter // Token bucket rate limiter for request throttling

	// Service endpoints - each service handles specific resource types
	Domain      *services.DomainService      // Domain analysis and WHOIS operations
	Host        *services.HostService        // Host scanning and port analysis
	Certificate *services.CertificateService // SSL/TLS certificate operations
	CVE         *services.CVEService         // Vulnerability and CVE data operations
}

// ClientOptions holds configuration options for creating a new client
type ClientOptions struct {
	Timeout        time.Duration // HTTP request timeout (default: 60s)
	RequestsPerSec float64       // Rate limit in requests per second (default: 10)
	BurstSize      int           // Maximum burst size for rate limiter (default: 20)
}

// NewClient creates a new API client with the provided authentication token.
// It initializes all service endpoints and sets up rate limiting.
func NewClient(token string) (*Client, error) {
	return NewClientWithOptions(token, ClientOptions{})
}

// NewClientWithOptions creates a new API client with custom configuration options.
// This allows fine-tuning of timeouts and rate limiting behavior.
func NewClientWithOptions(token string, opts ClientOptions) (*Client, error) {
	// Set default values for unspecified options
	if opts.Timeout <= 0 {
		opts.Timeout = 60 * time.Second
	}
	if opts.RequestsPerSec <= 0 {
		opts.RequestsPerSec = internal.DefaultRateLimit
	}
	if opts.BurstSize <= 0 {
		opts.BurstSize = internal.DefaultBurst
	}

	// Validate configuration
	if token == "" {
		return nil, errors.New("token is required")
	}

	// Configure HTTP client with timeout and other settings
	httpClient := &http.Client{
		Timeout: opts.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,              // Pool idle connections
			MaxIdleConnsPerHost: 10,               // Limit per-host connections
			IdleConnTimeout:     90 * time.Second, // Close idle connections after 90s
			DisableCompression:  false,            // Enable gzip compression
		},
	}

	// Initialize rate limiter using token bucket algorithm
	// This allows burst traffic up to BurstSize, then enforces steady rate
	rateLimiter := rate.NewLimiter(rate.Limit(opts.RequestsPerSec), opts.BurstSize)

	// Create main client instance
	client := &Client{
		baseURL:     internal.BaseURL,
		token:       token,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
	}

	// Initialize service endpoints with reference to this client
	client.Domain = services.NewDomainService(client)
	client.Host = services.NewHostService(client)
	client.Certificate = services.NewCertificateService(client)
	client.CVE = services.NewCVEService(client)

	return client, nil
}

// MakeRequest performs an authenticated HTTP request with rate limiting.
// It automatically handles authentication headers and enforces request rate limits.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - method: HTTP method (GET, POST, PUT, DELETE, etc.)
//   - endpoint: API endpoint path (will be appended to baseURL)
//   - body: Request body data (will be JSON encoded if not nil)
//
// Returns the HTTP response or an error if the request fails.
func (c *Client) MakeRequest(ctx context.Context, method, endpoint string, body interface{}) (*http.Response, error) {
	// Apply rate limiting before making the request
	// This will block until a token is available or context is cancelled
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait cancelled: %w", err)
	}

	// Construct full URL and make authenticated request
	fullURL := c.baseURL + endpoint
	return internal.MakeAuthenticatedRequest(ctx, c.httpClient, method, fullURL, c.token, body)
}

// SetRateLimit updates the client's rate limiting configuration.
// This is useful for adjusting limits based on API tier or runtime conditions.
//
// Parameters:
//   - requestsPerSec: New rate limit in requests per second
//   - burstSize: Maximum number of requests that can be made in a burst
func (c *Client) SetRateLimit(requestsPerSec float64, burstSize int) error {
	if requestsPerSec <= 0 {
		return fmt.Errorf("requests per second must be positive, got: %f", requestsPerSec)
	}
	if burstSize <= 0 {
		return fmt.Errorf("burst size must be positive, got: %d", burstSize)
	}

	// Update the rate limiter with new parameters
	c.rateLimiter.SetLimit(rate.Limit(requestsPerSec))
	c.rateLimiter.SetBurst(burstSize)

	return nil
}

// GetRateLimitStatus returns the current rate limiting status.
// This can be useful for monitoring or debugging rate limit behavior.
func (c *Client) GetRateLimitStatus() (limit rate.Limit, burst int) {
	// The rate package doesn't expose current token count, only limit and burst
	return c.rateLimiter.Limit(), c.rateLimiter.Burst()
}

// Close performs cleanup operations for the client.
// It closes idle connections and releases resources.
func (c *Client) Close() {
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

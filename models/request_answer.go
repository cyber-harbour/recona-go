package models

// RequestAnswer represents the complete response data from an HTTP request or network probe.
// This struct captures both successful responses and error conditions, along with proxy and redirect information.
// It's commonly used in network scanning, web crawling, or API testing scenarios where detailed response
// metadata is required for analysis.
type RequestAnswer struct {
	// IP is the resolved IP address of the target host that was actually contacted.
	// This may differ from the original hostname if DNS resolution was involved.
	IP string `json:"ip,omitempty"`

	// Host is the original hostname or domain name that was requested.
	// This preserves the original target identifier before any DNS resolution.
	Host string `json:"host,omitempty"`

	// RawResponse contains the complete HTTP response as a string, including headers and body.
	// This provides the full, unprocessed response for detailed analysis or debugging.
	RawResponse string `json:"raw_response,omitempty"`

	// RawResponseBytes contains the raw HTTP response as a byte array.
	// This is useful for binary content or when precise byte-level data is needed,
	// especially for non-text responses or when character encoding matters.
	RawResponseBytes []byte `json:"raw_response_bytes,omitempty"`

	// Headers contains the HTTP response headers as an array of strings.
	// Each string typically represents a header in "Key: Value" format.
	// This allows for structured access to response metadata.
	Headers []string `json:"headers,omitempty"`

	// StatusCode is the HTTP status code returned by the server (e.g., 200, 404, 500).
	// This indicates the success or failure status of the HTTP request.
	StatusCode int64 `json:"status_code,omitempty"`

	// Error contains any error message that occurred during the request processing.
	// This field is populated when the request fails due to network issues,
	// timeouts, DNS resolution failures, or other connectivity problems.
	Error string `json:"error,omitempty"`

	// ExternalRedirectURL contains the URL if the response included an HTTP redirect.
	// This captures redirect targets (from 3xx status codes) for further processing
	// or to track redirect chains in web crawling scenarios.
	ExternalRedirectURL string `json:"external_redirect_url,omitempty"`

	// ProxyURL is the URL or address of the proxy server used for this request.
	// This field is populated when the request was made through a proxy server.
	ProxyURL string `json:"proxy_url,omitempty"`

	// ProxyPort is the port number of the proxy server used for this request.
	// This works in conjunction with ProxyURL to fully specify the proxy endpoint.
	ProxyPort int64 `json:"proxy_port,omitempty"`

	// ProxyType specifies the type of proxy used (e.g., "HTTP", "SOCKS5", "HTTPS").
	// This indicates the proxy protocol that was employed for the request.
	ProxyType string `json:"proxy_type,omitempty"`
}

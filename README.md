# Recona API Client for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/cyber-harbour/recona-go.svg)](https://pkg.go.dev/github.com/cyber-harbour/recona-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cyber-harbour/recona-go)](https://goreportcard.com/report/github.cyber-harbour/recona-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

The official Go client library for the [Recona.io](https://recona.io/) API, designed to help cybersecurity professionals and developers integrate comprehensive internet asset intelligence into their applications.

[Recona](https://recona.io/) provides extensive cybersecurity intelligence and reconnaissance data to help security teams discover, analyze, and monitor digital assets across the internet.

## 🚀 Key Features

**Comprehensive Asset Intelligence:**
- **Domain Analysis**: WHOIS data, DNS records, subdomains, and domain reputation
- **Host Discovery**: Port scanning, service detection, and vulnerability assessment
- **Certificate Intelligence**: SSL/TLS certificate analysis and monitoring
- **Vulnerability Data**: CVE information with CVSS scores, exploit availability, and remediation guidance

**Advanced Capabilities:**
- Built-in rate limiting to respect API quotas
- Configurable HTTP client with connection pooling
- Context-aware requests with timeout support
- Comprehensive error handling and retry logic
- Type-safe API responses with detailed struct definitions

## 📦 Installation

```bash
go get github.com/cyber-harbour/recona-go
```

## ⚡ Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

	recona "github.com/cyber-harbour/recona-go"
)

func main() {
    // Initialize client with your API token
    client, err := recona.NewClient("your-api-token-here")
    if err != nil {
        log.Fatal(err)
    }
    // Example 1: Get domain information
    domainInfo, err := client.Domain.GetDetails(context.Background(), "example.com")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Domain: %s, Status: %s\n", domainInfo.Name, domainInfo.Status)

    // Example 2: Scan a host
    hostInfo, err := client.Host.GetDetails(context.Background(), "192.168.1.1")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Host: %s, Open Ports: %d\n", hostInfo.IP, len(hostInfo.Ports))

    // Example 3: Search for CVEs
	certificates, err := client.Certificate.Search(context.Background(), models.Search{
		Query: "parsed.validity.end.gte: 2025-09-15 13:21:18",
	})
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d Certificates\n", len(certificates))
}
```

## 🛠 Advanced Configuration

### Custom Client Options

```go
client := recona.NewClientWithOptions("your-token", recona.ClientOptions{
    Timeout:        30 * time.Second,  // Custom timeout
    RequestsPerSec: 10,                // Rate limit: 10 requests/second
    BurstSize:      5,                 // Allow bursts up to 5 requests
})
```

### Rate Limiting Management

```go
// Adjust rate limits at runtime
client.SetRateLimit(15.0, 10) // 15 req/sec, 10 burst

// Check current rate limit status
limit, burst := client.GetRateLimitStatus()
fmt.Printf("Current limit: %v req/sec, Burst: %d\n", limit, burst)
```

## 📚 API Services

### Domain Service
```go
// Get comprehensive domain information
domain, err := client.Domain.GetDetails(ctx, "example.com")

// Search for domains by criteria
domains, err := client.Domain.Search(ctx, models.SearchParams{Query: "google.com", Limit: 100, Offset: 0})

```

### Host Service
```go
// Comprehensive host analysis
host, err := client.Host.GetDetails(ctx, "192.168.1.1")

// Search hosts by criteria
hosts, err := client.Host.Search(ctx, models.SearchParams{Query: "192.168.1.1", Limit: 100, Offset: 0})

```

### Certificate Service
```go
// SSL certificate analysis
cert, err := client.Certificate.GetDetails(ctx, "certificate_fingerprint")

// Search certificates by criteria
certificates, err := client.Certificate.Search(
	ctx, models.SearchParams{Query: "fingerprint_sha256.eq: certificate_fingerprint", Limit: 100, Offset: 0})
```

### CVE Service
```go
// Get CVE details
cve, err := client.CVE.GetDetails(ctx, "CVE-2023-12345")

// Search CVEs by criteria
cves, err := client.Certificate.Search(ctx, models.SearchParams{Query: "id.eq: CVE-2023-12345", Limit: 100, Offset: 0}

```

## 🔧 Examples

Account:
- [Check your API quotas](./examples/account/main.go)

Target info:
- [IP details](./examples/ip_details/main.go)
- [Domain details](./examples/domain_details/main.go)
- [Domain batch_search](./examples/domain_batch/main.go)

Search with params (up to 10 000 results):
- [Subdomains lookup](./examples/domain_subdomains/main.go)
- [Search domains by technology](./examples/domains_with_technology/main.go)
- [Search emails by domain name](./examples/emails_search/main.go)
- [Search IPv4 hosts with specific open port](./examples/ips_with_open_port/main.go)
- [Search IPv4 hosts by geolocation](./examples/ips_by_country/main.go)

### Running Examples

```bash
# Set your API token
export RECONA_API_TOKEN="your-api-token-here"

# Run domain analysis example
go run ./examples/domain_details/main.go --access_token token_here

# Run host scanning example
go run ./examples/ips_with_open_port/main.go --access_token token_here

```

## 🧪 Testing

Run the complete test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run tests with verbose output
go test -v ./...

```

## 📈 Performance and Best Practices

### Rate Limiting
- Default: 10 requests/second with burst capacity of 2
- Configure based on your API plan and usage patterns
- Monitor rate limit status to optimize performance

### Connection Management
- HTTP client uses connection pooling automatically
- Idle connections are managed and cleaned up
- Configure timeouts based on your network conditions

### Context Usage
- Always use context for request cancellation
- Set appropriate timeouts for long-running operations
- Handle context cancellation gracefully

```go
// Example of proper context usage
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := client.Domain.GetDetails(ctx, "example.com")
```

## 🔗 API Coverage

This client library provides complete coverage of the Recona API:

- ✅ **Domain API**: Full WHOIS, DNS, and subdomain functionality
- ✅ **Host API**: Port scanning, service detection, and analysis
- ✅ **Certificate API**: SSL/TLS certificate management and validation
- ✅ **CVE API**: Vulnerability data and exploit intelligence
- ✅ **AS API**: Autonomous System information and routing data
- ✅ **Account API**: Usage statistics and quota management

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🆘 Support

- **Documentation**: [API Reference](https://reconatest.io/docs/search-concept)
- **Email**: support@recona.io

## 🔐 Security

For security vulnerabilities, please email security@recona.io instead of opening a public issue.

---

Built with ❤️ by the Recona team for the cybersecurity community.
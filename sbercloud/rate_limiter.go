package sbercloud

import (
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimitedTransport wraps http.RoundTripper with rate limiting functionality.
// It uses the token bucket algorithm to limit the number of HTTP requests per second.
type RateLimitedTransport struct {
	// Transport is the underlying HTTP transport used to execute requests.
	// If nil, http.DefaultTransport will be used.
	Transport http.RoundTripper

	// Limiter controls the rate of requests.
	// If nil, no rate limiting is applied.
	Limiter *rate.Limiter
}

// RoundTrip implements the http.RoundTripper interface.
// It applies rate limiting before forwarding the request to the underlying transport.
func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// If limiter is configured, wait for permission to proceed
	if t.Limiter != nil {
		// Wait blocks until limiter permits one event or context is done.
		// This ensures we don't exceed the configured rate limit.
		err := t.Limiter.Wait(req.Context())
		if err != nil {
			// Return error if context was cancelled or deadline exceeded
			return nil, err
		}
	}

	// Proceed with the actual HTTP request
	return t.Transport.RoundTrip(req)
}

// NewRateLimitedTransport creates a new rate-limited HTTP transport.
//
// Parameters:
//   - rateLimit: maximum requests per second (0 = unlimited)
//   - baseTransport: underlying transport to wrap (nil = http.DefaultTransport)
//
// The rate limiter uses a token bucket algorithm where:
//   - Tokens are added at a constant rate (rateLimit per second)
//   - Each request consumes one token
//   - If no tokens available, the request waits until one becomes available
//
// Example:
//
//	// Limit to 80 requests per second
//	transport := NewRateLimitedTransport(80, nil)
//	client := &http.Client{Transport: transport}
func NewRateLimitedTransport(rateLimit int, baseTransport http.RoundTripper) *RateLimitedTransport {
	// Use default transport if none provided
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	transport := &RateLimitedTransport{
		Transport: baseTransport,
	}

	// Only create limiter if rate limit is specified and greater than 0
	if rateLimit > 0 {
		// rate.Limit(rateLimit) = requests per second
		// rateLimit = burst size (allows burst of N requests before rate limiting kicks in)
		transport.Limiter = rate.NewLimiter(rate.Limit(rateLimit), rateLimit)
	}

	return transport
}

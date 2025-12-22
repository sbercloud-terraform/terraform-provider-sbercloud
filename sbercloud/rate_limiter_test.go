package sbercloud

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestRateLimitedTransport_NoLimit(t *testing.T) {
	// Test that requests pass through without delay when rate_limit = 0
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with no limit (0)
	transport := NewRateLimitedTransport(0, nil)
	client := &http.Client{Transport: transport}

	// Send 10 requests rapidly
	for i := 0; i < 10; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}

	mu.Lock()
	count := requestCount
	mu.Unlock()

	if count != 10 {
		t.Errorf("Expected 10 requests, got %d", count)
	}
}

func TestRateLimitedTransport_WithLimit(t *testing.T) {
	// Test that rate limiting actually works and enforces delay
	requestTimes := []time.Time{}
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestTimes = append(requestTimes, time.Now())
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with limit of 5 req/s
	transport := NewRateLimitedTransport(5, nil)
	client := &http.Client{Transport: transport}

	start := time.Now()

	// Send 10 requests (should take ~2 seconds with 5 req/s limit)
	for i := 0; i < 10; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		resp.Body.Close()
	}

	duration := time.Since(start)

	mu.Lock()
	count := len(requestTimes)
	mu.Unlock()

	// Verify all requests completed
	if count != 10 {
		t.Errorf("Expected 10 requests, got %d", count)
	}

	// With burst=1 and rate=5 req/s: first request is immediate, then 9 requests at 200ms intervals
	// Total time: 0 + 9*200ms = 1.8s minimum
	// Should take at least 1.5 seconds (with some margin for the first request)
	if duration < 1500*time.Millisecond {
		t.Errorf("Requests completed too quickly: %v (expected >= 1.5s)", duration)
	}

	// Should not take more than 3 seconds (with reasonable margin)
	if duration > 3*time.Second {
		t.Errorf("Requests took too long: %v (expected <= 3s)", duration)
	}
}

func TestRateLimitedTransport_ContextCancellation(t *testing.T) {
	// Test that context cancellation is properly handled
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with very low limit to force waiting
	transport := NewRateLimitedTransport(1, nil)
	client := &http.Client{Transport: transport}

	// Create a context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// First request should succeed
	req1, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	resp1, err1 := client.Do(req1)
	if err1 != nil {
		t.Fatalf("First request failed: %v", err1)
	}
	resp1.Body.Close()

	// Second request should fail due to context timeout
	// (rate limiter will wait, but context will expire)
	req2, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	_, err2 := client.Do(req2)
	if err2 == nil {
		t.Error("Expected second request to fail due to context timeout")
	}
}

func TestNewRateLimitedTransport_NilTransport(t *testing.T) {
	// Test that nil transport defaults to http.DefaultTransport
	transport := NewRateLimitedTransport(10, nil)

	if transport.Transport == nil {
		t.Error("Expected transport to be set to default")
	}

	if transport.Limiter == nil {
		t.Error("Expected limiter to be created for rate > 0")
	}
}

func TestNewRateLimitedTransport_CustomTransport(t *testing.T) {
	// Test that custom transport is preserved
	customTransport := &http.Transport{
		MaxIdleConns: 100,
	}

	transport := NewRateLimitedTransport(10, customTransport)

	if transport.Transport != customTransport {
		t.Error("Expected custom transport to be preserved")
	}

	if transport.Limiter == nil {
		t.Error("Expected limiter to be created for rate > 0")
	}
}

func TestNewRateLimitedTransport_ZeroRate(t *testing.T) {
	// Test that zero rate does not create a limiter
	transport := NewRateLimitedTransport(0, nil)

	if transport.Limiter != nil {
		t.Error("Expected no limiter for rate = 0")
	}

	if transport.Transport == nil {
		t.Error("Expected transport to be set even with rate = 0")
	}
}

func TestRateLimitedTransport_ConcurrentRequests(t *testing.T) {
	// Test that rate limiter works correctly with concurrent requests
	requestCount := 0
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create transport with limit of 10 req/s
	transport := NewRateLimitedTransport(10, nil)
	client := &http.Client{Transport: transport}

	// Launch 20 concurrent goroutines making requests
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := client.Get(server.URL)
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	mu.Lock()
	count := requestCount
	mu.Unlock()

	// All 20 requests should have completed
	if count != 20 {
		t.Errorf("Expected 20 requests, got %d", count)
	}

	// Should take at least 1 second (20 requests / 10 req/s = 2s, but some margin)
	if duration < 1*time.Second {
		t.Errorf("Concurrent requests completed too quickly: %v", duration)
	}
}

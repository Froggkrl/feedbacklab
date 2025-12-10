package health

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Service defines the interface for health check operations.
type Service interface {
	Check(ctx context.Context) error
}

// HTTPHealthService performs health checks via HTTP requests.
type HTTPHealthService struct {
	URL     string
	Timeout time.Duration
}

// NewHTTPHealthService creates a new HTTPHealthService instance.
func NewHTTPHealthService(url string, timeout time.Duration) *HTTPHealthService {
	return &HTTPHealthService{
		URL:     url,
		Timeout: timeout,
	}
}

// Check performs a health check by making an HTTP request to the configured URL.
func (h *HTTPHealthService) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("health request failed: %w", err)
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

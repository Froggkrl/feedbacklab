package health

import "context"

// SelfHealthService provides a simple health check that always returns healthy.
type SelfHealthService struct{}

// NewSelfHealthService creates a new SelfHealthService instance.
func NewSelfHealthService() *SelfHealthService {
	return &SelfHealthService{}
}

// Check always returns nil, indicating the service is healthy.
func (s *SelfHealthService) Check(_ context.Context) error {
	return nil
}

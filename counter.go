package counter

import "context"

// Counter provides a counter service.
type Counter interface {
	// Incr increases counter for given id.
	Incr(ctx context.Context, id string) error
	// Get returns counter for given id.
	Get(ctx context.Context, id string) (int64, error)
}

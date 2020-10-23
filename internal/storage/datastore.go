package storage

import "context"

// Datastore is a data storage implementation.
type Datastore interface {
	ReportDropout(ctx context.Context) error
}

package watermark

import (
	"context"

	"github.com/justisGipson/go-service/internal"
)

// Contexts enable the microservices to handle the multiple concurrent requests
type Service interface {
	// get doc list
	Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error)
	Status(ctx context.Context, ticketID string) (internal.Status, error)
	Watermark(ctx context.Context, ticketID, mark string) (int, error)
	AddDocument(ctx context.Context, doc *internal.Document) (string, error)
	ServiceStatus(ctx context.Context) (int, error)
}

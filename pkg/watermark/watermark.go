package watermark

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/justisGipson/go-service/internal"

	"github.com/go-kit/kit/log"
	"github.com/lithammer/shortuuid/v3"
)

type watermarkService struct{}

func NewService() Service {
	return &watermarkService{}
}

func (w *watermarkService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	// query db using filters and return list of docs
	// throw err if filter (key) is invalid and also return error if no items found
	doc := internal.Document{
		Content: "book",
		Title:   "Harry Potter and the Sorcerer's Stone",
		Author:  "J.K. Rowling",
		Topic:   "Fiction and Magic",
	}
	return []internal.Document{doc}, nil
}

func (w *watermarkService) Status(_ context.Context, ticketID string) (internal.Status, error) {
	// query db using ticketID and return doc info
	// throw err if ticketID is invalid or no doc exists for ticketID
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(_ context.Context, ticketID, mark string) (int, error) {
	// update db entry with watermark field as non-empty
	// first check if watermark status is not in a `InProgress, Started, or Finished` state
	// if yes, return invalid request
	// throw error if no item found w/ ticketID
	return http.StatusOK, nil
}

func (w *watermarkService) AddDocument(_ context.Context, doc *internal.Document) (string, error) {
	// add doc entry to db by calling db service
	// return error if doc is invalid and/or the db theres invalid entry err
	newTicketID := shortuuid.New()
	return newTicketID, nil
}

func (w *watermarkService) ServiceStatus(_ context.Context) (int, error) {
	logger.Log("Checking the service health...")
	return http.StatusOK, nil
}

var logger = log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

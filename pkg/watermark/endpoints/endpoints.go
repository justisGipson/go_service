package endpoints

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/justisGipson/go-service/internal"
	"github.com/justisGipson/go-service/pkg/watermark"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type Set struct {
	GetEndpoint           endpoint.Endpoint
	AddDocumentEndpoint   endpoint.Endpoint
	StatusEndpoint        endpoint.Endpoint
	ServiceStatusEndpoint endpoint.Endpoint
	WatermarkEndpoint     endpoint.Endpoint
}

func NewEndpointSet(svc watermark.Service) Set {
	return Set{
		GetEndpoint:           MakeGetEndpoint(svc),
		AddDocumentEndpoint:   MakeAddDocumentEndpoint(svc),
		StatusEndpoint:        MakeStatusEndpoint(svc),
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
		WatermarkEndpoint:     MakeWatermarkEndpoint(svc),
	}
}

func MakeGetEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRequest)
		docs, err := svc.Get(ctx, req.Filters...)
		if err != nil {
			return GetResponse{docs, err.Error()}, nil
		}
		return GetResponse{docs, ""}, nil
	}
}

func MakeStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StatusRequest)
		status, err := svc.Status(ctx, req.TicketID)
		if err != nil {
			return StatusResponse{Status: status, Err: err.Error()}, nil
		}
		return StatusResponse{Status: status, Err: ""}, nil
	}
}

func MakeAddDocumentEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddDocumentRequest)
		ticketID, err := svc.AddDocument(ctx, req.Document)
		if err != nil {
			return AddDocumentResponse{TicketID: ticketID, Err: err.Error()}, nil
		}
		return AddDocumentResponse{TicketID: ticketID, Err: ""}, nil
	}
}

func MakeWatermarkEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(WatermarkRequest)
		code, err := svc.Watermark(ctx, req.TicketID, req.Mark)
		if err != nil {
			return WatermarkResponse{Code: code, Err: err.Error()}, nil
		}
		return WatermarkResponse{Code: code, Err: ""}, nil
	}
}

func MakeServiceStatusEndpoint(svc watermark.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code, err := svc.ServiceStatus(ctx)
		if err != nil {
			return ServiceStatusResponse{Code: code, Err: err.Error()}, nil
		}
		return ServiceStatusResponse{Code: code, Err: ""}, nil
	}
}

func (s *Set) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	res, err := s.GetEndpoint(ctx, GetRequest{Filters: filters})
	if err != nil {
		return []internal.Document{}, err
	}
	getRes := res.(GetResponse)
	if getRes.Err != "" {
		return []internal.Document{}, errors.New(getRes.Err)
	}
	return getRes.Documents, nil
}

func (s *Set) ServiceStatus(ctx context.Context) (int, error) {
	res, err := s.ServiceStatusEndpoint(ctx, ServiceStatusRequest{})
	svcStatusRes := res.(ServiceStatusResponse)
	if err != nil {
		return svcStatusRes.Code, err
	}
	if svcStatusRes.Err != "" {
		return svcStatusRes.Code, errors.New(svcStatusRes.Err)
	}
	return svcStatusRes.Code, nil
}

func (s *Set) AddDocument(ctx context.Context, doc *internal.Document) (string, error) {
	res, err := s.AddDocumentEndpoint(ctx, AddDocumentRequest{Document: doc})
	if err != nil {
		return "", err
	}
	adRes := res.(AddDocumentResponse)
	if adRes.Err != "" {
		return "", errors.New(adRes.Err)
	}
	return adRes.TicketID, nil
}

func (s *Set) Status(ctx context.Context, ticketID string) (internal.Status, error) {
	res, err := s.StatusEndpoint(ctx, StatusRequest{TicketID: ticketID})
	if err != nil {
		return internal.Failed, err
	}
	stsRes := res.(StatusResponse)
	if stsRes.Err != "" {
		return internal.Failed, errors.New(stsRes.Err)
	}
	return stsRes.Status, nil
}

func (s *Set) Watermark(ctx context.Context, ticketID, mark string) (int, error) {
	res, err := s.WatermarkEndpoint(ctx, WatermarkRequest{TicketID: ticketID, Mark: mark})
	wmRes := res.(WatermarkResponse)
	if err != nil {
		return wmRes.Code, err
	}
	if wmRes.Err != "" {
		return wmRes.Code, errors.New(wmRes.Err)
	}
	return wmRes.Code, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

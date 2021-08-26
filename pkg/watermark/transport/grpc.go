/// map of the protobuf payload to their requests and responses
/// create gRPC handler constructor which creates the set of grpcServers
/// and registers the appropriate endpoint to the decoders and encoders of requests and responses

package transport

import (
	"context"

	"github.com/justisGipson/go-service/api/vi/pb/watermark"
	"github.com/justisGipson/go-service/internal"
	"github.com/justisGipson/go-service/pkg/watermark/endpoints"

	grpcTransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	get           grpcTransport.Handler
	status        grpcTransport.Handler
	addDocument   grpcTransport.Handler
	watermark     grpcTransport.Handler
	serviceStatus grpcTransport.Handler
}

func NewGRPCServer(ep endpoints.Set) watermark.watermarkServer {
	return &grpcServer{
		get: grpcTransport.NewServer(
			ep.GetEndpoint,
			decodeGRPCGetRequest,
			decodeGRPCGetResponse,
		),

		status: grpcTransport.NewServer(
			ep.StatusEndpoint,
			decodeGRPCStatusRequest,
			decodeGRPCStatusResponse,
		),
		addDocument: grpcTransport.NewServer(
			ep.AddDocumentEndpoint,
			decodeGRPCAddDocumentRequest,
			decodeGRPCAddDocumentResponse,
		),
		watermark: grpcTransport.NewServer(
			ep.WatermarkEndpoint,
			decodeGRPCWatermarkRequest,
			decodeGRPCWatermarkResponse,
		),
		serviceStatus: grpcTransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeGRPCServiceStatusRequest,
			decodeGRPCServiceStatusResponse,
		),
	}
}

func (g *grpcServer) Get(ctx context.Context, r *watermark.GetRequest) (*watermark.GetReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.GetReply), nil
}

func (g *grpcServer) ServiceStatus(ctx context.Context, r *watermark.ServiceStatusRequest) (*watermark.ServiceStatusReply, error) {
	_, rep, err := g.get.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.ServiceStatusReply), nil
}

func (g *grpcServer) AddDocument(ctx context.Context, r *watermark.AddDocumentRequest) (*watermark.AddDocumentReply, error) {
	_, rep, err := g.addDocument.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.AddDocumentReply), nil
}

func (g *grpcServer) Status(ctx context.Context, r *watermark.StatusRequest) (*watermark.StatusReply, error) {
	_, rep, err := g.serviceStatus.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.StatusReply), nil
}

func (g *grpcServer) Watermark(ctx context.Context, r *watermark.WatermarkRequest) (*watermark.WatermarkReply, error) {
	_, rep, err := g.watermark.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return rep.(*watermark.WatermarkReply), nil
}

func decodeGRPCGetRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.GetRequest)
	var filters []internal.Filter
	for _, f := range req.Filters {
		filters = append(filters, internal.Filter{Key: f.key, Value: f.value})
	}
	return endpoints.GetRequest{Filters: filters}, nil
}

func decodeGRPCStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.StatusRequest)
	return endpoints.StatusRequest{TicketID: req.ticketID}, nil
}

func decodeGRPCWatermarkRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.WatermarkRequest)
	return endpoints.WatermarkRequest{TicketID: req.ticketID, Mark: req.mark}, nil
}

func decodeGRPCAddDocumentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*watermark.AddDocumentRequest)
	doc := &internal.Document{
		Content:   req.Document.Content,
		Title:     req.Document.Title,
		Author:    req.Document.Author,
		Topic:     req.Document.Topic,
		Watermark: req.Document.Watermark,
	}
	return endpoints.AddDocumentRequest{Document: doc}, nil
}

func decodeGRPCServiceStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return endpoints.ServiceStatusRequest{}, nil
}

func decodeGRPCGetResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.GetReply)
	var docs []internal.Document
	for _, d := range reply.Documents {
		doc := internal.Document{
			Content:   d.Content,
			Title:     d.Title,
			Author:    d.Author,
			Topic:     d.Topic,
			Watermark: d.Watermark,
		}
		docs = append(docs, doc)
	}
	return endpoints.GetResponse{Documents: docs, Err: reply.Err}, nil
}

func decodeGRPCStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.StatusReply)
	return endpoints.StatusResponse{Status: internal.Status(reply.Status), Err: reply.Err}, nil
}

func decodeGRPCWatermarkResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.WatermarkReply)
	return endpoints.WatermarkResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

func decodeGRPCAddDocumentResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.AddDocumentReply)
	return endpoints.AddDocumentResponse{TicketID: reply.TicketID, Err: reply.Err}, nil
}

func decodeGRPCServiceStatusResponse(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*watermark.ServiceStatusReply)
	return endpoints.ServiceStatusResponse{Code: int(reply.Code), Err: reply.Err}, nil
}

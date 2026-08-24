package server

import (
	"context"
	"fmt"

	"fluxa/pkg/plugin"
	pluginv1 "fluxa/proto/fluxa/v1"
)

// Server implementa o contrato gRPC de canal e pode ser alimentado pelo poller.
type Server struct {
	pluginv1.UnimplementedChannelPluginServerServer
	Send func(context.Context, plugin.OutboundMessageRequest) (plugin.OutboundMessageResponse, error)
}

func (s *Server) SendMessage(ctx context.Context, req *pluginv1.OutboundMessageRequest) (*pluginv1.OutboundMessageResponse, error) {
	if s.Send == nil {
		return nil, fmt.Errorf("telegram: envio não configurado")
	}
	resp, err := s.Send(ctx, plugin.OutboundMessageRequest{ChannelID: req.GetChannelId(), CustomerRef: req.GetCustomerRef(), Text: req.GetText(), Metadata: req.GetMetadata()})
	if err != nil {
		return nil, fmt.Errorf("telegram: enviando: %w", err)
	}
	return &pluginv1.OutboundMessageResponse{MessageId: resp.MessageID, Status: resp.Status}, nil
}
func (s *Server) HandleWebhook(context.Context, *pluginv1.WebhookPayloadRequest) (*pluginv1.WebhookPayloadResponse, error) {
	return &pluginv1.WebhookPayloadResponse{HttpStatusCode: 405, ResponseBody: "long-poll only"}, nil
}
func (s *Server) HealthCheck(context.Context, *pluginv1.Empty) (*pluginv1.HealthStatus, error) {
	return &pluginv1.HealthStatus{Healthy: true, Message: "ok"}, nil
}

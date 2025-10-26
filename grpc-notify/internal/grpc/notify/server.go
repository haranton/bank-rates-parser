package notify

import (
	"context"
	"fmt"
	"grpc-notify/internal/service"
	"grpc-notify/proto"

	"google.golang.org/grpc"
)

type serverAPI struct {
	proto.UnimplementedNotificationServiceServer
	srv *service.TelegramClient
}

func Register(gRPC *grpc.Server, srv *service.TelegramClient) {
	proto.RegisterNotificationServiceServer(gRPC, &serverAPI{srv: srv})
}

func (s *serverAPI) SendNotification(ctx context.Context, req *proto.NotificationRequest) (*proto.NotificationResponse, error) {
	fmt.Printf("[NOTIFICATION] Type: %s | Target: %s | Message: %s\n", req.Type, req.Target, req.Message)

	err := s.srv.SendMessage(req.Message)
	if err != nil {
		return nil, err
	}

	return &proto.NotificationResponse{
		Success: true,
		Error:   "",
	}, nil
}

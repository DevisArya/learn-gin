package grpcdelivery

import (
	"context"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/entity"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/pb"
	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/usecase"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationGrpc interface {
	pb.NotificationServiceServer
}

type NotificationGrpcImpl struct {
	pb.UnimplementedNotificationServiceServer
	notificationUC usecase.NotificationUseCase
	Log            *logrus.Logger
}

func NewNotificationGrpc(notificationUc usecase.NotificationUseCase, log *logrus.Logger) NotificationGrpc {
	return &NotificationGrpcImpl{
		notificationUC: notificationUc,
		Log:            log,
	}
}

func (g *NotificationGrpcImpl) ReceiveUserAction(ctx context.Context, reqPb *pb.UserActionRequest) (*pb.Empty, error) {

	notificationCreateReq := dto.CreateNotificationRequest{
		UserID:      reqPb.GetUserId(),
		Action:      entity.ActionType(reqPb.Action.String()),
		Description: reqPb.GetDescription(),
	}

	err := g.notificationUC.SaveNotification(ctx, &notificationCreateReq)

	if err != nil {
		g.Log.Errorf("Failed to save notification: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	g.Log.Infof("Successfully saved notification for user_id=%v", notificationCreateReq.UserID)

	return &pb.Empty{}, nil
}

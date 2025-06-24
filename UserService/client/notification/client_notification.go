package notification

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DevisArya/BE-challenge-syn/UserService/client/pb"
	"github.com/DevisArya/BE-challenge-syn/UserService/internal/dto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var stringToActionTypeEnum = map[string]pb.ActionType{
	"CREATE": pb.ActionType_CREATE,
	"UPDATE": pb.ActionType_UPDATE,
	"DELETE": pb.ActionType_DELETE,
}

type NotificationClient interface {
	SendUserNotification(ctx context.Context, request *dto.SendUserNotificationRequest) error
	Close() error
}

type NotificationClientImpl struct {
	client pb.NotificationServiceClient
	conn   *grpc.ClientConn
}

func NewNotificationClientGRPC(address string) (*NotificationClientImpl, error) {

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewNotificationServiceClient(conn)
	return &NotificationClientImpl{
		client: client,
		conn:   conn,
	}, nil
}

func (n *NotificationClientImpl) SendUserNotification(ctx context.Context, request *dto.SendUserNotificationRequest) error {

	actionEnum, ok := stringToActionTypeEnum[strings.ToUpper(request.Action)]
	if !ok {
		return fmt.Errorf("invalid action: %s", request.Action)
	}

	pbReq := &pb.UserActionRequest{
		UserId:      request.UserID,
		Action:      actionEnum,
		Description: request.Description,
		OccurredAt:  timestamppb.New(time.Now()),
	}

	_, err := n.client.ReceiveUserAction(ctx, pbReq)
	if err != nil {
		return err
	}

	return nil
}

func (n *NotificationClientImpl) Close() error {
	return n.conn.Close()
}

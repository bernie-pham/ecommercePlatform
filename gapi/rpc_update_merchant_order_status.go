package gapi

import (
	"context"
	"database/sql"

	db "github.com/bernie-pham/ecommercePlatform/db/sqlc"
	"github.com/bernie-pham/ecommercePlatform/pb"
	"github.com/bernie-pham/ecommercePlatform/session"
	"github.com/bernie-pham/ecommercePlatform/token"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	statusValue = map[db.OrderStatus]int{
		db.OrderStatusOpen:       0,
		db.OrderStatusApproved:   1,
		db.OrderStatusPrepared:   2,
		db.OrderStatusCanceled:   3,
		db.OrderStatusPicked:     4,
		db.OrderStatusOnDelivery: 5,
		db.OrderStatusDeliveried: 6,
	}
	statusString = map[int32]db.OrderStatus{
		0: db.OrderStatusOpen,
		1: db.OrderStatusApproved,
		2: db.OrderStatusPrepared,
		3: db.OrderStatusCanceled,
		4: db.OrderStatusPicked,
		5: db.OrderStatusOnDelivery,
		6: db.OrderStatusDeliveried,
	}
)

func (server *Server) UpdateOrderStatus(
	ctx context.Context,
	req *pb.UpdateMerchantOrderReq,
) (*pb.UpdateMerchantOrderResponse, error) {
	session_payload := ctx.Value(token.AuthPayloadKey).(session.SessionPayload)
	// fmt.Printf("session_payload: %v\n", session_payload)
	arg := db.UpdateMerchantOrderStatusParams{
		OrderStatus: statusString[int32(req.Status.Number())],
		ID:          req.MerchantOrderId,
		MerchantID:  session_payload.UserID,
	}
	// TODO: check order status, if backward status, it is invalid.

	order, err := server.store.UpdateMerchantOrderStatus(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "merchant order id is not found")
		}
		log.Error().
			Err(err).
			Msg("failed to update merchent order status")
		return nil, status.Errorf(codes.Internal, "failed to update merchant order status")
	}

	rsp := &pb.UpdateMerchantOrderResponse{
		MerchantOrder: &pb.MerchantOrder{
			Id:         order.ID,
			OrderId:    order.OrderID,
			Status:     pb.OrderStatus(statusValue[order.OrderStatus]),
			MerchantId: order.MerchantID,
			TotalPrice: order.TotalPrice,
			CreatedAt:  timestamppb.New(order.CreatedAt),
			UpdatedAt:  timestamppb.New(order.UpdatedAt),
		},
	}
	return rsp, nil
}

// func validateUpdateOrderStatusParams(
// 	req *pb.UpdateMerchantOrderReq,
// ) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if err := req.
// }

package service

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/goodgodx64/orderservice-go/pkg/api/grpc"
	"github.com/goodgodx64/orderservice-go/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	items map[string]*pb.Order
	mu    sync.Mutex
}

func New() *OrderService {
	return &OrderService{
		items: make(map[string]*pb.Order),
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id := uuid.New().String()
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "CreateOrder, creating order", zap.String("item", in.GetItem()), zap.Int32("quantity", in.GetQuantity()), zap.String("order_id", id))

	o.mu.Lock()
	defer o.mu.Unlock()

	o.items[id] = &pb.Order{Id: id, Item: in.GetItem(), Quantity: in.GetQuantity()}

	return &pb.CreateOrderResponse{Id: id}, nil
}

func (o *OrderService) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return nil, fmt.Errorf("order with id %s not found", in.GetId())
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "GetOrder, returning order", zap.String("order_id", in.GetId()))

	return &pb.GetOrderResponse{Order: o.items[in.GetId()]}, nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, in *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return nil, fmt.Errorf("order with id %s not found", in.GetId())
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "UpdateOrder, updating order", zap.String("order_id", in.GetId()), zap.String("new_item", in.GetItem()), zap.Int32("new_quantity", in.GetQuantity()))

	o.mu.Lock()
	defer o.mu.Unlock()

	updatedOrder := &pb.Order{Id: in.GetId(), Item: in.GetItem(), Quantity: in.GetQuantity()}
	o.items[in.GetId()] = updatedOrder

	return &pb.UpdateOrderResponse{Order: updatedOrder}, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, in *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return &pb.DeleteOrderResponse{Success: false}, fmt.Errorf("order with id %s not found", in.GetId())
	}
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "DeleteOrder, deleting order", zap.String("order_id", in.GetId()))

	o.mu.Lock()
	defer o.mu.Unlock()

	delete(o.items, in.GetId())
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (o *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	logger.GetLoggerFromCtx(ctx).Debug(ctx, "ListOrders, listing orders")

	var orders []*pb.Order

	o.mu.Lock()
	defer o.mu.Unlock()

	for _, v := range o.items {
		orders = append(orders, v)
	}
	return &pb.ListOrdersResponse{Orders: orders}, nil
}

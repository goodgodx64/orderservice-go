package service

import (
	"context"
	"fmt"
	pb "grpctask/pkg/api/test"
	"log"
	"sync"

	"github.com/google/uuid"
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
	log.Printf("Created order: %s %d", in.GetItem(), in.GetQuantity())
	id := uuid.NewString()
	o.mu.Lock()
	defer o.mu.Unlock()
	o.items[id] = &pb.Order{Id: id, Item: in.GetItem(), Quantity: in.GetQuantity()}
	return &pb.CreateOrderResponse{Id: id}, nil
}

func (o *OrderService) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return nil, fmt.Errorf("order with id %s not found", in.GetId())
	}
	log.Printf("Returning order with id: %s", in.GetId())
	return &pb.GetOrderResponse{Order: o.items[in.GetId()]}, nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, in *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return nil, fmt.Errorf("order with id %s not found", in.GetId())
	}
	o.mu.Lock()
	reqOrder := o.items[in.GetId()]
	log.Printf(`Updating order with id: %s, item: %s, quantity: %d`, reqOrder.GetId(), reqOrder.GetItem(), reqOrder.GetQuantity())

	updatedOrder := &pb.Order{Id: in.GetId(), Item: in.GetItem(), Quantity: in.GetQuantity()}
	defer o.mu.Unlock()
	o.items[in.GetId()] = updatedOrder
	return &pb.UpdateOrderResponse{Order: updatedOrder}, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, in *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if _, ok := o.items[in.GetId()]; !ok {
		return &pb.DeleteOrderResponse{Success: false}, fmt.Errorf("order with id %s not found", in.GetId())
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.items, in.GetId())
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (o *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	var orders []*pb.Order
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, v := range o.items {
		orders = append(orders, v)
	}
	return &pb.ListOrdersResponse{Orders: orders}, nil
}

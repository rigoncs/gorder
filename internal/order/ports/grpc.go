package ports

import (
	context "context"
	"github.com/rigoncs/gorder/common/genproto/orderpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) UpdateOder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

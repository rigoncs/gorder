package adapters

import (
	"context"
	"github.com/rigoncs/gorder/common/genproto/orderpb"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (g *OrderGRPC) UpdateOrder(ctx context.Context, request *orderpb.Order) error {
	_, err := g.client.UpdateOder(ctx, request)
	return err
}

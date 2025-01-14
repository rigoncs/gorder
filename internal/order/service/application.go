package service

import (
	"context"
	"github.com/rigoncs/gorder/order/adapters"
	"github.com/rigoncs/gorder/order/app"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	return app.Application{}
}

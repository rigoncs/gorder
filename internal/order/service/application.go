package service

import (
	"context"
	"github.com/rigoncs/gorder/common/metrics"
	"github.com/rigoncs/gorder/order/adapters"
	"github.com/rigoncs/gorder/order/app"
	"github.com/rigoncs/gorder/order/app/command"
	"github.com/rigoncs/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}

package service

import (
	"context"
	grpcClient "github.com/rigoncs/gorder/common/client"
	"github.com/rigoncs/gorder/common/metrics"
	"github.com/rigoncs/gorder/payment/adapters"
	"github.com/rigoncs/gorder/payment/app"
	"github.com/rigoncs/gorder/payment/app/command"
	"github.com/rigoncs/gorder/payment/domain"
	"github.com/rigoncs/gorder/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmemProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}

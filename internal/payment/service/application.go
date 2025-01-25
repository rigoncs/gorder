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
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	//memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	metricClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		Host:        viper.GetString("payment.metrics_export_addr"),
		ServiceName: viper.GetString("payment.service-name"),
	})
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logrus.StandardLogger(), metricClient),
		},
	}
}

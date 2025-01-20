package main

import (
	"context"
	"github.com/rigoncs/gorder/common/broker"
	_ "github.com/rigoncs/gorder/common/config"
	"github.com/rigoncs/gorder/common/logging"
	"github.com/rigoncs/gorder/common/server"
	"github.com/rigoncs/gorder/common/tracing"
	"github.com/rigoncs/gorder/payment/infrastructure/consumer"
	"github.com/rigoncs/gorder/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverType := viper.GetString("payment.server-to-run")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}

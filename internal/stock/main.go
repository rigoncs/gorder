package main

import (
	"context"
	_ "github.com/rigoncs/gorder/common/config"
	"github.com/rigoncs/gorder/common/discovery"
	"github.com/rigoncs/gorder/common/genproto/stockpb"
	"github.com/rigoncs/gorder/common/logging"
	"github.com/rigoncs/gorder/common/server"
	"github.com/rigoncs/gorder/common/tracing"
	"github.com/rigoncs/gorder/stock/ports"
	"github.com/rigoncs/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	logrus.Info(serverType)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)
	application := service.NewApplication(ctx)

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			srv := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, srv)
		})
	case "http":
	// 暂时不用
	default:
		panic("unexpected server type")
	}
}

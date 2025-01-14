package main

import (
	"github.com/rigoncs/gorder/common/genproto/stockpb"
	"github.com/rigoncs/gorder/common/server"
	"github.com/rigoncs/gorder/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			srv := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, srv)
		})
	case "http":
	// 暂时不用
	default:
		panic("unexpected server type")
	}
}

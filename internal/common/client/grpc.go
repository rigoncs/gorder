package client

import (
	"context"
	"github.com/rigoncs/gorder/common/discovery"
	"github.com/rigoncs/gorder/common/genproto/orderpb"
	"github.com/rigoncs/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for stock grpc")
	}
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error { return nil }, err
	}
	if grpcAddr == "" {
		logrus.Warn("empty grpc addr for order grpc")
	}
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}

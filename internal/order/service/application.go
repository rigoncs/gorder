package service

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rigoncs/gorder/common/broker"
	grpcClient "github.com/rigoncs/gorder/common/client"
	"github.com/rigoncs/gorder/common/metrics"
	"github.com/rigoncs/gorder/order/adapters"
	"github.com/rigoncs/gorder/order/adapters/grpc"
	"github.com/rigoncs/gorder/order/app"
	"github.com/rigoncs/gorder/order/app/command"
	"github.com/rigoncs/gorder/order/app/query"
	"github.com/rigoncs/gorder/order/infrastructure/mq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCh()
		_ = ch.Close()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	//orderRepo := adapters.NewMemoryOrderRepository()
	mongoClient := newMongoClient()
	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
	metricClient := metrics.TodoMetrics{}
	eventPublisher := mq.NewRabbitMQEventPublisher(ch)
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, eventPublisher, logrus.StandardLogger(), metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logrus.StandardLogger(), metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logrus.StandardLogger(), metricClient),
		},
	}
}

func newMongoClient() *mongo.Client {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return c
}

package service

import (
	"context"
	"github.com/rigoncs/gorder/common/metrics"
	"github.com/rigoncs/gorder/stock/adapters"
	"github.com/rigoncs/gorder/stock/app"
	"github.com/rigoncs/gorder/stock/app/query"
	"github.com/rigoncs/gorder/stock/infrastructure/integration"
	"github.com/rigoncs/gorder/stock/infrastructure/persistent"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	//stockRepo := adapters.NewMemoryStockRepository()
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripeAPI := integration.NewStripeAPI()
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeAPI, logrus.StandardLogger(), metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logrus.StandardLogger(), metricsClient),
		},
	}
}

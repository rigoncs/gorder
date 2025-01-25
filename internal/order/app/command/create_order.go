package command

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rigoncs/gorder/common/broker"
	"github.com/rigoncs/gorder/common/convertor"
	"github.com/rigoncs/gorder/common/decorator"
	"github.com/rigoncs/gorder/common/entity"
	"github.com/rigoncs/gorder/common/logging"
	"github.com/rigoncs/gorder/order/app/query"
	domain "github.com/rigoncs/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/status"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp.Channel,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	if stockGRPC == nil {
		panic("nil stockGRPC")
	}
	if channel == nil {
		panic("nil channel")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC, channel: channel},
		logger,
		metricClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	var err error
	defer logging.WhenCommandExecute(ctx, "CreateOrderHandler", cmd, err)

	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", broker.EventOrderCreated))
	defer span.End()

	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerID, validItems)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, pendingOrder)
	if err != nil {
		return nil, err
	}

	err = broker.PublishEvent(ctx, broker.PublishEventReq{
		Channel:  c.channel,
		Routing:  broker.Direct,
		Queue:    broker.EventOrderCreated,
		Exchange: "",
		Body:     o,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "publish event error q.Name=%s", broker.EventOrderCreated)
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, status.Convert(err).Err()
	}
	return convertor.NewItemConvertor().ProtosToEntities(resp.Items), nil
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	var res []*entity.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, entity.NewItemWithQuantity(id, quantity))
	}
	return res
}

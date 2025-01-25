package command

import (
	"context"
	"github.com/rigoncs/gorder/common/convertor"
	"github.com/rigoncs/gorder/common/decorator"
	"github.com/rigoncs/gorder/common/entity"
	"github.com/rigoncs/gorder/common/logging"
	"github.com/rigoncs/gorder/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *entity.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	var err error
	defer logging.WhenCommandExecute(ctx, "CreatePaymentHandler", cmd, err)

	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	newOrder, err := entity.NewValidOrder(
		cmd.Order.ID,
		cmd.Order.CustomerID,
		"waiting_for_payment",
		link,
		cmd.Order.Items,
	)
	if err != nil {
		return "", err
	}
	err = c.orderGRPC.UpdateOrder(ctx, convertor.NewOrderConvertor().EntityToProto(newOrder))
	return link, err
}

func NewCreatePaymentHandler(
	processor domain.Processor,
	orderGRPC OrderService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		createPaymentHandler{processor: processor, orderGRPC: orderGRPC},
		logger,
		metricClient)
}

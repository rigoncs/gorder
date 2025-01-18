package consumer

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rigoncs/gorder/common/broker"
	"github.com/rigoncs/gorder/order/app"
	"github.com/rigoncs/gorder/order/app/command"
	domain "github.com/rigoncs/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{app: app}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderPaied, true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.QueueBind(q.Name, "", broker.EventOrderPaied, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("error unmarshal msg.body into domain.order, err = %v", err)
		_ = msg.Nack(false, false)
		return
	}
	_, err := c.app.Commands.UpdateOrder.Handle(context.Background(), command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
		// TODO: retry
		return
	}
	_ = msg.Ack(false)
	logrus.Info("order consume paid event success!")
}

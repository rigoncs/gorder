package mq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rigoncs/gorder/common/broker"
	domain "github.com/rigoncs/gorder/order/domain/order"
)

type RabbitMQEventPublisher struct {
	Channel *amqp.Channel
}

func NewRabbitMQEventPublisher(channel *amqp.Channel) *RabbitMQEventPublisher {
	return &RabbitMQEventPublisher{Channel: channel}
}

func (r *RabbitMQEventPublisher) Publish(ctx context.Context, event domain.DomainEvent) error {
	return broker.PublishEvent(ctx, broker.PublishEventReq{
		Channel:  r.Channel,
		Routing:  broker.Direct,
		Queue:    event.Dest,
		Exchange: "",
		Body:     event.Data,
	})
}

func (r *RabbitMQEventPublisher) Broadcast(ctx context.Context, event domain.DomainEvent) error {
	return broker.PublishEvent(ctx, broker.PublishEventReq{
		Channel:  r.Channel,
		Routing:  broker.FanOut,
		Queue:    event.Dest,
		Exchange: "",
		Body:     event.Data,
	})
}

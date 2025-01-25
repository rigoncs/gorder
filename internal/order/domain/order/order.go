package order

import (
	"errors"
	"fmt"
	"github.com/rigoncs/gorder/common/consts"
	"github.com/rigoncs/gorder/common/entity"
	"slices"
)

func (o *Order) UpdatePaymentLink(paymentLink string) error {
	o.PaymentLink = paymentLink
	return nil
}

func (o *Order) UpdateItems(items []*entity.Item) error {
	o.Items = items
	return nil
}

func (o *Order) UpdateStatus(to string) error {
	if !o.isValidStatusTransition(to) {
		return fmt.Errorf("cannot transit from '%s' to '%s'", o.Status, to)
	}
	o.Status = to
	return nil
}

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if status == "" {
		return nil, errors.New("empty status")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func NewPendingOrder(customerId string, items []*entity.Item) (*Order, error) {
	if customerId == "" {
		return nil, errors.New("empty customerID")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	return &Order{
		CustomerID: customerId,
		Status:     consts.OrderStatusPending,
		Items:      items,
	}, nil
}

func (o *Order) isValidStatusTransition(to string) bool {
	switch o.Status {
	default:
		return false
	case consts.OrderStatusPending:
		return slices.Contains([]string{consts.OrderStatusPending}, to)
	case consts.OrderStatusWaitingForPayment:
		return slices.Contains([]string{consts.OrderStatusWaitingForPayment}, to)
	case consts.OrderStatusPaid:
		return slices.Contains([]string{consts.OrderStatusPaid}, to)
	}
}

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rigoncs/gorder/common"
	client "github.com/rigoncs/gorder/common/client/order"
	"github.com/rigoncs/gorder/order/app"
	"github.com/rigoncs/gorder/order/app/command"
	"github.com/rigoncs/gorder/order/app/dto"
	"github.com/rigoncs/gorder/order/app/query"
	"github.com/rigoncs/gorder/order/convertor"
	"net/http"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
		err  error
	)
	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	resp = convertor.NewOrderConvertor().EntityToClient(o)

}

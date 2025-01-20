// Package order provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package order

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	CustomerID string             `json:"customerID"`
	Items      []ItemWithQuantity `json:"items"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// Item defines model for Item.
type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PriceID  string `json:"priceID"`
	Quantity int32  `json:"quantity"`
}

// ItemWithQuantity defines model for ItemWithQuantity.
type ItemWithQuantity struct {
	Id       string `json:"id"`
	Quantity int32  `json:"quantity"`
}

// Order defines model for Order.
type Order struct {
	CustomerID  string `json:"customerID"`
	Id          string `json:"id"`
	Items       []Item `json:"items"`
	PaymentLink string `json:"paymentLink"`
	Status      string `json:"status"`
}

// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest

package app

import "github.com/rigoncs/gorder/stock/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	CheckIfItemsInStock query.CheckIfItemsInStockHandler
	GetItems            query.GetItemsHandler
}

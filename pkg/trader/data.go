package trader

import "time"

const (
	OrderTypeSell = "SELL"
)

type Ticker struct {
	Id           int64
	Symbol       string
	BestBidPrice float64
	BestBidQty   float64
	BestAskPrice float64
	BestAskQty   float64
}

type Order struct {
	Id         string
	Type       string
	Symbol     string
	Size       float64
	Price      float64
	IsFinished bool
	FinalPrice float64
}

type OrderFill struct {
	Id        string
	Order     *Order
	Ticker    *Ticker
	CreatedAt time.Time
	Size      float64
	Price     float64
}

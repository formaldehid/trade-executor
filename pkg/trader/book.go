package trader

import (
	"github.com/adshao/go-binance/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type book struct {
	cfg    *Config
	db     *DB
	stream *Stream
	ticker *Ticker
	order  *Order
	fills  []*OrderFill
}

type Book interface {
	GetStream() *Stream
	CreateSellOrder(symbol string, size float64, price float64) (Order, error)
	FillOrder(order *Order, ticker *Ticker, fills []*OrderFill) (bool, *OrderFill)
}

func NewBook(cfg *Config) (Book, error) {
	var st Stream
	var bk *book

	db, err := NewDB(cfg)
	if err != nil {
		return nil, err
	}

	handler := func(event *binance.WsBookTickerEvent) {
		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
				st.Stop()
			}
		}()
		func() {
			t, err := ParseWsBookTickerEvent(event)
			if err != nil {
				panic(err)
			}

			bk.ticker = &t

			err = db.InsertTicker(bk.ticker)
			if err != nil {
				panic(err)
			}

			success, of := bk.FillOrder(bk.order, bk.ticker, bk.fills)
			if of != nil {
				err = db.InsertOrderFill(of)
				if err != nil {
					panic(err)
				}
				const l string = "order fill: id=%s size=%f price=%f ticker=%d timestamp=%s"
				log.Infof(l, of.Id, of.Size, of.Price, of.Ticker.Id, of.CreatedAt)
			}
			if success == true {
				err = db.FinalizeOrder(bk.order)
				if err != nil {
					panic(err)
				}
				st.Stop()
			}
		}()
	}

	st, err = NewStream(cfg, handler)
	if err != nil {
		return nil, err
	}

	bk = &book{
		cfg:    cfg,
		db:     &db,
		stream: &st,
	}

	return bk, nil
}

func ParseWsBookTickerEvent(event *binance.WsBookTickerEvent) (Ticker, error) {
	bidPrice, err := strconv.ParseFloat(event.BestBidPrice, 64)
	if err != nil {
		return Ticker{}, err
	}

	bidQty, err := strconv.ParseFloat(event.BestBidQty, 64)
	if err != nil {
		return Ticker{}, err
	}

	askPrice, err := strconv.ParseFloat(event.BestAskPrice, 64)
	if err != nil {
		return Ticker{}, err
	}

	askQty, err := strconv.ParseFloat(event.BestAskQty, 64)
	if err != nil {
		return Ticker{}, err
	}

	t := Ticker{
		Id:           event.UpdateID,
		Symbol:       event.Symbol,
		BestBidPrice: bidPrice,
		BestBidQty:   bidQty,
		BestAskPrice: askPrice,
		BestAskQty:   askQty,
	}

	return t, nil
}

func (b *book) GetStream() *Stream {
	return b.stream
}

func (b *book) CreateSellOrder(symbol string, size float64, price float64) (Order, error) {
	b.order = &Order{
		Id:     uuid.NewString(),
		Type:   OrderTypeSell,
		Symbol: symbol,
		Size:   size,
		Price:  price,
	}

	err := (*b.db).InsertOrder(b.order)
	if err != nil {
		return Order{}, err
	}

	return *b.order, nil
}

func (b *book) FillOrder(order *Order, ticker *Ticker, fills []*OrderFill) (bool, *OrderFill) {
	var bestSize float64 = 0
	var bestPrice float64 = 0
	var sizeF float64 = 0
	var priceF float64 = 0
	var bestBidQty float64 = 0

	for i := 0; i < len(fills); i++ {
		if fills[i].Price > bestPrice {
			bestPrice = fills[i].Price
			bestSize = fills[i].Size
		} else if fills[i].Price == bestPrice {
			bestSize += fills[i].Size
		}

		sizeF += fills[i].Size
		priceF += fills[i].Price
	}

	sizeR := order.Size - sizeF
	if sizeR == 0 {
		return true, nil
	}

	if ticker.BestBidPrice >= order.Price && (bestPrice <= ticker.BestBidPrice || len(fills) == 0) {
		if len(fills) == 0 || bestPrice < ticker.BestBidPrice {
			bestBidQty = ticker.BestBidQty
		} else if bestPrice == ticker.BestBidPrice {
			bestBidQty = ticker.BestBidQty - bestSize
			if bestBidQty <= 0 {
				return false, nil
			}
		}

		if bestBidQty > 0 {
			if sizeR <= bestBidQty {
				of := &OrderFill{
					Id:        uuid.NewString(),
					Order:     order,
					Ticker:    ticker,
					CreatedAt: time.Now(),
					Size:      sizeR,
					Price:     ticker.BestBidPrice,
				}
				b.fills = append(fills, of)
				return true, of
			} else {
				of := &OrderFill{
					Id:        uuid.NewString(),
					Order:     order,
					Ticker:    ticker,
					CreatedAt: time.Now(),
					Size:      bestBidQty,
					Price:     ticker.BestBidPrice,
				}
				b.fills = append(fills, of)
				return false, of
			}
		}
	}

	return false, nil
}

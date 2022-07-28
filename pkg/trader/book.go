package trader

import (
	"github.com/adshao/go-binance/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	FillOrder(order *Order, ticker *Ticker, fills []*OrderFill) bool
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
				st.Done()
				db.CloseDB()
			}
		}()
		func() {
			t, err := ParseWsBookTickerEvent(event)
			if err != nil {
				log.Error(err)
				panic(err)
			}

			bk.ticker = &t

			err = db.InsertTicker(bk.ticker)
			if err != nil {
				log.Error(err)
				panic(err)
			}

			success := bk.FillOrder(bk.order, bk.ticker, bk.fills)
			if success == true {
				st.Done()
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

func (b *book) FillOrder(order *Order, ticker *Ticker, fills []*OrderFill) bool {
	var sizeF float64 = 0

	for i := 0; i < len(fills); i++ {
		sizeF += fills[i].Size
	}

	sizeR := order.Size - sizeF
	if sizeR == 0 {
		return true
	}

	return false
}

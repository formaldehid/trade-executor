package trader

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type bookTestSuite struct {
	suite.Suite
}

func TestBook(t *testing.T) {
	suite.Run(t, new(bookTestSuite))
}

func (s *bookTestSuite) TestFillOrder() {
	order := &Order{
		Id:     uuid.NewString(),
		Type:   "SELL",
		Symbol: "BNBUSDT",
		Size:   25,
		Price:  42,
	}

	ticker1 := &Ticker{
		Id:           400900217,
		Symbol:       "BNBUSDT",
		BestBidPrice: 40,
		BestBidQty:   10,
		BestAskPrice: 41,
		BestAskQty:   10,
	}

	ticker2 := &Ticker{
		Id:           400900223,
		Symbol:       "BNBUSDT",
		BestBidPrice: 42,
		BestBidQty:   5,
		BestAskPrice: 43,
		BestAskQty:   10,
	}

	ticker3 := &Ticker{
		Id:           400900235,
		Symbol:       "BNBUSDT",
		BestBidPrice: 42.5,
		BestBidQty:   30,
		BestAskPrice: 43,
		BestAskQty:   10,
	}

	var orderFills []*OrderFill

	cfg := NewConfig()
	cfg.Symbol = "BNBUSDT"

	book, _ := NewBook(cfg)

	res, _ := book.FillOrder(order, ticker1, orderFills)
	assert.Equal(s.T(), res, false)

	res, ofs2 := book.FillOrder(order, ticker2, orderFills)
	assert.Equal(s.T(), res, false)
	assert.Equal(s.T(), ofs2.Price, ticker2.BestBidPrice)
	assert.Equal(s.T(), ofs2.Size, ticker2.BestBidQty)
	assert.Equal(s.T(), ofs2.Ticker.Id, ticker2.Id)

	res, ofs3 := book.FillOrder(order, ticker3, append(orderFills, ofs2))
	assert.Equal(s.T(), res, true)
	assert.Equal(s.T(), ofs3.Price, 42.5)
	assert.Equal(s.T(), ofs3.Size, 20.0)
	assert.Equal(s.T(), ofs3.Ticker.Id, ticker3.Id)
}

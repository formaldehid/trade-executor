package trader

import (
	"github.com/adshao/go-binance/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

type stream struct {
	doneC <-chan struct{}
	stopC <-chan struct{}
}

type Stream interface {
	Stop()
	Done()
}

func NewStream(cfg *Config, handler func(event *binance.WsBookTickerEvent)) (Stream, error) {
	errorHandler := func(err error) {
		log.Error(err)
	}

	doneC, stopC, err := binance.WsBookTickerServe(strings.ToUpper(cfg.Symbol), handler, errorHandler)
	if err != nil {
		return nil, err
	}

	s := &stream{
		doneC: doneC,
		stopC: stopC,
	}

	return s, nil
}

func (s *stream) Stop() {
	<-s.stopC
}

func (s *stream) Done() {
	<-s.doneC
}

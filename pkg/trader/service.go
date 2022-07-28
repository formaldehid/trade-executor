package trader

import log "github.com/sirupsen/logrus"

type service struct {
	cfg  *Config
	book *Book
}

type Service interface {
	Listen()
	Sell(symbol string, size float64, price float64) error
}

func NewService(cfg *Config) (Service, error) {
	book, err := NewBook(cfg)
	if err != nil {
		return nil, err
	}

	s := &service{
		cfg:  cfg,
		book: &book,
	}

	return s, nil
}

func (s *service) Listen() {
	(*(*s.book).GetStream()).Done()
}

func (s *service) Sell(symbol string, size float64, price float64) error {
	order, err := (*s.book).CreateSellOrder(symbol, size, price)
	if err != nil {
		return err
	}

	log.Infof("order created: id=%s", order.Id)

	return nil
}

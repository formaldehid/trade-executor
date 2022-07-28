package main

import (
	"github.com/formaldehid/trade-executor/pkg/trader"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.Infof("log level: %d", log.GetLevel())
}

func main() {
	app := &cli.App{
		Name:  "sell",
		Usage: "create market sell order",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "symbol",
				Required: true,
				Usage:    "ticker symbol",
			},
			&cli.Float64Flag{
				Name:     "size",
				Required: true,
				Usage:    "size of trade",
			},
			&cli.Float64Flag{
				Name:     "price",
				Required: true,
				Usage:    "price of trade",
			},
		},
		Action: func(c *cli.Context) error {
			symbol := c.String("symbol")
			size := c.Float64("size")
			price := c.Float64("price")

			log.Infof("symbol: %s", symbol)
			log.Infof("size: %f", size)
			log.Infof("price: %f", price)

			cfg := trader.NewConfig()
			cfg.Symbol = symbol

			service, err := trader.NewService(cfg)
			if err != nil {
				log.Error(err)
				return err
			}

			err = service.Sell(symbol, size, price)
			if err != nil {
				log.Error(err)
				return err
			}

			service.Listen()

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package freestuff

import (
	"github.com/caarlos0/env"
	"log"
)

type config struct {
	PriceThreshold float64 `env:"FREESTUFF_PRICE_THRESHOLD" envDefault:"10"`
	Webhook        string  `env:"FREESTUFF_WEBHOOK"`
}

var Config config

func init() {
	if err := env.Parse(&Config); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

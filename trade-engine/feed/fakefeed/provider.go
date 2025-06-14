package fakefeed

import (
	"math/rand/v2"
)

type Provider struct {
	tickers map[string]int64 // map[symbol]ticker
}

func NewProvider(tickers map[string]int64) *Provider {
	return &Provider{
		tickers: tickers,
	}
}

func (p *Provider) RetrievePrices() (map[string]Quote, error) {
	quotes := make(map[string]Quote)

	for symbol, lastPrice := range p.tickers {
		// make price a random variation between 0.1% and 0.5% and floor it to 2 decimal places
		change := int64(float64(lastPrice)*rand.Float64()*0.04 + 0.01*(rand.Float64()*2-1))

		// get 1 or -1 to determine if the price should go up or down
		direction := int64(rand.IntN(2)*2 - 1)

		lastPrice += change * direction

		p.tickers[symbol] = lastPrice

		quotes[symbol] = Quote{
			LastTradePrice: uint64(lastPrice),
		}
	}

	return quotes, nil
}

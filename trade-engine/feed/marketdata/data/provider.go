package data

import "context"

type MarketPricesProvider interface {
	GetMarketPriceBySymbol(ctx context.Context, symbol string) (MarketPrice, error)

	UpsertMarketPrice(ctx context.Context, marketPrice MarketPrice) (*MarketPrice, error)
	GetSymbols(ctx context.Context) ([]Symbols, error)
}

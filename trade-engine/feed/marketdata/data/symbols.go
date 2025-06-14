package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Symbols struct {
	bun.BaseModel      `bun:"table:symbols,alias:s"`
	ID                 string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Symbol             string    `bun:"symbol"`
	Exchange           string    `bun:"exchange"`
	LastTradePrice     uint64    `bun:"last_trade_price"`
	LastTradeTimestamp time.Time `bun:"last_trade_timestamp"`
}

// GetSymbols is a function that returns all symbols.
func (dp DataProvider) GetSymbols(ctx context.Context) ([]Symbols, error) {
	var symbols []Symbols

	err := dp.db.NewSelect().Model(&symbols).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return symbols, nil
}

func (dp DataProvider) GetSymbol(ctx context.Context, symbol string) (*Symbols, error) {
	var symbols Symbols

	err := dp.db.NewSelect().Model(&symbols).Where("symbol = ?", symbol).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &symbols, nil
}
func (dp DataProvider) UpsertSymbol(ctx context.Context, symbol Symbols) (*Symbols, error) {
	_, err := dp.db.NewInsert().Model(&symbol).On("CONFLICT (symbol) DO UPDATE").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

func (dp DataProvider) UpdateSymbol(ctx context.Context, symbol Symbols) (*Symbols, error) {
	_, err := dp.db.NewUpdate().Model(&symbol).Where("symbol = ?", symbol.Symbol).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

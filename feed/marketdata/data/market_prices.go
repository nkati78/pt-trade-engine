package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type MarketPrice struct {
	bun.BaseModel `bun:"table:market_prices,alias:mp"`

	ID             string    `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Symbol         string    `bun:"symbol"`
	Price          uint64    `bun:"price"`
	StartingPrice  uint64    `bun:"starting_price"`
	TradeDate      string    `bun:"trade_date,omitzero"`
	TodayHigh      uint64    `bun:"today_high"`
	TodayLow       uint64    `bun:"today_low"`
	YesterdayClose uint64    `bun:"yesterday_close"`
	YesterdayOpen  uint64    `bun:"yesterday_open"`
	YesterdayHigh  uint64    `bun:"yesterday_high"`
	YesterdayLow   uint64    `bun:"yesterday_low"`
	UpdatedAt      time.Time `bun:"updated_at"`

	ReferenceID string `bun:"reference_id"`
}

func (dp DataProvider) GetMarketPrice(ctx context.Context, symbol string) (*MarketPrice, error) {
	var marketPrice MarketPrice

	err := dp.db.NewSelect().Model(&marketPrice).Where("symbol = ?", symbol).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &marketPrice, nil
}

func (dp DataProvider) UpsertMarketPrice(ctx context.Context, marketPrice MarketPrice) (*MarketPrice, error) {
	_, err := dp.db.NewInsert().Model(&marketPrice).On("CONFLICT (symbol) DO UPDATE").Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &marketPrice, nil
}

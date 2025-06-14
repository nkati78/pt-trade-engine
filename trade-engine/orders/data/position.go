package data

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Position struct {
	bun.BaseModel `bun:"table:positions,alias:p"`

	ID         string `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	AvgPrice   uint64 `bun:"average_price"`
	Quantity   uint32 `bun:"quantity"`
	Direction  string `bun:"direction"`
	ProfitLoss int64  `bun:"profit_loss"`
	Symbol     string `bun:"symbol"`
	UserID     string `bun:"user_id"`
	OrderID    string `bun:"order_id"`
	Status     string `bun:"status"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

// GetUserPositions is a function that returns all orders for a user.
func (dp DataProvider) GetUserPositions(ctx context.Context, userID string) ([]Position, error) {
	var orders []Position

	err := dp.db.NewSelect().Model(&orders).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (dp DataProvider) CreatePosition(ctx context.Context, position Position) (*Position, error) {
	_, err := dp.db.NewInsert().Model(&position).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &position, nil
}

func (dp DataProvider) UpdatePosition(ctx context.Context, position Position) (*Position, error) {
	_, err := dp.db.NewUpdate().Model(&position).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &position, nil
}

func (dp DataProvider) UpdatePositionByOrderID(ctx context.Context, position Position) (*Position, error) {
	_, err := dp.db.NewUpdate().Model(&position).Where("order_id = ?", position.OrderID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &position, nil
}

func (dp DataProvider) UpdatePositionsBySymbol(ctx context.Context, symbol string, newPrice uint64) error {
	_, err := dp.db.NewUpdate().Model((*Position)(nil)).Set("profit_loss = ABS((? - average_price) * quantity)", newPrice).Where("symbol = ?", symbol).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

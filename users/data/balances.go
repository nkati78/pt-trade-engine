package data

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

type Balance struct {
	bun.BaseModel `bun:"table:balances,alias:b"`
	ID            string  `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID        string  `bun:"user_id"`
	Balance       float64 `bun:"balance"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

// GetUserBalance is a function that returns all balances for a user.
func (dp DataProvider) GetUserBalance(ctx context.Context, userID string) (*Balance, error) {
	var balance Balance

	err := dp.db.NewSelect().Model(&balance).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (dp DataProvider) CreateBalance(ctx context.Context, userID string) (*Balance, error) {
	balance := Balance{
		UserID:  userID,
		Balance: 10_000,
	}

	_, err := dp.db.NewInsert().Model(&balance).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

// TODO: Josh, added userID here for balance update
func (dp DataProvider) UpdateUserBalance(ctx context.Context, userID string, balance Balance) (*Balance, error) {
	_, err := dp.db.NewUpdate().Model(&balance).Where("user_id = ?", userID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

package data

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID       string `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Price    uint64 `bun:"price"`
	Quantity uint32 `bun:"quantity"`
	Side     string `bun:"side"`
	Type     string `bun:"type"`
	Symbol   string `bun:"symbol"`
	UserID   string `bun:"user_id"`
	Status   string `bun:"status"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

// GetUserOrders is a function that returns all orders for a user.
func (dp DataProvider) GetUserOrders(ctx context.Context, userID string) ([]Order, error) {
	var orders []Order

	err := dp.db.NewSelect().Model(&orders).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (dp DataProvider) GetOrdersBySymbol(ctx context.Context, symbol string) ([]Order, error) {
	var orders []Order

	err := dp.db.NewSelect().Model(&orders).Where("symbol = ?", symbol).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (dp DataProvider) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var order Order

	err := dp.db.NewSelect().Model(&order).Where("id = ?", orderID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (dp DataProvider) CreateOrder(ctx context.Context, order Order) (*Order, error) {
	res, err := dp.db.NewInsert().Model(&order).Exec(ctx)
	fmt.Println("HERE BE RESPONSE res: ", res)
	if err != nil {
		fmt.Println("AGH! Error creating order:", err)
		return nil, err
	}

	return &order, nil
}

func (dp DataProvider) UpdateOrder(ctx context.Context, order Order) (*Order, error) {
	_, err := dp.db.NewUpdate().Model(&order).OmitZero().WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *Order) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		o.CreatedAt = time.Now()
	case *bun.UpdateQuery:
		o.UpdatedAt = time.Now()
	}
	return nil
}

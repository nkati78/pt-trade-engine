package data

import "context"

type OrderProvider interface {
	GetUserOrders(ctx context.Context, userID string) ([]Order, error)
	GetOrdersBySymbol(ctx context.Context, symbol string) ([]Order, error)

	GetOrder(ctx context.Context, orderID string) (*Order, error)
	CreateOrder(ctx context.Context, order Order) (*Order, error)
	UpdateOrder(ctx context.Context, order Order) (*Order, error)

	GetUserPositions(ctx context.Context, userID string) ([]Position, error)
	CreatePosition(ctx context.Context, position Position) (*Position, error)
	UpdatePosition(ctx context.Context, position Position) (*Position, error)
	UpdatePositionByOrderID(ctx context.Context, position Position) (*Position, error)
	UpdatePositionsBySymbol(ctx context.Context, symbol string, newPrice uint64) error
}

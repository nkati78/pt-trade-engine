package orders

import (
	"time"

	"github.com/paper-thesis/trade-engine/orders/data"

	"github.com/google/uuid"
)

type TradeSide string

const (
	Buy  TradeSide = "buy"
	Sell TradeSide = "sell"
)

type OrderType string

const (
	Market OrderType = "market"
	Limit  OrderType = "limit"
)

type OrderStatus string

const (
	Open     OrderStatus = "open"
	Closed   OrderStatus = "closed"
	Filled   OrderStatus = "filled"
	Canceled OrderStatus = "canceled"
)

type Order struct {
	OrderID        string
	Price          uint64
	Quantity       uint32
	QuantityFilled uint32
	Filled         bool
	UserID         string
	Symbol         string
	Timestamp      time.Time
	Side           TradeSide
	Type           OrderType
	FilledTime     time.Time
	Status         OrderStatus
}

type Fill struct {
	FillID       string
	OrderID      string
	FillPrice    int64
	FillQuantity uint32
	FillTime     time.Time
}

func NewOrder(price uint64, quantity uint32, userID string, symbol string, side TradeSide, orderType OrderType) *Order {
	return &Order{
		OrderID:   uuid.New().String(),
		Price:     price,
		Quantity:  quantity,
		UserID:    userID,
		Symbol:    symbol,
		Timestamp: time.Now(),
		Side:      side,
		Type:      orderType,
		Status:    Open,
	}
}

type OrderBook []*Order

var orderBooks = make(map[string]OrderBook)

func (o *Order) Fill(quantity uint32) {
	o.QuantityFilled += quantity
	if o.QuantityFilled == o.Quantity {
		o.Filled = true
		o.FilledTime = time.Now()
		o.Status = Closed
	}
}

func OrderBookInsert(orders []*Order, order *Order) []*Order {
	if len(orders) == 0 {
		return append(orders, order)
	}

	for i, o := range orders {
		if o.Timestamp.Compare(order.Timestamp) > 0 {
			return append(orders[:i], append([]*Order{order}, orders[i:]...)...)
		}
	}

	return append(orders, order)
}

func OrderBookRemove(orders []*Order, orderID string) []*Order {
	for i, o := range orders {
		if o.OrderID == orderID {
			return append(orders[:i], orders[i+1:]...)
		}
	}

	return orders
}

func GetOrderByPrice(orders []*Order, price uint64) *Order {
	for _, o := range orders {
		if o.Price == price {
			return o
		}
	}

	return nil
}

func (o Order) ToDB() data.Order {
	return data.Order{
		ID:       o.OrderID,
		Price:    o.Price,
		Quantity: o.Quantity,
		Side:     string(o.Side),
		Type:     string(o.Type),
		Symbol:   o.Symbol,
		UserID:   o.UserID,
		Status:   string(o.Status),
	}
}

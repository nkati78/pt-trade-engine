package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paper-thesis/trade-engine/orders/data"
	"github.com/paper-thesis/trade-engine/users"
	userData "github.com/paper-thesis/trade-engine/users"
)

type OrderRequest struct {
	Price    uint64 `json:"price"`
	Quantity uint32 `json:"quantity"`
	Side     string `json:"side"`
	Type     string `json:"type"`
	Symbol   string `json:"symbol"`
}

type OrderResponse struct {
	OrderID   string `json:"order_id"`
	Timestamp string `json:"timestamp"`
}

func NewOrderService(dal data.OrderProvider, userService users.UserService) OrderService {
	return OrderService{
		dal:         dal,
		userService: userService,
	}
}

type OrderService struct {
	dal         data.OrderProvider
	userService userData.UserService
}

func (os OrderService) CreateOrder(ctx context.Context, userID string, orderRequest OrderRequest) (*OrderResponse, error) {
	if len(userID) == 0 {
		return nil, errors.New("user is not authenticated")
	}

	if orderRequest.Type == "" {
		return nil, errors.New("order type is required")
	}

	order := NewOrder(
		orderRequest.Price,
		orderRequest.Quantity,
		userID,
		orderRequest.Symbol,
		TradeSide(orderRequest.Side),
		OrderType(orderRequest.Type),
	)

	// Insert the order into the order book
	orders := OrderBookInsert(orderBooks[order.Symbol], order)
	orderBooks[order.Symbol] = orders
	order.Status = Open

	fmt.Println("Order book updated: ", orderBooks)

	orderDB := order.ToDB()

	createdOrder, err := os.dal.CreateOrder(ctx, orderDB)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{OrderID: createdOrder.ID, Timestamp: createdOrder.CreatedAt.Format("2006-01-02T15:04:05")}, nil
}

func (os OrderService) GetOrderBook(symbol string) OrderBook {
	return orderBooks[symbol]
}

func (os OrderService) GetOrderByUserID(ctx context.Context, userID string) ([]*Order, error) {
	orders, err := os.dal.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, nil
	}

	userOrders := make([]*Order, 0)
	for _, order := range orders {
		userOrders = append(userOrders, &Order{
			OrderID:   order.ID,
			Price:     order.Price,
			Quantity:  order.Quantity,
			UserID:    order.UserID,
			Symbol:    order.Symbol,
			Timestamp: order.CreatedAt,
			Status:    OrderStatus(order.Status),
			Side:      TradeSide(order.Side),
			Type:      OrderType(order.Type),
		})
	}

	return userOrders, nil
}

func (os OrderService) CancelOrder(orderID string) {
	for symbol, orders := range orderBooks {
		for _, order := range orders {
			if order.OrderID == orderID {
				orders = OrderBookRemove(orders, orderID)
				orderBooks[symbol] = orders
				return
			}
		}
	}
}

func (os OrderService) FillOrder(ctx context.Context, order *Order, price uint64) error {
	order.Filled = true
	order.FilledTime = time.Now()
	order.Status = Filled

	order.Price = price

	_, err := os.dal.UpdateOrder(ctx, order.ToDB())
	if err != nil {
		return err
	}

	position := data.Position{
		Quantity:   order.Quantity,
		Direction:  string(order.Side),
		ProfitLoss: 0,
		Symbol:     order.Symbol,
		UserID:     order.UserID,
		OrderID:    order.OrderID,
		Status:     string(Open),
		AvgPrice:   price,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err = os.dal.CreatePosition(ctx, position)
	if err != nil {
		return err
	}

	OrderBookRemove(orderBooks[order.Symbol], order.OrderID)

	//TODO: Where I am currently working Josh
	//balance, err := os.ud.GetUserBalance(ctx, order.UserID)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println(balance)
	//fmt.Println("ron")
	//
	//totalCost := float64(order.Quantity) * float64(price)
	//balance.Balance -= totalCost
	//
	//_, err = os.bd.UpdateUserBalance(ctx, order.UserID, *balance)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (os OrderService) UpdatePositionsBySymbol(ctx context.Context, symbol string, newPrice uint64) error {
	return os.dal.UpdatePositionsBySymbol(ctx, symbol, newPrice)
}

func (os OrderService) GetPositionsByUserID(ctx context.Context, userID string) ([]*Position, error) {
	positions, err := os.dal.GetUserPositions(ctx, userID)
	if err != nil {
		fmt.Println("Error getting user positions: ", err)
		return nil, err
	}

	userPositions := make([]*Position, 0)
	for _, position := range positions {
		userPositions = append(userPositions, &Position{
			ID:         position.ID,
			Quantity:   position.Quantity,
			Direction:  position.Direction,
			AvgPrice:   position.AvgPrice,
			ProfitLoss: position.ProfitLoss,
			Symbol:     position.Symbol,
			UserID:     position.UserID,
			OrderID:    position.OrderID,
			Status:     position.Status,
			CreatedAt:  position.CreatedAt,
			UpdatedAt:  position.UpdatedAt,
		})
	}

	return userPositions, nil
}

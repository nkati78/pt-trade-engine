package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paper-thesis/trade-engine/orders"
	"github.com/paper-thesis/trade-engine/security"
)

type OrderHandler struct {
	orderService orders.OrderService
}

type OrderResponse struct {
	OrderID   string `json:"orderId"`
	Price     uint64 `json:"price"`
	Quantity  uint32 `json:"quantity"`
	Side      string `json:"side"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	UpdatedAt string `json:"updatedAt"`
}

func NewOrderHandler(orderService orders.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (oh OrderHandler) HandleCreateOrder(c *gin.Context) {
	var orderRequest orders.OrderRequest

	userID := c.GetString(security.UserCtxKey)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("user is not authenticated"))
		return
	}

	// parse the request body
	err := c.BindJSON(&orderRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	order, err := oh.orderService.CreateOrder(c, userID, orderRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	orderResponse := OrderResponse{
		OrderID:   order.OrderID,
		Timestamp: order.Timestamp,
	}

	c.IndentedJSON(http.StatusCreated, orderResponse)
}

type OrderListResponse struct {
	Orders []orders.Order `json:"orders"`
}

func (oh OrderHandler) HandleGetOrders(c *gin.Context) {
	userID := c.GetString(security.UserCtxKey)
	order, err := oh.orderService.GetOrderByUserID(c, userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.IndentedJSON(http.StatusOK, order)
}

func (oh OrderHandler) HandleGetUserPositions(c *gin.Context) {
	userID := c.GetString(security.UserCtxKey)
	positions, err := oh.orderService.GetPositionsByUserID(c, userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.IndentedJSON(http.StatusOK, positions)
}

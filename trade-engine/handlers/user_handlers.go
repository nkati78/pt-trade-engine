package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/paper-thesis/trade-engine/orders"
	"github.com/paper-thesis/trade-engine/security"
	"github.com/paper-thesis/trade-engine/users"
	"github.com/paper-thesis/trade-engine/users/models"
)

type UserHandler struct {
	userService  users.UserService
	orderService orders.OrderService
	auth         security.Auth
}

func NewUserHandler(userService users.UserService, orderService orders.OrderService, auth security.Auth) *UserHandler {
	return &UserHandler{
		userService:  userService,
		orderService: orderService,
		auth:         auth,
	}
}

type CreateUserHandlerResponse struct {
	Token string                     `json:"token"`
	User  *models.CreateUserResponse `json:"user"`
}

// CreateUser creates a new user
func (uh UserHandler) CreateUser(c *gin.Context) (HTTPStatusCode, interface{}) {
	var userRequest models.CreateUserRequest

	err := c.BindJSON(&userRequest)
	if err != nil {
		return HTTPStatusBadRequest, HTTPError{Message: "Invalid request"}
	}

	user, err := uh.userService.CreateUser(c, userRequest)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	token, err := uh.auth.GenerateToken(user.ID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusCreated, CreateUserHandlerResponse{
		Token: token,
		User:  user,
	}
}

type LoginHandlerResponse struct {
	Token string      `json:"token"`
	User  *users.User `json:"user"`
}

// Login logs in a user
func (uh UserHandler) Login(c *gin.Context) (HTTPStatusCode, interface{}) {
	var loginRequest models.LoginUserRequest

	err := c.BindJSON(&loginRequest)
	if err != nil {
		return HTTPStatusBadRequest, HTTPError{Message: "Invalid request"}
	}

	user, err := uh.userService.Login(c, loginRequest)
	if err != nil {
		if err == users.IncorrectEmailOrPass {
			return HTTPStatusUnauthorized, HTTPError{Message: "Incorrect email or password"}
		}

		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	token, err := uh.auth.GenerateToken(user.ID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, LoginHandlerResponse{
		Token: token,
		User:  user,
	}
}

func (uh UserHandler) GetBalance(c *gin.Context) (HTTPStatusCode, interface{}) {
	userID := c.GetString(security.UserCtxKey)

	balance, err := uh.userService.GetUserBalance(c, userID)
	if err != nil {
		fmt.Println(err)
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	// get open positions for the user and calculate the total value
	positions, err := uh.orderService.GetPositionsByUserID(c, userID)
	if err != nil {
		fmt.Println(err)
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	totalValue := 0.0
	for _, position := range positions {
		if position.Status != string(orders.Open) {
			continue
		}

		totalValue += float64(position.ProfitLoss)
	}

	balance.Balance += totalValue

	return HTTPStatusOK, balance
}

func (uh UserHandler) GetUser(c *gin.Context) (HTTPStatusCode, interface{}) {
	userID := c.GetString(security.UserCtxKey)

	user, err := uh.userService.GetUser(c, userID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, user
}

func (uh UserHandler) GetUserWatchList(c *gin.Context) (HTTPStatusCode, interface{}) {
	userID := c.GetString(security.UserCtxKey)

	watchList, err := uh.userService.GetUserWatchList(c, userID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, watchList
}

func (uh UserHandler) AddSymbolWatchList(c *gin.Context) (HTTPStatusCode, interface{}) {
	userID := c.GetString(security.UserCtxKey)
	var watchListInput models.WatchListInput

	err := c.BindJSON(&watchListInput)
	if err != nil {
		return HTTPStatusBadRequest, HTTPError{Message: "Invalid request"}
	}

	watchtListResponse, err := uh.userService.CreateWatchlistSymbol(c, watchListInput.Symbol, userID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	fmt.Println(watchtListResponse)

	return HTTPStatusOK, watchtListResponse
}

func (uh UserHandler) RemoveSymbolWatchList(c *gin.Context) (HTTPStatusCode, interface{}) {
	userID := c.GetString(security.UserCtxKey)
	var watchListInput models.WatchListInput

	err := c.BindJSON(&watchListInput)
	if err != nil {
		return HTTPStatusBadRequest, HTTPError{Message: "Invalid request, missing symbol"}
	}

	err = uh.userService.DeleteWatchlistSymbol(c, watchListInput.Symbol, userID)
	if err != nil {
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, nil
}

package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/paper-thesis/trade-engine/users/data"
	"github.com/paper-thesis/trade-engine/users/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dal data.DataProvider
}

type LoginError error

var (
	IncorrectEmailOrPass LoginError = errors.New("incorrect email or password")
)

func NewUserService(dal data.DataProvider) UserService {
	return UserService{dal: dal}
}

func (us UserService) CreateUser(ctx context.Context, request models.CreateUserRequest) (*models.CreateUserResponse, error) {
	user := data.UserDB{
		Username:   request.Username,
		Email:      request.Email,
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		ProviderID: request.ProviderID,
	}

	existingUser, err := us.dal.GetUserByEmail(ctx, request.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	if request.Password != "" {
		// hash the password
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.PasswordHash = string(passwordHash)
		user.ConnectionType = data.UsernamePassword
	}

	userResult, err := us.dal.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	_, err = us.dal.CreateBalance(ctx, userResult.ID)
	if err != nil {
		return nil, err
	}

	return &models.CreateUserResponse{
		ID:        userResult.ID,
		Username:  userResult.Username,
		Email:     userResult.Email,
		FirstName: userResult.FirstName,
		LastName:  userResult.LastName,
	}, nil
}

func (us UserService) Login(ctx context.Context, request models.LoginUserRequest) (*User, error) {
	user, err := us.dal.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, IncorrectEmailOrPass
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return nil, IncorrectEmailOrPass
	}

	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (us UserService) GetUserBalance(ctx context.Context, userID string) (*models.BalanceResponse, error) {
	balance, err := us.dal.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.BalanceResponse{
		Balance: balance.Balance,
		UserID:  balance.UserID,
		ID:      balance.ID,
	}, nil
}

func (us UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	user, err := us.dal.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (us UserService) CreateWatchlistSymbol(ctx context.Context, symbol, userID string) (*WatchList, error) {
	count, err := us.dal.CountWatchList(ctx, userID)
	if err != nil {
		return nil, err
	}

	wl := data.WatchList{
		UserID:         userID,
		Symbol:         symbol,
		SequenceNumber: count + 1,
	}

	wlResult, err := us.dal.CreateWatchList(ctx, wl)
	if err != nil {
		return nil, err
	}

	return &WatchList{
		ID:             wlResult.ID,
		UserID:         wlResult.UserID,
		Symbol:         wlResult.Symbol,
		SequenceNumber: wlResult.SequenceNumber,
	}, nil
}

func (us UserService) GetUserWatchList(ctx context.Context, userID string) ([]WatchList, error) {
	wl, err := us.dal.GetUserWatchList(ctx, userID)
	if err != nil {
		return nil, err
	}

	var watchList []WatchList
	for _, w := range wl {
		watchList = append(watchList, WatchList{
			ID:             w.ID,
			UserID:         w.UserID,
			Symbol:         w.Symbol,
			SequenceNumber: w.SequenceNumber,
		})
	}

	return watchList, nil
}

func (us UserService) DeleteWatchlistSymbol(ctx context.Context, symbol, userID string) error {
	wl, err := us.dal.GetUserWatchList(ctx, userID)
	if err != nil {
		return err
	}

	for _, w := range wl {
		if w.Symbol == symbol {
			return us.dal.DeleteWatchList(ctx, w.ID)
		}
	}

	return nil
}

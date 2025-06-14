package data

import "context"

type UserProvider interface {
	GetUser(ctx context.Context, userID string) (*UserDB, error)
	CreateUser(ctx context.Context, user UserDB) (*UserDB, error)
	UpdateUser(ctx context.Context, user UserDB) (*UserDB, error)
	GetUserByUsername(ctx context.Context, username string) (*UserDB, error)
	GetUserByEmail(ctx context.Context, email string) (*UserDB, error)
	CreateWatchList(ctx context.Context, watchList WatchList) (*WatchList, error)
	GetUserWatchList(ctx context.Context, userID string) ([]WatchList, error)
	UpdateWatchList(ctx context.Context, watchList WatchList) (*WatchList, error)
	CountWatchList(ctx context.Context, userID string) (int, error)
	DeleteWatchList(ctx context.Context, id string) error
	GetUserBalance(ctx context.Context, userID string) (*Balance, error)
}

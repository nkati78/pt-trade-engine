package data

import (
	"context"
	"github.com/uptrace/bun"
	"time"
)

type DataProvider struct {
	db *bun.DB
}

type UserDB struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID             string         `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Username       string         `bun:"username"`
	PasswordHash   string         `bun:"password_hash"`
	Email          string         `bun:"email"`
	FirstName      string         `bun:"first_name"`
	LastName       string         `bun:"last_name"`
	ConnectionType ConnectionType `bun:"connection_type"`
	ProviderID     string         `bun:"provider_id"`

	CreatedAt time.Time `bun:"created_at"`
	UpdatedAt time.Time `bun:"updated_at"`
}

type ConnectionType string

const (
	UsernamePassword ConnectionType = "username-password"
	Facebook         ConnectionType = "facebook"
	Google           ConnectionType = "google"
)

func NewDataProvider(db *bun.DB) DataProvider {
	return DataProvider{db: db}
}

func (dp DataProvider) GetUser(ctx context.Context, userID string) (*UserDB, error) {
	var user UserDB

	err := dp.db.NewSelect().Model(&user).Where("id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dp DataProvider) GetUserByUsername(ctx context.Context, username string) (*UserDB, error) {
	var user UserDB

	err := dp.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dp DataProvider) GetUserByEmail(ctx context.Context, email string) (*UserDB, error) {
	var user UserDB

	err := dp.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dp DataProvider) CreateUser(ctx context.Context, user UserDB) (*UserDB, error) {
	_, err := dp.db.NewInsert().Model(&user).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dp DataProvider) UpdateUser(ctx context.Context, user UserDB) (*UserDB, error) {
	_, err := dp.db.NewUpdate().Model(&user).Where("id = ?", user.ID).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

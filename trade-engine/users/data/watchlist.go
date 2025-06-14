package data

import (
	"context"

	"github.com/uptrace/bun"
)

// watch list table
type WatchList struct {
	bun.BaseModel  `bun:"table:watchlist,alias:wl"`
	ID             string `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID         string `bun:"user_id"`
	Symbol         string `bun:"symbol"`
	SequenceNumber int    `bun:"sequence_number"`
}

// GetUserWatchList is a function that returns all watchlist for a user.
func (dp DataProvider) GetUserWatchList(ctx context.Context, userID string) ([]WatchList, error) {
	var watchList []WatchList

	err := dp.db.NewSelect().Model(&watchList).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return watchList, nil
}

func (dp DataProvider) CreateWatchList(ctx context.Context, watchList WatchList) (*WatchList, error) {
	_, err := dp.db.NewInsert().Model(&watchList).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &watchList, nil
}

func (dp DataProvider) CountWatchList(ctx context.Context, userID string) (int, error) {
	count, err := dp.db.NewSelect().Model((*WatchList)(nil)).Where("user_id = ?", userID).Count(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// update watchlist
func (dp DataProvider) AddWatchListSymbol(ctx context.Context, watchList WatchList) (*WatchList, error) {
	_, err := dp.db.NewUpdate().Model(&watchList).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &watchList, nil
}

func (dp DataProvider) DeleteWatchList(ctx context.Context, id string) error {
	_, err := dp.db.NewDelete().Model((*WatchList)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

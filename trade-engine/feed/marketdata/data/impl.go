package data

import "github.com/uptrace/bun"

type DataProvider struct {
	db *bun.DB
}

func NewDataProvider(db *bun.DB) DataProvider {
	return DataProvider{db: db}
}

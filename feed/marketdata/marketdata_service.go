package marketdata

import (
	"context"
	"fmt"
	"time"

	"github.com/paper-thesis/trade-engine/feed/marketdata/data"
)

type MarketDataService struct {
	dal data.DataProvider
}

func NewMarketDataService(dal data.DataProvider) MarketDataService {
	return MarketDataService{
		dal: dal,
	}
}

type MarketData struct {
	ID               string `json:"id"`
	Symbol           string `json:"symbol"`
	Price            uint64 `json:"price"`
	StartingPrice    uint64 `json:"startingPrice"`
	YesterdayClose   uint64 `json:"yesterdayClose"`
	YesterdayOpen    uint64 `json:"yesterdayOpen"`
	YesterdayHigh    uint64 `json:"yesterdayHigh"`
	YesterdayLow     uint64 `json:"yesterdayLow"`
	TodayOpen        uint64 `json:"todayOpen"`
	TodayHigh        uint64 `json:"todayHigh"`
	TodayLow         uint64 `json:"todayLow"`
	TradeDate        string `json:"tradeDate"`
	PriceChange      int64  `json:"priceChange"`
	PercentageChange string `json:"percentageChange"`
	UpdatedAt        string `json:"updated_at"`
}

type SymbolData struct {
	Symbol           string `json:"symbol"`
	Exchange         string `json:"exchange"`
	LastTradePrice   uint64 `json:"lastTradePrice"`
	PreviousDayClose uint64 `json:"previousDayClose"`
	TradeDate        string `json:"tradeDate"`
	OpenPrice        uint64 `json:"openPrice"`
}

func (m MarketDataService) GetMarketData(ctx context.Context, symbol string) (*MarketData, error) {
	// Get market data
	data, err := m.dal.GetMarketPrice(ctx, symbol)
	if err != nil {
		return nil, err
	}

	// calculate the price change
	priceChange := int64(data.Price - data.StartingPrice)

	var percentageChange string

	if data.StartingPrice == 0 {
		percentageChange = "100.00"
	} else {
		percentageChange = fmt.Sprintf("%.3f", float64(priceChange)/float64(data.StartingPrice))
		// round to 2 decimal places as a float

	}

	fmt.Println("trade date in db: ", data.TradeDate)

	return &MarketData{
		ID:               data.ID,
		Symbol:           data.Symbol,
		Price:            data.Price,
		StartingPrice:    data.StartingPrice,
		PriceChange:      priceChange,
		PercentageChange: percentageChange,
		TradeDate:        data.TradeDate,
		TodayHigh:        data.TodayHigh,
		TodayLow:         data.TodayLow,
		YesterdayClose:   data.YesterdayClose,
		YesterdayOpen:    data.YesterdayOpen,
		YesterdayHigh:    data.YesterdayHigh,
		YesterdayLow:     data.YesterdayLow,
		UpdatedAt:        data.UpdatedAt.String(),
	}, nil
}

func (m MarketDataService) UpsertMarketData(ctx context.Context, marketData MarketData) (*MarketData, error) {

	// Upsert market data
	data := data.MarketPrice{
		ID:             marketData.ID,
		Symbol:         marketData.Symbol,
		Price:          marketData.Price,
		TradeDate:      marketData.TradeDate,
		TodayHigh:      marketData.TodayHigh,
		TodayLow:       marketData.TodayLow,
		StartingPrice:  marketData.StartingPrice,
		YesterdayClose: marketData.YesterdayClose,
		YesterdayOpen:  marketData.YesterdayOpen,
		YesterdayHigh:  marketData.YesterdayHigh,
		UpdatedAt:      time.Now(),
	}

	//fmt.Println(data)

	res, err := m.dal.UpsertMarketPrice(ctx, data)
	if err != nil {
		return nil, err
	}

	fmt.Println("New Data: ", res)

	return &marketData, nil
}

func (m MarketDataService) GetSymbols(ctx context.Context) ([]SymbolData, error) {
	// Get symbols
	symbols, err := m.dal.GetSymbols(ctx)
	if err != nil {
		return nil, err
	}

	var symbolsList []SymbolData
	for _, symbol := range symbols {
		symbolsList = append(symbolsList, SymbolData{
			Symbol:   symbol.Symbol,
			Exchange: symbol.Exchange,
		})
	}

	return symbolsList, nil
}

func (m MarketDataService) GetSymbol(ctx context.Context, symbol string) (*SymbolData, error) {
	// Get symbols
	symbols, err := m.dal.GetSymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	return &SymbolData{
		Symbol:   symbols.Symbol,
		Exchange: symbols.Exchange,
	}, nil
}

func (m MarketDataService) UpsertSymbol(ctx context.Context, symbol SymbolData) (*SymbolData, error) {
	// Upsert symbol
	data := data.Symbols{
		Symbol:             symbol.Symbol,
		Exchange:           symbol.Exchange,
		LastTradePrice:     symbol.LastTradePrice,
		LastTradeTimestamp: time.Now(),
	}

	_, err := m.dal.UpsertSymbol(ctx, data)
	if err != nil {
		return nil, err
	}

	return &symbol, nil
}

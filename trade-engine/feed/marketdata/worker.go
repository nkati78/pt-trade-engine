package marketdata

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/paper-thesis/trade-engine/feed/fakefeed"
	"github.com/paper-thesis/trade-engine/orders"
)

// Worker is a struct that represents a worker that will be used to process the feed
type Worker struct {
	marketDataService MarketDataService
	orderService      orders.OrderService
}

// NewWorker creates a new worker
func NewWorker(marketDataService MarketDataService, orderService orders.OrderService) *Worker {
	return &Worker{
		marketDataService: marketDataService,
		orderService:      orderService,
	}
}

func getBearerToken() string {
	token := os.Getenv("BEARER_TOKEN")

	return token
}

// Start starts the worker
func (w *Worker) Start() {
	// marketdata := robinhood.NewProvider(getBearerToken(), map[string]string{"MSFT": "50810c35-d215-4866-9758-0ada4ac79ffa"})
	marketdata := fakefeed.NewProvider(map[string]int64{"MSFT": 23000, "AAPL": 23000, "GOOGL": 23000, "AMZN": 23000, "TSLA": 23000, "FB": 23000, "NVDA": 23000, "NFLX": 23000, "AMD": 23000, "INTC": 23000, "CSCO": 23000, "IBM": 23000, "ORCL": 23000, "QCOM": 23000, "TXN": 23000, "AVGO": 23000, "ADBE": 23000, "CRM": 23000})

	for {
		quotes, err := marketdata.RetrievePrices()
		if err != nil {
			panic(err)
		}

		for symbol, value := range quotes {
			// get market data from the database
			marketData, err := w.marketDataService.GetMarketData(context.Background(), symbol)
			if err != nil {
				fmt.Println("\n\n\n\nError getting market data: \n\n\n\n\n", err)
				marketData = &MarketData{
					Symbol: symbol,
					Price:  value.LastTradePrice,
				}
			}

			marketData.Price = value.LastTradePrice

			// first time we are getting the market data
			if marketData.StartingPrice == 0 {
				marketData.StartingPrice = value.LastTradePrice
				marketData.YesterdayOpen = value.LastTradePrice - 1000
				marketData.TodayHigh = value.LastTradePrice
				marketData.TodayLow = value.LastTradePrice
				marketData.YesterdayClose = value.LastTradePrice - 500
				marketData.YesterdayHigh = value.LastTradePrice + 400
				marketData.YesterdayLow = value.LastTradePrice - 1000
				marketData.TradeDate = time.Now().Format("2006-01-02")
			}

			// check if new trading day
			if marketData.TradeDate != time.Now().Format("2006-01-02") {
				fmt.Println("\n==================\nNew trading day\n==================\n")
				fmt.Println("Old trade date: ", marketData.TradeDate)
				fmt.Println("New trade date: ", time.Now().Format("2006-01-02"))
				marketData.TradeDate = time.Now().Format("2006-01-02")
				marketData.StartingPrice = value.LastTradePrice
				marketData.TodayHigh = value.LastTradePrice
				marketData.TodayLow = value.LastTradePrice
				marketData.YesterdayClose = marketData.Price
				marketData.YesterdayHigh = marketData.TodayHigh
				marketData.YesterdayLow = marketData.TodayLow
			}

			// check if we are todays high or low
			if value.LastTradePrice < marketData.TodayLow || marketData.TodayLow == 0 {
				marketData.TodayLow = value.LastTradePrice
			}

			if value.LastTradePrice > marketData.TodayHigh || marketData.TodayHigh == 0 {
				marketData.TodayHigh = value.LastTradePrice
			}

			// check if we are the starting price
			if marketData.StartingPrice == 0 {
				marketData.StartingPrice = value.LastTradePrice
			}

			_, err = w.marketDataService.UpsertMarketData(context.Background(), *marketData)
			if err != nil {
				fmt.Println(err)
			}

			newPrice := value.LastTradePrice

			orderBook := w.orderService.GetOrderBook(symbol)
			for _, order := range orderBook {
				fmt.Println("Checking order: ", order)
				if order.Type == orders.Market && order.Status == orders.Open {
					w.orderService.FillOrder(context.Background(), order, newPrice)
				}
			}

			err = w.orderService.UpdatePositionsBySymbol(context.Background(), symbol, newPrice)
			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

package marketdata

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/paper-thesis/trade-engine/feed/robinhood"
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
	fmt.Println("Starting market data worker with Robinhood feed...")
	marketdata := robinhood.NewProvider(getBearerToken(), map[string]string{"MSFT": "50810c35-d215-4866-9758-0ada4ac79ffa", "AAPL": "450dfc6d-5510-4d40-abfb-f633b7d9be3e", "IBM": "ae2e4ada-197d-42a4-825c-aff01cc3a8dd", "GOOGL": "54db869e-f7d5-45fb-88f1-8d7072d4c8b2", "NFLX": "81733743-965a-4d93-b87a-6973cb9efd34", "NVDA": "a4ecd608-e7b4-4ff3-afa5-f77ae7632dfb"})
	//marketdata := fakefeed.NewProvider(map[string]int64{"MSFT": 23000, "AAPL": 23000, "GOOGL": 23000, "AMZN": 23000, "TSLA": 23000, "FB": 23000, "NVDA": 23000, "NFLX": 23000, "AMD": 23000, "INTC": 23000, "CSCO": 23000, "IBM": 23000, "ORCL": 23000, "QCOM": 23000, "TXN": 23000, "AVGO": 23000, "ADBE": 23000, "CRM": 23000})

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
				marketData.StartingPrice = value.LastClosePrice
				marketData.TodayHigh = value.LastTradePrice
				marketData.TodayLow = value.LastTradePrice
				marketData.YesterdayClose = value.LastClosePrice
				marketData.TradeDate = value.TradingDay
				marketData.TodayOpen = value.LastClosePrice
				marketData.TradeDate = value.TradingDay
			}

			// check if new trading day
			if marketData.TradeDate != value.TradingDay {
				fmt.Println("\n==================\nNew trading day\n==================\n")
				fmt.Println("Old trade date: ", marketData.TradeDate)
				fmt.Println("New trade date: ", value.TradingDay)
				marketData.TradeDate = value.TradingDay
				marketData.StartingPrice = value.LastClosePrice
				marketData.TodayHigh = value.LastTradePrice
				marketData.TodayLow = value.LastTradePrice
				marketData.YesterdayClose = marketData.Price
				marketData.YesterdayHigh = marketData.TodayHigh
				marketData.YesterdayLow = marketData.TodayLow
				marketData.TodayOpen = value.LastClosePrice
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
				marketData.StartingPrice = value.LastClosePrice
			}

			_, err = w.marketDataService.UpsertMarketData(context.Background(), *marketData)
			if err != nil {
				fmt.Println(err)
			}

			newPrice := value.LastTradePrice

			orderBook := w.orderService.GetOrderBook(symbol)
			for _, order := range orderBook {
				// fmt.Println("Checking order: ", order)
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

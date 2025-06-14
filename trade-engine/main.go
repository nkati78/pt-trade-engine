package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/paper-thesis/trade-engine/db"
	"github.com/paper-thesis/trade-engine/feed/marketdata"
	feedData "github.com/paper-thesis/trade-engine/feed/marketdata/data"
	"github.com/paper-thesis/trade-engine/orders"
	"github.com/paper-thesis/trade-engine/orders/data"
	"github.com/paper-thesis/trade-engine/users"
	userData "github.com/paper-thesis/trade-engine/users/data"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func run() int {
	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dsn = "postgres://postgres:QYs1Ecdtv1xvycyo7bGX@paper-thesis.cje8aqmy09mu.us-east-1.rds.amazonaws.com:5432/postgres"
	}

	database := db.NewDatabaseConnection(dsn)

	m, err := migrate.New(
		"file://db/migrations",
		dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		fmt.Print(err)
	}

	userProvider := userData.NewDataProvider(database)
	userService := users.NewUserService(userProvider)

	orderProvider := data.NewDataProvider(database)
	orderService := orders.NewOrderService(orderProvider, userService)

	marketDataService := marketdata.NewMarketDataService(feedData.NewDataProvider(database))
	marketDataWorker := marketdata.NewWorker(marketDataService, orderService)

	go func() {
		marketDataWorker.Start()
	}()

	/*
		go func() {
			StartDataFeed()
		}()
	*/

	if err := StartServer(orderService, userService, marketDataService); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("Server closed under request")
			return 0
		} else {
			return 1
		}
	}

	log.Println("Server exiting")

	return 0
}

func main() {
	os.Exit(run())
}

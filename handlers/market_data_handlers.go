package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/paper-thesis/trade-engine/feed/marketdata"
)

type MarketDataHandler struct {
	marketDataService marketdata.MarketDataService
}

func NewMarketDataHandler(marketDataService marketdata.MarketDataService) *MarketDataHandler {
	return &MarketDataHandler{
		marketDataService: marketDataService,
	}
}

func (mh MarketDataHandler) HandleGetMarketData(c *gin.Context) (HTTPStatusCode, interface{}) {
	symbol := c.Param("symbol")

	if symbol == "" {
		return HTTPStatusBadRequest, HTTPError{Message: "Invalid request"}
	}

	marketPrice, err := mh.marketDataService.GetMarketData(c, symbol)
	if err != nil {
		fmt.Println(err)
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, marketPrice
}

func (mh MarketDataHandler) HandleGetSymbols(c *gin.Context) (HTTPStatusCode, interface{}) {
	symbols, err := mh.marketDataService.GetSymbols(c)
	if err != nil {
		fmt.Println(err)
		return HTTPStatusInternalServerError, HTTPError{Message: "Internal server error"}
	}

	return HTTPStatusOK, symbols
}

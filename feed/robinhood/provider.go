// robinhood api provider
package robinhood

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/paper-thesis/trade-engine/feed/fakefeed"
)

type Provider struct {
	tickers        map[string]string // map[symbol]ticker
	baseURL        string
	defaultHeaders http.Header
	client         *http.Client
	tickerIds      string
}

const (
	QUOTES_ROUTE = "/quotes/"
	BOUNDS       = "24_5"
	BASE_RH_URL  = "https://api.robinhood.com/marketdata"
)

func NewProvider(secret string, tickers map[string]string) *Provider {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	defaultHeaders := make(http.Header)

	defaultHeaders.Add("Content-Type", "application/json")
	defaultHeaders.Add("Authorization", "Bearer "+secret)
	defaultHeaders.Add("Accept", "*/*")
	defaultHeaders.Add("X-Timezone-Id", "America/Chicago")
	defaultHeaders.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")

	// tickerIds variable that is all the values in the tickers map joined by comma
	tickerIds := ""
	for _, id := range tickers {
		tickerIds += id + ","
	}
	tickerIds = tickerIds[:len(tickerIds)-1] // remove trailing comma

	return &Provider{
		defaultHeaders: defaultHeaders,
		tickers:        tickers,
		baseURL:        BASE_RH_URL,
		client:         client,
		tickerIds:      tickerIds,
	}
}

func (p *Provider) RetrievePrices() (map[string]fakefeed.Quote, error) {
	url := fmt.Sprintf("%s/?bounds=%s&ids=%s&include_inactive=true", p.baseURL+QUOTES_ROUTE, BOUNDS, p.tickerIds)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = p.defaultHeaders

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	quotesResponse := QuotesResponse{}
	if jsonErr := json.NewDecoder(res.Body).Decode(&quotesResponse); jsonErr != nil {
		return nil, jsonErr
	}

	quotes := make(map[string]fakefeed.Quote)
	for _, quote := range quotesResponse.Results {
		priceFloat, err := strconv.ParseFloat(quote.LastTradePrice, 64)
		var price uint64
		if err != nil {
			price = 0
			fmt.Println("Error parsing price for ", quote.Symbol, ": ", err)
		} else {
			price = uint64(priceFloat * 100) // Convert dollars to cents
		}

		var lastClosePrice uint64
		prevCloseFloat, err := strconv.ParseFloat(quote.PreviousClose, 64)
		if err != nil {
			lastClosePrice = 0
			fmt.Println("Error parsing previous close price for ", quote.Symbol, ": ", err)
		} else {
			lastClosePrice = uint64(prevCloseFloat * 100) // Convert dollars to cents
		}

		quotes[quote.Symbol] = fakefeed.Quote{
			LastTradePrice: price,
			LastClosePrice: lastClosePrice,
			TradingDay:     quote.VenueLastTradeTime.Format("2006-01-02"),
		}
	}

	return quotes, nil

}

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

	return &Provider{
		defaultHeaders: defaultHeaders,
		tickers:        tickers,
		baseURL:        BASE_RH_URL,
		client:         client,
	}
}

func (p *Provider) RetrievePrices() (map[string]fakefeed.Quote, error) {
	url := fmt.Sprintf("%s/%s/?bounds=%s&include_inactive=true", p.baseURL+QUOTES_ROUTE, p.tickers["MSFT"], BOUNDS)

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

	msftQuote := Quote{}
	if jsonErr := json.NewDecoder(res.Body).Decode(&msftQuote); jsonErr != nil {
		return nil, jsonErr
	}

	quotes := make(map[string]fakefeed.Quote)
	price, err := strconv.ParseUint(msftQuote.LastTradePrice, 10, 64)
	if err != nil {
		price = 0
		fmt.Println("Error parsing price: ", err)
	}
	quotes["MSFT"] = fakefeed.Quote{
		LastTradePrice: price,
	}

	return quotes, nil

}

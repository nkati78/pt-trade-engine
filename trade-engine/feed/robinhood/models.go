package robinhood

import "time"

type Quote struct {
	AskPrice                    string    `json:"ask_price"`
	AskSize                     int       `json:"ask_size"`
	VenueAskTime                time.Time `json:"venue_ask_time"`
	BidPrice                    string    `json:"bid_price"`
	BidSize                     int       `json:"bid_size"`
	VenueBidTime                time.Time `json:"venue_bid_time"`
	LastTradePrice              string    `json:"last_trade_price"`
	VenueLastTradeTime          time.Time `json:"venue_last_trade_time"`
	LastExtendedHoursTradePrice string    `json:"last_extended_hours_trade_price"`
	LastNonRegTradePrice        string    `json:"last_non_reg_trade_price"`
	VenueLastNonRegTradeTime    time.Time `json:"venue_last_non_reg_trade_time"`
	PreviousClose               string    `json:"previous_close"`
	AdjustedPreviousClose       string    `json:"adjusted_previous_close"`
	PreviousCloseDate           string    `json:"previous_close_date"`
	Symbol                      string    `json:"symbol"`
	TradingHalted               bool      `json:"trading_halted"`
	HasTraded                   bool      `json:"has_traded"`
	LastTradePriceSource        string    `json:"last_trade_price_source"`
	LastNonRegTradePriceSource  string    `json:"last_non_reg_trade_price_source"`
	UpdatedAt                   time.Time `json:"updated_at"`
	Instrument                  string    `json:"instrument"`
	InstrumentID                string    `json:"instrument_id"`
	State                       string    `json:"state"`
}

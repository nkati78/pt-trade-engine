package fakefeed

type Quote struct {
	LastTradePrice uint64 `json:"last_trade_price"`
	LastClosePrice uint64 `json:"last_close_price"`
	TradingDay     string `json:"trading_day"`
}

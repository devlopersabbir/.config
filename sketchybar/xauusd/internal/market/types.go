package market

import "encoding/json"

// BybitWsMessage represents a generic message from Bybit V5 WebSocket.
type BybitWsMessage struct {
	Topic string          `json:"topic"`
	Type  string          `json:"type"`
	Ts    int64           `json:"ts"`
	Data  json.RawMessage `json:"data"`
}

// BybitTickerData represents the ticker update payload.
type BybitTickerData struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	Bid1Price string `json:"bid1Price"`
	Ask1Price string `json:"ask1Price"`
}

// BybitTradeData represents an individual executed trade.
type BybitTradeData struct {
	Price string `json:"p"`
	Side  string `json:"S"`
	Size  string `json:"v"`
	Time  int64  `json:"T"`
}

// BybitKlineData represents real-time kline/candlestick updates.
type BybitKlineData struct {
	Start     int64  `json:"start"`
	End       int64  `json:"end"`
	Interval  string `json:"interval"`
	Open      string `json:"open"`
	Close     string `json:"close"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Volume    string `json:"volume"`
	Turnover  string `json:"turnover"`
	Confirm   bool   `json:"confirm"`
	Timestamp int64  `json:"timestamp"`
}

// BybitOrderbookData represents orderbook depth update.
type BybitOrderbookData struct {
	Symbol string     `json:"s"`
	Bids   [][]string `json:"b"`
	Asks   [][]string `json:"a"`
}

// CandleBar represents a single OHLC candlestick for charting.
type CandleBar struct {
	Open  float64
	High  float64
	Low   float64
	Close float64
}

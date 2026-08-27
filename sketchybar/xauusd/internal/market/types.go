package market

import "encoding/json"

// BinanceCombined represents the wrapper for combined Binance WebSocket streams.
type BinanceCombined struct {
	Stream string          `json:"stream"`
	Data   json.RawMessage `json:"data"`
}

// BinanceBookTicker represents the real-time best bid and ask order book stream.
type BinanceBookTicker struct {
	UpdateID int64  `json:"u"`
	Symbol   string `json:"s"`
	BidPrice string `json:"b"`
	BidQty   string `json:"B"`
	AskPrice string `json:"a"`
	AskQty   string `json:"A"`
}

// BinanceTrade represents an individual executed market trade.
type BinanceTrade struct {
	Event     string `json:"e"`
	Symbol    string `json:"s"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	TradeTime int64  `json:"T"`
}

// BinanceKlineStream represents real-time kline/candlestick updates.
type BinanceKlineStream struct {
	Event string `json:"e"`
	Kline struct {
		OpenTime float64 `json:"t"`
		Open     string  `json:"o"`
		High     string  `json:"h"`
		Low      string  `json:"l"`
		Close    string  `json:"c"`
		Volume   string  `json:"v"`
		IsClosed bool    `json:"x"`
	} `json:"k"`
}

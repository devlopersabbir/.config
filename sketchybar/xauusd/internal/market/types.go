package market

// CandleBar represents a single OHLC candlestick for charting.
type CandleBar struct {
	Timestamp int64   `json:"timestamp"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
}

// TVQuoteData represents parsed live quote fields from TradingView quote stream.
type TVQuoteData struct {
	Lp        *float64 `json:"lp,omitempty"`
	Ch        *float64 `json:"ch,omitempty"`
	Chp       *float64 `json:"chp,omitempty"`
	Bid       *float64 `json:"bid,omitempty"`
	Ask       *float64 `json:"ask,omitempty"`
	Open      *float64 `json:"open_price,omitempty"`
	High      *float64 `json:"high_price,omitempty"`
	Low       *float64 `json:"low_price,omitempty"`
	PrevClose *float64 `json:"prev_close_price,omitempty"`
}

// ScannerResponse represents TradingView CFD scanner response.
type ScannerResponse struct {
	TotalCount int `json:"totalCount"`
	Data       []struct {
		S string        `json:"s"`
		D []interface{} `json:"d"`
	} `json:"data"`
}

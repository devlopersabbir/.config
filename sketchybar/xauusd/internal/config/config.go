package config

import "time"

const (
	// Symbol is the target asset: Gold / US Dollar Perpetual (XAUUSDT)
	Symbol = "XAUUSDT"

	// WsURL is the Bybit V5 Public Linear WebSocket endpoint for sub-second live market ticks
	WsURL = "wss://stream.bybit.com/v5/public/linear"

	// OhlcURL is the Bybit V5 REST API for 5m candle analytical data & historical chart
	OhlcURL = "https://api.bybit.com/v5/market/kline?category=linear&symbol=XAUUSDT&interval=5&limit=30"

	// TickerURL is the Bybit V5 REST API for 24h ticker snapshot
	TickerURL = "https://api.bybit.com/v5/market/tickers?category=linear&symbol=XAUUSDT"

	// CacheFile is the local cache path for SketchyBar plugin scripts
	CacheFile = "/tmp/sketchybar_gold_price"

	// AnalysisRefresh is the background interval for refreshing 5m OHLC analytical data
	AnalysisRefresh = 30 * time.Second

	// MinRenderInterval is the throttling interval between SketchyBar CLI updates
	MinRenderInterval = 100 * time.Millisecond

	// NumCandleBars is the number of individual candlestick bars displayed in the chart
	NumCandleBars = 18

	// UI Theme Colors (Catppuccin Macchiato Palette)
	ColorGreen   = "0xffa6e3a1"
	ColorRed     = "0xfff38ba8"
	ColorNeutral = "0xffcdd6f4"
	ColorMuted   = "0xff737994"
	ColorGold    = "0xfff9e2af"
)

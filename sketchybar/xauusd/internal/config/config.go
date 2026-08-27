package config

import "time"

const (
	// Symbol is the target asset: Spot Gold backed 1:1 by London Good Delivery gold bullion
	Symbol = "PAXGUSDT"

	// WsURL is the Binance WebSocket combined stream for real-time order book ticks, live trades, and 15m klines
	WsURL = "wss://stream.binance.com:9443/stream?streams=paxgusdt@bookTicker/paxgusdt@trade/paxgusdt@kline_15m"

	// OhlcURL is the Binance REST API for 15m candle analytical data & historical sparkline
	OhlcURL = "https://api.binance.com/api/v3/klines?symbol=PAXGUSDT&interval=15m&limit=30"

	// CacheFile is the local cache path for SketchyBar plugin scripts
	CacheFile = "/tmp/sketchybar_gold_price"

	// AnalysisRefresh is the background interval for refreshing 15m OHLC analytical data
	AnalysisRefresh = 30 * time.Second

	// MinRenderInterval is the throttling interval between SketchyBar CLI updates to keep CPU at ~0%
	MinRenderInterval = 150 * time.Millisecond

	// UI Theme Colors (Catppuccin Macchiato Palette)
	ColorGreen   = "0xffa6e3a1"
	ColorRed     = "0xfff38ba8"
	ColorNeutral = "0xffcdd6f4"
)

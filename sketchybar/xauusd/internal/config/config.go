package config

import "time"

const (
	// Target Asset: Gold Spot / U.S. Dollar (OANDA:XAUUSD)
	Symbol        = "OANDA:XAUUSD"
	SymbolDisplay = "XAU/USD"

	// WsURL is the TradingView live real-time WebSocket endpoint
	WsURL = "wss://data.tradingview.com/socket.io/websocket"

	// ScannerURL is the TradingView REST CFD scanner endpoint for analytical snapshots & fallback
	ScannerURL = "https://scanner.tradingview.com/cfd/scan"

	// CacheFile is the local cache path for SketchyBar plugin scripts
	CacheFile = "/tmp/sketchybar_gold_price"

	// AnalysisRefresh is the background interval for refreshing analytical data
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

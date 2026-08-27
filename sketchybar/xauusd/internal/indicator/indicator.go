package indicator

import (
	"xauusd/internal/config"
	"xauusd/internal/market"
)

// Unicode blocks representing increasing heights for sparkline bars.
var sparkBlocks = []string{
	" ", "▂", "▃", "▄", "▅", "▆", "▇", "█",
}

// CandleBarView holds the rendered Unicode character and color for an individual candle bar.
type CandleBarView struct {
	Char  string
	Color string
}

// ComputeCandleBars calculates dynamic heights and bullish/bearish colors for each individual candle.
func ComputeCandleBars(candles []market.CandleBar, currentPrice, currentOpen float64, count int) []CandleBarView {
	if count <= 0 {
		count = config.NumCandleBars
	}

	result := make([]CandleBarView, count)

	if len(candles) == 0 {
		for i := 0; i < count; i++ {
			result[i] = CandleBarView{Char: " ", Color: config.ColorMuted}
		}
		return result
	}

	// Work with the most recent `count` candles
	working := candles
	if len(working) > count {
		working = working[len(working)-count:]
	}

	// Update the latest active bar with live price
	if len(working) > 0 && currentPrice > 0 {
		working[len(working)-1].Close = currentPrice
		if currentOpen > 0 {
			working[len(working)-1].Open = currentOpen
		}
	}

	// Find min and max price across the window
	min := working[0].Close
	max := working[0].Close

	for _, bar := range working {
		if bar.Close < min {
			min = bar.Close
		}
		if bar.Close > max {
			max = bar.Close
		}
		if bar.Open < min {
			min = bar.Open
		}
		if bar.Open > max {
			max = bar.Open
		}
	}

	rng := max - min

	// Pad with leading empty bars if we have fewer candles than count
	offset := count - len(working)
	for i := 0; i < offset; i++ {
		result[i] = CandleBarView{Char: " ", Color: config.ColorMuted}
	}

	for i, bar := range working {
		level := 3
		if rng > 0 {
			level = int(((bar.Close - min) / rng) * 7.0)
		}
		if level < 0 {
			level = 0
		}
		if level > 7 {
			level = 7
		}

		color := config.ColorGreen
		if bar.Close < bar.Open {
			color = config.ColorRed
		}

		result[offset+i] = CandleBarView{
			Char:  sparkBlocks[level],
			Color: color,
		}
	}

	return result
}

// CalculateTrend analyzes price momentum across historical candles to classify market trend.
func CalculateTrend(candles []market.CandleBar) (trend string, icon string, color string) {
	if len(candles) < 6 {
		return "RANGE", "↔", config.ColorMuted
	}

	half := len(candles) / 2
	var oldSum, newSum float64

	for i := 0; i < half; i++ {
		oldSum += candles[i].Close
	}
	for i := half; i < len(candles); i++ {
		newSum += candles[i].Close
	}

	oldAvg := oldSum / float64(half)
	newAvg := newSum / float64(len(candles)-half)

	if oldAvg == 0 {
		return "RANGE", "↔", config.ColorMuted
	}

	change := ((newAvg - oldAvg) / oldAvg) * 100.0

	if change > 0.03 {
		return "BULL", "↑", config.ColorGreen
	}
	if change < -0.03 {
		return "BEAR", "↓", config.ColorRed
	}

	return "RANGE", "↔", config.ColorMuted
}

// CalculateCandleChange calculates the percentage difference relative to the 15m candle open.
func CalculateCandleChange(price, candleOpen float64) float64 {
	if candleOpen <= 0 {
		return 0
	}
	return ((price - candleOpen) / candleOpen) * 100.0
}

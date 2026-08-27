package indicator

import "strings"

// Blocks used for rendering dynamic unicode sparkline bars.
var sparkBlocks = []string{
	" ", "▂", "▃", "▄", "▅", "▆", "▇", "█",
}

// BuildSparkline converts an array of price values into a unicode sparkline string.
func BuildSparkline(values []float64) string {
	if len(values) == 0 {
		return "────────"
	}

	// Keep the latest 24 points to fit the bar comfortably
	if len(values) > 24 {
		values = values[len(values)-24:]
	}

	min := values[0]
	max := values[0]

	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	rng := max - min
	var builder strings.Builder

	for _, v := range values {
		level := 3
		if rng > 0 {
			level = int(((v - min) / rng) * 7)
		}
		if level < 0 {
			level = 0
		}
		if level > 7 {
			level = 7
		}
		builder.WriteString(sparkBlocks[level])
	}

	return builder.String()
}

// CalculateTrend analyzes price history moving averages to classify market momentum.
func CalculateTrend(history []float64) (trend string, icon string) {
	if len(history) < 6 {
		return "RANGE", "↔"
	}

	half := len(history) / 2
	var oldSum, newSum float64

	for i := 0; i < half; i++ {
		oldSum += history[i]
	}
	for i := half; i < len(history); i++ {
		newSum += history[i]
	}

	oldAvg := oldSum / float64(half)
	newAvg := newSum / float64(len(history)-half)

	if oldAvg == 0 {
		return "RANGE", "↔"
	}

	change := ((newAvg - oldAvg) / oldAvg) * 100.0

	if change > 0.03 {
		return "BULL", "↑"
	}
	if change < -0.03 {
		return "BEAR", "↓"
	}

	return "RANGE", "↔"
}

// CalculateCandleChange calculates the percentage difference relative to the 15m candle open.
func CalculateCandleChange(price, candleOpen float64) float64 {
	if candleOpen <= 0 {
		return 0
	}
	return ((price - candleOpen) / candleOpen) * 100.0
}

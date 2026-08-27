package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"xauusd/internal/config"
	"xauusd/internal/indicator"
	"xauusd/internal/market"
)

// Renderer controls SketchyBar UI updates and local file caching.
type Renderer struct {
	state *market.State
}

// NewRenderer creates a new Renderer instance.
func NewRenderer(state *market.State) *Renderer {
	return &Renderer{
		state: state,
	}
}

// Start runs the throttled render loop to provide smooth sub-second updates without CPU overhead.
func (r *Renderer) Start(ctx context.Context) {
	ticker := time.NewTicker(config.MinRenderInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			snap := r.state.PrepareRender()
			if snap.ShouldRender {
				r.Render(snap)
			}
		}
	}
}

// Render formats and publishes the latest market state to SketchyBar with individual bar coloring.
func (r *Renderer) Render(snap market.Snapshot) {
	pctChange := snap.ChangePct
	if pctChange == 0 && snap.CandleOpen > 0 {
		pctChange = indicator.CalculateCandleChange(snap.Price, snap.CandleOpen)
	}

	// Determine direction arrow based on price change
	arrow := "▲"
	if pctChange < -0.0001 {
		arrow = "▼"
	} else if pctChange > 0.0001 {
		arrow = "▲"
	} else {
		arrow = "•"
	}

	// Change color: only % and arrow change color for bearish/bullish
	changeColor := config.ColorNeutral
	if pctChange > 0.0001 {
		changeColor = config.ColorGreen
	} else if pctChange < -0.0001 {
		changeColor = config.ColorRed
	}

	// Trend indicator and color
	trend, trendIcon, trendColor := indicator.CalculateTrend(snap.Candles)

	// Compute individual candle bars with green/red colors
	candleBars := indicator.ComputeCandleBars(snap.Candles, snap.Price, snap.CandleOpen, config.NumCandleBars)

	priceStr := fmt.Sprintf("$%.2f", snap.Price)
	changeStr := fmt.Sprintf("%+.2f%% %s", pctChange, arrow)
	trendStr := fmt.Sprintf("%s %s", trendIcon, trend)

	// Build atomic SketchyBar update command
	args := []string{
		"--set", "gold.price", "label=" + priceStr, "label.color=" + config.ColorNeutral,
		"--set", "gold.change", "label=" + changeStr, "label.color=" + changeColor,
	}

	for i, bar := range candleBars {
		itemName := fmt.Sprintf("gold.bar.%d", i+1)
		args = append(args, "--set", itemName, "label="+bar.Char, "label.color="+bar.Color)
	}

	args = append(args, "--set", "gold.trend", "label="+trendStr, "label.color="+trendColor)

	cmd := exec.Command("sketchybar", args...)
	_ = cmd.Run()

	// Write cache
	sparkText := ""
	for _, bar := range candleBars {
		sparkText += bar.Char
	}
	cacheLine := fmt.Sprintf("%s %s │ %s │ %s|%s", priceStr, changeStr, sparkText, trendStr, changeColor)
	_ = os.WriteFile(config.CacheFile, []byte(cacheLine), 0644)
}

// RenderDisconnected displays the disconnected / offline status on SketchyBar.
func (r *Renderer) RenderDisconnected() {
	price, _ := r.state.GetPrice()

	priceStr := "offline"
	if price > 0 {
		priceStr = fmt.Sprintf("$%.2f", price)
	}

	args := []string{
		"--set", "gold.price", "label=" + priceStr, "label.color=" + config.ColorRed,
		"--set", "gold.change", "label=---", "label.color=" + config.ColorMuted,
	}

	for i := 1; i <= config.NumCandleBars; i++ {
		itemName := fmt.Sprintf("gold.bar.%d", i)
		args = append(args, "--set", itemName, "label= ", "label.color="+config.ColorMuted)
	}

	args = append(args, "--set", "gold.trend", "label=---", "label.color="+config.ColorMuted)

	cmd := exec.Command("sketchybar", args...)
	_ = cmd.Run()
}

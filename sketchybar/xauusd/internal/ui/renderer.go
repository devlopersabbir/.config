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

// Render formats and publishes the latest market state to SketchyBar.
func (r *Renderer) Render(snap market.Snapshot) {
	candleChange := indicator.CalculateCandleChange(snap.Price, snap.CandleOpen)

	// Determine label color
	color := config.ColorNeutral
	if candleChange > 0.0001 {
		color = config.ColorGreen
	} else if candleChange < -0.0001 {
		color = config.ColorRed
	}

	// Determine direction arrow
	arrow := "▲"
	if snap.Direction == "DOWN" {
		arrow = "▼"
	} else if snap.Direction == "UP" {
		arrow = "▲"
	} else if candleChange < -0.0001 {
		arrow = "▼"
	}

	// Compute analytical indicators
	sparkline := indicator.BuildSparkline(snap.History)
	trend, trendIcon := indicator.CalculateTrend(snap.History)

	priceStr := fmt.Sprintf("$%.2f", snap.Price)
	changeStr := fmt.Sprintf("%+.2f%%", candleChange)

	label := fmt.Sprintf(
		"%s %s %s │ %s │ %s %s",
		priceStr,
		changeStr,
		arrow,
		sparkline,
		trendIcon,
		trend,
	)

	r.UpdateSketchyBar(label, color)
	r.WriteCache(label, color)
}

// RenderDisconnected displays the disconnected / offline status on SketchyBar.
func (r *Renderer) RenderDisconnected() {
	price, _ := r.state.GetPrice()

	if price <= 0 {
		r.UpdateSketchyBar("🥇 Connecting...", config.ColorNeutral)
		return
	}

	label := fmt.Sprintf("$%.2f │ offline", price)
	r.UpdateSketchyBar(label, config.ColorRed)
}

// UpdateSketchyBar invokes the sketchybar CLI to update the gold item.
func (r *Renderer) UpdateSketchyBar(label, color string) {
	cmd := exec.Command(
		"sketchybar",
		"--set", "gold",
		"label="+label,
		"label.color="+color,
	)
	_ = cmd.Run()
}

// WriteCache writes the current label and color to disk for shell plugins.
func (r *Renderer) WriteCache(label, color string) {
	data := label + "|" + color
	_ = os.WriteFile(config.CacheFile, []byte(data), 0644)
}

package analysis

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"xauusd/internal/config"
	"xauusd/internal/market"
)

// Analyzer manages periodic OHLC candle fetching for historical analysis.
type Analyzer struct {
	state  *market.State
	client *http.Client
}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer(state *market.State) *Analyzer {
	return &Analyzer{
		state: state,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Start begins the background analytical polling loop.
func (a *Analyzer) Start(ctx context.Context) {
	// Immediate initial fetch
	a.Refresh(ctx)

	ticker := time.NewTicker(config.AnalysisRefresh)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.Refresh(ctx)
		}
	}
}

// Refresh fetches the latest 15m OHLC klines from the REST API.
func (a *Analyzer) Refresh(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.OhlcURL, nil)
	if err != nil {
		log.Printf("OHLC request creation failed: %v", err)
		return
	}

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("OHLC fetch error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("OHLC HTTP error: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Format: [ [openTime, open, high, low, close, volume, ...], ... ]
	var rawKlines [][]interface{}
	if err := json.Unmarshal(body, &rawKlines); err != nil {
		log.Printf("OHLC decode error: %v", err)
		return
	}

	if len(rawKlines) == 0 {
		return
	}

	var history []float64
	for _, bar := range rawKlines {
		if len(bar) >= 5 {
			if closeStr, ok := bar[4].(string); ok {
				if cVal, err := strconv.ParseFloat(closeStr, 64); err == nil {
					history = append(history, cVal)
				}
			}
		}
	}

	lastBar := rawKlines[len(rawKlines)-1]
	var currentOpen float64
	if openStr, ok := lastBar[1].(string); ok {
		if val, err := strconv.ParseFloat(openStr, 64); err == nil {
			currentOpen = val
		}
	}

	a.state.SetAnalysisHistory(currentOpen, history)
	log.Printf("15m analysis refreshed: open=%.2f, history_points=%d", currentOpen, len(history))
}

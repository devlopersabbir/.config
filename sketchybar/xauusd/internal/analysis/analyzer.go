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

	var candles []market.CandleBar
	for _, bar := range rawKlines {
		if len(bar) >= 5 {
			openStr, ok1 := bar[1].(string)
			highStr, ok2 := bar[2].(string)
			lowStr, ok3 := bar[3].(string)
			closeStr, ok4 := bar[4].(string)

			if ok1 && ok2 && ok3 && ok4 {
				oVal, err1 := strconv.ParseFloat(openStr, 64)
				hVal, err2 := strconv.ParseFloat(highStr, 64)
				lVal, err3 := strconv.ParseFloat(lowStr, 64)
				cVal, err4 := strconv.ParseFloat(closeStr, 64)

				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					candles = append(candles, market.CandleBar{
						Open:  oVal,
						High:  hVal,
						Low:   lVal,
						Close: cVal,
					})
				}
			}
		}
	}

	a.state.SetAnalysisCandles(candles)
	log.Printf("5m analysis refreshed: %d candles loaded", len(candles))
}

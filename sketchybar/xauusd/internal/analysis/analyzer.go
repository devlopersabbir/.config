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

// Analyzer manages periodic OHLC candle fetching from Bybit V5 REST API.
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

// Refresh fetches the latest 5m OHLC klines from Bybit V5 REST API.
func (a *Analyzer) Refresh(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.OhlcURL, nil)
	if err != nil {
		log.Printf("Bybit OHLC request creation failed: %v", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("Bybit OHLC fetch error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Bybit OHLC HTTP error: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var bybitResp struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Symbol string     `json:"symbol"`
			List   [][]string `json:"list"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &bybitResp); err != nil {
		log.Printf("Bybit OHLC decode error: %v", err)
		return
	}

	rawList := bybitResp.Result.List
	if len(rawList) == 0 {
		return
	}

	// Bybit returns klines in reverse chronological order (newest first).
	// We reverse them so the oldest is at index 0 and latest is at the end.
	var candles []market.CandleBar
	for i := len(rawList) - 1; i >= 0; i-- {
		bar := rawList[i]
		if len(bar) >= 5 {
			o, err1 := strconv.ParseFloat(bar[1], 64)
			h, err2 := strconv.ParseFloat(bar[2], 64)
			l, err3 := strconv.ParseFloat(bar[3], 64)
			c, err4 := strconv.ParseFloat(bar[4], 64)

			if err1 == nil && err2 == nil && err3 == nil && err4 == nil && o > 0 && c > 0 {
				candles = append(candles, market.CandleBar{
					Open:  o,
					High:  h,
					Low:   l,
					Close: c,
				})
			}
		}
	}

	if len(candles) > 0 {
		a.state.SetAnalysisCandles(candles)
		log.Printf("Bybit 5m analysis refreshed: %d candles loaded for %s", len(candles), config.Symbol)
	}
}

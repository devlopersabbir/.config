package analysis

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"xauusd/internal/config"
	"xauusd/internal/market"
)

// Analyzer manages periodic OHLC candle & CFD scanner data from TradingView REST API.
type Analyzer struct {
	state  *market.State
	client *http.Client
}

// NewAnalyzer creates a new Analyzer instance.
func NewAnalyzer(state *market.State) *Analyzer {
	return &Analyzer{
		state: state,
		client: &http.Client{
			Timeout: 6 * time.Second,
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

// Refresh fetches the latest market snapshot from TradingView CFD scanner.
func (a *Analyzer) Refresh(ctx context.Context) {
	payload := []byte(`{
		"symbols": {"tickers": ["` + config.Symbol + `"]},
		"columns": ["close", "change", "change_abs", "open", "high", "low", "Recommend.All", "open|5", "high|5", "low|5", "close|5"]
	}`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.ScannerURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("TradingView scanner request creation failed: %v", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		log.Printf("TradingView scanner fetch error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("TradingView scanner HTTP error: %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var sr market.ScannerResponse
	if err := json.Unmarshal(body, &sr); err != nil {
		log.Printf("TradingView scanner decode error: %v", err)
		return
	}

	if len(sr.Data) == 0 || len(sr.Data[0].D) < 6 {
		return
	}

	d := sr.Data[0].D
	closePrice := toFloat(d[0])
	chPct := toFloat(d[1])
	chAbs := toFloat(d[2])
	openPrice := toFloat(d[3])

	if closePrice > 0 {
		a.state.UpdateQuote(&closePrice, &chAbs, &chPct, nil, nil, &openPrice, nil)
		log.Printf("TradingView analysis refreshed: %s spot price $%.3f (%+.2f%%)", config.Symbol, closePrice, chPct)
	}
}

func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

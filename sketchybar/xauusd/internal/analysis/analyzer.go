package analysis

import (
	"bytes"
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

// Analyzer manages periodic OHLC candle & CFD scanner data from REST APIs.
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
	// Immediate initial fetch for spot quote and candlestick history
	a.Refresh(ctx)
	a.RefreshCandles(ctx)

	ticker := time.NewTicker(config.AnalysisRefresh)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.Refresh(ctx)
			a.RefreshCandles(ctx)
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
	}
}

// RefreshCandles fetches 30 historical 5m OHLC candlestick bars from REST endpoints.
func (a *Analyzer) RefreshCandles(ctx context.Context) {
	candles := a.fetchBybitCandles(ctx)
	if len(candles) == 0 {
		candles = a.fetchBinanceCandles(ctx)
	}

	if len(candles) > 0 {
		a.state.SetCandles(candles)
	}
}

func (a *Analyzer) fetchBybitCandles(ctx context.Context) []market.CandleBar {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.KlineRESTURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var res struct {
		RetCode int `json:"retCode"`
		Result  struct {
			List [][]string `json:"list"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &res); err != nil || len(res.Result.List) == 0 {
		return nil
	}

	// Bybit returns newest first, so reverse to chronological order
	raw := res.Result.List
	candles := make([]market.CandleBar, len(raw))
	for i := range raw {
		item := raw[len(raw)-1-i]
		ts, _ := strconv.ParseInt(item[0], 10, 64)
		open, _ := strconv.ParseFloat(item[1], 64)
		high, _ := strconv.ParseFloat(item[2], 64)
		low, _ := strconv.ParseFloat(item[3], 64)
		close, _ := strconv.ParseFloat(item[4], 64)
		vol, _ := strconv.ParseFloat(item[5], 64)

		candles[i] = market.CandleBar{
			Timestamp: ts,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    vol,
		}
	}

	return candles
}

func (a *Analyzer) fetchBinanceCandles(ctx context.Context) []market.CandleBar {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.KlineBackupRESTURL, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var raw [][]interface{}
	if err := json.Unmarshal(body, &raw); err != nil || len(raw) == 0 {
		return nil
	}

	candles := make([]market.CandleBar, len(raw))
	for i, item := range raw {
		if len(item) < 6 {
			continue
		}
		ts := int64(item[0].(float64))
		open, _ := strconv.ParseFloat(item[1].(string), 64)
		high, _ := strconv.ParseFloat(item[2].(string), 64)
		low, _ := strconv.ParseFloat(item[3].(string), 64)
		close, _ := strconv.ParseFloat(item[4].(string), 64)
		vol, _ := strconv.ParseFloat(item[5].(string), 64)

		candles[i] = market.CandleBar{
			Timestamp: ts,
			Open:      open,
			High:      high,
			Low:       low,
			Close:     close,
			Volume:    vol,
		}
	}

	return candles
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

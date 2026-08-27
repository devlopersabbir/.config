package market

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"xauusd/internal/config"
)

// State represents the thread-safe global market state.
type State struct {
	mu sync.RWMutex

	Price         float64
	PrevPrice     float64
	Change        float64
	ChangePct     float64
	Bid           float64
	Ask           float64
	Direction     string // UP / DOWN / FLAT
	CandleOpen    float64
	Candles       []CandleBar
	Connected     bool

	LastTick      time.Time
	LastRender    time.Time
	PendingRender bool
}

// NewState initializes a new State instance with disk-cached candles.
func NewState() *State {
	s := &State{
		Direction: "FLAT",
	}

	// Restore cached candles for 0ms instant startup rendering
	cachedCandles := loadCandlesFromDisk()
	if len(cachedCandles) > 0 {
		s.Candles = cachedCandles
		s.CandleOpen = cachedCandles[len(cachedCandles)-1].Open
		s.Price = cachedCandles[len(cachedCandles)-1].Close
		s.Connected = true
		s.PendingRender = true
	}

	return s
}

// UpdateQuote updates the state with live quote fields received from TradingView.
func (s *State) UpdateQuote(lp, ch, chp, bid, ask, open, prevClose *float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.LastTick = now
	s.Connected = true

	if bid != nil && *bid > 0 {
		s.Bid = *bid
	}
	if ask != nil && *ask > 0 {
		s.Ask = *ask
	}
	if ch != nil {
		s.Change = *ch
	}
	if chp != nil {
		s.ChangePct = *chp
	}
	if open != nil && *open > 0 {
		s.CandleOpen = *open
	}

	if lp != nil && *lp > 0 {
		prev := s.Price
		s.Price = *lp

		if prev > 0 {
			if *lp > prev {
				s.Direction = "UP"
			} else if *lp < prev {
				s.Direction = "DOWN"
			}
		} else if s.Change > 0 {
			s.Direction = "UP"
		} else if s.Change < 0 {
			s.Direction = "DOWN"
		}

		s.PrevPrice = prev

		// Update the active candle's close with the latest trade
		if len(s.Candles) > 0 {
			lastIdx := len(s.Candles) - 1
			s.Candles[lastIdx].Close = *lp
			if *lp > s.Candles[lastIdx].High {
				s.Candles[lastIdx].High = *lp
			}
			if *lp < s.Candles[lastIdx].Low {
				s.Candles[lastIdx].Low = *lp
			}
		}
	}

	s.PendingRender = true
}

// SetCandles updates or replaces historical OHLC analysis candles and saves them to disk.
func (s *State) SetCandles(candles []CandleBar) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(candles) > 0 {
		s.Candles = candles
		if s.Price > 0 {
			s.Candles[len(s.Candles)-1].Close = s.Price
		}
		if s.CandleOpen <= 0 {
			s.CandleOpen = candles[len(candles)-1].Open
		}
		saveCandlesToDisk(s.Candles)
	}
	s.PendingRender = true
}

// UpdateLastCandle updates or appends the active candle bar.
func (s *State) UpdateLastCandle(candle CandleBar) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Candles) == 0 {
		s.Candles = append(s.Candles, candle)
	} else {
		lastIdx := len(s.Candles) - 1
		if s.Candles[lastIdx].Timestamp == candle.Timestamp {
			s.Candles[lastIdx] = candle
		} else if candle.Timestamp > s.Candles[lastIdx].Timestamp {
			s.Candles = append(s.Candles, candle)
		} else {
			s.Candles[lastIdx] = candle
		}
	}

	if candle.Close > 0 && s.Price == 0 {
		s.Price = candle.Close
	}

	saveCandlesToDisk(s.Candles)
	s.PendingRender = true
}

// SetConnected marks the connection status.
func (s *State) SetConnected(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Connected = connected
}

// Snapshot retrieves the current state values needed for rendering.
type Snapshot struct {
	Price        float64
	Change       float64
	ChangePct    float64
	CandleOpen   float64
	Direction    string
	Connected    bool
	Candles      []CandleBar
	ShouldRender bool
}

// PrepareRender checks if a render is pending and returns a snapshot while clearing PendingRender.
func (s *State) PrepareRender() Snapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	shouldRender := s.PendingRender && s.Price > 0
	if shouldRender {
		s.PendingRender = false
		s.LastRender = time.Now()
	}

	candlesCopy := make([]CandleBar, len(s.Candles))
	copy(candlesCopy, s.Candles)

	return Snapshot{
		Price:        s.Price,
		Change:       s.Change,
		ChangePct:    s.ChangePct,
		CandleOpen:   s.CandleOpen,
		Direction:    s.Direction,
		Connected:    s.Connected,
		Candles:      candlesCopy,
		ShouldRender: shouldRender,
	}
}

// GetPrice returns the current price and connection status.
func (s *State) GetPrice() (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Price, s.Connected
}

func loadCandlesFromDisk() []CandleBar {
	data, err := os.ReadFile(config.CandleCacheFile)
	if err != nil {
		return nil
	}
	var candles []CandleBar
	if err := json.Unmarshal(data, &candles); err != nil {
		return nil
	}
	return candles
}

func saveCandlesToDisk(candles []CandleBar) {
	if len(candles) == 0 {
		return
	}
	data, err := json.Marshal(candles)
	if err != nil {
		return
	}
	_ = os.WriteFile(config.CandleCacheFile, data, 0644)
}

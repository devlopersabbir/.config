package market

import (
	"sync"
	"time"
)

// State represents the thread-safe global market state.
type State struct {
	mu sync.RWMutex

	Price      float64
	PrevPrice  float64
	Bid        float64
	Ask        float64
	Direction  string // UP / DOWN / FLAT
	CandleOpen float64
	Candles    []CandleBar
	Connected  bool

	LastTick      time.Time
	LastRender    time.Time
	PendingRender bool
}

// NewState initializes a new State instance.
func NewState() *State {
	return &State{
		Direction: "FLAT",
	}
}

// UpdateBookTick updates the state with the latest bid and ask prices.
func (s *State) UpdateBookTick(bid, ask float64) {
	now := time.Now()
	mid := (bid + ask) / 2.0

	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.Price
	s.Bid = bid
	s.Ask = ask
	s.Price = mid
	s.LastTick = now
	s.Connected = true

	if prev > 0 {
		if mid > prev {
			s.Direction = "UP"
		} else if mid < prev {
			s.Direction = "DOWN"
		}
	} else if s.CandleOpen > 0 {
		if mid >= s.CandleOpen {
			s.Direction = "UP"
		} else {
			s.Direction = "DOWN"
		}
	}

	s.PrevPrice = prev
	if len(s.Candles) > 0 {
		s.Candles[len(s.Candles)-1].Close = mid
		if mid > s.Candles[len(s.Candles)-1].High {
			s.Candles[len(s.Candles)-1].High = mid
		}
		if mid < s.Candles[len(s.Candles)-1].Low {
			s.Candles[len(s.Candles)-1].Low = mid
		}
	}
	s.PendingRender = true
}

// UpdateTradeTick updates the state with the latest executed trade price.
func (s *State) UpdateTradeTick(price float64) {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.Price
	s.Price = price
	s.LastTick = now
	s.Connected = true

	if prev > 0 {
		if price > prev {
			s.Direction = "UP"
		} else if price < prev {
			s.Direction = "DOWN"
		}
	} else if s.CandleOpen > 0 {
		if price >= s.CandleOpen {
			s.Direction = "UP"
		} else {
			s.Direction = "DOWN"
		}
	}

	s.PrevPrice = prev
	if len(s.Candles) > 0 {
		s.Candles[len(s.Candles)-1].Close = price
		if price > s.Candles[len(s.Candles)-1].High {
			s.Candles[len(s.Candles)-1].High = price
		}
		if price < s.Candles[len(s.Candles)-1].Low {
			s.Candles[len(s.Candles)-1].Low = price
		}
	}
	s.PendingRender = true
}

// UpdateCandle updates the current 15m candle open price and close in real-time.
func (s *State) UpdateCandle(open, close float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if open > 0 {
		s.CandleOpen = open
	}
	if len(s.Candles) > 0 {
		if open > 0 {
			s.Candles[len(s.Candles)-1].Open = open
		}
		if close > 0 {
			s.Candles[len(s.Candles)-1].Close = close
		}
	}
	s.PendingRender = true
}

// SetAnalysisCandles sets or updates historical OHLC analysis data.
func (s *State) SetAnalysisCandles(candles []CandleBar) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(candles) >= 18 {
		s.Candles = candles
		s.CandleOpen = candles[len(candles)-1].Open
	} else if len(candles) > 0 {
		for _, c := range candles {
			if len(s.Candles) > 0 {
				s.Candles[len(s.Candles)-1] = c
			} else {
				s.Candles = append(s.Candles, c)
			}
		}
		if len(s.Candles) > 0 {
			s.CandleOpen = s.Candles[len(s.Candles)-1].Open
		}
	}
	s.PendingRender = true
}

// SetConnected marks the WebSocket connection status.
func (s *State) SetConnected(connected bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Connected = connected
}

// Snapshot retrieves the current state values needed for rendering.
type Snapshot struct {
	Price        float64
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

	shouldRender := s.PendingRender && s.Connected && s.Price > 0
	if shouldRender {
		s.PendingRender = false
		s.LastRender = time.Now()
	}

	candlesCopy := make([]CandleBar, len(s.Candles))
	copy(candlesCopy, s.Candles)

	return Snapshot{
		Price:        s.Price,
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

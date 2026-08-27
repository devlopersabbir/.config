package streamer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"xauusd/internal/config"
	"xauusd/internal/market"
)

// Streamer handles the persistent real-time WebSocket connection to Binance.
type Streamer struct {
	state *market.State
}

// NewStreamer creates a new Streamer instance.
func NewStreamer(state *market.State) *Streamer {
	return &Streamer{
		state: state,
	}
}

// Start manages the continuous reconnection loop.
func (s *Streamer) Start(ctx context.Context) {
	backoff := 1 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		log.Println("Connecting to Binance Gold WebSocket stream...")

		err := s.connectAndStream(ctx)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			log.Printf("WebSocket stream error: %v", err)
		}

		s.state.SetConnected(false)
		log.Printf("Reconnecting in %s...", backoff)

		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}

		if backoff < 15*time.Second {
			backoff *= 2
		}
	}
}

func (s *Streamer) connectAndStream(ctx context.Context) error {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 5 * time.Second

	conn, _, err := dialer.DialContext(ctx, config.WsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}
	defer conn.Close()

	log.Println("WebSocket connected to Binance Gold stream.")
	s.state.SetConnected(true)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		var combined market.BinanceCombined
		if err := json.Unmarshal(msg, &combined); err != nil {
			continue
		}

		if strings.Contains(combined.Stream, "bookTicker") {
			var bt market.BinanceBookTicker
			if err := json.Unmarshal(combined.Data, &bt); err == nil {
				bid, _ := strconv.ParseFloat(bt.BidPrice, 64)
				ask, _ := strconv.ParseFloat(bt.AskPrice, 64)
				if bid > 0 && ask > 0 {
					s.state.UpdateBookTick(bid, ask)
				}
			}
		} else if strings.Contains(combined.Stream, "trade") {
			var tr market.BinanceTrade
			if err := json.Unmarshal(combined.Data, &tr); err == nil {
				tPrice, _ := strconv.ParseFloat(tr.Price, 64)
				if tPrice > 0 {
					s.state.UpdateTradeTick(tPrice)
				}
			}
		} else if strings.Contains(combined.Stream, "kline_5m") {
			var ks market.BinanceKlineStream
			if err := json.Unmarshal(combined.Data, &ks); err == nil {
				openP, _ := strconv.ParseFloat(ks.Kline.Open, 64)
				closeP, _ := strconv.ParseFloat(ks.Kline.Close, 64)
				s.state.UpdateCandle(openP, closeP)
			}
		}
	}
}

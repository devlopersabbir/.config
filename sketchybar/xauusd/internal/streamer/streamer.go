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

// Streamer handles the persistent real-time WebSocket connection to Bybit V5.
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

		log.Printf("Connecting to Bybit V5 WebSocket stream (%s)...", config.Symbol)

		err := s.connectAndStream(ctx)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			log.Printf("Bybit WebSocket stream error: %v", err)
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
		return fmt.Errorf("bybit websocket dial failed: %w", err)
	}
	defer conn.Close()

	log.Printf("WebSocket connected to Bybit V5 Linear Stream (%s).", config.Symbol)
	s.state.SetConnected(true)

	// Subscribe to topics
	subMsg := map[string]interface{}{
		"op": "subscribe",
		"args": []string{
			"tickers." + config.Symbol,
			"publicTrade." + config.Symbol,
			"kline.5." + config.Symbol,
			"orderbook.1." + config.Symbol,
		},
	}
	if err := conn.WriteJSON(subMsg); err != nil {
		return fmt.Errorf("subscribe failed: %w", err)
	}

	// Ping ticker to keep Bybit connection alive
	pingTicker := time.NewTicker(20 * time.Second)
	defer pingTicker.Stop()

	// Goroutine for periodic Bybit ping
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-pingTicker.C:
				_ = conn.WriteJSON(map[string]string{"op": "ping"})
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		var wsMsg market.BybitWsMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			continue
		}

		topic := strings.ToLower(wsMsg.Topic)

		if strings.HasPrefix(topic, "tickers") {
			var ticker market.BybitTickerData
			if err := json.Unmarshal(wsMsg.Data, &ticker); err == nil {
				if lp, err := strconv.ParseFloat(ticker.LastPrice, 64); err == nil && lp > 0 {
					s.state.UpdateTradeTick(lp)
				}
				bid, err1 := strconv.ParseFloat(ticker.Bid1Price, 64)
				ask, err2 := strconv.ParseFloat(ticker.Ask1Price, 64)
				if err1 == nil && err2 == nil && bid > 0 && ask > 0 {
					s.state.UpdateBookTick(bid, ask)
				}
			}
		} else if strings.HasPrefix(topic, "publictrade") {
			var trades []market.BybitTradeData
			if err := json.Unmarshal(wsMsg.Data, &trades); err == nil {
				for _, tr := range trades {
					if p, err := strconv.ParseFloat(tr.Price, 64); err == nil && p > 0 {
						s.state.UpdateTradeTick(p)
					}
				}
			}
		} else if strings.HasPrefix(topic, "orderbook") {
			var ob market.BybitOrderbookData
			if err := json.Unmarshal(wsMsg.Data, &ob); err == nil {
				var bid, ask float64
				if len(ob.Bids) > 0 && len(ob.Bids[0]) > 0 {
					bid, _ = strconv.ParseFloat(ob.Bids[0][0], 64)
				}
				if len(ob.Asks) > 0 && len(ob.Asks[0]) > 0 {
					ask, _ = strconv.ParseFloat(ob.Asks[0][0], 64)
				}
				if bid > 0 && ask > 0 {
					s.state.UpdateBookTick(bid, ask)
				}
			}
		} else if strings.HasPrefix(topic, "kline") {
			var klines []market.BybitKlineData
			if err := json.Unmarshal(wsMsg.Data, &klines); err == nil && len(klines) > 0 {
				latest := klines[len(klines)-1]
				openP, _ := strconv.ParseFloat(latest.Open, 64)
				closeP, _ := strconv.ParseFloat(latest.Close, 64)
				if openP > 0 && closeP > 0 {
					s.state.UpdateCandle(openP, closeP)
				}
			}
		}
	}
}

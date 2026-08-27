package streamer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"xauusd/internal/config"
	"xauusd/internal/market"
)

// Streamer handles the persistent real-time WebSocket connection to TradingView.
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

		err := s.connectAndStream(ctx)

		if ctx.Err() != nil {
			return
		}

		if err != nil {
			log.Printf("TradingView stream notice: %v", err)
		}

		s.state.SetConnected(false)

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

func randomSessionID(prefix string) string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

func formatMessage(m string, p []interface{}) string {
	payload, _ := json.Marshal(map[string]interface{}{
		"m": m,
		"p": p,
	})
	return fmt.Sprintf("~m~%d~m~%s", len(payload), string(payload))
}

func (s *Streamer) connectAndStream(ctx context.Context) error {
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 6 * time.Second

	header := http.Header{}
	header.Set("Origin", "https://www.tradingview.com")
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	conn, _, err := dialer.DialContext(ctx, config.WsURL, header)
	if err != nil {
		return fmt.Errorf("tradingview websocket dial failed: %w", err)
	}
	defer conn.Close()

	s.state.SetConnected(true)

	// Read initial connection handshake
	_, _, err = conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("initial handshake read failed: %w", err)
	}

	qs := randomSessionID("qs_")
	cs := randomSessionID("cs_")

	sendMessage := func(m string, p []interface{}) error {
		msg := formatMessage(m, p)
		return conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}

	// Initialize TradingView session & subscribe to symbol and series
	if err := sendMessage("set_auth_token", []interface{}{"unauthorized_user_token"}); err != nil {
		return err
	}
	if err := sendMessage("chart_create_session", []interface{}{cs, ""}); err != nil {
		return err
	}
	if err := sendMessage("quote_create_session", []interface{}{qs}); err != nil {
		return err
	}
	if err := sendMessage("quote_set_fields", []interface{}{
		qs, "lp", "ch", "chp", "open_price", "high_price", "low_price", "prev_close_price", "bid", "ask",
	}); err != nil {
		return err
	}
	if err := sendMessage("quote_add_symbols", []interface{}{qs, config.Symbol}); err != nil {
		return err
	}
	if err := sendMessage("resolve_symbol", []interface{}{
		cs, "sds_sym_1", "={\"symbol\":\"" + config.Symbol + "\",\"adjustment\":\"splits\"}",
	}); err != nil {
		return err
	}
	if err := sendMessage("create_series", []interface{}{
		cs, "sds_1", "s1", "sds_sym_1", "5", 30, "",
	}); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		str := string(raw)

		// Handle heartbeat ping
		if strings.HasPrefix(str, "~m~") && strings.Contains(str, "~h~") {
			_ = conn.WriteMessage(websocket.TextMessage, raw)
			continue
		}

		parts := strings.Split(str, "~m~")
		for _, part := range parts {
			if len(part) == 0 {
				continue
			}
			if strings.HasPrefix(part, "~h~") {
				resp := fmt.Sprintf("~m~%d~m~%s", len(part), part)
				_ = conn.WriteMessage(websocket.TextMessage, []byte(resp))
				continue
			}
			if strings.HasPrefix(part, "{") {
				var pkt struct {
					M string          `json:"m"`
					P json.RawMessage `json:"p"`
				}
				if err := json.Unmarshal([]byte(part), &pkt); err == nil {
					s.handlePacket(pkt.M, pkt.P)
				}
			}
		}
	}
}

func (s *Streamer) handlePacket(method string, payload json.RawMessage) {
	switch method {
	case "qsd":
		var params []json.RawMessage
		if err := json.Unmarshal(payload, &params); err == nil && len(params) >= 2 {
			var qData struct {
				N string              `json:"n"`
				V market.TVQuoteData `json:"v"`
			}
			if err := json.Unmarshal(params[1], &qData); err == nil {
				v := qData.V
				s.state.UpdateQuote(v.Lp, v.Ch, v.Chp, v.Bid, v.Ask, v.Open, v.PrevClose)
			}
		}

	case "timescale_update":
		var params []json.RawMessage
		if err := json.Unmarshal(payload, &params); err == nil && len(params) >= 2 {
			var seriesData map[string]struct {
				S []struct {
					I int       `json:"i"`
					V []float64 `json:"v"`
				} `json:"s"`
			}
			if err := json.Unmarshal(params[1], &seriesData); err == nil {
				if s1, ok := seriesData["sds_1"]; ok && len(s1.S) > 0 {
					var candles []market.CandleBar
					for _, item := range s1.S {
						if len(item.V) >= 5 {
							vol := 0.0
							if len(item.V) > 5 {
								vol = item.V[5]
							}
							candles = append(candles, market.CandleBar{
								Timestamp: int64(item.V[0]),
								Open:      item.V[1],
								High:      item.V[2],
								Low:       item.V[3],
								Close:     item.V[4],
								Volume:    vol,
							})
						}
					}
					if len(candles) > 0 {
						s.state.SetCandles(candles)
					}
				}
			}
		}

	case "du":
		var params []json.RawMessage
		if err := json.Unmarshal(payload, &params); err == nil && len(params) >= 2 {
			var seriesData map[string]struct {
				S []struct {
					I int       `json:"i"`
					V []float64 `json:"v"`
				} `json:"s"`
			}
			if err := json.Unmarshal(params[1], &seriesData); err == nil {
				if s1, ok := seriesData["sds_1"]; ok && len(s1.S) > 0 {
					item := s1.S[len(s1.S)-1]
					if len(item.V) >= 5 {
						vol := 0.0
						if len(item.V) > 5 {
							vol = item.V[5]
						}
						s.state.UpdateLastCandle(market.CandleBar{
							Timestamp: int64(item.V[0]),
							Open:      item.V[1],
							High:      item.V[2],
							Low:       item.V[3],
							Close:     item.V[4],
							Volume:    vol,
						})
					}
				}
			}
		}
	}
}


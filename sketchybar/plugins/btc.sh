#!/usr/bin/env bash

# Fetch Live BTC Price and 15-minute percentage price change
CACHE_FILE="/tmp/sketchybar_btc_price"

fetch_btc() {
  local raw
  raw=$(curl -s --connect-timeout 1.5 --max-time 2 "https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=15m&limit=1" 2>/dev/null)

  if [ -n "$raw" ] && [ "$raw" != "[]" ]; then
    local open_price curr_price
    open_price=$(echo "$raw" | jq -r '.[0][1] // empty' 2>/dev/null)
    curr_price=$(echo "$raw" | jq -r '.[0][4] // empty' 2>/dev/null)

    if [ -n "$open_price" ] && [ -n "$curr_price" ]; then
      local arrow color formatted sign pct_val
      read -r arrow color formatted sign pct_val <<< $(awk -v open="$open_price" -v curr="$curr_price" 'BEGIN {
        diff = curr - open
        pct = (diff / open) * 100
        if (diff >= 0) {
          arr = "▲"
          col = "0xffa6e3a1" # Green
          s = "+"
        } else {
          arr = "▼"
          col = "0xfff38ba8" # Red
          s = ""
        }
        printf "%s %s $%\x27.0f %s %.2f%%\n", arr, col, curr, s, pct
      }')

      local full_label="$formatted ${sign}${pct_val} $arrow"
      echo "$full_label|$color" > "$CACHE_FILE"
      echo "$full_label|$color"
      return
    fi
  fi

  # Fallback to cache if request timed out
  if [ -f "$CACHE_FILE" ]; then
    cat "$CACHE_FILE"
  else
    echo "--- |0xffcdd6f4"
  fi
}

DATA=$(fetch_btc)
LABEL=$(echo "$DATA" | cut -d'|' -f1)
COLOR=$(echo "$DATA" | cut -d'|' -f2)

if [ -n "$NAME" ]; then
  sketchybar --set "$NAME" label="$LABEL" label.color="$COLOR"
else
  echo "BTC: $LABEL (Color: $COLOR)"
fi

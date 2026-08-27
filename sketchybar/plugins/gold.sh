#!/usr/bin/env bash

CACHE_FILE="/tmp/sketchybar_gold_price"

if [ -f "$CACHE_FILE" ]; then

    DATA=$(cat "$CACHE_FILE")

    LABEL="${DATA%%|*}"
    COLOR="${DATA#*|}"

    sketchybar --set "$NAME" \
        label="$LABEL" \
        label.color="$COLOR"

else

    sketchybar --set "$NAME" \
        label="🥇 Connecting..." \
        label.color="0xffcdd6f4"

fi

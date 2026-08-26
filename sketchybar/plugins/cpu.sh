#!/usr/bin/env bash

# Calculate instantaneous CPU usage percentage
CPU_USAGE=$(top -l 1 -n 0 | awk '/CPU usage/ { sub("%","",$3); sub("%","",$5); printf "%.0f", $3 + $5 }')
CPU_NORM=$(awk -v u="$CPU_USAGE" 'BEGIN { printf "%.2f", u / 100 }')

sketchybar --push "$NAME" "$CPU_NORM" \
           --set "$NAME" label="${CPU_USAGE}%"

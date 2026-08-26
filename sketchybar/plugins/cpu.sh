#!/usr/bin/env bash

# Calculate instantaneous CPU usage percentage
CPU_USAGE=$(top -l 1 -n 0 | awk '/CPU usage/ { sub("%","",$3); sub("%","",$5); printf "%.0f%%", $3 + $5 }')

if [ -z "$CPU_USAGE" ]; then
  CPU_USAGE="0%"
fi

if [ -n "$NAME" ]; then
  sketchybar --set "$NAME" label="$CPU_USAGE"
else
  echo "$CPU_USAGE"
fi

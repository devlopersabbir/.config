#!/usr/bin/env bash

# Calculate primary Disk usage percentage
DISK_USAGE=$(df -H / | awk 'NR==2 { print $5 }')

if [ -z "$DISK_USAGE" ]; then
  DISK_USAGE="0%"
fi

if [ -n "$NAME" ]; then
  sketchybar --set "$NAME" label="$DISK_USAGE"
else
  echo "$DISK_USAGE"
fi

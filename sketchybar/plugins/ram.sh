#!/usr/bin/env bash

# Calculate physical RAM usage percentage
RAM_USAGE=$(top -l 1 -n 0 | awk '/PhysMem/ {
  gsub(/[^0-9\.]/, "", $2)
  u = $2
  if (index($0, "M used") > 0) u = u / 1024

  match($0, /[0-9\.]+M unused|[0-9\.]+G unused/)
  free_str = substr($0, RSTART, RLENGTH)
  gsub(/[^0-9\.]/, "", free_str)
  f = free_str
  if (index($0, "M unused") > 0) f = f / 1024

  total = u + f
  if (total > 0) {
    printf "%.0f%%", (u / total) * 100
  } else {
    printf "0%%"
  }
}')

if [ -z "$RAM_USAGE" ]; then
  RAM_USAGE="0%"
fi

if [ -n "$NAME" ]; then
  sketchybar --set "$NAME" label="$RAM_USAGE"
else
  echo "$RAM_USAGE"
fi

#!/usr/bin/env bash

# Launcher for SketchyBar Live Stream Daemon
PIDFILE="/tmp/sketchybar_live_stream.pid"
SCRIPT_PATH="$HOME/.config/sketchybar/plugins/live_stream.py"

# Stop existing instance if running
if [ -f "$PIDFILE" ]; then
  old_pid=$(cat "$PIDFILE" 2>/dev/null)
  if [ -n "$old_pid" ] && kill -0 "$old_pid" 2>/dev/null; then
    kill "$old_pid" 2>/dev/null
    sleep 0.2
  fi
  rm -f "$PIDFILE"
fi

# Start live daemon in background
nohup python3 "$SCRIPT_PATH" >/dev/null 2>&1 &
echo $! > "$PIDFILE"

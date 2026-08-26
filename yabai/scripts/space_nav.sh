#!/usr/bin/env bash

# Ensure 10 spaces exist in macOS / yabai
TARGET_SPACES=10

ensure_spaces() {
  local current_spaces
  current_spaces=$(yabai -m query --spaces 2>/dev/null | jq 'length' 2>/dev/null)
  
  if [ -z "$current_spaces" ] || [ "$current_spaces" -eq 0 ]; then
    current_spaces=1
  fi

  if [ "$current_spaces" -lt "$TARGET_SPACES" ]; then
    local needed=$((TARGET_SPACES - current_spaces))
    
    # Try yabai native create first
    local yabai_created=0
    for ((i=1; i<=needed; i++)); do
      if yabai -m space --create 2>/dev/null; then
        yabai_created=$((yabai_created + 1))
      else
        break
      fi
    done

    # If yabai native failed (e.g., SIP enabled without scripting-addition), use AppleScript Mission Control UI
    if [ "$yabai_created" -lt "$needed" ]; then
      local remaining=$((needed - yabai_created))
      osascript -e "
      tell application \"Mission Control\" to launch
      delay 0.5
      tell application \"System Events\"
          tell process \"Dock\"
              repeat $remaining times
                  try
                      click (button 1 of group 2 of group 1 of group 1)
                      delay 0.12
                  end try
              end repeat
          end tell
          delay 0.3
          key code 53
      end tell
      " >/dev/null 2>&1
    fi

    # Reload Sketchybar so it immediately reflects all 10 spaces
    if command -v sketchybar >/dev/null 2>&1; then
      sketchybar --reload >/dev/null 2>&1 || true
    fi
  fi
}

focus_space() {
  local target="$1"
  if [ -n "$target" ]; then
    yabai -m space --focus "$target" 2>/dev/null || true
  fi
}

move_window_to_space() {
  local target="$1"
  if [ -n "$target" ]; then
    yabai -m window --space "$target" 2>/dev/null && yabai -m space --focus "$target" 2>/dev/null || true
  fi
}

case "$1" in
  ensure)
    ensure_spaces
    ;;
  focus)
    focus_space "$2"
    ;;
  move)
    move_window_to_space "$2"
    ;;
  *)
    ensure_spaces
    ;;
esac

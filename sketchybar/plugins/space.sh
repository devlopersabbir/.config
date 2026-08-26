#!/usr/bin/env bash

# When sketchybar triggers space script:
# $SELECTED holds "true" if the space is currently focused, "false" otherwise.
# $NAME holds the space item name (e.g. space.1).

if [ "$SELECTED" = "true" ]; then
  sketchybar --animate sin 10 --set "$NAME" \
    icon.color=0xff11111b \
    background.color=0xff7aa2f7 \
    background.border_color=0xff7aa2f7
else
  sketchybar --animate sin 10 --set "$NAME" \
    icon.color=0xff737994 \
    background.color=0x25313244 \
    background.border_color=0x00000000
fi

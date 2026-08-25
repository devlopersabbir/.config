#!/usr/bin/env zsh
# ~/.config/zsh/banner.zsh — Terminal welcome panel (compact)

white=$'\033[1;97m'
dim=$'\033[38;5;240m'
faint=$'\033[38;5;236m'
accent=$'\033[38;5;80m'
green=$'\033[38;5;114m'
yellow=$'\033[38;5;180m'
red=$'\033[38;5;203m'
pink=$'\033[38;5;210m'
r=$'\033[0m'

clear

# --- time-aware greeting ---
h=$(date +%H)
if   (( h < 5 ));  then greet="still up, night owl"
elif (( h < 12 )); then greet="good morning"
elif (( h < 17 )); then greet="good afternoon"
elif (( h < 21 )); then greet="good evening"
else                    greet="wind down soon"
fi

# --- stats ---
uptime_str=$(uptime | sed -E 's/.*up ([^,]*),.*/\1/' | xargs)
disk_pct=$(df -h / | awk 'NR==2 {gsub("%","",$5); print $5}')

if [[ "$(uname -s)" == "Darwin" ]]; then
  mem_pct=$(vm_stat | awk '
    /Pages active/ {active=$3} /Pages wired/ {wired=$4} /Pages free/ {free=$3}
    END { gsub("\\.","",active); gsub("\\.","",wired); gsub("\\.","",free)
          total=active+wired+free; printf "%.0f", (active+wired)/total*100 }')
else
  mem_pct=$(free | awk '/Mem:/ {printf "%.0f", (($2-$7)/$2)*100}')
fi

git_line=""
if git rev-parse --is-inside-work-tree &> /dev/null; then
  branch=$(git branch --show-current 2>/dev/null)
  changes=$(git status --porcelain 2>/dev/null | wc -l | xargs)
  git_line="  ${pink}⎇${r}  $branch ${dim}·${r} ${changes} changed"
fi

# --- inline bar ---
bar() {
  local pct=$1 width=10
  local filled=$(( pct * width / 100 ))
  local empty=$(( width - filled ))
  local color=$green
  (( pct >= 60 )) && color=$yellow
  (( pct >= 85 )) && color=$red
  local out="${color}"
  for ((i=0; i<filled; i++)); do out+="━"; done
  out+="${faint}"
  for ((i=0; i<empty; i++)); do out+="━"; done
  out+="${r}"
  echo "$out"
}

spine="${accent}▌${r}"

echo ""
echo "  ${spine}  ${white}Sabbir Hossain${r}  ${dim}·${r}  ${dim}${greet}${r}"
echo "  ${spine}  ${accent}◷${r} $(date '+%I:%M %p') ${dim}$(date '+%a, %d %b')${r}  ${accent}↑${r} up ${uptime_str}"
echo "  ${spine}  ${dim}disk${r}   $(bar $disk_pct) ${dim}${disk_pct}%${r}   ${dim}mem${r}  $(bar $mem_pct) ${dim}${mem_pct}%${r}"
[[ -n "$git_line" ]] && echo "  ${spine}${git_line}"
echo "  ${spine}  ${accent}✦${r} ${white}ready to build${r}"
echo ""

#!/usr/bin/env bash

# ==============================================================================
#  macOS Dotfiles & Development Environment Automated Uninstaller
#  Author: Sabbir Hossain Shuvo (@devlopersabbir)
# ==============================================================================

# --- Styling & Colors ---
BOLD='\033[1m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

info()    { echo -e "${BLUE}${BOLD}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}${BOLD}[SUCCESS]${NC} $1"; }
warn()    { echo -e "${YELLOW}${BOLD}[WARN]${NC} $1"; }
error()   { echo -e "${RED}${BOLD}[ERROR]${NC} $1"; }
step()    { echo -e "\n${CYAN}${BOLD}==>${NC} ${BOLD}$1${NC}"; }

# Helper function for y/N confirmation (Default is NO)
confirm() {
    local prompt="$1"
    local response
    read -r -p "$(echo -e "${YELLOW}${BOLD}?${NC} ${prompt} [y/N]: ")" response
    case "${response}" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

echo -e "${RED}${BOLD}"
cat << "EOF"
 ██╗   ██╗███╗   ██╗██╗███╗   ██╗███████╗████████╗ █████╗ ██╗     ██╗     
 ██║   ██║████╗  ██║██║████╗  ██║██╔════╝╚══██╔══╝██╔══██╗██║     ██║     
 ██║   ██║██╔██╗ ██║██║██╔██╗ ██║███████╗   ██║   ███████║██║     ██║     
 ██║   ██║██║╚██╗██║██║██║╚██╗██║╚════██║   ██║   ██╔══██║██║     ██║     
 ╚██████╔╝██║ ╚████║██║██║ ╚████║███████║   ██║   ██║  ██║███████╗███████╗
  ╚═════╝ ╚═╝  ╚═══╝╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚══════╝
  Interactive Tool & Environment Uninstaller
EOF
echo -e "${NC}"

echo -e "${BOLD}This script will guide you through removing tools, packages, and configurations.${NC}"
echo -e "${BOLD}Every step requires your explicit confirmation (y/N). Default is No.${NC}\n"

# Ensure Homebrew is in PATH
if [[ -x /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [[ -x /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
fi

# 1. Yabai & Skhd
step "Window Management (Yabai & Skhd)"
if confirm "Stop and uninstall Yabai (Tiling Window Manager)?"; then
    info "Stopping Yabai service..."
    yabai --stop-service 2>/dev/null || true
    brew services stop yabai 2>/dev/null || true
    if command -v brew &>/dev/null && brew list --formula | grep -q "^yabai\$"; then
        info "Uninstalling yabai via Homebrew..."
        brew uninstall yabai
    fi
    success "Yabai uninstalled"
else
    info "Skipping Yabai"
fi

if confirm "Stop and uninstall Skhd (Hotkey Daemon)?"; then
    info "Stopping Skhd service..."
    skhd --stop-service 2>/dev/null || true
    brew services stop skhd 2>/dev/null || true
    if command -v brew &>/dev/null && brew list --formula | grep -q "^skhd\$"; then
        info "Uninstalling skhd via Homebrew..."
        brew uninstall skhd
    fi
    success "Skhd uninstalled"
else
    info "Skipping Skhd"
fi

# 2. SketchyBar & Gold Market Daemon
step "Status Bar (SketchyBar & XAUUSD Market Daemon)"
if confirm "Stop and uninstall SketchyBar and background market daemons?"; then
    info "Stopping SketchyBar service and background daemons..."
    brew services stop sketchybar 2>/dev/null || true
    pkill -f sketchybar 2>/dev/null || true
    pkill -f xauusd 2>/dev/null || true

    # Clean compiled binary and caches
    rm -f "${HOME}/.config/sketchybar/xauusd/xauusd"
    rm -f /tmp/sketchybar_* /tmp/xauusd*

    if command -v brew &>/dev/null && brew list --formula | grep -q "^sketchybar\$"; then
        info "Uninstalling sketchybar via Homebrew..."
        brew uninstall sketchybar
    fi
    success "SketchyBar uninstalled and caches cleared"
else
    info "Skipping SketchyBar"
fi

# 3. Neovim & Packer Plugins
step "Editor (Neovim & Packer.nvim)"
if confirm "Remove Neovim plugins and Packer packages (~/.local/share/nvim)?"; then
    info "Removing Neovim packer plugins and cache..."
    rm -rf "${HOME}/.local/share/nvim"
    rm -rf "${HOME}/.cache/nvim"
    rm -rf "${HOME}/.local/state/nvim"
    success "Neovim plugin directories removed"
fi

if confirm "Uninstall Neovim (brew package)?"; then
    if command -v brew &>/dev/null && brew list --formula | grep -q "^neovim\$"; then
        info "Uninstalling neovim via Homebrew..."
        brew uninstall neovim
        success "Neovim package uninstalled"
    fi
else
    info "Skipping Neovim uninstallation"
fi

# 4. Tmux
step "Terminal Multiplexer (Tmux)"
if confirm "Uninstall Tmux and remove ~/.tmux.conf link?"; then
    rm -f "${HOME}/.tmux.conf"
    if command -v brew &>/dev/null && brew list --formula | grep -q "^tmux\$"; then
        info "Uninstalling tmux via Homebrew..."
        brew uninstall tmux
    fi
    success "Tmux uninstalled and ~/.tmux.conf removed"
else
    info "Skipping Tmux"
fi

# 5. Oh My Zsh & Custom Plugins
step "Shell (Oh My Zsh & Custom Plugins)"
if confirm "Remove custom Zsh plugins (zsh-autosuggestions, zsh-syntax-highlighting)?"; then
    ZSH_CUSTOM="${ZSH_CUSTOM:-${HOME}/.oh-my-zsh/custom}"
    rm -rf "${ZSH_CUSTOM}/plugins/zsh-autosuggestions"
    rm -rf "${ZSH_CUSTOM}/plugins/zsh-syntax-highlighting"
    success "Zsh custom plugins removed"
fi

if confirm "Completely uninstall Oh My Zsh framework (~/.oh-my-zsh)?"; then
    if [[ -d "${HOME}/.oh-my-zsh" ]]; then
        if [[ -f "${HOME}/.oh-my-zsh/tools/uninstall.sh" ]]; then
            info "Running official Oh My Zsh uninstaller..."
            env ZSH="${HOME}/.oh-my-zsh" sh "${HOME}/.oh-my-zsh/tools/uninstall.sh" --unattended || rm -rf "${HOME}/.oh-my-zsh"
        else
            rm -rf "${HOME}/.oh-my-zsh"
        fi
        success "Oh My Zsh removed"
    fi
else
    info "Skipping Oh My Zsh"
fi

# 6. Flameshot (Screenshot Utility)
step "Screenshots (Flameshot)"
if confirm "Uninstall Flameshot screenshot application?"; then
    if command -v brew &>/dev/null && brew list --cask 2>/dev/null | grep -q "^flameshot\$"; then
        info "Uninstalling flameshot cask..."
        brew uninstall --cask flameshot
        success "Flameshot uninstalled"
    fi
else
    info "Skipping Flameshot"
fi

# 7. Nerd Fonts
step "Nerd Fonts"
if confirm "Uninstall Nerd Fonts (JetBrainsMono & Hack Nerd Fonts)?"; then
    if command -v brew &>/dev/null; then
        brew uninstall --cask font-jetbrains-mono-nerd-font 2>/dev/null || true
        brew uninstall --cask font-hack-nerd-font 2>/dev/null || true
        success "Nerd Fonts uninstalled"
    fi
else
    info "Skipping Nerd Fonts"
fi

# 8. Extra CLI Utilities
step "Additional Developer CLI Tools"
EXTRA_PKGS=(gh ripgrep tree-sitter htop neofetch)

for pkg in "${EXTRA_PKGS[@]}"; do
    if command -v brew &>/dev/null && brew list --formula | grep -q "^${pkg}\$"; then
        if confirm "Uninstall CLI utility '${pkg}'?"; then
            brew uninstall "${pkg}"
            success "Uninstalled ${pkg}"
        fi
    fi
done

# 9. Homebrew Taps (Optional Cleanup)
step "Homebrew Third-Party Taps"
if confirm "Remove custom Homebrew taps (koekeishiya/formulae, felixkratz/formulae)?"; then
    brew untap koekeishiya/formulae 2>/dev/null || true
    brew untap felixkratz/formulae 2>/dev/null || true
    success "Homebrew taps removed"
fi

# 10. Summary
step "Uninstallation Process Finished"
echo -e "${GREEN}${BOLD}✓ Selected uninstallation tasks have been completed.${NC}\n"
echo -e "${CYAN}To re-install or reconfigure everything in the future, run:${NC}"
echo -e "  ${BOLD}bash ~/.config/scripts/setup.sh${NC}\n"

#!/usr/bin/env bash
set -e

# ==============================================================================
#  macOS Dotfiles & Development Environment Automated Setup
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

CONFIG_DIR="${HOME}/.config"

echo -e "${CYAN}${BOLD}"
cat << "EOF"
  ██████╗  ██████╗ ███╗   ██╗███████╗██╗ ██████╗ 
 ██╔════╝ ██╔═══██╗████╗  ██║██╔════╝██║██╔════╝ 
 ██║  ███╗██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗
 ██║   ██║██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║
 ╚██████╔╝╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝
  ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝ 
  Complete macOS Dotfiles & Environment Setup
EOF
echo -e "${NC}"

# 1. macOS Environment Check
step "Checking System Compatibility"
if [[ "$(uname)" != "Darwin" ]]; then
    error "This setup script is optimized for macOS. Aborting."
    exit 1
fi
success "Running on macOS ($(uname -m))"

# 2. Xcode Command Line Tools
step "Checking Xcode Command Line Tools"
if ! xcode-select -p &>/dev/null; then
    info "Installing Xcode Command Line Tools..."
    xcode-select --install
    echo -e "${YELLOW}Please finish the Xcode tools prompt and press Enter to continue...${NC}"
    read -r
else
    success "Xcode Command Line Tools already installed"
fi

# 3. Homebrew Installation & Taps
step "Setting up Homebrew & Repositories"
if ! command -v brew &>/dev/null; then
    info "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    
    # Configure PATH for Apple Silicon / Intel
    if [[ -f /opt/homebrew/bin/brew ]]; then
        eval "$(/opt/homebrew/bin/brew shellenv)"
    elif [[ -f /usr/local/bin/brew ]]; then
        eval "$(/usr/local/bin/brew shellenv)"
    fi
else
    success "Homebrew is installed"
fi

# Ensure brew is available in current shell session
if [[ -x /opt/homebrew/bin/brew ]]; then
    eval "$(/opt/homebrew/bin/brew shellenv)"
elif [[ -x /usr/local/bin/brew ]]; then
    eval "$(/usr/local/bin/brew shellenv)"
fi

info "Updating Homebrew formulae..."
brew update --quiet || true

info "Adding third-party Homebrew taps..."
brew tap koekeishiya/formulae --quiet 2>/dev/null || true
brew tap felixkratz/formulae --quiet 2>/dev/null || true

# 4. Core CLI & Window Management Formulae
step "Installing Packages & Window Management Tools"

PACKAGES=(
    git
    gh
    neovim
    tmux
    zsh
    go
    node
    ripgrep
    tree-sitter
    htop
    neofetch
    yabai
    skhd
    sketchybar
)

for pkg in "${PACKAGES[@]}"; do
    if brew list --formula | grep -q "^${pkg}\$"; then
        success "Formula ${pkg} is already installed"
    else
        info "Installing ${pkg}..."
        brew install "${pkg}" --quiet || warn "Failed to install ${pkg}"
    fi
done

# 5. Fonts & GUI Applications
step "Installing Nerd Fonts & Applications"

CASKS=(
    font-jetbrains-mono-nerd-font
    font-hack-nerd-font
    flameshot
)

for cask in "${CASKS[@]}"; do
    if brew list --cask 2>/dev/null | grep -q "^${cask}\$"; then
        success "Cask ${cask} is already installed"
    else
        info "Installing cask ${cask}..."
        brew install --cask "${cask}" --quiet || warn "Failed to install cask ${cask}"
    fi
done

# 6. Oh My Zsh Setup
step "Configuring Oh My Zsh"
if [[ ! -d "${HOME}/.oh-my-zsh" ]]; then
    info "Installing Oh My Zsh..."
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended
else
    success "Oh My Zsh already installed"
fi

# Install Oh My Zsh Plugins
ZSH_CUSTOM="${ZSH_CUSTOM:-${HOME}/.oh-my-zsh/custom}"

if [[ ! -d "${ZSH_CUSTOM}/plugins/zsh-autosuggestions" ]]; then
    info "Installing zsh-autosuggestions..."
    git clone https://github.com/zsh-users/zsh-autosuggestions "${ZSH_CUSTOM}/plugins/zsh-autosuggestions" --quiet
fi

if [[ ! -d "${ZSH_CUSTOM}/plugins/zsh-syntax-highlighting" ]]; then
    info "Installing zsh-syntax-highlighting..."
    git clone https://github.com/zsh-users/zsh-syntax-highlighting.git "${ZSH_CUSTOM}/plugins/zsh-syntax-highlighting" --quiet
fi

# Ensure alias vim="nvim" exists in ~/.zshrc
if [[ -f "${HOME}/.zshrc" ]]; then
    if ! grep -q 'alias vim="nvim"' "${HOME}/.zshrc"; then
        echo 'alias vim="nvim"' >> "${HOME}/.zshrc"
    fi
fi
success "Zsh environment configured"

# 7. Neovim & Packer.nvim Setup
step "Configuring Neovim Plugins & Packer"
PACKER_DIR="${HOME}/.local/share/nvim/site/pack/packer/start/packer.nvim"
if [[ ! -d "${PACKER_DIR}" ]]; then
    info "Cloning packer.nvim..."
    git clone --depth 1 https://github.com/wbthomason/packer.nvim "${PACKER_DIR}" --quiet
else
    success "packer.nvim already installed"
fi

info "Synchronizing and compiling Neovim plugins (headless)..."
nvim --headless -c 'autocmd User PackerComplete quitall' -c 'PackerSync' 2>/dev/null || true
success "Neovim plugins synchronized"

# 8. Tmux Configuration
step "Configuring Tmux"
if [[ -f "${CONFIG_DIR}/tmux/.tmux.conf" ]]; then
    ln -sf "${CONFIG_DIR}/tmux/.tmux.conf" "${HOME}/.tmux.conf"
    success "Linked ~/.tmux.conf -> ~/.config/tmux/.tmux.conf"
fi

# 9. Compile SketchyBar Gold Market Daemon
step "Compiling SketchyBar Go Market Daemon (XAUUSD)"
if [[ -d "${CONFIG_DIR}/sketchybar/xauusd" ]]; then
    info "Building xauusd binary with Go..."
    (cd "${CONFIG_DIR}/sketchybar/xauusd" && go build -o xauusd main.go)
    success "xauusd daemon compiled successfully"
fi

# 10. Start & Restart Services
step "Starting Window Management Services"

info "Restarting Yabai..."
brew services restart yabai 2>/dev/null || yabai --restart-service 2>/dev/null || true

info "Restarting Skhd..."
brew services restart skhd 2>/dev/null || skhd --restart-service 2>/dev/null || true

info "Restarting Sketchybar..."
brew services restart sketchybar 2>/dev/null || sketchybar --reload 2>/dev/null || true

# 11. Final Instructions
step "Setup Complete!"
echo -e "${GREEN}${BOLD}✓ All dependencies, plugins, configs, and services have been installed and configured!${NC}\n"

echo -e "${YELLOW}${BOLD}IMPORTANT: macOS Accessibility Permissions Required${NC}"
echo -e "To allow Yabai, Skhd, and Flameshot to manage windows and hotkeys:"
echo -e "1. Open ${BOLD}System Settings -> Privacy & Security -> Accessibility${NC}"
echo -e "2. Enable ${BOLD}yabai${NC}, ${BOLD}skhd${NC}, and ${BOLD}Flameshot${NC}"
echo -e "3. Under ${BOLD}Privacy & Security -> Screen Recording${NC}, enable ${BOLD}Flameshot${NC}\n"

echo -e "${CYAN}Reload anytime using:${NC}"
echo -e "  • Yabai config:     ${BOLD}alt - r${NC}  or  ${BOLD}yabai --reload-config${NC}"
echo -e "  • Skhd config:      ${BOLD}skhd --reload${NC}"
echo -e "  • Sketchybar:       ${BOLD}sketchybar --reload${NC}"
echo -e "  • Shell:            ${BOLD}source ~/.zshrc${NC}\n"

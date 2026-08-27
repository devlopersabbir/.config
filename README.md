# ⚡ macOS Dotfiles & Development Environment

> Automated, high-performance tiling window management, terminal workflow, real-time status bar, and developer tools for macOS.

---

## 🚀 Quick Setup (Single Command)

To set up and configure everything on a fresh machine (including Homebrew packages, Yabai, Skhd, SketchyBar, Neovim plugins via Packer, Tmux, Oh My Zsh, and fonts):

### If already cloned to `~/.config`:
```bash
bash ~/.config/scripts/setup.sh
```

### From a fresh terminal (Clone & Install):
```bash
git clone https://github.com/devlopersabbir/.config.git ~/.config && bash ~/.config/scripts/setup.sh
```

---

## 🛠 What Gets Installed & Configured

| Component | Description |
| :--- | :--- |
| **[Yabai](https://github.com/koekeishiya/yabai)** | Tiling window manager configured for 10 spaces, BSP layouts, smooth floating/tiling toggles, and multi-monitor management. |
| **[Skhd](https://github.com/koekeishiya/skhd)** | Simple hotkey daemon for rapid app launches, space switching, focus navigation, and smooth window resizing. |
| **[SketchyBar](https://github.com/FelixKratz/SketchyBar)** | Modern dark glassmorphism status bar featuring active space indicators, system resource monitors (CPU, RAM, Disk), and a live real-time XAU/USD (Gold) market streamer with dynamic 18-bar candlestick chart and momentum trend analysis. |
| **[Neovim](https://neovim.io/)** | Lua-based setup powered by `packer.nvim`, `lsp-zero`, Mason, Treesitter, Telescope, Supermaven AI, and Antigravity IDE integration. |
| **[Tmux](https://github.com/tmux/tmux)** | Terminal multiplexer with custom statusline, vi-style pane navigation (`h/j/k/l`), and resizing. |
| **[Zsh & Oh My Zsh](https://ohmyz.sh/)** | Zsh shell with `zsh-autosuggestions`, `zsh-syntax-highlighting`, and developer aliases. |
| **[Flameshot](https://flameshot.org/)** | Native screenshot capture mapped to <kbd>Cmd</kbd> + <kbd>Shift</kbd> + <kbd>X</kbd>. |
| **Nerd Fonts** | JetBrainsMono Nerd Font & Hack Nerd Font for glyphs, status icons, and UI rendering. |

---

## ⌨️ Keybinding Cheatsheet

### 🪟 Window Resizing
| Shortcut | Action |
| :--- | :--- |
| <kbd>Ctrl</kbd> + <kbd>L</kbd> | **Expand Width** (smooth fluid widening from left & right) |
| <kbd>Ctrl</kbd> + <kbd>H</kbd> | **Shrink Width** (smooth fluid narrowing from left & right) |
| <kbd>Ctrl</kbd> + <kbd>K</kbd> | **Increase Height** (expand top & bottom) |
| <kbd>Ctrl</kbd> + <kbd>J</kbd> | **Decrease Height** (shrink top & bottom) |
| <kbd>Alt</kbd> + <kbd>Shift</kbd> + <kbd>H/J/K/L</kbd> | Alternative resize controls |

### 🎯 Window Focus & Movement
| Shortcut | Action |
| :--- | :--- |
| <kbd>Cmd</kbd> + <kbd>H</kbd> / <kbd>J</kbd> / <kbd>K</kbd> / <kbd>L</kbd> | Move window focus (West / South / North / East) |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>H</kbd> / <kbd>J</kbd> / <kbd>K</kbd> / <kbd>L</kbd> | Swap window position with neighbor |
| <kbd>Ctrl</kbd> + <kbd>Cmd</kbd> + <kbd>H</kbd> / <kbd>J</kbd> / <kbd>K</kbd> / <kbd>L</kbd> | Warp / split window into target direction |
| <kbd>Ctrl</kbd> + <kbd>Alt</kbd> + <kbd>H</kbd> / <kbd>L</kbd> | Focus external display (West / East) |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>S</kbd> / <kbd>G</kbd> | Move window to display (West / East) |

### 🧭 Spaces & Navigation
| Shortcut | Action |
| :--- | :--- |
| <kbd>Alt</kbd> + <kbd>1</kbd> .. <kbd>0</kbd> | Focus workspace / space 1 through 10 |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>1</kbd> .. <kbd>0</kbd> | Send active window to space 1..10 and focus it |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>U</kbd> / <kbd>I</kbd> | Move active window to previous / next space |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>E</kbd> | Balance all windows to occupy equal space |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>R</kbd> | Rotate layout clockwise (90°) |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>Y</kbd> | Mirror layout along Y-axis |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>M</kbd> | Toggle window fullscreen zoom (within space) |
| <kbd>Shift</kbd> + <kbd>Cmd</kbd> + <kbd>T</kbd> | Toggle window floating mode |

### 🚀 Application Shortcuts
| Shortcut | App |
| :--- | :--- |
| <kbd>Alt</kbd> + <kbd>G</kbd> | **Ghostty Terminal** |
| <kbd>Alt</kbd> + <kbd>S</kbd> | **Safari** |
| <kbd>Alt</kbd> + <kbd>T</kbd> | **Telegram** |
| <kbd>Alt</kbd> + <kbd>O</kbd> | **Obsidian** |
| <kbd>Alt</kbd> + <kbd>F</kbd> | **Final Cut Pro** |
| <kbd>Alt</kbd> + <kbd>Q</kbd> | **QuickTime Player** |
| <kbd>Cmd</kbd> + <kbd>Shift</kbd> + <kbd>X</kbd> | **Flameshot** (Interactive Screenshot GUI) |

---

## ⚙️ Configuration & Service Reloads

| Command / Hotkey | Purpose |
| :--- | :--- |
| <kbd>Alt</kbd> + <kbd>R</kbd> or `yabai --reload-config` | Reload Yabai configuration |
| `skhd --reload` | Reload Skhd keybindings |
| `sketchybar --reload` | Reload SketchyBar & restart Gold daemon |
| `source ~/.zshrc` | Reload Zsh shell environment |

---

## 🔒 Required macOS Permissions

After running the setup script, ensure macOS permissions are enabled:

1. **Accessibility**:
   - Open **System Settings -> Privacy & Security -> Accessibility**
   - Enable **yabai**, **skhd**, and **Flameshot**.
2. **Screen Recording**:
   - Open **System Settings -> Privacy & Security -> Screen Recording**
   - Enable **Flameshot** for screenshot captures.

---

## 📁 Repository Structure

```text
.config/
├── README.md              # Documentation and keybinding reference
├── scripts/
│   └── setup.sh           # One-command automated installation script
├── yabai/
│   ├── yabairc            # Yabai window manager configuration
│   └── scripts/           # Space creation and window helper scripts
├── skhd/
│   └── skhdrc             # Hotkey bindings and shortcuts
├── sketchybar/
│   ├── sketchybarrc       # Status bar structure, styling & items
│   ├── plugins/           # System metric scripts (CPU, RAM, Disk, Clock)
│   └── xauusd/            # Real-time Go market daemon for live Gold (XAUUSD)
├── nvim/
│   ├── init.lua           # Neovim entry point
│   ├── lua/sabbir/        # Plugin manager (packer), keymaps, LSP & AI configs
│   └── after/             # Plugin configuration overrides
├── tmux/
│   └── .tmux.conf         # Tmux multiplexer layout & styling
├── flameshot/             # Flameshot configuration
└── zsh/                   # Zsh environment scripts & banner
```

---

## 📄 License
Personal configuration maintained by [Sabbir Hossain Shuvo](https://github.com/devlopersabbir). Feel free to customize and adapt for your own workflow!

-- set leader key to space

vim.g.mapleader = " "

local keymap = vim.keymap -- for conciseness

-- Terminal keybindings (terminal.lua file location at nvim/after/plugin/terminal.lua)
local terminal = require("sabbir.terminal")

-- use leader pv for go to the file explorer
keymap.set("n", "<leader>pv", vim.cmd.Ex)

-- use jk to exit insert mode
keymap.set("i", "jk", "<ESC>", { desc = "Exit insert mode with jk" })

vim.keymap.set("v", "J", ":m '>+1<CR>gv=gv")
vim.keymap.set("v", "K", ":m '<-2<CR>gv=gv")

-- nvim builtin terminal (NOTE: this is commented for terminal config file)
-- keymap.set("n", "<leader>tt", "<cmd>botright 15split | terminal<CR>")
-- keymap.set(
--     "t",
--     "<leader>xx",
--     [[<C-\><C-n>:b#<CR>]]
-- )
-- Terminal

-- Toggle terminal
keymap.set({ "n", "t" }, "<C-t>", function()
    terminal.toggle()
end, { desc = "Toggle Terminal" })

-- Focus editor
keymap.set("t", "<C-j>", function()
    terminal.focus_editor()
end, { desc = "Focus Editor" })

-- Focus terminal
keymap.set("n", "<C-k>", function()
    terminal.focus_terminal()
end, { desc = "Focus Terminal" })

-- Toggle terminal height
keymap.set({ "n", "t" }, "<C-\\>", function()
    terminal.toggle_height()
end, { desc = "Maximize Terminal" })

-- Exit terminal insert mode
keymap.set("t", "<Esc>", [[<C-\><C-n>]], {
    desc = "Terminal Normal Mode",
})

-- clear search highlights ---- new start from here....
keymap.set("n", "<leader>nh", ":nohl<CR>", { desc = "Clear search highlights" })

-- delete single character without copying into register
keymap.set("n", "x", '"_x')

-- window management
keymap.set("n", "<leader>sv", "<C-w>v", { desc = "Split window vertically" })     -- split window vertically
keymap.set("n", "<leader>sh", "<C-w>s", { desc = "Split window horizontally" })   -- split window horizontally
keymap.set("n", "<leader>se", "<C-w>=", { desc = "Make splits equal size" })      -- make split windows equal width & height
keymap.set("n", "<leader>sx", "<cmd>close<CR>", { desc = "Close current split" }) -- close current split window
keymap.set("n", "<leader>sm", "<cmd>MaximizerToggle<CR>", { desc = "Maximize/minimize a split" })

-- focus windows
keymap.set("n", "<leader>h", "<C-w>h", { desc = "Move to the window on the left" })
keymap.set("n", "<leader>j", "<C-w>j", { desc = "Move to the window below" })
keymap.set("n", "<leader>k", "<C-w>k", { desc = "Move to the window above" })
keymap.set("n", "<leader>l", "<C-w>l", { desc = "Move to the window on the right" })
keymap.set("n", "<leader>b", "<C-w>t", { desc = "Move to the TOP window" })
keymap.set("n", "<leader>t", "<C-w>b", { desc = "Move to the BOTTOM window" })

-- tab management
keymap.set("n", "<leader>to", "<cmd>tabnew<CR>", { desc = "Tab: New" })
keymap.set("n", "<leader>tf", "<cmd>tabnew %<CR>", { desc = "Tab: Current Buffer in New Tab" })
keymap.set("n", "<S-l>", "<cmd>tabnext<CR>", { desc = "Tab: Next" })
keymap.set("n", "<S-h>", "<cmd>tabprevious<CR>", { desc = "Tab: Previous" })
keymap.set("n", "<S-q>", "<cmd>tabclose<CR>", { desc = "Tab: Close" })

-- nvim-tree remap
keymap.set("n", "<leader>ee", "<cmd>NvimTreeToggle<CR>", { desc = "Toggle file explorer" })                         -- toggle file explorer
keymap.set("n", "<leader>eo", "<cmd>NvimTreeFocus<CR>", { desc = "Focus file explorer" })                           -- focus file explorer
keymap.set("n", "<leader>ef", "<cmd>NvimTreeFindFileToggle<CR>", { desc = "Toggle file explorer on current file" }) -- toggle file explorer on current file
keymap.set("n", "<leader>ec", "<cmd>NvimTreeCollapse<CR>", { desc = "Collapse file explorer" })                     -- collapse file explorer
keymap.set("n", "<leader>er", "<cmd>NvimTreeRefresh<CR>", { desc = "Refresh file explorer" })

-- ai
keymap.set("n", "<leader>ag", "<cmd>Antigravity<CR>", { desc = "Toggle Antigravity AI" })
keymap.set("v", "<leader>ag", function() require("antigravity").ask_selection() end,
    { desc = "Ask Antigravity with selection" })

-- Define a custom on_attach function to remap keys
local function on_attach(bufnr)
    local api = require("nvim-tree.api")

    local function opts(desc)
        return { desc = "nvim-tree: " .. desc, buffer = bufnr, noremap = true, silent = true, nowait = true }
    end

    -- Load the default mappings first
    api.config.mappings.default_on_attach(bufnr)

    -- Open in current tab (same window layout, keeping sidebar)
    vim.keymap.set("n", "<CR>", api.node.open.edit, opts("Open"))
    vim.keymap.set("n", "o", api.node.open.edit, opts("Open"))

    -- Open in new tab: 't', 'T', '<C-t>', or '<S-CR>'
    vim.keymap.set("n", "t", api.node.open.tab, opts("Open: New Tab"))
    vim.keymap.set("n", "T", api.node.open.tab, opts("Open: New Tab"))

    -- Open in split windows (bonus useful shortcuts)
    vim.keymap.set("n", "v", api.node.open.vertical, opts("Open: Vertical Split"))
    vim.keymap.set("n", "s", api.node.open.horizontal, opts("Open: Horizontal Split"))

    -- global Shift+H / Shift+L tab switching inside nvim-tree buffer
    vim.keymap.set("n", "<S-h>", "<cmd>tabprevious<CR>", opts("Tab: Previous"))
    vim.keymap.set("n", "H", "<cmd>tabprevious<CR>", opts("Tab: Previous"))
    vim.keymap.set("n", "<S-l>", "<cmd>tabnext<CR>", opts("Tab: Next"))
    vim.keymap.set("n", "L", "<cmd>tabnext<CR>", opts("Tab: Next"))
end

-- OR setup with some options
require("nvim-tree").setup({
    on_attach = on_attach,
    sort = {
        sorter = "case_sensitive",
    },
    view = {
        width = 28,
        relativenumber = true,
    },
    renderer = {
        group_empty = true,
        icons = {
            show = {
                file = true,
                folder = true,
                folder_arrow = true,
                git = true,
                modified = true,
            },
        },
    },
    filters = {
        dotfiles = false,
    },

    -- disable window_picker for
    -- explorer to work well with
    -- window splits
    actions = {
        open_file = {
            window_picker = {
                enable = false,
            },
        },
    },

    git = {
        ignore = false,
    },
})

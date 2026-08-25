-- This file can be loaded by calling `lua require('plugins')` from your init.vim

-- Only required if you have packer configured as `opt`
vim.cmd([[packadd packer.nvim]])

return require("packer").startup(function(use)
    -- Packer can manage itself
    use("wbthomason/packer.nvim")

    use({ "catppuccin/nvim", as = "catppuccin" })
    use("folke/tokyonight.nvim")
    use({ "AlexvZyl/nordic.nvim" })

    use({
        "nvim-telescope/telescope.nvim",
        branch = "0.1.x",
        requires = { { "nvim-lua/plenary.nvim" } },
    })

    use({
        "nvim-treesitter/nvim-treesitter",
        run = function()
            local ts_update = require("nvim-treesitter.install").update({ with_sync = true })
            ts_update()
        end,
    })

    use("nvim-treesitter/nvim-treesitter-context")

    use({
        "nvim-tree/nvim-tree.lua",
        requires = {
            "nvim-tree/nvim-web-devicons", -- optional
        },
    })

    use("nvim-lua/plenary.nvim")

    -- vim maximizser
    use("szw/vim-maximizer")

    -- auto pairs
    use({
        "windwp/nvim-autopairs",
        event = "InsertEnter",
        config = function()
            require("nvim-autopairs").setup({})
        end,
    })

    -- comment
    use("numToStr/Comment.nvim")

    -- formation
    use("stevearc/conform.nvim")


    use({
        "VonHeikemen/lsp-zero.nvim",
        branch = "v4.x",
        requires = {
            { "neovim/nvim-lspconfig" },
            { "williamboman/mason.nvim" },
            { "williamboman/mason-lspconfig.nvim" },

            -- Autocomplete
            { "hrsh7th/nvim-cmp" },
            { "saadparwaiz1/cmp_luasnip" },
            { "hrsh7th/cmp-path" },
            { "hrsh7th/cmp-nvim-lsp" },

            -- Snippets
            { "L3MON4D3/LuaSnip" },
            { "rafamadriz/friendly-snippets" },
        },
    })

    -- Ai code completation with avante
    use({
        "supermaven-inc/supermaven-nvim",
        config = function()
            require("supermaven-nvim").setup({})
        end,
    })
    -- dressing sidebar
    use({
        "stevearc/dressing.nvim",
        config = function()
            require("dressing").setup()
        end,
    })
    -- ai
    use({
        "NakLast/antigravity-cli.nvim",
        config = function()
            local cmd = vim.fn.exepath("agy")
            if cmd == "" then
                local local_bin = vim.fn.expand("~/.local/bin/agy")
                if vim.fn.executable(local_bin) == 1 then
                    cmd = local_bin
                else
                    cmd = "agy"
                end
            end
            require("antigravity").setup({
                cmd = cmd,
                position = "right",
                width_ratio = 0.35,
                height_ratio = 0.4,
                border = "rounded",
                style = "vsplit",
            })
        end,
    })
end)

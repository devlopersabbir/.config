local conform = require("conform")

conform.setup({
    formatters_by_ft = {
        javascript = { "prettier" },
        typescript = { "prettier" },
        javascriptreact = { "prettier" },
        typescriptreact = { "prettier" },

        json = { "prettier" },
        html = { "prettier" },
        css = { "prettier" },
        yaml = { "prettier" },
        markdown = { "prettier" },
        graphql = { "prettier" },

        lua = { "stylua" },

        go = { "goimports", "gofmt" },

        python = { "isort", "black" },

        dockerfile = { "prettier" }, -- simple + works fine
    },

    format_on_save = {
        lsp_fallback = true, -- 🔥 IMPORTANT (fixes your issue)
        async = false,
        timeout_ms = 2000,
    },
})

vim.keymap.set({ "n", "v" }, "<leader>fm", function()
    conform.format({
        lsp_fallback = true,
        async = false,
        timeout_ms = 2000,
    })
end, { desc = "Format file or range" })

vim.keymap.set({ "n", "v" }, "<leader>mp", function()
    conform.format({
        lsp_fallback = true,
        async = false,
        timeout_ms = 2000,
    })
end, { desc = "Make Pretty (Format)" })

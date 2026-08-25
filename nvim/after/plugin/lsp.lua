vim.api.nvim_create_autocmd("LspAttach", {
	callback = function(event)
		local km = vim.keymap.set
		local opts = { buffer = event.buf, remap = false }

		km("n", "gd", vim.lsp.buf.definition, opts)
		km("n", "gr", vim.lsp.buf.references, opts)
		km("n", "gi", vim.lsp.buf.implementation, opts)
		km("n", "K", vim.lsp.buf.hover, opts)

		km("n", "<leader>rn", vim.lsp.buf.rename, opts)

		km("n", "<leader>ca", vim.lsp.buf.code_action, opts)
	end,
})

-- =============================
-- List of servers
-- =============================
local servers = {
	"gopls", -- IMPORTANT
	"ts_ls",
	"lua_ls",
}

-- =============================
-- Special Lua config
-- =============================
vim.lsp.config("lua_ls", {
	settings = {
		Lua = {
			diagnostics = {
				globals = { "vim" },
			},
			workspace = {
				checkThirdParty = false,
			},
			telemetry = {
				enable = false,
			},
		},
	},
})

-- =============================
-- Go config
-- =============================
vim.lsp.config("gopls", {
	settings = {
		gopls = {
			completeUnimported = true,
			usePlaceholders = true,
			staticcheck = true,
			gofumpt = true,
		},
	},
})

-- =============================
-- Enable all servers
-- =============================
for _, server in ipairs(servers) do
	vim.lsp.enable(server)
end

-- =============================
-- Auto import on save
-- =============================
vim.api.nvim_create_autocmd("BufWritePre", {
	pattern = "*.go",
	callback = function()
		vim.lsp.buf.code_action({
			context = {
				only = { "source.organizeImports" },
			},
			apply = true,
		})
	end,
})

-- =============================
-- Diagnostics config
-- =============================
vim.diagnostic.config({
	virtual_text = true,
	signs = true,
	underline = true,
	update_in_insert = false,
	severity_sort = true,
})

-- =============================
-- Mason
-- =============================
require("mason").setup()

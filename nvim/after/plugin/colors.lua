require("tokyonight").setup({
	style = "moon",
	transparent = true,
	terminal_colors = true,

	styles = {
		sidebars = "transparent",
		floats = "transparent",
	},
})

vim.cmd.colorscheme("tokyonight")

local groups = {
	"Normal",
	"NormalNC",
	"NormalFloat",
	"SignColumn",
	"EndOfBuffer",
	"LineNr",
	"FoldColumn",
	"StatusLine",
	"StatusLineNC",
	"FloatBorder",
	"WinSeparator",

	"TelescopeNormal",
	"TelescopeBorder",
	"TelescopePromptNormal",
	"TelescopePromptBorder",
	"TelescopeResultsNormal",
	"TelescopeResultsBorder",
	"TelescopePreviewNormal",
	"TelescopePreviewBorder",

	"NvimTreeNormal",
	"NvimTreeNormalNC",
	"NvimTreeEndOfBuffer",
}

for _, group in ipairs(groups) do
	vim.api.nvim_set_hl(0, group, { bg = "NONE" })
end

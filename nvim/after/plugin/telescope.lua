local builtin = require("telescope.builtin")
local keymap = vim.keymap

-- Find
keymap.set("n", "<leader>ff", builtin.find_files, {
	desc = "Find Files",
})

keymap.set("n", "<leader>fg", builtin.live_grep, {
	desc = "Live Grep",
})

keymap.set("n", "<leader>fw", function()
	builtin.live_grep({
		default_text = vim.fn.expand("<cword>"),
	})
end, {
	desc = "Search Word Under Cursor",
})

keymap.set("n", "<leader>fb", builtin.buffers, {
	desc = "Buffers",
})

keymap.set("n", "<leader>fr", builtin.oldfiles, {
	desc = "Recent Files",
})

keymap.set("n", "<leader>fh", builtin.help_tags, {
	desc = "Help Tags",
})

----------------------------------------------------
-- Git
----------------------------------------------------

keymap.set("n", "<leader>gf", builtin.git_files, {
	desc = "Git Files",
})

keymap.set("n", "<leader>gc", builtin.git_commits, {
	desc = "Git Commits",
})

keymap.set("n", "<leader>gb", builtin.git_branches, {
	desc = "Git Branches",
})

keymap.set("n", "<leader>gs", builtin.git_status, {
	desc = "Git Status",
})
local comment = require("Comment")

comment.setup({
    pre_hook = function(ctx)
        local cs = vim.bo.commentstring
        if not cs or cs == "" then
            pcall(function()
                cs = vim.filetype.get_option(vim.bo.filetype, "commentstring")
            end)
        end
        if not cs or cs == "" then
            cs = "# %s"
        end
        return cs
    end,
})

-- Visual mode mappings: allow both `gc` and `gcc` to toggle linewise comment on selected lines
vim.keymap.set("x", "gcc", "<Plug>(comment_toggle_linewise_visual)", { desc = "Toggle comment on selection (linewise)" })
vim.keymap.set("x", "gbc", "<Plug>(comment_toggle_blockwise_visual)", { desc = "Toggle comment on selection (blockwise)" })


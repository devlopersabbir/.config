local status_ok, antigravity = pcall(require, "antigravity")
if not status_ok then
    return
end

local cmd = vim.fn.exepath("agy")
if cmd == "" then
    local local_bin = vim.fn.expand("~/.local/bin/agy")
    if vim.fn.executable(local_bin) == 1 then
        cmd = local_bin
    else
        cmd = "agy"
    end
end

antigravity.setup({
    cmd = cmd,
    position = "right", -- Opens on the right side (like VS Code secondary sidebar)
    style = "vsplit",   -- "vsplit" (vertical sidebar split) or "float" (floating right panel)
    width_ratio = 0.30, -- Width of the sidebar
    height_ratio = 0.8,
    border = "rounded",
})

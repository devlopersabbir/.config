local M = {}

local state = {
	buf = nil,
	win = nil,
	height = 15,
	maximized = false,
	previous_win = nil,
}

function M.toggle()
	-- Hide terminal if it's focused
	if state.win and vim.api.nvim_win_is_valid(state.win) and vim.api.nvim_get_current_win() == state.win then
		vim.api.nvim_win_close(state.win, true)
		state.win = nil

		if state.previous_win and vim.api.nvim_win_is_valid(state.previous_win) then
			vim.api.nvim_set_current_win(state.previous_win)
		end

		return
	end

	state.previous_win = vim.api.nvim_get_current_win()

	if state.win and vim.api.nvim_win_is_valid(state.win) then
		vim.api.nvim_set_current_win(state.win)
		vim.cmd("startinsert")
		return
	end

	vim.cmd("botright " .. state.height .. "split")

	if state.buf and vim.api.nvim_buf_is_valid(state.buf) then
		vim.api.nvim_win_set_buf(0, state.buf)
	else
		vim.cmd("terminal")
		state.buf = vim.api.nvim_get_current_buf()
	end

	state.win = vim.api.nvim_get_current_win()

	vim.cmd("startinsert")
end

function M.focus_terminal()
	if state.win and vim.api.nvim_win_is_valid(state.win) then
		state.previous_win = vim.api.nvim_get_current_win()
		vim.api.nvim_set_current_win(state.win)
		vim.cmd("startinsert")
	end
end

function M.focus_editor()
	vim.cmd("stopinsert")

	if state.previous_win and vim.api.nvim_win_is_valid(state.previous_win) then
		vim.api.nvim_set_current_win(state.previous_win)
	else
		vim.cmd("wincmd k")
	end
end

function M.toggle_height()
	if not (state.win and vim.api.nvim_win_is_valid(state.win)) then
		return
	end

	if state.maximized then
		vim.api.nvim_win_set_height(state.win, state.height)
	else
		vim.api.nvim_win_set_height(state.win, vim.o.lines - 3)
	end

	state.maximized = not state.maximized
end

return M

package view

import (
	"fmt"

	"github.com/arit-pal/gedit/internal/editor"
	"github.com/gdamore/tcell/v2"
)

// Draw renders the entire editor UI based on the current state.
func Draw(state *editor.State) {
	state.Screen.Clear()
	_, height := state.Screen.Size()
	textHeight := height - 1

	highlightStyle := tcell.StyleDefault.Reverse(true)

	for y := 0; y < textHeight; y++ {
		fileRow := y + state.RowOffset
		if fileRow >= 0 && fileRow < len(state.Content) {
			line := state.Content[fileRow]
			for x, r := range line {
				style := tcell.StyleDefault
				if fileRow == state.LastMatchY && x >= state.LastMatchX && x < state.LastMatchX+len(state.SearchQuery) {
					style = highlightStyle
				}
				state.Screen.SetContent(x, y, r, nil, style)
			}
		}
	}

	drawStatusBar(state)
	state.Screen.ShowCursor(state.CursorX, state.CursorY-state.RowOffset)
}

func drawStatusBar(state *editor.State) {
	width, height := state.Screen.Size()
	style := tcell.StyleDefault.Reverse(true)

	// Clear the status bar
	for i := 0; i < width; i++ {
		state.Screen.SetContent(i, height-1, ' ', nil, style)
	}

	// If there's a transient message, show it.
	if state.StatusMessage != "" {
		for i, r := range []rune(state.StatusMessage) {
			state.Screen.SetContent(i, height-1, r, nil, style)
		}
	} else {
		// Otherwise, show permanent info.
		dirtyIndicator := ""
		if state.IsDirty {
			dirtyIndicator = " [Modified]"
		}

		// Left side: Filename and dirty status
		fileInfo := fmt.Sprintf("'%s'%s - %d lines", state.FileName, dirtyIndicator, len(state.Content))
		for i, r := range []rune(fileInfo) {
			state.Screen.SetContent(i, height-1, r, nil, style)
		}

		// Right side: Cursor position
		cursorInfo := fmt.Sprintf("%d:%d", state.CursorY+1, state.CursorX+1)
		for i, r := range []rune(cursorInfo) {
			state.Screen.SetContent(width-len(cursorInfo)+i, height-1, r, nil, style)
		}
	}
}

// UpdateScrolling adjusts the rowOffset based on the cursor's position.
func UpdateScrolling(state *editor.State) {
	_, height := state.Screen.Size()
	textHeight := height - 1

	if state.CursorY < state.RowOffset {
		state.RowOffset = state.CursorY
	}

	if state.CursorY >= state.RowOffset+textHeight {
		state.RowOffset = state.CursorY - textHeight + 1
	}
}

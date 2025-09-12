package input

import (
	"strings"

	"github.com/arit-pal/gedit/internal/editor"
	"github.com/arit-pal/gedit/internal/file"
	"github.com/arit-pal/gedit/internal/view"
	"github.com/gdamore/tcell/v2"
)

// HandleKeyEvent processes keyboard events and modifies the editor state.
func HandleKeyEvent(ev *tcell.EventKey, state *editor.State) (quit bool) {
	state.StatusMessage = "" // Clear status message on any key press

	switch ev.Key() {
	case tcell.KeyCtrlX:
		return true
	case tcell.KeyCtrlS:
		if err := file.Save(state.FileName, state.Content); err != nil {
			state.StatusMessage = "Could not save file: " + err.Error()
		} else {
			state.StatusMessage = "File saved successfully!"
		}
	case tcell.KeyCtrlF:
		search(state)
	case tcell.KeyTab:
		line := state.Content[state.CursorY]
		spaces := []rune{' ', ' ', ' ', ' '}
		newLine := append(line[:state.CursorX], append(spaces, line[state.CursorX:]...)...)
		state.Content[state.CursorY] = newLine
		state.CursorX += 4
	case tcell.KeyUp:
		if state.CursorY > 0 {
			state.CursorY--
		}
	case tcell.KeyDown:
		if state.CursorY < len(state.Content)-1 {
			state.CursorY++
		}
	case tcell.KeyLeft:
		if state.CursorX > 0 {
			state.CursorX--
		}
	case tcell.KeyRight:
		if state.CursorY < len(state.Content) && state.CursorX < len(state.Content[state.CursorY]) {
			state.CursorX++
		}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		// Smart Backspace Logic.
		if state.CursorX > 0 {
			isSoftTab := false
			// Check if we are on a tab boundary and the preceding characters are all spaces.
			if state.CursorX%4 == 0 {
				line := state.Content[state.CursorY]
				if state.CursorX >= 4 {
					subSlice := line[state.CursorX-4 : state.CursorX]
					isAllSpaces := true
					for _, r := range subSlice {
						if r != ' ' {
							isAllSpaces = false
							break
						}
					}
					isSoftTab = isAllSpaces
				}
			}

			if isSoftTab {
				// Delete all 4 spaces of the soft tab
				line := state.Content[state.CursorY]
				state.Content[state.CursorY] = append(line[:state.CursorX-4], line[state.CursorX:]...)
				state.CursorX -= 4
			} else {
				// Default: delete one character
				line := state.Content[state.CursorY]
				state.Content[state.CursorY] = append(line[:state.CursorX-1], line[state.CursorX:]...)
				state.CursorX--
			}
		} else if state.CursorY > 0 {
			// Join with the line above
			newCursorX := len(state.Content[state.CursorY-1])
			state.Content[state.CursorY-1] = append(state.Content[state.CursorY-1], state.Content[state.CursorY]...)
			state.Content = append(state.Content[:state.CursorY], state.Content[state.CursorY+1:]...)
			state.CursorY--
			state.CursorX = newCursorX
		}
	case tcell.KeyDelete:
		line := state.Content[state.CursorY]
		if state.CursorX < len(line) {
			state.Content[state.CursorY] = append(line[:state.CursorX], line[state.CursorX+1:]...)
		} else if state.CursorY < len(state.Content)-1 {
			state.Content[state.CursorY] = append(line, state.Content[state.CursorY+1]...)
			state.Content = append(state.Content[:state.CursorY+1], state.Content[state.CursorY+2:]...)
		}
	case tcell.KeyEnter:
		line := state.Content[state.CursorY]
		restOfLine := line[state.CursorX:]
		state.Content[state.CursorY] = line[:state.CursorX]
		newLine := restOfLine
		state.Content = append(state.Content[:state.CursorY+1], append([][]rune{newLine}, state.Content[state.CursorY+1:]...)...)
		state.CursorY++
		state.CursorX = 0
	case tcell.KeyRune:
		line := state.Content[state.CursorY]
		newLine := append(line[:state.CursorX], append([]rune{ev.Rune()}, line[state.CursorX:]...)...)
		state.Content[state.CursorY] = newLine
		state.CursorX++
	}

	// Clamp cursor X position
	if state.CursorY < len(state.Content) && state.CursorX > len(state.Content[state.CursorY]) {
		state.CursorX = len(state.Content[state.CursorY])
	}

	view.UpdateScrolling(state)
	return false
}

func search(state *editor.State) {
	state.LastMatchX, state.LastMatchY = -1, -1
	query, ok := promptUser(state, "Search: ")
	if !ok {
		return
	}
	state.SearchQuery = query

	for y := state.CursorY; y < len(state.Content); y++ {
		line := string(state.Content[y])
		searchFromX := 0
		if y == state.CursorY {
			searchFromX = state.CursorX + 1
		}
		if x := strings.Index(line[searchFromX:], state.SearchQuery); x != -1 {
			state.CursorX = searchFromX + x
			state.CursorY = y
			state.LastMatchX = state.CursorX
			state.LastMatchY = state.CursorY
			view.UpdateScrolling(state)
			state.StatusMessage = ""
			return
		}
	}
	state.StatusMessage = "Search term not found"
}

func promptUser(state *editor.State, prompt string) (string, bool) {
	var input []rune
	for {
		state.StatusMessage = prompt + string(input)
		view.Draw(state)
		state.Screen.Show()
		ev := state.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return string(input), true
			case tcell.KeyEscape:
				state.StatusMessage = ""
				return "", false
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
			case tcell.KeyRune:
				input = append(input, ev.Rune())
			}
		}
	}
}

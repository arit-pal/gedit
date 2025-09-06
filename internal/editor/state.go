package editor

import "github.com/gdamore/tcell/v2"

// State holds all the data representing the current state of the editor.
type State struct {
	Screen        tcell.Screen
	Content       [][]rune
	CursorX       int
	CursorY       int
	RowOffset     int
	FileName      string
	StatusMessage string
	SearchQuery   string
	LastMatchX    int
	LastMatchY    int
}

// NewState creates the initial state of the editor.
func NewState(screen tcell.Screen, fileName string) *State {
	return &State{
		Screen:     screen,
		FileName:   fileName,
		LastMatchX: -1,
		LastMatchY: -1,
	}
}

package editor

import (
	"bufio"
	"os"

	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen  tcell.Screen
	content [][]rune
	cursorX int
	cursorY int
}

func NewEditor(fileName string) (*Editor, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)

	editor := &Editor{screen: screen}
	if err := editor.loadFile(fileName); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return editor, nil
}

func (e *Editor) loadFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		e.content = append(e.content, []rune(line))
	}

	return scanner.Err()
}

func (e *Editor) draw() {
	e.screen.Clear()
	for y, line := range e.content {
		for x, r := range line {
			e.screen.SetContent(x, y, r, nil, tcell.StyleDefault)
		}
	}
	e.screen.ShowCursor(e.cursorX, e.cursorY)
}

func (e *Editor) handleKeyEvent(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyUp:
		if e.cursorY > 0 {
			e.cursorY--
		}
	case tcell.KeyDown:
		if e.cursorY < len(e.content)-1 {
			e.cursorY++
		}
	case tcell.KeyLeft:
		if e.cursorX > 0 {
			e.cursorX--
		}
	case tcell.KeyRight:
		if e.cursorY < len(e.content) && e.cursorX < len(e.content[e.cursorY]) {
			e.cursorX++
		}
	case tcell.KeyRune:
		switch ev.Rune() {
		case 'x':
			return true
		}
	}

	if e.cursorY < len(e.content) && e.cursorX > len(e.content[e.cursorY]) {
		e.cursorX = len(e.content[e.cursorY])
	}

	return false
}

func (e *Editor) Start() {
	defer e.Finish()

	for {
		e.draw()
		e.screen.Show()

		ev := e.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.screen.Sync()
		case *tcell.EventKey:
			if quit := e.handleKeyEvent(ev); quit {
				return
			}
		}
	}
}

func (e *Editor) Finish() {
	e.screen.Fini()
}

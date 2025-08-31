package editor

import (
	"bufio"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen        tcell.Screen
	content       [][]rune
	cursorX       int
	cursorY       int
	rowOffset     int
	fileName      string
	statusMessage string
	searchQuery   string
	lastMatchX    int
	lastMatchY    int
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

	editor := &Editor{screen: screen, fileName: fileName, lastMatchX: -1, lastMatchY: -1}
	if err := editor.loadFile(fileName); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return editor, nil
}

func (e *Editor) saveFile() error {
	var lines []string
	for _, line := range e.content {
		lines = append(lines, string(line))
	}

	err := os.WriteFile(e.fileName, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		e.statusMessage = "Could not save file: " + err.Error()
		return err
	}

	e.statusMessage = "File saved successfully!"
	return nil
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

func (e *Editor) updateScrolling() {
	_, height := e.screen.Size()
	textHeight := height - 1

	if e.cursorY < e.rowOffset {
		e.rowOffset = e.cursorY
	}

	if e.cursorY >= e.rowOffset+textHeight {
		e.rowOffset = e.cursorY - textHeight + 1
	}
}

func (e *Editor) drawStatusBar() {
	width, height := e.screen.Size()
	style := tcell.StyleDefault.Reverse(true)

	for i := 0; i < width; i++ {
		e.screen.SetContent(i, height-1, ' ', nil, style)
	}

	for i, r := range []rune(e.statusMessage) {
		e.screen.SetContent(i, height-1, r, nil, style)
	}
}

func (e *Editor) promptUser(prompt string) (string, bool) {
	var input []rune
	for {
		e.statusMessage = prompt + string(input)
		e.draw()
		e.screen.Show()

		ev := e.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				return string(input), true
			case tcell.KeyEscape:
				e.statusMessage = ""
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

func (e *Editor) search() {
	e.lastMatchX, e.lastMatchY = -1, -1

	query, ok := e.promptUser("Search: ")
	if !ok {
		return
	}
	e.searchQuery = query

	for y := e.cursorY; y < len(e.content); y++ {
		line := string(e.content[y])
		searchFromX := 0
		if y == e.cursorY {
			searchFromX = e.cursorX + 1
		}

		if x := strings.Index(line[searchFromX:], e.searchQuery); x != -1 {
			e.cursorX = searchFromX + x
			e.cursorY = y
			e.lastMatchX = e.cursorX
			e.lastMatchY = e.cursorY
			e.updateScrolling()
			e.statusMessage = ""
			return
		}
	}

	e.statusMessage = "Search term not found"
}

func (e *Editor) draw() {
	e.screen.Clear()
	_, height := e.screen.Size()
	textHeight := height - 1

	highlightStyle := tcell.StyleDefault.Reverse(true)

	for y := 0; y < textHeight; y++ {
		fileRow := y + e.rowOffset
		if fileRow >= 0 && fileRow < len(e.content) {
			line := e.content[fileRow]
			for x, r := range line {
				style := tcell.StyleDefault
				if fileRow == e.lastMatchY && x >= e.lastMatchX && x < e.lastMatchX+len(e.searchQuery) {
					style = highlightStyle
				}
				e.screen.SetContent(x, y, r, nil, style)
			}
		}
	}

	e.drawStatusBar()
	e.screen.ShowCursor(e.cursorX, e.cursorY-e.rowOffset)
}

func (e *Editor) handleKeyEvent(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlX:
		return true
	case tcell.KeyCtrlS:
		e.saveFile()
	case tcell.KeyCtrlF:
		e.search()
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
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if e.cursorX > 0 {
			line := e.content[e.cursorY]
			e.content[e.cursorY] = append(line[:e.cursorX-1], line[e.cursorX:]...)
			e.cursorX--
		} else if e.cursorY > 0 {
			newCursorX := len(e.content[e.cursorY-1])
			e.content[e.cursorY-1] = append(e.content[e.cursorY-1], e.content[e.cursorY]...)
			e.content = append(e.content[:e.cursorY], e.content[e.cursorY+1:]...)
			e.cursorY--
			e.cursorX = newCursorX
		}
	case tcell.KeyEnter:
		line := e.content[e.cursorY]
		restOfLine := line[e.cursorX:]
		e.content[e.cursorY] = line[:e.cursorX]
		newLine := restOfLine
		e.content = append(e.content[:e.cursorY+1], append([][]rune{newLine}, e.content[e.cursorY+1:]...)...)
		e.cursorY++
		e.cursorX = 0
	case tcell.KeyRune:
		line := e.content[e.cursorY]
		newLine := append(line[:e.cursorX], append([]rune{ev.Rune()}, line[e.cursorX:]...)...)
		e.content[e.cursorY] = newLine
		e.cursorX++
	}

	if e.cursorY < len(e.content) && e.cursorX > len(e.content[e.cursorY]) {
		e.cursorX = len(e.content[e.cursorY])
	}

	e.updateScrolling()
	return false
}

func (e *Editor) Start() {
	e.updateScrolling()
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

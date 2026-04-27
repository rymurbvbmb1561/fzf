package fzf

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/junegunn/fzf/src/tui"
)

// Terminal represents the interactive terminal UI for fzf
type Terminal struct {
	mutex    sync.Mutex
	input    []rune
	cursor   int
	result   []string
	matches  []string
	selected map[int]bool
	eventBox *EventBox
	options  *Options
	tui      tui.Renderer
	prompt   string
	query    string
	count    int
	readChan chan string
	quitChan chan bool
}

// NewTerminal creates a new Terminal instance with the given options and event box
func NewTerminal(opts *Options, eventBox *EventBox) *Terminal {
	return &Terminal{
		input:    []rune{},
		cursor:   0,
		result:   []string{},
		matches:  []string{},
		selected: make(map[int]bool),
		eventBox: eventBox,
		options:  opts,
		prompt:   opts.Prompt,
		readChan: make(chan string, 256),
		quitChan: make(chan bool),
	}
}

// Start initializes the terminal and begins the event loop
func (t *Terminal) Start() {
	go t.render()
	t.loop()
}

// UpdateList updates the displayed list of matches
func (t *Terminal) UpdateList(items []string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.matches = items
	t.count = len(items)
}

// GetQuery returns the current search query
func (t *Terminal) GetQuery() string {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return string(t.input)
}

// render handles the visual rendering of the terminal UI
func (t *Terminal) render() {
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.draw()
		case <-t.quitChan:
			return
		}
	}
}

// draw renders the current state to the terminal
func (t *Terminal) draw() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// Clear current line and redraw prompt with current input
	fmt.Fprintf(os.Stderr, "\r\033[K%s%s",
		t.prompt,
		string(t.input),
	)
}

// loop handles user input events
func (t *Terminal) loop() {
	for {
		select {
		case line, ok := <-t.readChan:
			if !ok {
				return
			}
			t.handleInput(line)
		case <-t.quitChan:
			return
		}
	}
}

// handleInput processes a single input event
func (t *Terminal) handleInput(input string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	switch input {
	case "\r", "\n": // Enter
		t.accept()
	case "\x7f", "\b": // Backspace
		if len(t.input) > 0 && t.cursor > 0 {
			t.input = append(t.input[:t.cursor-1], t.input[t.cursor:]...)
			t.cursor--
		}
	case "\x03", "\x1b": // Ctrl-C or Escape
		t.quit()
	default:
		runes := []rune(input)
		t.input = append(t.input[:t.cursor], append(runes, t.input[t.cursor:]...)...)
		t.cursor += len(runes)
		t.query = strings.TrimSpace(string(t.input))
		t.eventBox.Set(EvtSearchNew, t.query)
	}
}

// accept finalizes selection and exits
func (t *Terminal) accept() {
	if len(t.matches) > 0 {
		t.result = []string{t.matches[0]}
	}
	t.quit()
}

// quit signals the terminal to stop
func (t *Terminal) quit() {
	close(t.quitChan)
}

// Results returns the final selected results
func (t *Terminal) Results() []string {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.result
}

package fzf

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

// Reader reads input from various sources (stdin, command output, filesystem)
// and feeds items into the event channel for processing.
type Reader struct {
	chan_  chan<- []byte
	eventBox *EventBox
	delimNil bool
	pauseMutex sync.Mutex
	paused     bool
}

// NewReader creates a new Reader instance.
func NewReader(ch chan<- []byte, eventBox *EventBox, delimNil bool) *Reader {
	return &Reader{
		chan_:    ch,
		eventBox: eventBox,
		delimNil: delimNil,
		paused:   false,
	}
}

// ReadSource reads from the given io.Reader, splitting on newline or null byte.
func (r *Reader) ReadSource(reader io.Reader, limit int64) bool {
	br := bufio.NewReaderSize(reader, 64*1024)
	delim := byte('\n')
	if r.delimNil {
		delim = 0
	}

	count := 0
	for {
		r.pauseMutex.Lock()
		paused := r.paused
		r.pauseMutex.Unlock()

		if paused {
			time.Sleep(50 * time.Millisecond)
			continue
		}

		line, err := br.ReadBytes(delim)
		if len(line) > 0 {
			// Strip the delimiter from the end
			if len(line) > 0 && line[len(line)-1] == delim {
				line = line[:len(line)-1]
			}
			// Strip carriage return on Windows-style line endings
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if len(line) > 0 {
				copy := make([]byte, len(line))
				copy_n := copy(copy, line)
				r.chan_ <- copy[:copy_n]
				count++
				if limit > 0 && int64(count) >= limit {
					return true
				}
			}
		}
		if err != nil {
			break
		}
	}
	return count > 0
}

// ReadStdin reads from standard input.
func (r *Reader) ReadStdin() bool {
	return r.ReadSource(os.Stdin, 0)
}

// Pause temporarily halts reading, useful during interactive resize or reload.
func (r *Reader) Pause(paused bool) {
	r.pauseMutex.Lock()
	r.paused = paused
	r.pauseMutex.Unlock()
}

// EventBox is a simple thread-safe event notification mechanism.
type EventBox struct {
	mutex  sync.Mutex
	cond   *sync.Cond
	events map[EventType]interface{}
}

// EventType represents a type of event in the event loop.
type EventType int

const (
	EvtReadNew EventType = iota
	EvtReadFin
	EvtSearchNew
	EvtSearchProgress
	EvtSearchFin
	EvtClose
)

// NewEventBox creates a new EventBox.
func NewEventBox() *EventBox {
	eb := &EventBox{events: make(map[EventType]interface{})}
	eb.cond = sync.NewCond(&eb.mutex)
	return eb
}

// Set records an event, optionally with an associated value.
func (eb *EventBox) Set(event EventType, value interface{}) {
	eb.mutex.Lock()
	eb.events[event] = value
	eb.mutex.Unlock()
	eb.cond.Broadcast()
}

// Wait blocks until at least one event is available, then calls the handler.
func (eb *EventBox) Wait(handler func(events map[EventType]interface{})) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	for len(eb.events) == 0 {
		eb.cond.Wait()
	}
	handler(eb.events)
	eb.events = make(map[EventType]interface{})
}

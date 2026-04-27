package fzf

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Layout options
const (
	LayoutDefault  = "default"
	LayoutReverse  = "reverse"
	LayoutReverseList = "reverse-list"
)

// Options holds all configuration options for fzf
type Options struct {
	// Input/Output
	Input     chan string
	Output    chan string
	PrintQuery bool

	// Search behavior
	Fuzzy     bool
	Extended  bool
	Case      Case
	Normalize bool
	Nth       []Range
	WithNth   []Range
	Delimiter Delimiter

	// Display
	Height    heightSpec
	Layout    string
	Reverse   bool
	Border    bool
	Prompt    string
	Pointer   string
	Marker    string
	Header    []string
	HeaderLines int

	// Behavior
	Multi     int
	NoMouse   bool
	Keybinds  map[string]string
	Expect    map[string]string
	ToggleSort bool

	// Performance
	Sort      int
	Tac       bool
	Sync      bool

	// Scripting
	Filter    *string
	Query     string
	Select1   bool
	Exit0     bool

	// Preview
	Preview   previewOpts
}

// Case sensitivity mode
type Case int

const (
	CaseSmart  Case = iota
	CaseIgnore
	CaseRespect
)

// heightSpec represents a height specification (absolute or percentage)
type heightSpec struct {
	size    float64
	percent bool
}

// Delimiter for splitting fields
type Delimiter struct {
	str *string
}

// Range represents a field range for --nth and --with-nth
type Range struct {
	begin int
	end   int
}

// previewOpts holds preview window configuration
type previewOpts struct {
	command  string
	position string
	size     float64
	percent  bool
	hidden   bool
	wrap     bool
}

// defaultOptions returns Options with sensible defaults
func defaultOptions() *Options {
	return &Options{
		Fuzzy:     true,
		Extended:  true,
		Case:      CaseSmart,
		Normalize: true,
		Layout:    LayoutReverse, // personal preference: reverse layout feels more natural
		Prompt:    "❯ ",         // nicer prompt character
		Pointer:   "▶",          // arrow pointer instead of >
		Marker:    "✓",          // checkmark for selected items
		Multi:     0,
		Sort:      1000,
		Keybinds:  make(map[string]string),
		Expect:    make(map[string]string),
		Height:    heightSpec{size: 50, percent: true}, // default to 50% height
		Preview:   previewOpts{position: "right", size: 50, percent: true},
	}
}

// parseHeight parses a height string like "40%" or "20"
func parseHeight(s string) (heightSpec, error) {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil || v < 0 || v > 100 {
			return heightSpec{}, fmt.Errorf("invalid height: %s", s)
		}
		return heightSpec{size: v, percent: true}, nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v < 0 {
		return heightSpec{}, fmt.Errorf("invalid height: %s", s)
	}
	return heightSpec{size: v, percent: false}, nil
}

// environ returns the value of an environment variable with a fallback
func environ(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

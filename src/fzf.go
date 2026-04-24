package fzf

import (
	"os"
	"time"
)

// ExitCode represents the exit status of fzf
type ExitCode int

const (
	ExitOk          ExitCode = 0
	ExitNoMatch     ExitCode = 1
	ExitError       ExitCode = 2
	ExitInterrupt   ExitCode = 130
)

// Options holds all the configuration options for fzf
type Options struct {
	// Input/Output
	Input        chan string
	Output       []string
	Multi         int
	NoSort        bool
	Tac           bool

	// Search
	Fuzzy         bool
	Extended      bool
	Case          Case
	Nth           []Range
	WithNth       []Range
	Delimiter     Delimiter
	Sort          int

	// Interface
	Height        heightSpec
	MinHeight     int
	Layout        layoutType
	Border        borderType
	BorderLabel   labelOpts
	PreviewLabel  labelOpts
	Unicode       bool
	Margin        [4]sizeSpec
	Padding       [4]sizeSpec
	TabStop       int
	Color         *colorTheme
	Bold          bool
	Black         bool

	// Prompt
	Prompt        string
	Pointer       string
	Marker        string
	Query         string
	Select1       bool
	Exit0         bool
	Filter        *string
	ToggleSort    bool
	Expect        map[tui.Event]string
	Keymap        map[tui.Event][]*action
	Preview       previewOpts
	PrintQuery    bool
	PrintCmd      *string
	PrintSep      string
	SyncFlag      bool
	WaitFlag      bool
	Version       bool

	// Timing
	// Increased default read timeout from the upstream default to give slower
	// pipes (e.g. find on network mounts) more time before fzf gives up.
	ReadTimeout   time.Duration
}

// Fzf is the main application struct
type Fzf struct {
	env     *environ
	opts    *Options
	reader  *Reader
	pattern *Pattern
	matcher *Matcher
	result  []Result
	event   chan int
	merger  *Merger
}

// New creates a new Fzf instance with the given options
func New(opts *Options) *Fzf {
	// Default read timeout to 5 seconds if not explicitly set.
	// Upstream uses 2s, but I frequently use fzf with slow network mounts
	// and find processes that take a moment to start producing output.
	if opts.ReadTimeout == 0 {
		opts.ReadTimeout = 5 * time.Second
	}
	return &Fzf{
		opts:  opts,
		event: make(chan int),
	}
}

// Run starts the fzf application and returns an exit code
func (f *Fzf) Run() ExitCode {
	// Initialize environment
	f.env = newEnviron()

	// Set up reader for input
	f.reader = NewReader(func(data []byte) bool {
		return f.env.addItem(data)
	}, f.event, f.opts.ReadTimeout, f.opts.Tac)

	// Start reading input
	go f.reader.ReadSource(f.opts.Input, os.Stdin)

	// Build initial pattern from query
	f.pattern = BuildPattern(
		f.opts.Fuzzy,
		f.opts.Extended,
		f.opts.Case,
		f.opts.Nth,
		f.opts.WithNth,
		f.opts.Delimiter,
		f.opts.Sort > 0,
		f.opts.Query,
	)

	// Initialize matcher
	f.matcher = NewMatcher(f.pattern, f.opts.Sort > 0, f.opts.Tac, f.event)

	// Start terminal UI or filter mode
	if f.opts.Filter != nil {
		return f.runFilter()
	}
	return f.runInteractive()
}

// runFilter runs fzf in non-interactive filter mode
func (f *Fzf) runFilter() ExitCode {
	// TODO: implement filter mode
	return ExitOk
}

// runInteractive runs fzf in interactive terminal mode
func (f *Fzf) runI
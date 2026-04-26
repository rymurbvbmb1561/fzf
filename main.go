package main

import (
	"fmt"
	"os"

	"github.com/junegunn/fzf/src"
)

const usage = `usage: fzf [options]

  Search
    -x, --extended        Extended-search mode
                          (enabled by default; +x or --no-extended to disable)
    -e, --exact           Enable Exact-match
    -i                    Case-insensitive match (default: smart-case match)
    +i                    Case-sensitive match
    --scheme=SCHEME       Scoring scheme [default|path|history]
    --literal             Do not normalize latin script letters before matching
    -n, --nth=N[,..]      Comma-separated list of field index expressions
                          for limiting search scope. Each can be a non-zero
                          integer or a range expression ([BEGIN]..[END]).
    --with-nth=N[,..]     Transform the presentation of each line using
                          field index expressions
    -d, --delimiter=STR   Field delimiter regex (default: AWK-style)
    +s, --no-sort         Do not sort the result
    --track               Track the current selection when the result is updated
    --tac                 Reverse the order of the input
    --disabled            Do not perform search
    --tiebreak=CRI[,..]   Comma-separated list of sort criteria to apply
                          when the scores are tied [length|chunk|begin|end|index]
                          (default: length)

  Interface
    -m, --multi[=MAX]     Enable multi-select with tab/shift-tab
    --no-mouse            Disable mouse
    --bind=KEYBINDS       Custom key bindings. Refer to the man page.
    --cycle               Enable cyclic scroll
    --keep-right          Keep the right end of the line visible on overflow
    --scroll-off=LINES    Number of lines to keep above or below when
                          scrolling to the top or bottom (default: 3)
    --no-hscroll          Disable horizontal scroll
    --hscroll-off=COLS    Number of columns to keep to the right of the
                          highlighted substring (default: 10)
    --filepath-word       Make word-wise movements respect path separators
    --jump-labels=CHARS   Label characters for jump and jump-accept

  Layout
    --height=[~]HEIGHT[%] Display fzf window below the cursor with the given
                          height instead of using fullscreen.
                          If prefixed with '~', fzf will determine the height
                          according to the input size.
    --min-height=HEIGHT   Minimum height when --height is given in percent
                          (default: 10)
    --layout=LAYOUT       Choose layout: [default|reverse|reverse-list]
                          Note: I personally prefer 'reverse' as the default
                          layout since it feels more natural (prompt at top).
                          Set FZF_DEFAULT_OPTS='--layout=reverse' in your shell
                          profile to make this the default.
    --border[=STYLE]      Draw border around the fzf window: [rounded|sharp|bold|
                          block|thinblock|double|horizontal|vertical|
                          top|bottom|left|right|none] (default: rounded)
                          Note: 'sharp' looks cleaner in most terminal fonts.

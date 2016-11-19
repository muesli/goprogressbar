package goprogressbar

import (
	"fmt"
	"math"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const progressBarFormat = "[#>-]"

// ProgressBar is a helper for printing a progress bar
type ProgressBar struct {
	Text             string
	RightAlignedText string
	Total            int64
	Current          int64
	Width            uint
}

// NewProgressBar returns a new progress bar
func NewProgressBar(text string, total, current int64, width uint) *ProgressBar {
	return &ProgressBar{
		Text:    text,
		Total:   total,
		Current: current,
		Width:   width,
	}
}

// Print writes the progress bar to stdout
func (p *ProgressBar) Print() {
	pct := float64(p.Current) / float64(p.Total)
	if p.Total == 0 {
		if p.Current == 0 {
			// When both Total and Current are 0, show a full progressbar
			pct = 1
		} else {
			pct = 0
		}
	}

	// percentage is bound between 0 and 1
	pct = math.Min(1, math.Max(0, pct))

	clearCurrentLine()

	pcts := fmt.Sprintf("%.2f%%", pct*100)
	for len(pcts) < 7 {
		pcts = " " + pcts
	}

	tiWidth, _, _ := terminal.GetSize(int(syscall.Stdin))
	barWidth := uint(math.Min(float64(p.Width), float64(tiWidth)/3.0))

	size := int(barWidth) - len(pcts) - 4
	fill := int(math.Max(2, math.Floor((float64(size)*pct)+.5)))
	if size < 16 {
		barWidth = 0
	}

	text := p.Text
	maxTextWidth := int(tiWidth) - 3 - int(barWidth) - len(p.RightAlignedText)
	if maxTextWidth < 0 {
		maxTextWidth = 0
	}
	if len(p.Text) > maxTextWidth {
		if len(p.Text)-maxTextWidth+3 < len(p.Text) {
			text = "..." + p.Text[len(p.Text)-maxTextWidth+3:]
		} else {
			text = ""
		}
	}

	// Print text
	s := fmt.Sprintf("%s%s  %s ",
		text,
		strings.Repeat(" ", maxTextWidth-len(text)),
		p.RightAlignedText)
	fmt.Print(s)

	if barWidth > 0 {
		progChar := progressBarFormat[2]
		if p.Current == p.Total {
			progChar = progressBarFormat[1]
		}

		// Print progress bar
		fmt.Printf("%c%s%c%s%c %s",
			progressBarFormat[0],
			strings.Repeat(string(progressBarFormat[1]), fill-1),
			progChar,
			strings.Repeat(string(progressBarFormat[3]), size-fill),
			progressBarFormat[4],
			pcts)
	}
}

func clearCurrentLine() {
	fmt.Print("\033[2K\r")
}

func MoveCursorUp(lines uint) {
	fmt.Printf("\033[%dA", lines)
}

func MoveCursorDown(lines uint) {
	fmt.Printf("\033[%dB", lines)
}

/*
 * goprogressbar
 *     Copyright (c) 2016-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package goprogressbar

import (
	"fmt"
	"math"
	"strings"
	"syscall"
	"time"

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

	lastPrintTime time.Time
}

// MultiProgressBar is a helper for printing multiple progress bars
type MultiProgressBar struct {
	ProgressBars []*ProgressBar

	lastPrintTime time.Time
}

func (p *ProgressBar) percentage() float64 {
	pct := float64(p.Current) / float64(p.Total)
	if p.Total == 0 {
		if p.Current == 0 {
			// When both Total and Current are 0, show a full progressbar
			pct = 1
		} else {
			pct = 0
		}
	}

	return pct
}

// UpdateRequired returns true when this progressbar wants an update regardless
// of fps limitation
func (p *ProgressBar) UpdateRequired() bool {
	return p.Current == 0 || p.Current == p.Total
}

// LazyPrint writes the progress bar to stdout if a significant update occurred
func (p *ProgressBar) LazyPrint() {
	now := time.Now()
	if p.UpdateRequired() || now.Sub(p.lastPrintTime) > time.Second/25 {
		// Max out at 25fps
		p.lastPrintTime = now
		p.Print()
	}
}

// Print writes the progress bar to stdout
func (p *ProgressBar) Print() {
	// percentage is bound between 0 and 1
	pct := math.Min(1, math.Max(0, p.percentage()))

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

// AddProgressBar adds another progress bar to the multi struct
func (mp *MultiProgressBar) AddProgressBar(p *ProgressBar) {
	mp.ProgressBars = append(mp.ProgressBars, p)

	if len(mp.ProgressBars) > 1 {
		fmt.Println()
	}
	mp.Print()
}

// Print writes all progress bars to stdout
func (mp *MultiProgressBar) Print() {
	moveCursorUp(uint(len(mp.ProgressBars)))

	for _, p := range mp.ProgressBars {
		moveCursorDown(1)
		p.Print()
	}
}

// LazyPrint writes all progress bars to stdout if a significant update occurred
func (mp *MultiProgressBar) LazyPrint() {
	forced := false
	for _, p := range mp.ProgressBars {
		if p.UpdateRequired() {
			forced = true
			break
		}
	}

	now := time.Now()
	if !forced {
		forced = now.Sub(mp.lastPrintTime) > time.Second/25
	}

	if forced {
		// Max out at 20fps
		mp.lastPrintTime = now

		moveCursorUp(uint(len(mp.ProgressBars)))
		for _, p := range mp.ProgressBars {
			moveCursorDown(1)
			p.Print()
		}
	}
}

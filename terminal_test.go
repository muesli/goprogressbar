/*
 * goprogressbar
 *     Copyright (c) 2016-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package goprogressbar

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCursorMovement(t *testing.T) {
	buf := &bytes.Buffer{}
	Stdout = buf

	moveCursorUp(5)
	if buf.String() != "\033[5A" {
		t.Errorf("Unexpected cursor up movement behaviour")
	}
	buf.Reset()

	moveCursorDown(5)
	if buf.String() != "\033[5B" {
		t.Errorf("Unexpected cursor down movement behaviour")
	}
	buf.Reset()
}

func TestPercentageSpecialValues2(t *testing.T) {
	p := ProgressBar{Current: 0, Total: 0}
	if p.percentage() != 1 {
		t.Errorf("percentage should be 1 when both current and total are 0, got: %f", p.percentage())
	}

	p = ProgressBar{Current: 100, Total: 0}
	if p.percentage() != 0 {
		t.Errorf("percentage should be 0 when current is greater than 0 but the total is unknown (0), got: %f", p.percentage())
	}
}

func TestProgressBarOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	Stdout = buf

	p := ProgressBar{Text: "Test", Current: 0, Total: 100, Width: 60}
	p.RightAlignedText = fmt.Sprintf("%d of %d", p.Current, p.Total)
	p.Print()
	if buf.String() != "\033[2K\rTest                           0 of 100 [#>---------------------------]   0.00%" {
		t.Errorf("Unexpected progressbar print behaviour")
	}
	buf.Reset()

	p.Current = 10
	p.RightAlignedText = fmt.Sprintf("%d of %d", p.Current, p.Total)
	p.Print()
	if buf.String() != "\033[2K\rTest                          10 of 100 [##>--------------------------]  10.00%" {
		t.Errorf("Unexpected progressbar print behaviour")
	}
	buf.Reset()

	p.Current = 100
	p.RightAlignedText = fmt.Sprintf("%d of %d", p.Current, p.Total)
	p.Print()
	if buf.String() != "\033[2K\rTest                         100 of 100 [#############################] 100.00%" {
		t.Errorf("Unexpected progressbar print behaviour")
	}
	buf.Reset()
}

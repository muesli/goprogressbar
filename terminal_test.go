/*
 * goprogressbar
 *     Copyright (c) 2016-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package goprogressbar

import (
	"bytes"
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

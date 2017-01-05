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
	"os"
	"testing"
)

func TestPercentageBound(t *testing.T) {
	p := ProgressBar{Current: -1, Total: 100}
	if p.percentage() != 0 {
		t.Errorf("percentage should be bound to 0, got: %f", p.percentage())
	}

	p = ProgressBar{Current: 200, Total: 100}
	if p.percentage() != 1 {
		t.Errorf("percentage should be bound to 1, got: %f", p.percentage())
	}
}

func TestPercentageSpecialValues(t *testing.T) {
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

func TestMultiProgressBarOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	Stdout = buf

	f, _ := os.Create("/tmp/foo.output")
	p1 := ProgressBar{Text: "Test1", Current: 23, Total: 100, Width: 60}
	p1.RightAlignedText = fmt.Sprintf("%d of %d", p1.Current, p1.Total)
	p2 := ProgressBar{Text: "Test2", Current: 69, Total: 100, Width: 60}
	p2.RightAlignedText = fmt.Sprintf("%d of %d", p2.Current, p2.Total)

	mp := MultiProgressBar{}
	mp.AddProgressBar(&p1)
	mp.AddProgressBar(&p2)

	if buf.String() != "\033[1A\033[1B\033[2K\rTest1                         23 of 100 [######>----------------------]  23.00%"+
		"\033[2A\033[1B\033[2K\rTest1                         23 of 100 [######>----------------------]  23.00%"+
		"\033[1B\033[2K\rTest2                         69 of 100 [###################>---------]  69.00%" {
		t.Errorf("Unexpected multi progressbar print behaviour")
	}
	f.Write(buf.Bytes())
	f.Close()
	buf.Reset()
}

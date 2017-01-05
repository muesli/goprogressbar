/*
 * goprogressbar
 *     Copyright (c) 2016-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE
 */

package goprogressbar

import "testing"

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

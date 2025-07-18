
package eeslism

import (
	"math"
	"testing"
)

func TestEOP(t *testing.T) {
	// Setup
	p := []*P_MENN{
		{wa: 0, wb: 90},
		{wa: 90, wb: 90},
		{wa: 45, wb: 45},
	}

	// Execute
	EOP(len(p), p)

	// Verify
	// Case 1
	expected1 := XYZ{X: -0, Y: -1, Z: 0}
	if math.Abs(p[0].e.X-expected1.X) > 1e-9 || math.Abs(p[0].e.Y-expected1.Y) > 1e-9 || math.Abs(p[0].e.Z-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got %v", expected1, p[0].e)
	}

	// Case 2
	expected2 := XYZ{X: -1, Y: -0, Z: 0}
	if math.Abs(p[1].e.X-expected2.X) > 1e-9 || math.Abs(p[1].e.Y-expected2.Y) > 1e-9 || math.Abs(p[1].e.Z-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got %v", expected2, p[1].e)
	}

	// Case 3
	expected3 := XYZ{X: -0.5, Y: -0.5, Z: 0.7071}
	if math.Abs(p[2].e.X-expected3.X) > 1e-9 || math.Abs(p[2].e.Y-expected3.Y) > 1e-9 || math.Abs(p[2].e.Z-expected3.Z) > 1e-9 {
		t.Errorf("Test Case 3 Failed: Expected %v, but got %v", expected3, p[2].e)
	}
}

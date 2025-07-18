
package eeslism

import (
	"math"
	"testing"
)

func TestHOUSEN(t *testing.T) {
	// Test case 1: A simple XY plane polygon
	lp1 := []*P_MENN{
		{
			P: []XYZ{{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0, 1, 0}},
		},
	}
	HOUSEN(lp1)
	expected1 := XYZ{0, 0, 1}
	if math.Abs(lp1[0].e.X-expected1.X) > 1e-9 || math.Abs(lp1[0].e.Y-expected1.Y) > 1e-9 || math.Abs(lp1[0].e.Z-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got %v", expected1, lp1[0].e)
	}

	// Test case 2: A simple XZ plane polygon
	lp2 := []*P_MENN{
		{
			P: []XYZ{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}},
		},
	}
	HOUSEN(lp2)
	expected2 := XYZ{0, -1, 0}
	if math.Abs(lp2[0].e.X-expected2.X) > 1e-9 || math.Abs(lp2[0].e.Y-expected2.Y) > 1e-9 || math.Abs(lp2[0].e.Z-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got %v", expected2, lp2[0].e)
	}
}

func TestHOUSEN2(t *testing.T) {
	// Test case 1: A simple XY plane triangle
	p0_1 := XYZ{0, 0, 0}
	p1_1 := XYZ{1, 0, 0}
	p2_1 := XYZ{0, 1, 0}
	var e1 XYZ
	HOUSEN2(&p0_1, &p1_1, &p2_1, &e1)
	expected1 := XYZ{0, 0, 1}
	if math.Abs(e1.X-expected1.X) > 1e-9 || math.Abs(e1.Y-expected1.Y) > 1e-9 || math.Abs(e1.Z-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got %v", expected1, e1)
	}

	// Test case 2: A simple XZ plane triangle
	p0_2 := XYZ{0, 0, 0}
	p1_2 := XYZ{1, 0, 0}
	p2_2 := XYZ{0, 0, 1}
	var e2 XYZ
	HOUSEN2(&p0_2, &p1_2, &p2_2, &e2)
	expected2 := XYZ{0, -1, 0}
	if math.Abs(e2.X-expected2.X) > 1e-9 || math.Abs(e2.Y-expected2.Y) > 1e-9 || math.Abs(e2.Z-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got %v", expected2, e2)
	}
}

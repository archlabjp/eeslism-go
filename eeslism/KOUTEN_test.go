
package eeslism

import (
	"math"
	"testing"
)

func TestKOUTEN(t *testing.T) {
	// Test case 1: Simple intersection with XY plane
	Qx1, Qy1, Qz1 := 0.0, 0.0, 1.0
	ls1, ms1, ns1 := 0.0, 0.0, -1.0
	lp1 := XYZ{0, 0, 0}
	E1 := XYZ{0, 0, 1}
	var Px1, Py1, Pz1 float64
	KOUTEN(Qx1, Qy1, Qz1, ls1, ms1, ns1, &Px1, &Py1, &Pz1, lp1, E1)
	expected1 := XYZ{0, 0, 0}
	if math.Abs(Px1-expected1.X) > 1e-9 || math.Abs(Py1-expected1.Y) > 1e-9 || math.Abs(Pz1-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got (%f, %f, %f)", expected1, Px1, Py1, Pz1)
	}

	// Test case 2: Intersection with a tilted plane
	Qx2, Qy2, Qz2 := 0.0, 0.0, 0.0
	ls2, ms2, ns2 := 1.0, 1.0, 1.0
	lp2 := XYZ{1, 1, 1}
	E2 := XYZ{1, 0, 0}
	var Px2, Py2, Pz2 float64
	KOUTEN(Qx2, Qy2, Qz2, ls2, ms2, ns2, &Px2, &Py2, &Pz2, lp2, E2)
	expected2 := XYZ{1, 1, 1}
	if math.Abs(Px2-expected2.X) > 1e-9 || math.Abs(Py2-expected2.Y) > 1e-9 || math.Abs(Pz2-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got (%f, %f, %f)", expected2, Px2, Py2, Pz2)
	}
}

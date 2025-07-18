
package eeslism

import (
	"math"
	"testing"
)

func TestINOROUT(t *testing.T) {
	// Test case 1: Point inside the triangle
	P0_1 := XYZ{0, 0, 0}
	P1_1 := XYZ{2, 0, 0}
	P2_1 := XYZ{1, 2, 0}
	Px_1, Py_1, Pz_1 := 1.0, 0.5, 0.0
	var S1, T1 float64
	INOROUT(Px_1, Py_1, Pz_1, P0_1, P1_1, P2_1, &S1, &T1)
	if S1 < 0 || S1 > 1 || T1 < 0 || T1 > 1 || S1+T1 > 1 {
		t.Errorf("Test Case 1 Failed: Point should be inside. S=%f, T=%f", S1, T1)
	}

	// Test case 2: Point outside the triangle
	P0_2 := XYZ{0, 0, 0}
	P1_2 := XYZ{2, 0, 0}
	P2_2 := XYZ{1, 2, 0}
	Px_2, Py_2, Pz_2 := 3.0, 1.0, 0.0
	var S2, T2 float64
	INOROUT(Px_2, Py_2, Pz_2, P0_2, P1_2, P2_2, &S2, &T2)
	if !(S2 < 0 || S2 > 1 || T2 < 0 || T2 > 1 || S2+T2 > 1) {
		t.Errorf("Test Case 2 Failed: Point should be outside. S=%f, T=%f", S2, T2)
	}

	// Test case 3: Point on an edge
	P0_3 := XYZ{0, 0, 0}
	P1_3 := XYZ{2, 0, 0}
	P2_3 := XYZ{1, 2, 0}
	Px_3, Py_3, Pz_3 := 1.0, 0.0, 0.0
	var S3, T3 float64
	INOROUT(Px_3, Py_3, Pz_3, P0_3, P1_3, P2_3, &S3, &T3)
	if math.Abs(S3-0.5) > 1e-9 || math.Abs(T3-0.0) > 1e-9 {
		t.Errorf("Test Case 3 Failed: Point should be on an edge. S=%f, T=%f", S3, T3)
	}
}

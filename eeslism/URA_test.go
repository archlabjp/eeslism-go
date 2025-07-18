package eeslism

import (
	"math"
	"testing"
)

func TestURA(t *testing.T) {
	// Setup
	LP := []*P_MENN{
		{
			polyd: 4,
			P:     []XYZ{{1, 1, 1}, {2, 1, 1}, {2, 2, 1}, {1, 2, 1}},
		},
	}
	OP := []*P_MENN{
		{
			e:     XYZ{0, 0, 1},
			polyd: 4,
			P:     []XYZ{{0, 0, 0}, {1, 0, 0}, {1, 1, 0}, {0, 1, 0}},
		},
	}
	bektSlice := []*bekt{
		{
			ps: make([][]float64, 1),
		},
	}
	bektSlice[0].ps[0] = make([]float64, 4)

	// Execute
	URA(1, 1, LP, bektSlice, OP)

	// Verify
	expected := -1.0
	for i := 0; i < 4; i++ {
		if math.Abs(bektSlice[0].ps[0][i]-expected) < 1e-9 {
			t.Errorf("Test Case Failed: Expected t[0].ps[0][%d] to be %f, but got %f", i, expected, bektSlice[0].ps[0][i])
		}
	}
}

func TestURA_M(t *testing.T) {
	// Test case 1: Vector pointing in the same direction
	s1 := URA_M(0, 0, 1, 0)
	if math.Abs(s1-1.0) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected s to be 1.0, but got %f", s1)
	}

	// Test case 2: Vector pointing in the opposite direction
	s2 := URA_M(0, 0, -1, 0)
	if math.Abs(s2-(-1.0)) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected s to be -1.0, but got %f", s2)
	}

	// Test case 3: Vector is perpendicular
	s3 := URA_M(1, 0, 0, 0)
	if math.Abs(s3-0.0) > 1e-9 {
		t.Errorf("Test Case 3 Failed: Expected s to be 0.0, but got %f", s3)
	}
}

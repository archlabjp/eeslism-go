
package eeslism

import (
	"math"
	"testing"
)

func TestYOGEN(t *testing.T) {
	// Test case 1: Vectors pointing in the same direction
	var s1 float64
	YOGEN(0, 0, 0, 1, 0, 0, &s1, XYZ{1, 0, 0})
	if math.Abs(s1-1.0) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected s to be 1.0, but got %f", s1)
	}

	// Test case 2: Vectors pointing in opposite directions
	var s2 float64
	YOGEN(0, 0, 0, 1, 0, 0, &s2, XYZ{-1, 0, 0})
	if math.Abs(s2 - (-1.0)) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected s to be -1.0, but got %f", s2)
	}

	// Test case 3: Perpendicular vectors
	var s3 float64
	YOGEN(0, 0, 0, 1, 0, 0, &s3, XYZ{0, 1, 0})
	if math.Abs(s3-0.0) > 1e-9 {
		t.Errorf("Test Case 3 Failed: Expected s to be 0.0, but got %f", s3)
	}

	// Test case 4: Zero vector
	var s4 float64
	YOGEN(0, 0, 0, 0, 0, 0, &s4, XYZ{1, 1, 1})
	if s4 != -777 {
		t.Errorf("Test Case 4 Failed: Expected s to be -777, but got %f", s4)
	}
}

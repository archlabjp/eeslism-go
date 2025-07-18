
package eeslism

import (
	"math"
	"testing"
)

func TestZAHYOU(t *testing.T) {
	// Test case 1: No rotation
	Op1 := XYZ{1, 2, 3}
	G1 := XYZ{0, 0, 0}
	var op1_res XYZ
	ZAHYOU(Op1, G1, &op1_res, 0, 0)
	expected1 := XYZ{1, 2, 3}
	if math.Abs(op1_res.X-expected1.X) > 1e-9 || math.Abs(op1_res.Y-expected1.Y) > 1e-9 || math.Abs(op1_res.Z-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got %v", expected1, op1_res)
	}

	// Test case 2: 90-degree rotation around Z-axis
	Op2 := XYZ{1, 0, 0}
	G2 := XYZ{0, 0, 0}
	var op2_res XYZ
	ZAHYOU(Op2, G2, &op2_res, 90, 0)
	expected2 := XYZ{0, 1, 0}
	if math.Abs(op2_res.X-expected2.X) > 1e-9 || math.Abs(op2_res.Y-expected2.Y) > 1e-9 || math.Abs(op2_res.Z-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got %v", expected2, op2_res)
	}
}

func TestR_ZAHYOU(t *testing.T) {
	// Test case 1: No rotation
	Op1 := XYZ{1, 2, 3}
	G1 := XYZ{0, 0, 0}
	var op1_res XYZ
	R_ZAHYOU(Op1, G1, &op1_res, 0, 0)
	expected1 := XYZ{1, 2, 3}
	if math.Abs(op1_res.X-expected1.X) > 1e-9 || math.Abs(op1_res.Y-expected1.Y) > 1e-9 || math.Abs(op1_res.Z-expected1.Z) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected %v, but got %v", expected1, op1_res)
	}

	// Test case 2: 90-degree rotation around Z-axis
	Op2 := XYZ{0, 1, 0}
	G2 := XYZ{0, 0, 0}
	var op2_res XYZ
	R_ZAHYOU(Op2, G2, &op2_res, 90, 0)
	expected2 := XYZ{1, 0, 0}
	if math.Abs(op2_res.X-expected2.X) > 1e-9 || math.Abs(op2_res.Y-expected2.Y) > 1e-9 || math.Abs(op2_res.Z-expected2.Z) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected %v, but got %v", expected2, op2_res)
	}
}

func TestZAHYOU_R_ZAHYOU_Inverse(t *testing.T) {
	// Test that R_ZAHYOU is the inverse of ZAHYOU
	Op := XYZ{1, 2, 3}
	G := XYZ{4, 5, 6}
	wa := 30.0
	wb := 60.0

	var transformed, reversed XYZ
	ZAHYOU(Op, G, &transformed, wa, wb)
	R_ZAHYOU(transformed, G, &reversed, wa, wb)

	if math.Abs(Op.X-reversed.X) > 1e-4 || math.Abs(Op.Y-reversed.Y) > 1e-4 || math.Abs(Op.Z-reversed.Z) > 1e-4 {
		t.Errorf("Inverse Test Failed: Original %v, but got %v after transform and reverse", Op, reversed)
	}
}

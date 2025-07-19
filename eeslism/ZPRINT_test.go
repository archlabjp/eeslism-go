package eeslism

import (
	"os"
	"strings"
	"testing"
)

func TestZPRINT(t *testing.T) {
	// Create test data
	testFileName := "tmp_zprinttest_output"

	// Clean up test file after test
	defer func() {
		os.Remove(testFileName + ".gchi")
	}()

	// Create test OP data
	op := []P_MENN{
		{
			opname: "TestOP1",
			polyd:  3,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 1.0, Y: 0.0, Z: 0.0},
				{X: 0.5, Y: 1.0, Z: 0.0},
			},
			wd: 1,
			opw: []WD_MENN{
				{
					opwname: "TestWindow1",
					P: []XYZ{
						{X: 0.1, Y: 0.1, Z: 0.0},
						{X: 0.9, Y: 0.1, Z: 0.0},
						{X: 0.9, Y: 0.9, Z: 0.0},
						{X: 0.1, Y: 0.9, Z: 0.0},
					},
				},
			},
		},
	}

	// Create test LP data
	lp := []P_MENN{
		{
			opname: "TestLP1",
			polyd:  4,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 1.0},
				{X: 2.0, Y: 0.0, Z: 1.0},
				{X: 2.0, Y: 2.0, Z: 1.0},
				{X: 0.0, Y: 2.0, Z: 1.0},
			},
			e: XYZ{X: 0.0, Y: 0.0, Z: 1.0}, // Normal vector
		},
	}

	// Call ZPRINT function
	ZPRINT(lp, op, 1, 1, testFileName)

	// Check if file was created
	outputFile := testFileName + ".gchi"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("ZPRINT() should create file %s", outputFile)
		return
	}

	// Read and verify file content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
		return
	}

	contentStr := string(content)

	// Check for OP data in output
	if !strings.Contains(contentStr, "op[0] TestOP1") {
		t.Errorf("Output should contain OP name 'TestOP1'")
	}

	if !strings.Contains(contentStr, "P[0] X=0.000000 Y=0.000000 Z=0.000000") {
		t.Errorf("Output should contain OP vertex coordinates")
	}

	if !strings.Contains(contentStr, "op[0] opw[0] TestWindow1") {
		t.Errorf("Output should contain window name 'TestWindow1'")
	}

	// Check for LP data in output
	if !strings.Contains(contentStr, "lp[0] TestLP1") {
		t.Errorf("Output should contain LP name 'TestLP1'")
	}

	if !strings.Contains(contentStr, "e.X=0.000000 e.Y=0.000000 e.Z=1.000000") {
		t.Errorf("Output should contain LP normal vector")
	}

	if !strings.Contains(contentStr, "P[0] X=0.000000 Y=0.000000 Z=1.000000") {
		t.Errorf("Output should contain LP vertex coordinates")
	}
}

func TestZPRINT_EmptyData(t *testing.T) {
	testFileName := "tmp_zprinttest_empty"

	// Clean up test file after test
	defer func() {
		os.Remove(testFileName + ".gchi")
	}()

	// Test with empty arrays
	var op []P_MENN
	var lp []P_MENN

	// Should not panic with empty data
	ZPRINT(lp, op, 0, 0, testFileName)

	// Check if file was created (should be empty or minimal content)
	outputFile := testFileName + ".gchi"
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("ZPRINT() should create file %s even with empty data", outputFile)
		return
	}

	// File should exist but be mostly empty
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
		return
	}

	// Content should be minimal (just newlines or empty)
	if len(content) > 10 { // Allow for some minimal content
		t.Errorf("Output file should be mostly empty for empty input data, got %d bytes", len(content))
	}
}

func TestZPRINT_MultiplePolygons(t *testing.T) {
	testFileName := "tmp_zprinttest_multiple"

	// Clean up test file after test
	defer func() {
		os.Remove(testFileName + ".gchi")
	}()

	// Create test data with multiple polygons
	op := []P_MENN{
		{
			opname: "OP1",
			polyd:  3,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 1.0, Y: 0.0, Z: 0.0},
				{X: 0.5, Y: 1.0, Z: 0.0},
			},
			wd: 0, // No windows
		},
		{
			opname: "OP2",
			polyd:  4,
			P: []XYZ{
				{X: 2.0, Y: 0.0, Z: 0.0},
				{X: 3.0, Y: 0.0, Z: 0.0},
				{X: 3.0, Y: 1.0, Z: 0.0},
				{X: 2.0, Y: 1.0, Z: 0.0},
			},
			wd: 0, // No windows
		},
	}

	lp := []P_MENN{
		{
			opname: "LP1",
			polyd:  3,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 2.0},
				{X: 1.0, Y: 0.0, Z: 2.0},
				{X: 0.5, Y: 1.0, Z: 2.0},
			},
			e: XYZ{X: 0.0, Y: 0.0, Z: -1.0},
		},
		{
			opname: "LP2",
			polyd:  4,
			P: []XYZ{
				{X: 2.0, Y: 0.0, Z: 2.0},
				{X: 3.0, Y: 0.0, Z: 2.0},
				{X: 3.0, Y: 1.0, Z: 2.0},
				{X: 2.0, Y: 1.0, Z: 2.0},
			},
			e: XYZ{X: 0.0, Y: 0.0, Z: -1.0},
		},
	}

	// Call ZPRINT function
	ZPRINT(lp, op, 2, 2, testFileName)

	// Read and verify file content
	outputFile := testFileName + ".gchi"
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
		return
	}

	contentStr := string(content)

	// Check that both OP polygons are present
	if !strings.Contains(contentStr, "op[0] OP1") {
		t.Errorf("Output should contain first OP polygon")
	}
	if !strings.Contains(contentStr, "op[1] OP2") {
		t.Errorf("Output should contain second OP polygon")
	}

	// Check that both LP polygons are present
	if !strings.Contains(contentStr, "lp[0] LP1") {
		t.Errorf("Output should contain first LP polygon")
	}
	if !strings.Contains(contentStr, "lp[1] LP2") {
		t.Errorf("Output should contain second LP polygon")
	}
}

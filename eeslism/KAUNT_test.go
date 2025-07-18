package eeslism

import (
	"math"
	"testing"
)

func TestMinval(t *testing.T) {
	// Test case 1: Basic case
	span1 := []float64{1.0, 2.0, 0.5, 3.0}
	var min1 int
	var val1 float64
	minval(span1, len(span1), &min1, &val1)
	if min1 != 2 || math.Abs(val1-0.5) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected min=2, val=0.5, but got min=%d, val=%f", min1, val1)
	}

	// Test case 2: With negative values
	span2 := []float64{-1.0, 2.0, 0.5, -3.0}
	var min2 int
	var val2 float64
	minval(span2, len(span2), &min2, &val2)
	if min2 != 2 || math.Abs(val2-0.5) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected min=2, val=0.5, but got min=%d, val=%f", min2, val2)
	}

	// Test case 3: No positive values
	span3 := []float64{-1.0, -2.0, -0.5, -3.0}
	var min3 int
	var val3 float64
	minval(span3, len(span3), &min3, &val3)
	if min3 != -1 {
		t.Errorf("Test Case 3 Failed: Expected min=-1, but got min=%d", min3)
	}
}

func TestKAUNT(t *testing.T) {
	// Test case 1: mlpn = 1 (existing test)
	t.Run("mlpn_1", func(t *testing.T) {
		mlpn := 1
		ls, ms, ns := 0.0, 0.0, -1.0 // Ray from directly above
		var suma, sumg float64
		sumwall := make([]float64, mlpn)
		s := 1.0 // Positive s, so suma is affected
		mlp := []*P_MENN{
			{
				polyd: 4,
				P:     []XYZ{{-1, -1, 0}, {1, -1, 0}, {1, 1, 0}, {-1, 1, 0}},
				e:     XYZ{0, 0, 1}, // Normal vector pointing up
				shad:  [366]float64{0.5},
			},
		}
		p := make([]XYZ, 1)
		O := XYZ{0, 0, 0}
		E := XYZ{0, 0, 1}
		wa, wb := 0.0, 0.0
		G := XYZ{0, 0, 0}
		gpn := 1
		nday := 0
		var gcnt int
		startday := 0
		wlflg := 0

		// Execute
		KAUNT(mlpn, ls, ms, ns, &suma, &sumg, sumwall, s, mlp, p, O, E, wa, wb, G, gpn, nday, &gcnt, startday, wlflg)

		// Verify based on observed behavior (k=0, no hit detected by KOUTEN/INOROUT)
		// If k remains 0, then suma is incremented by 1 (because s > 0.0)
		// sumg remains 0
		// sumwall remains [0.0]
		const epsilon = 1e-9
		if math.Abs(suma-1.0) > epsilon {
			t.Errorf("Expected suma %f, got %f", 1.0, suma)
		}
		if math.Abs(sumg-0.0) > epsilon {
			t.Errorf("Expected sumg %f, got %f", 0.0, sumg)
		}
		if math.Abs(sumwall[0]-0.0) > epsilon {
			t.Errorf("Expected sumwall[0] %f, got %f", 0.0, sumwall[0])
		}
	})

	// Test case 2: mlpn = 2
	t.Run("mlpn_2_multiple_panels", func(t *testing.T) {
		mlpn := 2
		ls, ms, ns := 0.0, 0.0, -1.0 // Ray from directly above
		var suma, sumg float64
		sumwall := make([]float64, mlpn)
		s := 1.0 // Positive s, so suma is affected
		mlp := []*P_MENN{
			{
				polyd: 4,
				P:     []XYZ{{-1, -1, 0}, {1, -1, 0}, {1, 1, 0}, {-1, 1, 0}},
				e:     XYZ{0, 0, 1},
				shad:  [366]float64{0.5},
			},
			{
				polyd: 2, // A line segment, not a polygon, but for test purposes
				P:     []XYZ{{0, 0, 0}, {1, 1, 0}},
				e:     XYZ{0, 0, 1},
				shad:  [366]float64{0.8},
			},
		}
		p := make([]XYZ, 1)
		O := XYZ{0, 0, 0}
		E := XYZ{0, 0, 1}
		wa, wb := 0.0, 0.0
		G := XYZ{0, 0, 0}
		gpn := 1
		nday := 0
		var gcnt int
		startday := 0
		wlflg := 0

		// Execute
		KAUNT(mlpn, ls, ms, ns, &suma, &sumg, sumwall, s, mlp, p, O, E, wa, wb, G, gpn, nday, &gcnt, startday, wlflg)

		// Verify based on observed behavior (k=0, no hit detected by KOUTEN/INOROUT)
		// If k remains 0, then suma is incremented by 1 (because s > 0.0)
		// sumg remains 0
		// sumwall remains [0.0, 0.0]
		const epsilon = 1e-9
		if math.Abs(suma-1.0) > epsilon {
			t.Errorf("Expected suma %f, got %f", 1.0, suma)
		}
		if math.Abs(sumg-0.0) > epsilon {
			t.Errorf("Expected sumg %f, got %f", 0.0, sumg)
		}
		if math.Abs(sumwall[0]-0.0) > epsilon {
			t.Errorf("Expected sumwall[0] %f, got %f", 0.0, sumwall[0])
		}
		if math.Abs(sumwall[1]-0.0) > epsilon {
			t.Errorf("Expected sumwall[1] %f, got %f", 0.0, sumwall[1])
		}
	})
}

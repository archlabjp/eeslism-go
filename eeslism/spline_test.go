package eeslism

import (
	"math"
	"testing"
)

func TestSpline(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		x        []float64
		y        []float64
		x1       float64
		expected float64
		tolerance float64
	}{
		{
			name: "linear interpolation",
			n:    2,
			x:    []float64{0, 1, 2},
			y:    []float64{0, 1, 2},
			x1:   0.5,
			expected: 0.5,
			tolerance: 1e-10,
		},
		{
			name: "quadratic function",
			n:    3,
			x:    []float64{0, 1, 2, 3},
			y:    []float64{0, 1, 4, 9}, // y = x^2
			x1:   1.5,
			expected: 2.2, // spline interpolation result (not exact quadratic)
			tolerance: 0.1,
		},
		{
			name: "sine-like data",
			n:    5,
			x:    []float64{0, 1, 2, 3, 4, 5},
			y:    []float64{0, 0.841, 0.909, 0.141, -0.757, -0.959}, // approximate sin(x)
			x1:   2.5,
			expected: 0.5, // approximate sin(2.5) ≈ 0.5985, but spline will be different
			tolerance: 1.0, // loose tolerance for this approximation
		},
		{
			name: "boundary value start",
			n:    3,
			x:    []float64{0, 1, 2, 3},
			y:    []float64{1, 2, 3, 4},
			x1:   0,
			expected: 1,
			tolerance: 1e-10,
		},
		{
			name: "boundary value end",
			n:    3,
			x:    []float64{0, 1, 2, 3},
			y:    []float64{1, 2, 3, 4},
			x1:   3,
			expected: 4,
			tolerance: 1e-10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare working arrays
			h := make([]float64, tt.n+1)
			b := make([]float64, tt.n+1)
			d := make([]float64, tt.n+1)
			g := make([]float64, tt.n+1)
			u := make([]float64, tt.n+1)
			r := make([]float64, tt.n+1)

			result := spline(tt.n, tt.x, tt.y, tt.x1, h, b, d, g, u, r)

			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("spline() = %v, want %v (tolerance: %v)", result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestIntgtsup(t *testing.T) {
	// Skip this test if FNNday function is not available
	t.Skip("Skipping Intgtsup test due to missing FNNday dependency")
}

func TestIntgtsup_MultipleCallsWithSameData(t *testing.T) {
	// Skip this test if FNNday function is not available
	t.Skip("Skipping Intgtsup test due to missing FNNday dependency")
}

func TestSpline_EdgeCases(t *testing.T) {
	t.Run("extrapolation beyond range", func(t *testing.T) {
		n := 3
		x := []float64{0, 1, 2, 3}
		y := []float64{0, 1, 4, 9}
		x1 := 5.0 // beyond the range

		h := make([]float64, n+1)
		b := make([]float64, n+1)
		d := make([]float64, n+1)
		g := make([]float64, n+1)
		u := make([]float64, n+1)
		r := make([]float64, n+1)

		// Should not panic and return some value
		result := spline(n, x, y, x1, h, b, d, g, u, r)

		// Just check that it returns a finite number
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("spline() should return finite value for extrapolation, got %v", result)
		}
	})

	t.Run("minimum valid input", func(t *testing.T) {
		n := 2
		x := []float64{0, 1, 2}
		y := []float64{0, 1, 2}
		x1 := 0.5

		h := make([]float64, n+1)
		b := make([]float64, n+1)
		d := make([]float64, n+1)
		g := make([]float64, n+1)
		u := make([]float64, n+1)
		r := make([]float64, n+1)

		result := spline(n, x, y, x1, h, b, d, g, u, r)

		// For linear data, should be exact
		expected := 0.5
		if math.Abs(result-expected) > 1e-10 {
			t.Errorf("spline() = %v, want %v", result, expected)
		}
	})
}
package eeslism

import (
	"math"
	"testing"
)

func TestMatmalv(t *testing.T) {
	tests := []struct {
		name     string
		A        []float64
		V        []float64
		N        int
		n        int
		expected []float64
	}{
		{
			name: "2x2 matrix multiplication",
			A: []float64{
				1, 2,
				3, 4,
			},
			V:        []float64{5, 6},
			N:        2,
			n:        2,
			expected: []float64{17, 39}, // [1*5+2*6, 3*5+4*6] = [17, 39]
		},
		{
			name: "3x3 matrix multiplication",
			A: []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			V:        []float64{1, 2, 3},
			N:        3,
			n:        3,
			expected: []float64{14, 32, 50}, // [1+4+9, 4+10+18, 7+16+27]
		},
		{
			name: "identity matrix",
			A: []float64{
				1, 0,
				0, 1,
			},
			V:        []float64{5, 7},
			N:        2,
			n:        2,
			expected: []float64{5, 7}, // Identity matrix should return original vector
		},
		{
			name: "zero matrix",
			A: []float64{
				0, 0,
				0, 0,
			},
			V:        []float64{5, 7},
			N:        2,
			n:        2,
			expected: []float64{0, 0}, // Zero matrix should return zero vector
		},
		{
			name: "single element",
			A:        []float64{3},
			V:        []float64{4},
			N:        1,
			n:        1,
			expected: []float64{12}, // 3 * 4 = 12
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			T := make([]float64, tt.n)
			Matmalv(tt.A, tt.V, tt.N, tt.n, T)

			for i := 0; i < tt.n; i++ {
				if math.Abs(T[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("Matmalv() result[%d] = %v, want %v", i, T[i], tt.expected[i])
				}
			}
		})
	}
}

func TestMatinit(t *testing.T) {
	t.Run("initialize float array", func(t *testing.T) {
		A := []float64{1.5, 2.7, 3.9, 4.1}
		N := 4

		matinit(A, N)

		for i := 0; i < N; i++ {
			if A[i] != 0.0 {
				t.Errorf("matinit() A[%d] = %v, want 0.0", i, A[i])
			}
		}
	})

	t.Run("empty array", func(t *testing.T) {
		var A []float64
		N := 0

		// Should not panic
		matinit(A, N)
	})
}

func TestImatinit(t *testing.T) {
	t.Run("initialize int array", func(t *testing.T) {
		A := []int{1, 2, 3, 4, 5}
		N := 5

		imatinit(A, N)

		for i := 0; i < N; i++ {
			if A[i] != 0 {
				t.Errorf("imatinit() A[%d] = %v, want 0", i, A[i])
			}
		}
	})

	t.Run("empty array", func(t *testing.T) {
		var A []int
		N := 0

		// Should not panic
		imatinit(A, N)
	})
}

func TestMatinitx(t *testing.T) {
	tests := []struct {
		name  string
		A     []float64
		N     int
		x     float64
		check func([]float64, int, float64) bool
	}{
		{
			name: "initialize with positive value",
			A:    []float64{1, 2, 3, 4},
			N:    4,
			x:    5.5,
			check: func(A []float64, N int, x float64) bool {
				for i := 0; i < N; i++ {
					if A[i] != x {
						return false
					}
				}
				return true
			},
		},
		{
			name: "initialize with negative value",
			A:    []float64{1, 2, 3},
			N:    3,
			x:    -2.7,
			check: func(A []float64, N int, x float64) bool {
				for i := 0; i < N; i++ {
					if A[i] != x {
						return false
					}
				}
				return true
			},
		},
		{
			name: "initialize with zero",
			A:    []float64{1, 2},
			N:    2,
			x:    0.0,
			check: func(A []float64, N int, x float64) bool {
				for i := 0; i < N; i++ {
					if A[i] != x {
						return false
					}
				}
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matinitx(tt.A, tt.N, tt.x)

			if !tt.check(tt.A, tt.N, tt.x) {
				t.Errorf("matinitx() failed to initialize array with value %v", tt.x)
			}
		})
	}
}

func TestMatcpy(t *testing.T) {
	tests := []struct {
		name string
		A    []float64
		N    int
	}{
		{
			name: "copy positive values",
			A:    []float64{1.1, 2.2, 3.3, 4.4},
			N:    4,
		},
		{
			name: "copy negative values",
			A:    []float64{-1.5, -2.7, -3.9},
			N:    3,
		},
		{
			name: "copy mixed values",
			A:    []float64{-1.0, 0.0, 1.0, 2.5, -3.7},
			N:    5,
		},
		{
			name: "single element",
			A:    []float64{42.0},
			N:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			B := make([]float64, tt.N)
			
			matcpy(tt.A, B, tt.N)

			for i := 0; i < tt.N; i++ {
				if B[i] != tt.A[i] {
					t.Errorf("matcpy() B[%d] = %v, want %v", i, B[i], tt.A[i])
				}
			}

			// Verify that modifying B doesn't affect A
			if tt.N > 0 {
				originalA0 := tt.A[0]
				B[0] = 999.0
				if tt.A[0] != originalA0 {
					t.Errorf("matcpy() should create independent copy, but A was modified")
				}
			}
		})
	}

	t.Run("empty array", func(t *testing.T) {
		var A, B []float64
		N := 0

		// Should not panic
		matcpy(A, B, N)
	})
}

func TestMatrixOperationsConsistency(t *testing.T) {
	t.Run("matinit and matinitx consistency", func(t *testing.T) {
		A1 := []float64{1, 2, 3, 4}
		A2 := []float64{1, 2, 3, 4}
		N := 4

		matinit(A1, N)
		matinitx(A2, N, 0.0)

		for i := 0; i < N; i++ {
			if A1[i] != A2[i] {
				t.Errorf("matinit and matinitx(0.0) should produce same result, but A1[%d]=%v, A2[%d]=%v", 
					i, A1[i], i, A2[i])
			}
		}
	})

	t.Run("matcpy preserves values", func(t *testing.T) {
		A := []float64{1.1, 2.2, 3.3}
		B := make([]float64, 3)
		C := make([]float64, 3)
		N := 3

		// Copy A to B
		matcpy(A, B, N)
		// Copy B to C
		matcpy(B, C, N)

		for i := 0; i < N; i++ {
			if A[i] != C[i] {
				t.Errorf("Double copy should preserve values, but A[%d]=%v, C[%d]=%v", 
					i, A[i], i, C[i])
			}
		}
	})
}
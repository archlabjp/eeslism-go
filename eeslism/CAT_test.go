package eeslism

import (
	"math"
	"testing"
)

func TestCAT(t *testing.T) {
	tests := []struct {
		name     string
		inputA   float64
		inputB   float64
		inputC   float64
		expectedA float64
		expectedB float64
		expectedC float64
	}{
		{
			name:      "positive values",
			inputA:    1.23456789,
			inputB:    2.34567890,
			inputC:    3.45678901,
			expectedA: 1.2346,
			expectedB: 2.3457,
			expectedC: 3.4568,
		},
		{
			name:      "negative values",
			inputA:    -1.23456789,
			inputB:    -2.34567890,
			inputC:    -3.45678901,
			expectedA: -1.2346,
			expectedB: -2.3457,
			expectedC: -3.4568,
		},
		{
			name:      "zero values",
			inputA:    0.0,
			inputB:    0.0,
			inputC:    0.0,
			expectedA: 0.0,
			expectedB: 0.0,
			expectedC: 0.0,
		},
		{
			name:      "negative zero handling",
			inputA:    -0.00001,
			inputB:    -0.00001,
			inputC:    -0.00001,
			expectedA: 0.0,
			expectedB: 0.0,
			expectedC: 0.0,
		},
		{
			name:      "rounding up",
			inputA:    1.23455,
			inputB:    2.34565,
			inputC:    3.45675,
			expectedA: 1.2346,
			expectedB: 2.3457,
			expectedC: 3.4568,
		},
		{
			name:      "rounding down",
			inputA:    1.23454,
			inputB:    2.34564,
			inputC:    3.45674,
			expectedA: 1.2345,
			expectedB: 2.3456,
			expectedC: 3.4567,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, b, c := tt.inputA, tt.inputB, tt.inputC
			CAT(&a, &b, &c)

			if math.Abs(a-tt.expectedA) > 1e-10 {
				t.Errorf("CAT() a = %v, want %v", a, tt.expectedA)
			}
			if math.Abs(b-tt.expectedB) > 1e-10 {
				t.Errorf("CAT() b = %v, want %v", b, tt.expectedB)
			}
			if math.Abs(c-tt.expectedC) > 1e-10 {
				t.Errorf("CAT() c = %v, want %v", c, tt.expectedC)
			}
		})
	}
}

func TestCAT_NegativeZeroHandling(t *testing.T) {
	// Test specific case for -0.0 handling
	a, b, c := -0.0, -0.0, -0.0
	CAT(&a, &b, &c)

	// Check that -0.0 is converted to 0.0
	if math.Signbit(a) {
		t.Errorf("CAT() should convert -0.0 to 0.0 for a, got %v", a)
	}
	if math.Signbit(b) {
		t.Errorf("CAT() should convert -0.0 to 0.0 for b, got %v", b)
	}
	if math.Signbit(c) {
		t.Errorf("CAT() should convert -0.0 to 0.0 for c, got %v", c)
	}
}
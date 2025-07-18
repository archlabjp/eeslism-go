package eeslism

import (
	"math"
	"testing"
)

func TestPRA(t *testing.T) {
	tests := []struct {
		name     string
		ls, ms, ns float64
		x, y, z  float64
		expectedU float64
	}{
		{
			name:     "Case 2: ls is non-zero",
			ls:       1.0,
			ms:       0.0,
			ns:       0.0,
			x:        5.0,
			y:        0.0,
			z:        0.0,
			expectedU: 5.0, // U = x / ls
		},
		{
			name:     "Case 3: ls is zero, ms is non-zero",
			ls:       0.0,
			ms:       2.0,
			ns:       0.0,
			x:        0.0,
			y:        10.0,
			z:        0.0,
			expectedU: 5.0, // U = y / ms
		},
		{
			name:     "Case 4: ls and ms are zero, ns is non-zero",
			ls:       0.0,
			ms:       0.0,
			ns:       3.0,
			x:        0.0,
			y:        0.0,
			z:        15.0,
			expectedU: 5.0, // U = z / ns
		},
		{
			name:     "Case 5: Mixed non-zero components",
			ls:       1.0,
			ms:       2.0,
			ns:       3.0,
			x:        5.0,
			y:        10.0,
			z:        15.0,
			expectedU: 5.0, // U should be consistent across non-zero components
		},
		{
			name:     "Case 6: Negative components",
			ls:       -1.0,
			ms:       -2.0,
			ns:       -3.0,
			x:        -5.0,
			y:        -10.0,
			z:        -15.0,
			expectedU: 5.0,
		},
		{
			name:     "Case 7: Division by zero (ls) - should use ms or ns",
			ls:       0.0,
			ms:       1.0,
			ns:       1.0,
			x:        1.0,
			y:        1.0,
			z:        1.0,
			expectedU: 1.0, // U = y / ms
		},
		{
			name:     "Case 8: Division by zero (ms) - should use ls or ns",
			ls:       1.0,
			ms:       0.0,
			ns:       1.0,
			x:        1.0,
			y:        1.0,
			z:        1.0,
			expectedU: 1.0, // U = x / ls
		},
		{
			name:     "Case 9: Division by zero (ns) - should use ls or ms",
			ls:       1.0,
			ms:       1.0,
			ns:       0.0,
			x:        1.0,
			y:        1.0,
			z:        1.0,
			expectedU: 1.0, // U = x / ls
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			U := 0.0
			PRA(&U, tt.ls, tt.ms, tt.ns, tt.x, tt.y, tt.z)

			// Use a small epsilon for float comparisons
			const epsilon = 1e-9
			if math.Abs(U-tt.expectedU) > epsilon {
				t.Errorf("Expected U=%f, got %f", tt.expectedU, U)
			}
		})
	}
}


package eeslism

import (
	"math"
	"testing"
)

func TestRAND(t *testing.T) {
	const numTests = 1000
	for i := 0; i < numTests; i++ {
		var a, v float64
		RAND(&a, &v)

		if a < 0 || a > 2*math.Pi {
			t.Errorf("Test %d: Azimuth 'a' out of range [0, 2*Pi], got %f", i, a)
		}

		if v < 0 || v > math.Pi/2 {
			t.Errorf("Test %d: Elevation 'v' out of range [0, Pi/2], got %f", i, v)
		}
	}
}

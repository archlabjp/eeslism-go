
package eeslism

import (
	"math"
	"testing"
)

func TestCINC(t *testing.T) {
	// Test case 1: Sun is directly in front of the surface
	op1 := &P_MENN{wa: 0, wb: 90}
	ls1, ms1, ns1 := 0.0, 1.0, 0.0
	var co1 float64
	CINC(op1, ls1, ms1, ns1, &co1)
	if math.Abs(co1 - (-1.0)) > 1e-9 {
		t.Errorf("Test Case 1 Failed: Expected co to be -1.0, but got %f", co1)
	}

	// Test case 2: Sun is parallel to the surface
	op2 := &P_MENN{wa: 90, wb: 90}
	ls2, ms2, ns2 := 0.0, 1.0, 0.0
	var co2 float64
	CINC(op2, ls2, ms2, ns2, &co2)
	if math.Abs(co2-0.0) > 1e-9 {
		t.Errorf("Test Case 2 Failed: Expected co to be 0.0, but got %f", co2)
	}

	// Test case 3: Sun is at a 45-degree angle
	op3 := &P_MENN{wa: 45, wb: 90}
	ls3, ms3, ns3 := 0.0, 1.0, 0.0
	var co3 float64
	CINC(op3, ls3, ms3, ns3, &co3)
	if math.Abs(co3-(-math.Sin(45*math.Pi/180))) > 1e-9 {
		t.Errorf("Test Case 3 Failed: Expected co to be %f, but got %f", math.Sin(45*math.Pi/180), co3)
	}
}

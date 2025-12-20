package eeslism

import (
	"testing"
)

func TestFNXtrSimple(t *testing.T) {
	// Initialize psychrometric constants
	Psyint()

	// Test FNXtr at T=15°C, RH=50%
	xr := FNXtr(15.0, 50.0)
	t.Logf("FNXtr(15, 50) = %f kg/kg", xr)

	// Expected: around 0.00538 kg/kg
	if xr < 0.004 || xr > 0.006 {
		t.Errorf("FNXtr returned %f, expected around 0.00538", xr)
	}

	// Test FNXtr at T=11.5°C, RH=50%
	xr2 := FNXtr(11.5, 50.0)
	t.Logf("FNXtr(11.5, 50) = %f kg/kg", xr2)

	// Expected: around 0.0042 kg/kg
	if xr2 < 0.003 || xr2 > 0.005 {
		t.Errorf("FNXtr returned %f, expected around 0.0042", xr2)
	}

	// Test FNPws at T=15°C (saturation pressure)
	pws := FNPws(15.0)
	t.Logf("FNPws(15) = %f kPa", pws)

	// Expected: around 1.7 kPa
	if pws < 1.5 || pws > 2.0 {
		t.Errorf("FNPws returned %f, expected around 1.7", pws)
	}
}

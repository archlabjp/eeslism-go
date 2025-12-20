
package eeslism

import (
	"testing"
)

func TestRdpnlldsptr(t *testing.T) {
	// Setup
	load := ControlSWType('t')
	key := []string{"", "Tout"}
	Rdpnl := &RDPNL{}
	var idmrk byte

	// Execute
	vptr, err := rdpnlldsptr(&load, key, Rdpnl, &idmrk)

	// Verify
	if err != nil {
		t.Errorf("rdpnlldsptr failed with error: %v", err)
	}
	if vptr.Ptr == nil {
		t.Errorf("vptr.Ptr should not be nil")
	}
	if idmrk != 't' {
		t.Errorf("idmrk should be 't', but got %c", idmrk)
	}
}

func TestRdpnlldsschd(t *testing.T) {
	t.Run("LoadtNotNil_TosetAboveLimit", func(t *testing.T) {
		// Test Loadt != nil with Toset > TEMPLIMIT
		loadt := ControlSWType('t')
		Rdpnl := &RDPNL{
			Loadt: &loadt,
			Toset: TEMPLIMIT + 1, // Above TEMPLIMIT
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}},
			},
		}

		rdpnlldsschd(Rdpnl)

		if Rdpnl.cmp.Elouts[0].Control != ON_SW {
			t.Errorf("Expected Control=ON_SW when Toset > TEMPLIMIT, but got %d", Rdpnl.cmp.Elouts[0].Control)
		}
	})

	t.Run("LoadtNotNil_TosetBelowLimit", func(t *testing.T) {
		// Test Loadt != nil with Toset <= TEMPLIMIT (OFF case)
		loadt := ControlSWType('t')
		Rdpnl := &RDPNL{
			Loadt: &loadt,
			Toset: TEMPLIMIT - 1, // Below TEMPLIMIT (-100)
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}},
			},
		}

		rdpnlldsschd(Rdpnl)

		if Rdpnl.cmp.Elouts[0].Control != OFF_SW {
			t.Errorf("Expected Control=OFF_SW when Toset <= TEMPLIMIT, but got %d", Rdpnl.cmp.Elouts[0].Control)
		}
	})

	t.Run("LoadtNotNil_EoControlOff", func(t *testing.T) {
		// Test Loadt != nil but Eo.Control == OFF_SW (should not change)
		loadt := ControlSWType('t')
		Rdpnl := &RDPNL{
			Loadt: &loadt,
			Toset: 25.0,
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: OFF_SW}}, // Already OFF
			},
		}

		rdpnlldsschd(Rdpnl)

		// Should remain OFF_SW
		if Rdpnl.cmp.Elouts[0].Control != OFF_SW {
			t.Errorf("Expected Control to remain OFF_SW, but got %d", Rdpnl.cmp.Elouts[0].Control)
		}
	})

	t.Run("LoadtNil_NoChange", func(t *testing.T) {
		// Test Loadt == nil (no load pointer - should not modify anything)
		Rdpnl := &RDPNL{
			Loadt: nil,
			Toset: 25.0,
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}},
			},
		}

		origControl := Rdpnl.cmp.Elouts[0].Control
		rdpnlldsschd(Rdpnl)

		if Rdpnl.cmp.Elouts[0].Control != origControl {
			t.Errorf("Control should not change when Loadt=nil, but got %d", Rdpnl.cmp.Elouts[0].Control)
		}
	})
}

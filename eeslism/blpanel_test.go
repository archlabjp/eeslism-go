
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
	// Setup
	loadt := ControlSWType('t')
	Rdpnl := &RDPNL{
		Loadt: &loadt,
		Toset: TEMPLIMIT + 1,
		cmp: &COMPNT{
			Elouts: []*ELOUT{{Control: 0}},
		},
	}

	// Execute
	rdpnlldsschd(Rdpnl)

	// Verify
	if Rdpnl.cmp.Elouts[0].Control != ON_SW {
		t.Errorf("Rdpnl.cmp.Elouts[0].Control should be ON_SW, but got %d", Rdpnl.cmp.Elouts[0].Control)
	}
}

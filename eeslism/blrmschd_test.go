
package eeslism

import (
	"testing"
)

func TestWindowschdlr(t *testing.T) {
	// Setup
	isw := []ControlSWType{'B'}
	windows := []*WINDOW{
		{Rwall: 0.1},
		{Rwall: 0.2},
	}
	ds := []*RMSRF{
		{
			ble:    BLE_Window,
			Nfn:    2,
			fnsw:   0,
			fnmrk:  [10]rune{'A', 'B'},
			fnd:    [10]int{0, 1},
			Ctlif:  &CTLIF{},
			ifwin:  &WINDOW{Rwall: 0.3},
		},
	}

	// Execute
	Windowschdlr(isw, windows, ds)

	// Verify
	if ds[0].Rwall != 0.2 {
		t.Errorf("Windowschdlr failed: Rwall should be 0.2, but got %f", ds[0].Rwall)
	}
}

func TestQischdlr(t *testing.T) {
	// Setup
	hmsch := 1.0
	hmwksch := 1.0
	lightsch := 1.0
	assch := 1.0
	alsch := 1.0
	Room := &ROOM{
		Name:      "TestRoom",
		Nhm:       1,
		Hmsch:     &hmsch,
		Hmwksch:   &hmwksch,
		Light:     100,
		Lightsch:  &lightsch,
		Apsc:      50,
		Apsr:      50,
		Assch:     &assch,
		Apl:       20,
		Alsch:     &alsch,
		cmp:       &COMPNT{Elouts: []*ELOUT{{}}},
	}

	// Execute
	Room.Qischdlr()

	// Verify
	if Room.Hc == 0 {
		t.Errorf("Hc should not be zero")
	}
}

func TestVtschdlr(t *testing.T) {
	// Setup
	visc := 1.0
	vesc := 1.0
	rooms := []*ROOM{
		{
			Gvi:  1.0,
			Visc: &visc,
			Gve:  2.0,
			Vesc: &vesc,
		},
	}

	// Execute
	Vtschdlr(rooms)

	// Verify
	if rooms[0].Gvent != 3.0 {
		t.Errorf("Gvent should be 3.0, but got %f", rooms[0].Gvent)
	}
}

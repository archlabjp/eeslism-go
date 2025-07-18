package eeslism

import (
	"testing"
)

func TestVentdata(t *testing.T) {
	// Setup
	fi := NewEeTokens("Living Vent=(1.0,sch1) Inf=(0.5,sch2) *\n")
	Schdl := &SCHDL{
		Sch: []SCH{{name: "sch1"}, {name: "sch2"}},
		Val: []float64{0.8, 0.6},
	}

	Room := []*ROOM{
		{Name: "Living"},
	}
	Simc := &SIMCONTL{}

	// Execute
	Ventdata(fi, Schdl, Room, Simc)

	// Verify
	if Room[0].Gve != 1.0 {
		t.Errorf("Gve should be 1.0, but got %f", Room[0].Gve)
	}
	if *Room[0].Vesc != 0.8 {
		t.Errorf("Vesc should be 0.8, but got %f", *Room[0].Vesc)
	}
	if Room[0].Gvi != 0.5 {
		t.Errorf("Gvi should be 0.5, but got %f", Room[0].Gvi)
	}
	if *Room[0].Visc != 0.6 {
		t.Errorf("Visc should be 0.6, but got %f", *Room[0].Visc)
	}
}

func TestAichschdlr(t *testing.T) {
	// Setup
	val := []float64{0.5, 0.2}
	rooms := []*ROOM{
		{
			Nachr: 2,
			achr: []*ACHIR{
				{sch: 0},
				{sch: 1},
			},
		},
	}

	// Execute
	Aichschdlr(val, rooms)

	// Verify
	if rooms[0].achr[0].Gvr != 0.5 {
		t.Errorf("Gvr for achr[0] should be 0.5, but got %f", rooms[0].achr[0].Gvr)
	}
	if rooms[0].achr[1].Gvr != 0.2 {
		t.Errorf("Gvr for achr[1] should be 0.2, but got %f", rooms[0].achr[1].Gvr)
	}
}

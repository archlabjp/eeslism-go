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

func TestAichschdlr_ZeroNegativeValues(t *testing.T) {
	// Test: val <= 0.0 should set Gvr to 0.0
	val := []float64{0.0, -0.5, 0.3}
	rooms := []*ROOM{
		{
			Nachr: 3,
			achr: []*ACHIR{
				{sch: 0, Gvr: 1.0}, // Should become 0.0
				{sch: 1, Gvr: 2.0}, // Should become 0.0 (negative value)
				{sch: 2, Gvr: 0.0}, // Should become 0.3
			},
		},
	}

	Aichschdlr(val, rooms)

	// Verify zero value case
	if rooms[0].achr[0].Gvr != 0.0 {
		t.Errorf("Gvr for achr[0] should be 0.0 (val=0.0), but got %f", rooms[0].achr[0].Gvr)
	}
	// Verify negative value case
	if rooms[0].achr[1].Gvr != 0.0 {
		t.Errorf("Gvr for achr[1] should be 0.0 (val=-0.5), but got %f", rooms[0].achr[1].Gvr)
	}
	// Verify positive value case
	if rooms[0].achr[2].Gvr != 0.3 {
		t.Errorf("Gvr for achr[2] should be 0.3, but got %f", rooms[0].achr[2].Gvr)
	}
}

func TestAichschdlr_MultipleRooms(t *testing.T) {
	// Test: multiple rooms with different achr configurations
	val := []float64{0.5, 0.0}
	rooms := []*ROOM{
		{
			Nachr: 1,
			achr: []*ACHIR{
				{sch: 0},
			},
		},
		{
			Nachr: 2,
			achr: []*ACHIR{
				{sch: 0},
				{sch: 1},
			},
		},
		{
			Nachr: 0, // No inter-room airflow
			achr:  []*ACHIR{},
		},
	}

	Aichschdlr(val, rooms)

	if rooms[0].achr[0].Gvr != 0.5 {
		t.Errorf("Room[0].achr[0].Gvr should be 0.5, but got %f", rooms[0].achr[0].Gvr)
	}
	if rooms[1].achr[0].Gvr != 0.5 {
		t.Errorf("Room[1].achr[0].Gvr should be 0.5, but got %f", rooms[1].achr[0].Gvr)
	}
	if rooms[1].achr[1].Gvr != 0.0 {
		t.Errorf("Room[1].achr[1].Gvr should be 0.0, but got %f", rooms[1].achr[1].Gvr)
	}
}

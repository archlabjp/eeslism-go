package eeslism

import (
	"testing"
)

func TestResidata(t *testing.T) {
	// Setup
	fi := NewEeTokens("Living H=(2,sch1,sch2) comfrt=(sch3,sch4,sch5) ;\n*")
	Schdl := &SCHDL{
		Sch: []SCH{{name: "sch1"}, {name: "sch2"}, {name: "sch3"}, {name: "sch4"}, {name: "sch5"}},
		Val: []float64{0.8, 1.2, 1.0, 1.0, 0.1},
	}
	Room := []*ROOM{
		{Name: "Living"},
	}
	var pmvpri int
	Simc := &SIMCONTL{}

	// Execute
	Residata(fi, Schdl, Room, &pmvpri, Simc)

	// Verify
	if Room[0].Nhm != 2.0 {
		t.Errorf("Nhm should be 2.0, but got %f", Room[0].Nhm)
	}
	if *Room[0].Hmsch != 0.8 {
		t.Errorf("Hmsch should be 0.8, but got %f", *Room[0].Hmsch)
	}
	if *Room[0].Hmwksch != 1.2 {
		t.Errorf("Hmwksch should be 1.2, but got %f", *Room[0].Hmwksch)
	}
	if *Room[0].Metsch != 1.0 {
		t.Errorf("Metsch should be 1.0, but got %f", *Room[0].Metsch)
	}
	if *Room[0].Closch != 1.0 {
		t.Errorf("Closch should be 1.0, but got %f", *Room[0].Closch)
	}
	if *Room[0].Wvsch != 0.1 {
		t.Errorf("Wvsch should be 0.1, but got %f", *Room[0].Wvsch)
	}
	if pmvpri != 1 {
		t.Errorf("pmvpri should be 1, but got %d", pmvpri)
	}
}

func TestAppldata(t *testing.T) {
	// Setup
	fi := NewEeTokens("Living L=(100,x,sch1) As=(50,50,sch2) Al=(20,sch3) AE=(10,sch4) AG=(5,sch5) ;\n*")
	Schdl := &SCHDL{
		Sch: []SCH{{name: "sch1"}, {name: "sch2"}, {name: "sch3"}, {name: "sch4"}, {name: "sch5"}},
		Val: []float64{0.8, 0.6, 0.4, 0.2, 0.1},
	}
	Room := []*ROOM{
		{Name: "Living"},
	}
	Simc := &SIMCONTL{}

	// Execute
	Appldata(fi, Schdl, Room, Simc)

	// Verify
	if Room[0].Light != 100.0 {
		t.Errorf("Light should be 100.0, but got %f", Room[0].Light)
	}
	if *Room[0].Lightsch != 0.8 {
		t.Errorf("Lightsch should be 0.8, but got %f", *Room[0].Lightsch)
	}
	if Room[0].Apsc != 50.0 {
		t.Errorf("Apsc should be 50.0, but got %f", Room[0].Apsc)
	}
	if Room[0].Apsr != 50.0 {
		t.Errorf("Apsr should be 50.0, but got %f", Room[0].Apsr)
	}
	if *Room[0].Assch != 0.6 {
		t.Errorf("Assch should be 0.6, but got %f", *Room[0].Assch)
	}
	if Room[0].Apl != 20.0 {
		t.Errorf("Apl should be 20.0, but got %f", Room[0].Apl)
	}
	if *Room[0].Alsch != 0.4 {
		t.Errorf("Alsch should be 0.4, but got %f", *Room[0].Alsch)
	}
	if Room[0].AE != 10.0 {
		t.Errorf("AE should be 10.0, but got %f", Room[0].AE)
	}
	if *Room[0].AEsch != 0.2 {
		t.Errorf("AEsch should be 0.2, but got %f", *Room[0].AEsch)
	}
	if Room[0].AG != 5.0 {
		t.Errorf("AG should be 5.0, but got %f", Room[0].AG)
	}
	if *Room[0].AGsch != 0.1 {
		t.Errorf("AGsch should be 0.1, but got %f", *Room[0].AGsch)
	}
}

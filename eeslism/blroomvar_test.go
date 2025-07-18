
package eeslism

import (
	"testing"
)

func TestRoomelm(t *testing.T) {
	// Setup
	room1 := &ROOM{
		Nachr: 1,
		achr:  []*ACHIR{{room: &ROOM{cmp: &COMPNT{Elouts: []*ELOUT{{}, {}}}}}},
		Ntr:   1,
		trnx:  []*TRNX{{nextroom: &ROOM{cmp: &COMPNT{Elouts: []*ELOUT{{}}}}}},
		Nrp:   1,
		rmpnl: []*RPANEL{{
			pnl: &RDPNL{
				cmp: &COMPNT{
					Elins: []*ELIN{{Upo: &ELOUT{}}},
				},
			},
		}},
		cmp: &COMPNT{
			Elins:  []*ELIN{{}, {}, {}, {}},
			Elouts: []*ELOUT{{Ni: 2}},
		},
	}
	Room := []*ROOM{room1}
	_Rdpnl := []*RDPNL{}

	// Execute
	Roomelm(Room, _Rdpnl)

	// Verify (simple check)
	if Room[0].cmp.Elins[0].Upo == nil {
		t.Errorf("Roomelm failed: Upo should be assigned")
	}
}

func TestRoomvar(t *testing.T) {
	// Setup
	room1 := &ROOM{
		RMt:   10.0,
		RMC:   20.0,
		RMx:   30.0,
		RMXC:  40.0,
		Nachr: 1,
		achr:  []*ACHIR{{Gvr: 1.0}},
		Ntr:   1,
		ARN:   []float64{0.5},
		Nrp:   1,
		RMP:   []float64{0.2},
		Nasup: 1,
		cmp: &COMPNT{
			Elins: []*ELIN{
				{Lpath: &PLIST{G: 2.0}},
				{Lpath: &PLIST{G: 2.0}},
				{Lpath: &PLIST{G: 2.0}},
				{Lpath: &PLIST{G: 2.0}},
				{Lpath: &PLIST{G: 2.0}},
			},
			Elouts: []*ELOUT{{Coeffin: make([]float64, 4)}, {Coeffin: make([]float64, 2)}},
		},
	}
	Room := []*ROOM{room1}
	_Rdpnl := []*RDPNL{}

	// Execute
	Roomvar(Room, _Rdpnl)

	// Verify
	elout0 := Room[0].cmp.Elouts[0]
	if elout0.Coeffo == 0 {
		t.Errorf("Roomvar failed: Coeffo should not be zero")
	}
	elout1 := Room[0].cmp.Elouts[1]
	if elout1.Coeffo == 0 {
		t.Errorf("Roomvar failed: Coeffo should not be zero")
	}
}

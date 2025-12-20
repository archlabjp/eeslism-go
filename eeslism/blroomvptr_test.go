
package eeslism

import (
	"testing"
)

func TestRoomvptr(t *testing.T) {
	// Setup
	Room := &ROOM{
		Tr:   25.0,
		xr:   0.01,
		RH:   50.0,
		PMV:  0.5,
		Tsav: 24.0,
		Tot:  24.5,
		hr:   50.0,
		N:    1,
		rsrf: []*RMSRF{
			{
				Name:  "Wall1",
				Ts:    23.0,
				Tmrt:  23.5,
				Tcole: 22.0,
			},
		},
	}

	// Test cases
	testCases := []struct {
		name     string
		key      []string
		expected interface{}
	}{
		{"Tr", []string{"", "Tr"}, &Room.Tr},
		{"xr", []string{"", "xr"}, &Room.xr},
		{"RH", []string{"", "RH"}, &Room.RH},
		{"PMV", []string{"", "PMV"}, &Room.PMV},
		{"Tsav", []string{"", "Tsav"}, &Room.Tsav},
		{"Tot", []string{"", "Tot"}, &Room.Tot},
		{"hr", []string{"", "hr"}, &Room.hr},
		{"Ts", []string{"", "Wall1", "Ts"}, &Room.rsrf[0].Ts},
		{"Tmrt", []string{"", "Wall1", "Tmrt"}, &Room.rsrf[0].Tmrt},
		{"Te", []string{"", "Wall1", "Te"}, &Room.rsrf[0].Tcole},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vptr, err := roomvptr(len(tc.key), tc.key, Room)
			if err != nil {
				t.Errorf("roomvptr failed with error: %v", err)
			}
			if vptr.Ptr != tc.expected {
				t.Errorf("Expected pointer to %v, but got %v", tc.expected, vptr.Ptr)
			}
		})
	}
}

func TestRoomldptr(t *testing.T) {
	// Setup
	var load ControlSWType
	Room := &ROOM{
		rmld: &RMLOAD{},
		N:    1,
		rsrf: []*RMSRF{
			{
				Name: "Wall1",
			},
		},
	}

	// Test cases
	testCases := []struct {
		name    string
		key     []string
		idmrk   byte
		tropt   rune
		hmopt   rune
		exp_err bool
	}{
		{"Tr", []string{"", "Tr"}, 't', 'a', 0, false},
		{"Tot", []string{"", "Tot"}, 't', 'o', 0, false},
		{"RH", []string{"", "RH"}, 'x', 0, 'r', false},
		{"Tdp", []string{"", "Tdp"}, 'x', 0, 'd', false},
		{"xr", []string{"", "xr"}, 'x', 0, 'x', false},
		{"Ts", []string{"", "Wall1", "Ts"}, 't', 0, 0, false},
		{"Invalid", []string{"", "Invalid"}, 0, 0, 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var idmrk byte
			_, err := roomldptr(&load, tc.key, Room, &idmrk)

			if (err != nil) != tc.exp_err {
				t.Errorf("Expected error: %v, but got: %v", tc.exp_err, err)
			}
			if err == nil {
				if idmrk != tc.idmrk {
					t.Errorf("Expected idmrk: %c, but got: %c", tc.idmrk, idmrk)
				}
				if tc.tropt != 0 && Room.rmld.tropt != tc.tropt {
					t.Errorf("Expected tropt: %c, but got: %c", tc.tropt, Room.rmld.tropt)
				}
				if tc.hmopt != 0 && Room.rmld.hmopt != tc.hmopt {
					t.Errorf("Expected hmopt: %c, but got: %c", tc.hmopt, Room.rmld.hmopt)
				}
			}
		})
	}
}

func TestRoomldschd(t *testing.T) {
	t.Run("LoadtBranch_TsetAboveLimit", func(t *testing.T) {
		// Test loadt != nil with Tset > TEMPLIMIT
		loadt := ON_SW
		eo := &ELOUT{Control: ON_SW}
		eo.Eldobj = eo // Set Eldobj to itself (Eo == Eo.Eldobj)
		Room := &ROOM{
			rmld: &RMLOAD{
				loadt: &loadt,
				Tset:  25.0, // Above TEMPLIMIT (-100)
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{eo, {Control: ON_SW}},
			},
		}

		roomldschd(Room)

		if Room.cmp.Elouts[0].Control != LOAD_SW {
			t.Errorf("Expected Control=LOAD_SW, got %v", Room.cmp.Elouts[0].Control)
		}
		if Room.cmp.Elouts[0].Sysv != Room.rmld.Tset {
			t.Errorf("Expected Sysv=%f, got %f", Room.rmld.Tset, Room.cmp.Elouts[0].Sysv)
		}
		if Room.Tr != Room.rmld.Tset {
			t.Errorf("Expected Tr=%f, got %f", Room.rmld.Tset, Room.Tr)
		}
	})

	t.Run("LoadtBranch_TsetBelowLimit_WithVAV", func(t *testing.T) {
		// Test loadt != nil with Tset <= TEMPLIMIT and VAVcontrl set
		loadt := ON_SW
		eo := &ELOUT{Control: ON_SW}
		eo.Eldobj = eo
		vavEo := &ELOUT{Control: ON_SW}
		vavCmp := &COMPNT{Control: ON_SW, Elouts: []*ELOUT{vavEo}}
		vav := &VAV{Cmp: vavCmp}
		Room := &ROOM{
			rmld: &RMLOAD{
				loadt: &loadt,
				Tset:  -999.0, // Below TEMPLIMIT (-100)
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{eo, {Control: ON_SW}},
			},
			VAVcontrl: vav,
		}

		roomldschd(Room)

		// VAV should be turned off
		if Room.VAVcontrl.Cmp.Control != OFF_SW {
			t.Errorf("Expected VAVcontrl.Cmp.Control=OFF_SW, got %v", Room.VAVcontrl.Cmp.Control)
		}
		if Room.VAVcontrl.Cmp.Elouts[0].Control != OFF_SW {
			t.Errorf("Expected VAVcontrl Elouts[0].Control=OFF_SW, got %v", Room.VAVcontrl.Cmp.Elouts[0].Control)
		}
	})

	t.Run("LoadxBranch_RH_XsetPositive", func(t *testing.T) {
		// Test loadx != nil with hmopt='r' (relative humidity) and Xset > 0
		loadx := ON_SW
		eo1 := &ELOUT{Control: ON_SW}
		eo1.Eldobj = eo1
		Room := &ROOM{
			Tr: 25.0, // Room temperature for RH->x conversion
			rmld: &RMLOAD{
				loadx: &loadx,
				Xset:  50.0, // 50% RH
				hmopt: 'r',
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}, eo1},
			},
		}

		roomldschd(Room)

		if Room.cmp.Elouts[1].Control != LOAD_SW {
			t.Errorf("Expected Eo[1].Control=LOAD_SW, got %v", Room.cmp.Elouts[1].Control)
		}
		// Sysv should be calculated from FNXtr(25.0, 50.0)
		expectedX := FNXtr(25.0, 50.0)
		if Room.cmp.Elouts[1].Sysv != expectedX {
			t.Errorf("Expected Sysv=FNXtr(25,50)=%f, got %f", expectedX, Room.cmp.Elouts[1].Sysv)
		}
	})

	t.Run("LoadxBranch_Tdp_XsetAboveLimit", func(t *testing.T) {
		// Test loadx != nil with hmopt='d' (dew point) and Xset > TEMPLIMIT
		loadx := ON_SW
		eo1 := &ELOUT{Control: ON_SW}
		eo1.Eldobj = eo1
		Room := &ROOM{
			rmld: &RMLOAD{
				loadx: &loadx,
				Xset:  15.0, // Dew point temperature
				hmopt: 'd',
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}, eo1},
			},
		}

		roomldschd(Room)

		if Room.cmp.Elouts[1].Control != LOAD_SW {
			t.Errorf("Expected Eo[1].Control=LOAD_SW, got %v", Room.cmp.Elouts[1].Control)
		}
	})

	t.Run("LoadxBranch_xr_XsetPositive", func(t *testing.T) {
		// Test loadx != nil with hmopt='x' (absolute humidity) and Xset > 0
		loadx := ON_SW
		eo1 := &ELOUT{Control: ON_SW}
		eo1.Eldobj = eo1
		Room := &ROOM{
			rmld: &RMLOAD{
				loadx: &loadx,
				Xset:  0.010, // Absolute humidity [kg/kg]
				hmopt: 'x',
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}, eo1},
			},
		}

		roomldschd(Room)

		if Room.cmp.Elouts[1].Control != LOAD_SW {
			t.Errorf("Expected Eo[1].Control=LOAD_SW, got %v", Room.cmp.Elouts[1].Control)
		}
		if Room.cmp.Elouts[1].Sysv != Room.rmld.Xset {
			t.Errorf("Expected Sysv=%f, got %f", Room.rmld.Xset, Room.cmp.Elouts[1].Sysv)
		}
	})

	t.Run("NoLoadPointers_NilRmld", func(t *testing.T) {
		// Test with rmld = nil
		Room := &ROOM{
			rmld: nil,
			cmp: &COMPNT{
				Elouts: []*ELOUT{{Control: ON_SW}, {Control: ON_SW}},
			},
		}

		// Should not panic
		roomldschd(Room)

		// Controls should remain unchanged
		if Room.cmp.Elouts[0].Control != ON_SW {
			t.Errorf("Control should remain ON_SW, got %v", Room.cmp.Elouts[0].Control)
		}
	})

	t.Run("Eldobj_ControlOFF", func(t *testing.T) {
		// Test when Eldobj.Control == OFF_SW (should not set LOAD_SW)
		loadt := ON_SW
		eldobj := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{Control: ON_SW}
		eo.Eldobj = eldobj // Different Eldobj with OFF control
		Room := &ROOM{
			rmld: &RMLOAD{
				loadt: &loadt,
				Tset:  25.0,
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{eo, {Control: ON_SW}},
			},
		}

		roomldschd(Room)

		// Should NOT be set to LOAD_SW because Eldobj.Control == OFF_SW
		if Room.cmp.Elouts[0].Control != ON_SW {
			t.Errorf("Control should remain ON_SW when Eldobj.Control==OFF_SW, got %v", Room.cmp.Elouts[0].Control)
		}
	})
}

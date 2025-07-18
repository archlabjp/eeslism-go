
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

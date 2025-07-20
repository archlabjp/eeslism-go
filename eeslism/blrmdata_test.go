package eeslism

import (
	"testing"
)

func TestRoomDataStructures(t *testing.T) {
	// Test basic room data structures
	t.Run("ROOM structure", func(t *testing.T) {
		room := &ROOM{
			Name:  "TestRoom",
			N:     3,
			rsrf:  make([]*RMSRF, 3),
			Tr:    22.0,
			xr:    0.008,
			RH:    50.0,
			Ntr:   0,
			Nrp:   0,
			Nasup: 0,
		}

		if room.Name != "TestRoom" {
			t.Errorf("Room name = %s, want TestRoom", room.Name)
		}
		if room.N != 3 {
			t.Errorf("Room N = %d, want 3", room.N)
		}
		if len(room.rsrf) != 3 {
			t.Errorf("rsrf length = %d, want 3", len(room.rsrf))
		}
	})

	t.Run("RMSRF structure", func(t *testing.T) {
		surface := &RMSRF{
			typ:  RMSRFType_H,
			A:    15.0,
			ali:  8.0,
			alo:  25.0,
			alic: 3.0,
			alir: 5.0,
			Ts:   20.0,
			Te:   15.0,
		}

		if surface.typ != RMSRFType_H {
			t.Errorf("Surface typ = %c, want %c", surface.typ, RMSRFType_H)
		}
		if surface.A != 15.0 {
			t.Errorf("Surface A = %f, want 15.0", surface.A)
		}
	})

	t.Run("RMVLS structure", func(t *testing.T) {
		rmvls := &RMVLS{
			Room:  make([]*ROOM, 2),
			Sd:    make([]*RMSRF, 4),
			Mw:    make([]*MWALL, 4),
			Emrk:  make([]rune, 2),
			Snbk:  make([]*SNBK, 0),
			Qrm:   make([]*QRM, 0),
			Rdpnl: make([]*RDPNL, 0),
		}

		if len(rmvls.Room) != 2 {
			t.Errorf("RMVLS Room length = %d, want 2", len(rmvls.Room))
		}
		if len(rmvls.Sd) != 4 {
			t.Errorf("RMVLS Sd length = %d, want 4", len(rmvls.Sd))
		}
	})
}

func TestRoomDataValidation(t *testing.T) {
	// Test room data validation
	tests := []struct {
		name string
		room *ROOM
		valid bool
	}{
		{
			name: "Valid room",
			room: &ROOM{
				Name: "ValidRoom",
				N:    2,
				Tr:   22.0,
				xr:   0.008,
				RH:   50.0,
			},
			valid: true,
		},
		{
			name: "Invalid temperature",
			room: &ROOM{
				Name: "InvalidTempRoom",
				N:    2,
				Tr:   -60.0, // Invalid temperature (below -50)
				xr:   0.008,
				RH:   50.0,
			},
			valid: false,
		},
		{
			name: "Invalid humidity",
			room: &ROOM{
				Name: "InvalidHumidRoom",
				N:    2,
				Tr:   22.0,
				xr:   -0.001, // Invalid humidity
				RH:   50.0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			tempValid := tt.room.Tr >= -50.0 && tt.room.Tr <= 100.0
			humidValid := tt.room.xr >= 0.0 && tt.room.xr <= 0.030
			rhValid := tt.room.RH >= 0.0 && tt.room.RH <= 100.0

			isValid := tempValid && humidValid && rhValid

			if isValid != tt.valid {
				t.Errorf("Room validation = %t, want %t", isValid, tt.valid)
			}
		})
	}
}

func TestSurfaceDataValidation(t *testing.T) {
	// Test surface data validation
	tests := []struct {
		name    string
		surface *RMSRF
		valid   bool
	}{
		{
			name: "Valid surface",
			surface: &RMSRF{
				typ:  RMSRFType_H,
				A:    15.0,
				ali:  8.0,
				alic: 3.0,
				alir: 5.0,
			},
			valid: true,
		},
		{
			name: "Invalid area",
			surface: &RMSRF{
				typ:  RMSRFType_H,
				A:    -5.0, // Invalid area
				ali:  8.0,
				alic: 3.0,
				alir: 5.0,
			},
			valid: false,
		},
		{
			name: "Invalid heat transfer coefficient",
			surface: &RMSRF{
				typ:  RMSRFType_H,
				A:    15.0,
				ali:  -2.0, // Invalid coefficient
				alic: 3.0,
				alir: 5.0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			areaValid := tt.surface.A > 0.0
			aliValid := tt.surface.ali >= 0.0
			alicValid := tt.surface.alic >= 0.0
			alirValid := tt.surface.alir >= 0.0

			isValid := areaValid && aliValid && alicValid && alirValid

			if isValid != tt.valid {
				t.Errorf("Surface validation = %t, want %t", isValid, tt.valid)
			}
		})
	}
}

func TestRoomDataConsistency(t *testing.T) {
	// Test consistency between room and surface data
	room := &ROOM{
		Name: "ConsistencyTest",
		N:    2,
		rsrf: make([]*RMSRF, 2),
	}

	// Create surfaces
	for i := 0; i < 2; i++ {
		room.rsrf[i] = &RMSRF{
			typ:  RMSRFType_H,
			A:    10.0 + float64(i)*5.0,
			ali:  8.0,
			alic: 3.0,
			alir: 5.0,
			room: room, // Back reference to room
		}
	}

	// Verify consistency
	if len(room.rsrf) != room.N {
		t.Errorf("Surface count (%d) doesn't match room.N (%d)", len(room.rsrf), room.N)
	}

	for i, surface := range room.rsrf {
		if surface.room != room {
			t.Errorf("Surface %d doesn't have correct room reference", i)
		}
		if surface.A <= 0.0 {
			t.Errorf("Surface %d has invalid area: %f", i, surface.A)
		}
	}
}
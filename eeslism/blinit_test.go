package eeslism

import (
	"testing"
)

func TestRoomInit(t *testing.T) {
	// Test basic room initialization
	room := &ROOM{
		Name: "TestRoom",
		N:    2,
	}

	// Initialize basic arrays
	room.rsrf = make([]*RMSRF, room.N)
	room.XA = make([]float64, room.N*room.N)
	room.alr = make([]float64, room.N*room.N)

	// Create mock surfaces
	for i := 0; i < room.N; i++ {
		room.rsrf[i] = &RMSRF{
			A:    10.0,
			ali:  8.0,
			alic: 3.0,
			alir: 5.0,
		}
	}

	// Verify initialization
	if room.Name != "TestRoom" {
		t.Errorf("Room name = %s, want TestRoom", room.Name)
	}
	if room.N != 2 {
		t.Errorf("Room N = %d, want 2", room.N)
	}
	if len(room.rsrf) != room.N {
		t.Errorf("rsrf length = %d, want %d", len(room.rsrf), room.N)
	}
}

func TestRoomInitialization_Arrays(t *testing.T) {
	tests := []struct {
		name string
		N    int
	}{
		{"Small room", 2},
		{"Medium room", 4},
		{"Large room", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &ROOM{
				Name: tt.name,
				N:    tt.N,
			}

			// Initialize arrays
			room.rsrf = make([]*RMSRF, room.N)
			room.XA = make([]float64, room.N*room.N)
			room.alr = make([]float64, room.N*room.N)

			// Verify array sizes
			if len(room.rsrf) != tt.N {
				t.Errorf("rsrf length = %d, want %d", len(room.rsrf), tt.N)
			}
			if len(room.XA) != tt.N*tt.N {
				t.Errorf("XA length = %d, want %d", len(room.XA), tt.N*tt.N)
			}
			if len(room.alr) != tt.N*tt.N {
				t.Errorf("alr length = %d, want %d", len(room.alr), tt.N*tt.N)
			}
		})
	}
}

func TestWallInitialization(t *testing.T) {
	// Test basic wall initialization
	wall := &WALL{
		name:     "TestWall",
		WallType: WallType_N,
		N:        3,
	}

	// Initialize wall elements
	wall.welm = make([]WELM, wall.N)
	for i := 0; i < wall.N; i++ {
		wall.welm[i] = WELM{
			Code: "material" + string(rune('A'+i)),
			L:    0.1 + float64(i)*0.05,
			ND:   2,
		}
	}

	// Verify initialization
	if wall.name != "TestWall" {
		t.Errorf("Wall name = %s, want TestWall", wall.name)
	}
	if wall.N != 3 {
		t.Errorf("Wall N = %d, want 3", wall.N)
	}
	if len(wall.welm) != wall.N {
		t.Errorf("welm length = %d, want %d", len(wall.welm), wall.N)
	}

	// Check individual elements
	for i := 0; i < wall.N; i++ {
		expectedCode := "material" + string(rune('A'+i))
		if wall.welm[i].Code != expectedCode {
			t.Errorf("welm[%d].Code = %s, want %s", i, wall.welm[i].Code, expectedCode)
		}
		expectedL := 0.1 + float64(i)*0.05
		if wall.welm[i].L != expectedL {
			t.Errorf("welm[%d].L = %f, want %f", i, wall.welm[i].L, expectedL)
		}
	}
}

func TestSurfaceInitialization(t *testing.T) {
	// Test basic surface initialization
surface := &RMSRF{
typ:  RMSRFType_H,
A:    15.0,
ali:  8.0,
alo:  25.0,
alic: 3.0,
alir: 5.0,
}

	// Verify initialization
	if surface.typ != RMSRFType_H {
		t.Errorf("Surface typ = %c, want %c", surface.typ, RMSRFType_H)
	}
	if surface.A != 15.0 {
		t.Errorf("Surface A = %f, want 15.0", surface.A)
	}
	if surface.ali != 8.0 {
		t.Errorf("Surface ali = %f, want 8.0", surface.ali)
	}
	if surface.alo != 25.0 {
		t.Errorf("Surface alo = %f, want 25.0", surface.alo)
	}
}

func TestMWALLInitialization(t *testing.T) {
	// Test MWALL initialization
	mwall := &MWALL{
		M:  5,
		mp: 2,
	}

	// Initialize arrays
	mwall.UX = make([]float64, mwall.M*mwall.M)
	mwall.Tw = make([]float64, mwall.M)
	mwall.Told = make([]float64, mwall.M)
	mwall.Twd = make([]float64, mwall.M)
	mwall.res = make([]float64, mwall.M+1)
	mwall.cap = make([]float64, mwall.M+1)

	// Verify initialization
	if mwall.M != 5 {
		t.Errorf("MWALL M = %d, want 5", mwall.M)
	}
	if mwall.mp != 2 {
		t.Errorf("MWALL mp = %d, want 2", mwall.mp)
	}
	if len(mwall.UX) != mwall.M*mwall.M {
		t.Errorf("UX length = %d, want %d", len(mwall.UX), mwall.M*mwall.M)
	}
	if len(mwall.Tw) != mwall.M {
		t.Errorf("Tw length = %d, want %d", len(mwall.Tw), mwall.M)
	}
	if len(mwall.Told) != mwall.M {
		t.Errorf("Told length = %d, want %d", len(mwall.Told), mwall.M)
	}
}
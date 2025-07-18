package eeslism

import (
	"testing"
)

// TestWALL tests the WALL structure and related functions
func TestWALL(t *testing.T) {
	t.Run("WALL creation", func(t *testing.T) {
		wall := NewWall()
		
		if wall == nil {
			t.Fatal("NewWall() returned nil")
		}
		
		// Check default values
		if wall.name != "" {
			t.Errorf("Expected empty name, got %s", wall.name)
		}
		if wall.N != 0 {
			t.Errorf("Expected N=0, got %d", wall.N)
		}
		if wall.M != 0 {
			t.Errorf("Expected M=0, got %d", wall.M)
		}
		if wall.Ei != 0.9 {
			t.Errorf("Expected Ei=0.9, got %f", wall.Ei)
		}
		if wall.Eo != 0.9 {
			t.Errorf("Expected Eo=0.9, got %f", wall.Eo)
		}
		if wall.as != 0.7 {
			t.Errorf("Expected as=0.7, got %f", wall.as)
		}
		if wall.WallType != WallType_N {
			t.Errorf("Expected WallType_N, got %c", wall.WallType)
		}
	})

	t.Run("WALL with layers", func(t *testing.T) {
		wall := NewWall()
		wall.name = "TestWall"
		wall.ble = BLE_ExternalWall
		wall.N = 3
		wall.welm = []WELM{
			{Code: "CONC", L: 0.15, ND: 3, Cond: 1.6, Cro: 1900000},
			{Code: "INS", L: 0.05, ND: 1, Cond: 0.04, Cro: 30000},
			{Code: "GYPS", L: 0.012, ND: 1, Cond: 0.22, Cro: 830000},
		}
		
		if wall.name != "TestWall" {
			t.Errorf("Expected name 'TestWall', got %s", wall.name)
		}
		if wall.ble != BLE_ExternalWall {
			t.Errorf("Expected BLE_ExternalWall, got %c", wall.ble)
		}
		if wall.N != 3 {
			t.Errorf("Expected N=3, got %d", wall.N)
		}
		if len(wall.welm) != 3 {
			t.Errorf("Expected 3 layers, got %d", len(wall.welm))
		}
	})
}

// TestRMSRF tests the RMSRF structure
func TestRMSRF(t *testing.T) {
	t.Run("RMSRF creation", func(t *testing.T) {
		room := &ROOM{Name: "TestRoom"}
		
		rmsrf := &RMSRF{
			Name:  "TestSurface",
			ble:   BLE_ExternalWall,
			typ:   RMSRFType_H,
			mwtype: RMSRFMwType_I,
			room:  room,
			A:     10.0,
			Eo:    0.9,
			as:    0.7,
			ali:   8.0,
			alo:   25.0,
		}
		
		if rmsrf.Name != "TestSurface" {
			t.Errorf("Expected name 'TestSurface', got %s", rmsrf.Name)
		}
		if rmsrf.ble != BLE_ExternalWall {
			t.Errorf("Expected BLE_ExternalWall, got %c", rmsrf.ble)
		}
		if rmsrf.typ != RMSRFType_H {
			t.Errorf("Expected RMSRFType_H, got %c", rmsrf.typ)
		}
		if rmsrf.mwtype != RMSRFMwType_I {
			t.Errorf("Expected RMSRFMwType_I, got %c", rmsrf.mwtype)
		}
		if rmsrf.room != room {
			t.Error("Room reference not set correctly")
		}
		if rmsrf.A != 10.0 {
			t.Errorf("Expected A=10.0, got %f", rmsrf.A)
		}
	})
}

// TestWINDOW tests the WINDOW structure
func TestWINDOW(t *testing.T) {
	t.Run("WINDOW creation", func(t *testing.T) {
		window := NewWINDOW()
		
		if window == nil {
			t.Fatal("NewWINDOW() returned nil")
		}
		
		// Check default values
		if window.Name != "" {
			t.Errorf("Expected empty name, got %s", window.Name)
		}
		if window.Cidtype != "N" {
			t.Errorf("Expected Cidtype='N', got %s", window.Cidtype)
		}
		if window.Ei != 0.9 {
			t.Errorf("Expected Ei=0.9, got %f", window.Ei)
		}
		if window.Eo != 0.9 {
			t.Errorf("Expected Eo=0.9, got %f", window.Eo)
		}
		if window.RStrans != false {
			t.Errorf("Expected RStrans=false, got %t", window.RStrans)
		}
	})

	t.Run("WINDOW with properties", func(t *testing.T) {
		window := NewWINDOW()
		window.Name = "TestWindow"
		window.K = 3.0
		window.Rwall = 0.33
		window.tgtn = 0.8
		window.Bn = 0.1
		
		if window.Name != "TestWindow" {
			t.Errorf("Expected name 'TestWindow', got %s", window.Name)
		}
		if window.K != 3.0 {
			t.Errorf("Expected K=3.0, got %f", window.K)
		}
		if window.Rwall != 0.33 {
			t.Errorf("Expected Rwall=0.33, got %f", window.Rwall)
		}
		if window.tgtn != 0.8 {
			t.Errorf("Expected tgtn=0.8, got %f", window.tgtn)
		}
		if window.Bn != 0.1 {
			t.Errorf("Expected Bn=0.1, got %f", window.Bn)
		}
	})
}
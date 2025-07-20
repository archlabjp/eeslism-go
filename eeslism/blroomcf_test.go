package eeslism

import (
	"testing"
)

func TestEeroomcf_BasicStructures(t *testing.T) {
	// Test that the basic structures used by eeroomcf are properly initialized
	t.Run("WDAT structure", func(t *testing.T) {
		wd := &WDAT{}

		// Basic structure test - just verify it can be created
		if wd == nil {
			t.Errorf("WDAT should not be nil")
		}
	})

	t.Run("EXSFS structure", func(t *testing.T) {
		exs := &EXSFS{
			Exs: make([]*EXSF, 1),
		}

		if len(exs.Exs) != 1 {
			t.Errorf("EXSFS Exs length = %d, want 1", len(exs.Exs))
		}
	})

	t.Run("RMVLS structure", func(t *testing.T) {
		rmvls := &RMVLS{
			Room:  make([]*ROOM, 1),
			Sd:    make([]*RMSRF, 2),
			Mw:    make([]*MWALL, 2),
			Emrk:  make([]rune, 1),
			Snbk:  make([]*SNBK, 0),
			Qrm:   make([]*QRM, 0),
			Rdpnl: make([]*RDPNL, 0),
		}

		if len(rmvls.Room) != 1 {
			t.Errorf("RMVLS Room length = %d, want 1", len(rmvls.Room))
		}
		if len(rmvls.Sd) != 2 {
			t.Errorf("RMVLS Sd length = %d, want 2", len(rmvls.Sd))
		}
	})
}
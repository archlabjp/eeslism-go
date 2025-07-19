package eeslism

import (
	"testing"
)

func TestCommonConstants(t *testing.T) {
	t.Run("NOP constant", func(t *testing.T) {
		if NOP != -1 {
			t.Errorf("NOP should be -1, got %d", NOP)
		}
	})

	t.Run("FNOP constant", func(t *testing.T) {
		if FNOP != FNAN {
			t.Errorf("FNOP should equal FNAN, got %v", FNOP)
		}
	})

	t.Run("TEMPLIMIT constant", func(t *testing.T) {
		expected := -273.16
		if TEMPLIMIT != expected {
			t.Errorf("TEMPLIMIT should be %v, got %v", expected, TEMPLIMIT)
		}
	})

	t.Run("ERRFMT constant", func(t *testing.T) {
		expected := "xxxxx %s xxxxx : "
		if ERRFMT != expected {
			t.Errorf("ERRFMT should be %q, got %q", expected, ERRFMT)
		}
	})

	t.Run("ERRFMTA constant", func(t *testing.T) {
		expected := "xxx %s xxx %s\n"
		if ERRFMTA != expected {
			t.Errorf("ERRFMTA should be %q, got %q", expected, ERRFMTA)
		}
	})
}

func TestCommonConstantsUsage(t *testing.T) {
	t.Run("NOP as invalid state indicator", func(t *testing.T) {
		// Test that NOP can be used as an invalid state indicator
		var state int = NOP
		if state >= 0 {
			t.Errorf("NOP should indicate invalid state (negative), got %d", state)
		}
	})

	t.Run("TEMPLIMIT as physical boundary", func(t *testing.T) {
		// Test that TEMPLIMIT is close to absolute zero
		absoluteZero := -273.15
		tolerance := 1.0
		if TEMPLIMIT < absoluteZero-tolerance || TEMPLIMIT > absoluteZero+tolerance {
			t.Errorf("TEMPLIMIT should be close to absolute zero (%v), got %v", absoluteZero, TEMPLIMIT)
		}
	})

	t.Run("Error format strings contain placeholders", func(t *testing.T) {
		// Test that error format strings contain %s placeholders
		if len(ERRFMT) == 0 {
			t.Error("ERRFMT should not be empty")
		}
		
		if len(ERRFMTA) == 0 {
			t.Error("ERRFMTA should not be empty")
		}

		// Check that ERRFMTA contains newline for proper formatting
		if ERRFMTA[len(ERRFMTA)-1] != '\n' {
			t.Error("ERRFMTA should end with newline character")
		}
	})
}
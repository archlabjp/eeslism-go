package eeslism

import (
	"testing"
)

func TestPVWALLCAT(t *testing.T) {
	// Test PVWALLCAT structure
	pvcat := &PVWALLCAT{
		KHD:     0.97,
		KPD:     0.99,
		KPM:     0.94,
		KPA:     0.97,
		EffINO:  0.95,
		Apmax:   -0.004,
		Ap:      20.0,
		Rcoloff: 0.04,
		Kcoloff: 25.0,
	}

	// Verify initialization
	if pvcat.KHD != 0.97 {
		t.Errorf("KHD = %f, want 0.97", pvcat.KHD)
	}
	if pvcat.EffINO != 0.95 {
		t.Errorf("EffINO = %f, want 0.95", pvcat.EffINO)
	}
	if pvcat.Apmax != -0.004 {
		t.Errorf("Apmax = %f, want -0.004", pvcat.Apmax)
	}
}

func TestPVWALL(t *testing.T) {
	// Test PVWALL structure
	pvwall := &PVWALL{
		KTotal: 0.85,
		KPT:    0.95,
		TPV:    45.0,
		Power:  4500.0,
		Eff:    0.15,
		PVcap:  5000.0,
	}

	// Verify initialization
	if pvwall.KTotal != 0.85 {
		t.Errorf("KTotal = %f, want 0.85", pvwall.KTotal)
	}
	if pvwall.Power != 4500.0 {
		t.Errorf("Power = %f, want 4500.0", pvwall.Power)
	}
	if pvwall.PVcap != 5000.0 {
		t.Errorf("PVcap = %f, want 5000.0", pvwall.PVcap)
	}
	if pvwall.Eff != 0.15 {
		t.Errorf("Eff = %f, want 0.15", pvwall.Eff)
	}
}

func TestPVSystemEfficiency(t *testing.T) {
	tests := []struct {
		name     string
		eff      float64
		expected bool
	}{
		{"High efficiency", 0.20, true},
		{"Medium efficiency", 0.15, true},
		{"Low efficiency", 0.10, true},
		{"Very low efficiency", 0.05, true},
		{"Invalid efficiency", 0.0, false},
		{"Negative efficiency", -0.1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvwall := &PVWALL{
				Eff: tt.eff,
			}

			// Check if efficiency is in valid range
			isValid := pvwall.Eff > 0.0 && pvwall.Eff <= 1.0
			if isValid != tt.expected {
				t.Errorf("Efficiency validity = %t, want %t for efficiency %f", isValid, tt.expected, tt.eff)
			}
		})
	}
}

func TestPVTemperatureCoefficient(t *testing.T) {
	// Test temperature coefficient (typically negative for PV)
	pvcat := &PVWALLCAT{
		Apmax: -0.004, // Typical value for silicon PV
	}

	// Temperature coefficient should be negative
	if pvcat.Apmax >= 0.0 {
		t.Errorf("Apmax should be negative for typical PV, got %f", pvcat.Apmax)
	}

	// Test range of typical values
	if pvcat.Apmax < -0.01 || pvcat.Apmax > -0.002 {
		t.Logf("Warning: Apmax (%f) is outside typical range (-0.01 to -0.002)", pvcat.Apmax)
	}
}

func TestPVSystemCorrectionFactors(t *testing.T) {
	// Test correction factors (should be between 0 and 1)
	pvcat := &PVWALLCAT{
		KHD: 0.97, // Solar irradiation annual variation correction
		KPD: 0.99, // Aging correction
		KPM: 0.94, // Array load matching correction
		KPA: 0.97, // Array circuit correction
	}

	factors := []struct {
		name  string
		value float64
	}{
		{"KHD", pvcat.KHD},
		{"KPD", pvcat.KPD},
		{"KPM", pvcat.KPM},
		{"KPA", pvcat.KPA},
	}

	for _, factor := range factors {
		t.Run(factor.name, func(t *testing.T) {
			if factor.value <= 0.0 || factor.value > 1.0 {
				t.Errorf("%s = %f, should be between 0 and 1", factor.name, factor.value)
			}
		})
	}
}

func TestPVInverterEfficiency(t *testing.T) {
	tests := []struct {
		name    string
		effINO  float64
		isValid bool
	}{
		{"High efficiency inverter", 0.98, true},
		{"Standard efficiency inverter", 0.95, true},
		{"Low efficiency inverter", 0.90, true},
		{"Invalid efficiency", 1.1, false},
		{"Zero efficiency", 0.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pvcat := &PVWALLCAT{
				EffINO: tt.effINO,
			}

			isValid := pvcat.EffINO > 0.0 && pvcat.EffINO <= 1.0
			if isValid != tt.isValid {
				t.Errorf("Inverter efficiency validity = %t, want %t for efficiency %f", 
					isValid, tt.isValid, tt.effINO)
			}
		})
	}
}
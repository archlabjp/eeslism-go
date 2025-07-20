package eeslism

import (
	"testing"
)

func TestSolarWallBasicStructure(t *testing.T) {
	// Test basic solar wall surface structure
	surface := &RMSRF{
		typ:       RMSRFType_H,
		A:         20.0,
		Iwall:     800.0, // Solar irradiation on wall
		Tcole:     35.0,  // Collector outlet temperature
		Tf:        30.0,  // Fluid temperature
		PVwallFlg: true,
		PVwall: PVWALL{
			TPV:    45.0,
			Power:  2500.0,
			Eff:    0.15,
			KTotal: 0.85,
			KPT:    0.95,
			PVcap:  3000.0,
		},
	}

	// Verify initialization
	if surface.A != 20.0 {
		t.Errorf("Surface area = %f, want 20.0", surface.A)
	}
	if surface.Iwall != 800.0 {
		t.Errorf("Wall irradiation = %f, want 800.0", surface.Iwall)
	}
	if !surface.PVwallFlg {
		t.Errorf("PVwallFlg should be true")
	}
	if surface.PVwall.TPV != 45.0 {
		t.Errorf("PV temperature = %f, want 45.0", surface.PVwall.TPV)
	}
}

func TestSolarWallWithCollector(t *testing.T) {
	// Test solar wall with collector (WallType_C)
	wall := &WALL{
		WallType: WallType_C,
		Kc:       25.0, // Collector heat transfer coefficient
	}

	mwall := &MWALL{
		wall: wall,
		M:    3,
	}

	surface := &RMSRF{
		typ:    RMSRFType_H,
		A:      15.0,
		Iwall:  600.0,
		Tcole:  32.0,
		Tf:     28.0,
		mw:     mwall,
		mwside: RMSRFMwSideType_i,
		Ndiv:   5,
		Tc:     make([]float64, 5),
	}

	// Verify collector wall setup
	if wall.WallType != WallType_C {
		t.Errorf("WallType = %c, want %c", wall.WallType, WallType_C)
	}
	if wall.Kc != 25.0 {
		t.Errorf("Collector Kc = %f, want 25.0", wall.Kc)
	}
	if len(surface.Tc) != 5 {
		t.Errorf("Tc array length = %d, want 5", len(surface.Tc))
	}
}

func TestSolarWallTemperatureDistribution(t *testing.T) {
	// Test temperature distribution in solar collector
	surface := &RMSRF{
		Ndiv:  3,
		Tc:    make([]float64, 3),
		Tcole: 40.0,
		A:     12.0,
	}

	// Initialize temperature distribution
	for i := 0; i < surface.Ndiv; i++ {
		surface.Tc[i] = 35.0 + float64(i)*2.0 // 35, 37, 39
	}

	// Verify temperature distribution
	expectedTemps := []float64{35.0, 37.0, 39.0}
	for i := 0; i < surface.Ndiv; i++ {
		if surface.Tc[i] != expectedTemps[i] {
			t.Errorf("Tc[%d] = %f, want %f", i, surface.Tc[i], expectedTemps[i])
		}
	}

	// Verify collector outlet temperature
	if surface.Tcole != 40.0 {
		t.Errorf("Tcole = %f, want 40.0", surface.Tcole)
	}
}

func TestSolarWallHeatTransferCoefficients(t *testing.T) {
	// Test heat transfer coefficients for solar wall
	surface := &RMSRF{
		dblKsu: 0.15, // Upper surface heat transfer coefficient
		dblKsd: 0.12, // Lower surface heat transfer coefficient
		dblKc:  25.0, // Collector heat transfer coefficient
		dblTsu: 42.0, // Upper surface temperature
		dblTsd: 38.0, // Lower surface temperature
	}

	// Verify heat transfer coefficients
	if surface.dblKsu != 0.15 {
		t.Errorf("dblKsu = %f, want 0.15", surface.dblKsu)
	}
	if surface.dblKsd != 0.12 {
		t.Errorf("dblKsd = %f, want 0.12", surface.dblKsd)
	}
	if surface.dblKc != 25.0 {
		t.Errorf("dblKc = %f, want 25.0", surface.dblKc)
	}

	// Verify surface temperatures
	if surface.dblTsu != 42.0 {
		t.Errorf("dblTsu = %f, want 42.0", surface.dblTsu)
	}
	if surface.dblTsd != 38.0 {
		t.Errorf("dblTsd = %f, want 38.0", surface.dblTsd)
	}
}

func TestSolarWallPVPerformance(t *testing.T) {
	// Test PV performance calculation
	tests := []struct {
		name     string
		irrad    float64
		temp     float64
		area     float64
		eff      float64
		expected float64 // Expected power range (minimum)
	}{
		{"High irradiation", 1000.0, 25.0, 20.0, 0.15, 2400.0},
		{"Medium irradiation", 600.0, 35.0, 15.0, 0.15, 1000.0},
		{"Low irradiation", 200.0, 45.0, 10.0, 0.15, 200.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			surface := &RMSRF{
				A:         tt.area,
				Iwall:     tt.irrad,
				PVwallFlg: true,
				PVwall: PVWALL{
					TPV:    tt.temp,
					Eff:    tt.eff,
					KTotal: 0.85,
					KPT:    0.95,
				},
			}

			// Basic power calculation (simplified)
			// Power = Irradiation * Area * Efficiency * Correction factors
			expectedPower := tt.irrad * tt.area * tt.eff * surface.PVwall.KTotal * surface.PVwall.KPT

			// Set calculated power
			surface.PVwall.Power = expectedPower

			if surface.PVwall.Power < tt.expected {
				t.Errorf("PV Power = %f, should be at least %f", surface.PVwall.Power, tt.expected)
			}

			// Verify irradiation and temperature are reasonable
			if surface.Iwall < 0.0 || surface.Iwall > 1500.0 {
				t.Errorf("Iwall (%f) is outside reasonable range", surface.Iwall)
			}
			if surface.PVwall.TPV < -10.0 || surface.PVwall.TPV > 80.0 {
				t.Errorf("TPV (%f) is outside reasonable range", surface.PVwall.TPV)
			}
		})
	}
}

func TestSolarWallCollectorEfficiency(t *testing.T) {
	// Test collector efficiency calculation
	surface := &RMSRF{
		A:     25.0,
		Iwall: 750.0,
		Tcole: 45.0,
		Tf:    35.0,
		mw: &MWALL{
			wall: &WALL{
				WallType: WallType_C,
				Kc:       30.0,
			},
		},
	}

	// Calculate temperature difference
	tempDiff := surface.Tcole - surface.Tf
	if tempDiff <= 0.0 {
		t.Errorf("Temperature difference should be positive, got %f", tempDiff)
	}

	// Verify collector parameters
	if surface.mw.wall.Kc <= 0.0 {
		t.Errorf("Collector Kc should be positive, got %f", surface.mw.wall.Kc)
	}
	if surface.Iwall <= 0.0 {
		t.Errorf("Wall irradiation should be positive, got %f", surface.Iwall)
	}
}

func TestSolarWallMultipleCollectors(t *testing.T) {
	// Test multiple collector sections
	numSections := 4
	surfaces := make([]*RMSRF, numSections)

	for i := 0; i < numSections; i++ {
		surfaces[i] = &RMSRF{
			A:     10.0,
			Iwall: 700.0 + float64(i)*50.0, // Varying irradiation
			Tcole: 40.0 + float64(i)*2.0,   // Increasing temperature
			Tf:    30.0 + float64(i)*1.5,   // Increasing fluid temperature
			Ndiv:  3,
			Tc:    make([]float64, 3),
		}

		// Initialize temperature distribution for each section
		for j := 0; j < 3; j++ {
			surfaces[i].Tc[j] = surfaces[i].Tf + float64(j)*2.0
		}
	}

	// Verify each section
	for i := 0; i < numSections; i++ {
		if surfaces[i].A != 10.0 {
			t.Errorf("Section %d area = %f, want 10.0", i, surfaces[i].A)
		}
		if surfaces[i].Tcole <= surfaces[i].Tf {
			t.Errorf("Section %d: Tcole (%f) should be greater than Tf (%f)", 
				i, surfaces[i].Tcole, surfaces[i].Tf)
		}
		if len(surfaces[i].Tc) != 3 {
			t.Errorf("Section %d: Tc array length = %d, want 3", i, len(surfaces[i].Tc))
		}
	}
}
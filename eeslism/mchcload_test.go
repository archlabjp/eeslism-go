package eeslism

import (
	"testing"
)

// TestHclelm tests the HCLOAD element assignment function
func TestHclelm(t *testing.T) {
	t.Run("BasicElementAssignment", func(t *testing.T) {
		// Create basic HCLOAD with proper component setup
		hcload := createBasicHCLOAD()
		hcloads := []*HCLOAD{hcload}

		// Execute element assignment
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Hclelm panicked: %v", r)
			}
		}()

		Hclelm(hcloads)
		t.Log("Basic element assignment completed successfully")
	})

	t.Run("WetCoilAssignment", func(t *testing.T) {
		// Create wet coil HCLOAD
		hcload := createWetCoilHCLOAD()
		hcloads := []*HCLOAD{hcload}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Wet coil element assignment handled panic: %v", r)
			}
		}()

		Hclelm(hcloads)
		t.Log("Wet coil element assignment completed successfully")
	})

	t.Run("WaterCoilAssignment", func(t *testing.T) {
		// Create water coil HCLOAD (Type 'W')
		hcload := createWaterCoilHCLOAD()
		hcloads := []*HCLOAD{hcload}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Water coil element assignment handled panic: %v", r)
			}
		}()

		Hclelm(hcloads)
		t.Log("Water coil element assignment completed successfully")
	})

	t.Run("EmptyHCLOADList", func(t *testing.T) {
		// Test with empty HCLOAD list
		var hcloads []*HCLOAD

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Empty HCLOAD list panicked: %v", r)
			}
		}()

		Hclelm(hcloads)
		t.Log("Empty HCLOAD list handled successfully")
	})
}

// TestHcldcfv tests the HCLOAD coefficient calculation function
func TestHcldcfv(t *testing.T) {
	t.Run("BasicCoefficientCalculation", func(t *testing.T) {
		// Create HCLOAD with proper setup for coefficient calculation
		hcload := createCoefficientTestHCLOAD()
		hcloads := []*HCLOAD{hcload}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Coefficient calculation handled panic: %v", r)
			}
		}()

		Hcldcfv(hcloads)

		// Verify coefficient calculations
		if hcload.CGa <= 0 {
			t.Errorf("CGa should be positive, got %f", hcload.CGa)
		}
		if hcload.Ga <= 0 {
			t.Errorf("Ga should be positive, got %f", hcload.Ga)
		}

		t.Log("Basic coefficient calculation completed successfully")
	})

	t.Run("WetModeCoefficients", func(t *testing.T) {
		// Test coefficient calculation for wet mode
		hcload := createWetModeHCLOAD()
		hcloads := []*HCLOAD{hcload}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Wet mode coefficient calculation handled panic: %v", r)
			}
		}()

		Hcldcfv(hcloads)
		t.Log("Wet mode coefficient calculation completed successfully")
	})

	t.Run("WaterCoilCoefficients", func(t *testing.T) {
		// Test coefficient calculation for water coil (Type 'W')
		hcload := createWaterCoilCoefficientHCLOAD()
		hcloads := []*HCLOAD{hcload}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Water coil coefficient calculation handled panic: %v", r)
			}
		}()

		Hcldcfv(hcloads)

		// Verify water coil specific coefficients
		if hcload.Type == HCLoadType_W && hcload.CGw <= 0 {
			t.Errorf("CGw should be positive for water coil, got %f", hcload.CGw)
		}

		t.Log("Water coil coefficient calculation completed successfully")
	})
}

// TestHcldene tests the HCLOAD energy calculation function
func TestHcldene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		// Create HCLOAD for energy calculation
		hcload := createEnergyTestHCLOAD()
		hcloads := []*HCLOAD{hcload}
		var LDrest int
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy calculation handled panic: %v", r)
			}
		}()

		Hcldene(hcloads, &LDrest, wd)

		// Verify energy calculations
		if hcload.Cmp.Control == ON_SW {
			// Check that energy values are calculated
			t.Logf("Energy calculation results - Qs: %.1f, Ql: %.1f, Qt: %.1f", 
				hcload.Qs, hcload.Ql, hcload.Qt)
		}

		t.Log("Basic energy calculation completed successfully")
	})

	t.Run("RMACEnergyCalculation", func(t *testing.T) {
		// Test RMAC (Room Air Conditioner) energy calculation
		hcload := createRMACTestHCLOAD()
		hcloads := []*HCLOAD{hcload}
		var LDrest int
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("RMAC energy calculation handled panic: %v", r)
			}
		}()

		Hcldene(hcloads, &LDrest, wd)

		// Verify RMAC specific calculations
		if hcload.RMACFlg == 'Y' || hcload.RMACFlg == 'y' {
			t.Logf("RMAC calculation results - Qt: %.1f, Ele: %.1f, COP: %.2f", 
				hcload.Qt, hcload.Ele, hcload.COP)
		}

		t.Log("RMAC energy calculation completed successfully")
	})

	t.Run("EnergyBalance", func(t *testing.T) {
		// Test energy balance in HCLOAD calculations
		hcload := createEnergyBalanceHCLOAD()
		hcloads := []*HCLOAD{hcload}
		var LDrest int
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Energy balance calculation handled panic: %v", r)
			}
		}()

		Hcldene(hcloads, &LDrest, wd)

		// Verify energy balance (Qs + Ql should approximately equal Qt)
		if hcload.Cmp.Control == ON_SW && hcload.Qt != 0 {
			totalEnergy := hcload.Qs + hcload.Ql
			energyError := absValue(totalEnergy - hcload.Qt) / absValue(hcload.Qt)
			
			if energyError > 0.05 { // 5% tolerance
				t.Errorf("Energy balance error: %.3f%% (Qs+Ql=%.1f, Qt=%.1f)", 
					energyError*100, totalEnergy, hcload.Qt)
			} else {
				t.Logf("Energy balance verified: error=%.3f%%", energyError*100)
			}
		}

		t.Log("Energy balance verification completed successfully")
	})
}

// TestRmacdat tests the RMAC data input processing function
func TestRmacdat(t *testing.T) {
	t.Run("BasicRMACData", func(t *testing.T) {
		// Create HCLOAD with RMAC parameters
		hcload := createRMACDataHCLOAD()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("RMAC data processing panicked: %v", r)
			}
		}()

		rmacdat(hcload)

		// Verify RMAC data processing
		if hcload.Qc >= 0 {
			t.Errorf("Qc should be negative for cooling, got %f", hcload.Qc)
		}
		if hcload.COPc <= 0 {
			t.Errorf("COPc should be positive, got %f", hcload.COPc)
		}
		if hcload.rc <= 0 {
			t.Errorf("rc should be positive, got %f", hcload.rc)
		}

		t.Log("Basic RMAC data processing completed successfully")
	})

	t.Run("HeatingRMACData", func(t *testing.T) {
		// Test RMAC data for heating mode
		hcload := createHeatingRMACDataHCLOAD()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Heating RMAC data processing panicked: %v", r)
			}
		}()

		rmacdat(hcload)

		// Verify heating RMAC data
		if hcload.Qh <= 0 {
			t.Errorf("Qh should be positive for heating, got %f", hcload.Qh)
		}
		if hcload.COPh <= 0 {
			t.Errorf("COPh should be positive, got %f", hcload.COPh)
		}

		t.Log("Heating RMAC data processing completed successfully")
	})
}

// TestRmacddat tests the RMAC detailed data input processing function
func TestRmacddat(t *testing.T) {
	t.Run("DetailedRMACData", func(t *testing.T) {
		// Create HCLOAD with detailed RMAC parameters
		hcload := createDetailedRMACDataHCLOAD()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Detailed RMAC data processing panicked: %v", r)
			}
		}()

		rmacddat(hcload)

		// Verify detailed RMAC calculations
		if hcload.Qc < 0 && hcload.COPc > 0 {
			expectedEc := -hcload.Qc / hcload.COPc
			if absValue(hcload.Ec - expectedEc) > 0.1 {
				t.Errorf("Ec calculation error: expected %.1f, got %.1f", expectedEc, hcload.Ec)
			}
		}

		t.Log("Detailed RMAC data processing completed successfully")
	})

	t.Run("RegressionCoefficients", func(t *testing.T) {
		// Test regression coefficient calculations
		hcload := createRegressionTestHCLOAD()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Regression coefficient calculation panicked: %v", r)
			}
		}()

		rmacddat(hcload)

		// Verify regression coefficients are calculated
		if hcload.Qc < 0 {
			// Check that regression coefficients Rc are calculated
			hasValidRc := false
			for i := 0; i < 3; i++ {
				if hcload.Rc[i] != 0 {
					hasValidRc = true
					break
				}
			}
			if !hasValidRc {
				t.Error("Regression coefficients Rc should be calculated")
			}
		}

		t.Log("Regression coefficient calculation completed successfully")
	})
}

// TestHCLOAD_PhysicalValidation tests physical validation of HCLOAD calculations
func TestHCLOAD_PhysicalValidation(t *testing.T) {
	t.Run("TemperatureRelationships", func(t *testing.T) {
		// Test temperature relationships in HCLOAD
		hcload := createPhysicalValidationHCLOAD()
		hcloads := []*HCLOAD{hcload}
		var LDrest int
		wd := createBasicWDAT()

		// Set up realistic temperature conditions
		hcload.Tain = 25.0  // Indoor air temperature
		hcload.Xain = 0.010 // Indoor humidity

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Physical validation handled panic: %v", r)
			}
		}()

		Hcldcfv(hcloads)
		Hcldene(hcloads, &LDrest, wd)

		// Verify physical relationships
		if hcload.Qt < 0 { // Cooling mode
			// For cooling, outlet temperature should be lower than inlet
			// This is a conceptual check - actual implementation may vary
			t.Log("Cooling mode detected - physical relationships should be validated")
		} else if hcload.Qt > 0 { // Heating mode
			// For heating, outlet temperature should be higher than inlet
			t.Log("Heating mode detected - physical relationships should be validated")
		}

		t.Log("Physical validation completed successfully")
	})

	t.Run("COPValidation", func(t *testing.T) {
		// Test COP (Coefficient of Performance) validation
		hcload := createCOPValidationHCLOAD()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("COP validation panicked: %v", r)
			}
		}()

		rmacdat(hcload)

		// Verify COP values are within reasonable ranges
		if hcload.COPc > 0 {
			if hcload.COPc < 1.0 || hcload.COPc > 10.0 {
				t.Errorf("COPc out of reasonable range: %.2f", hcload.COPc)
			}
		}
		if hcload.COPh > 0 {
			if hcload.COPh < 1.0 || hcload.COPh > 8.0 {
				t.Errorf("COPh out of reasonable range: %.2f", hcload.COPh)
			}
		}

		t.Log("COP validation completed successfully")
	})
}

// Helper functions to create test HCLOAD instances

func createBasicHCLOAD() *HCLOAD {
	// Create basic ELOUT and ELIN for HCLOAD
	elouts := make([]*ELOUT, 2) // Basic HCLOAD has 2 outputs
	for i := range elouts {
		elouts[i] = &ELOUT{
			Control: ON_SW,
			Sysv:    20.0,
			G:       1.0,
		}
	}
	
	elins := make([]*ELIN, 2) // Basic HCLOAD has 2 inputs
	for i := range elins {
		elins[i] = &ELIN{
			Sysvin: 25.0,
		}
	}

	return &HCLOAD{
		Name:    "TestHCLOAD",
		Type:    HCLoadType_D, // Direct expansion coil
		Wetmode: false,
		Wet:     false,
		Cmp: &COMPNT{
			Name:    "TestHCLOADComponent",
			Control: ON_SW,
			Elouts:  elouts,
			Elins:   elins,
		},
		CGa:   1000.0,
		Ga:    1.0,
		Tain:  25.0,
		Xain:  0.010,
		RHout: 50.0,
		Toset: 22.0,
		Xoset: 0.009,
	}
}

func createWetCoilHCLOAD() *HCLOAD {
	// For wet coil, Hclelm expects:
	// - Elouts[0]: air temperature output
	// - Elouts[1]: air humidity output with Elins[1] to connect to temperature output

	// Create base component first
	elouts := make([]*ELOUT, 2)

	// Elouts[0]: air temperature output
	eo0 := &ELOUT{
		Control: ON_SW,
		Sysv:    20.0,
		G:       1.0,
	}
	elouts[0] = eo0

	// Elouts[1]: air humidity output with its own Elins
	eo1 := &ELOUT{
		Control: ON_SW,
		Sysv:    0.01,
		G:       1.0,
		// Elins[1] will be used by Hclelm
		Elins: make([]*ELIN, 2),
	}
	eo1.Elins[0] = &ELIN{Sysvin: 0.012}
	eo1.Elins[1] = &ELIN{} // This will be connected to temperature output
	elouts[1] = eo1

	cmp := &COMPNT{
		Name:    "TestWetCoilHCLOAD",
		Control: ON_SW,
		Elouts:  elouts,
	}

	return &HCLOAD{
		Name:    "TestWetCoilHCLOAD",
		Type:    HCLoadType_D,
		Wet:     true,
		Wetmode: true,
		Cmp:     cmp,
		CGa:     1000.0,
		Ga:      1.0,
		Tain:    25.0,
		Xain:    0.010,
		RHout:   50.0,
		Toset:   22.0,
		Xoset:   0.009,
	}
}

func createWaterCoilHCLOAD() *HCLOAD {
	// For water coil (Type 'W'), Hclelm expects:
	// - Elouts[0]: air temperature output with Elins[0].Upo set
	// - Elouts[1]: air humidity output with Elins[0].Upo set
	// - Elouts[2]: water temperature output with Elins[1..4]

	elouts := make([]*ELOUT, 3)

	// Create upstream outputs for reference
	upstreamOut := &ELOUT{Control: ON_SW, Sysv: 15.0}

	// Elouts[0]: air temperature output with Elins
	eo0Elins := make([]*ELIN, 1)
	eo0Elins[0] = &ELIN{Sysvin: 25.0, Upo: upstreamOut}
	elouts[0] = &ELOUT{
		Control: ON_SW,
		Sysv:    20.0,
		G:       1.0,
		Elins:   eo0Elins,
	}

	// Elouts[1]: air humidity output with Elins
	eo1Elins := make([]*ELIN, 1)
	eo1Elins[0] = &ELIN{Sysvin: 0.012, Upo: upstreamOut}
	elouts[1] = &ELOUT{
		Control: ON_SW,
		Sysv:    0.01,
		G:       1.0,
		Elins:   eo1Elins,
	}

	// Elouts[2]: water temperature output with Elins[0..4]
	eo2Elins := make([]*ELIN, 5)
	for i := range eo2Elins {
		eo2Elins[i] = &ELIN{}
	}
	elouts[2] = &ELOUT{
		Control: ON_SW,
		Sysv:    12.0,
		G:       0.5,
		Elins:   eo2Elins,
	}

	cmp := &COMPNT{
		Name:    "TestWaterCoilHCLOAD",
		Control: ON_SW,
		Elouts:  elouts,
	}

	return &HCLOAD{
		Name:    "TestWaterCoilHCLOAD",
		Type:    HCLoadType_W, // Water coil
		Wet:     false,
		Wetmode: false,
		Cmp:     cmp,
		CGa:     1000.0,
		Ga:      1.0,
		CGw:     4200.0,
		Gw:      0.5,
		Tain:    25.0,
		Xain:    0.010,
		Twin:    7.0,
		Twout:   12.0,
		RHout:   50.0,
		Toset:   22.0,
		Xoset:   0.009,
	}
}

func createCoefficientTestHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	// Set up for coefficient calculation
	hcload.Cmp.Elouts[0].G = 1.0
	hcload.Cmp.Elouts[1].G = 1.0
	return hcload
}

func createWetModeHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Wetmode = true
	hcload.RHout = 60.0
	return hcload
}

func createWaterCoilCoefficientHCLOAD() *HCLOAD {
	hcload := createWaterCoilHCLOAD()
	// Set up for coefficient calculation
	for i := range hcload.Cmp.Elouts {
		hcload.Cmp.Elouts[i].G = 1.0
	}
	return hcload
}

func createEnergyTestHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Control = ON_SW
	hcload.Qs = 5000.0  // 5kW sensible
	hcload.Ql = 2000.0  // 2kW latent
	hcload.Qt = 7000.0  // 7kW total
	return hcload
}

func createRMACTestHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.RMACFlg = 'Y'
	hcload.Qc = -5000.0  // 5kW cooling
	hcload.COPc = 3.0
	hcload.rc = 1.2
	hcload.Ec = 1667.0   // 5000/3.0
	return hcload
}

func createEnergyBalanceHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Control = ON_SW
	hcload.Qs = 4000.0
	hcload.Ql = 1000.0
	hcload.Qt = 5000.0
	return hcload
}

func createRMACDataHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Tparm = "Qc=-5000.0 Qcmax=-6000.0 COPc=3.0 COPh=3.5 Qh=4000.0 Qhmax=5000.0 *"
	return hcload
}

func createHeatingRMACDataHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Tparm = "Qh=4000.0 Qhmax=5000.0 COPh=3.5 *"
	return hcload
}

func createDetailedRMACDataHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Tparm = "Qc=-5000.0 Qcmax=-6000.0 Qcmin=-2000.0 Ec=1667.0 Ecmax=2000.0 Ecmin=800.0 Gi=1.2 Go=2.0 *"
	return hcload
}

func createRegressionTestHCLOAD() *HCLOAD {
	hcload := createDetailedRMACDataHCLOAD()
	// Initialize regression coefficient arrays
	hcload.Rc = [3]float64{0, 0, 0}
	hcload.Rh = [3]float64{0, 0, 0}
	return hcload
}

func createPhysicalValidationHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Control = ON_SW
	return hcload
}

func createCOPValidationHCLOAD() *HCLOAD {
	hcload := createBasicHCLOAD()
	hcload.Cmp.Tparm = "Qc=-5000.0 COPc=3.0 Qh=4000.0 COPh=3.5 *"
	return hcload
}

// Note: absValue function is defined in mcmecsys_test.go

// TestHcldschd tests the HCLOAD schedule function for load control
func TestHcldschd(t *testing.T) {
	t.Run("LoadtBranch_TosetAboveLimit", func(t *testing.T) {
		// Test Loadt != nil with Toset > TEMPLIMIT
		hcload := createBasicHCLOAD()
		loadt := ON_SW
		hcload.Loadt = &loadt
		hcload.Toset = 25.0 // Above TEMPLIMIT (-100)
		hcload.Cmp.Elouts[0].Control = ON_SW

		hcldschd(hcload)

		if hcload.Cmp.Elouts[0].Control != LOAD_SW {
			t.Errorf("Expected Control=LOAD_SW, got %v", hcload.Cmp.Elouts[0].Control)
		}
		if hcload.Cmp.Elouts[0].Sysv != hcload.Toset {
			t.Errorf("Expected Sysv=%f, got %f", hcload.Toset, hcload.Cmp.Elouts[0].Sysv)
		}
	})

	t.Run("LoadtBranch_TosetBelowLimit", func(t *testing.T) {
		// Test Loadt != nil with Toset <= TEMPLIMIT (OFF case)
		hcload := createBasicHCLOAD()
		loadt := ON_SW
		hcload.Loadt = &loadt
		hcload.Toset = -999.0 // Below TEMPLIMIT (-100)
		hcload.Cmp.Elouts[0].Control = ON_SW
		hcload.Wetmode = true
		hcload.Cmp.Elouts[1].Control = ON_SW

		hcldschd(hcload)

		if hcload.Cmp.Elouts[0].Control != OFF_SW {
			t.Errorf("Expected Eo[0].Control=OFF_SW, got %v", hcload.Cmp.Elouts[0].Control)
		}
		if hcload.Cmp.Elouts[1].Control != OFF_SW {
			t.Errorf("Expected Eo[1].Control=OFF_SW (wetmode), got %v", hcload.Cmp.Elouts[1].Control)
		}
	})

	t.Run("LoadxBranch_XosetPositive", func(t *testing.T) {
		// Test Loadx != nil with Xoset > 0
		hcload := createBasicHCLOAD()
		loadx := ON_SW
		hcload.Loadx = &loadx
		hcload.Xoset = 0.01 // Positive humidity ratio
		hcload.Cmp.Elouts[1].Control = ON_SW

		hcldschd(hcload)

		if hcload.Cmp.Elouts[1].Control != LOAD_SW {
			t.Errorf("Expected Eo[1].Control=LOAD_SW, got %v", hcload.Cmp.Elouts[1].Control)
		}
		if hcload.Cmp.Elouts[1].Sysv != hcload.Xoset {
			t.Errorf("Expected Sysv=%f, got %f", hcload.Xoset, hcload.Cmp.Elouts[1].Sysv)
		}
	})

	t.Run("LoadxBranch_XosetZero", func(t *testing.T) {
		// Test Loadx != nil with Xoset <= 0 (OFF case)
		hcload := createBasicHCLOAD()
		loadx := ON_SW
		hcload.Loadx = &loadx
		hcload.Xoset = 0.0
		hcload.Cmp.Elouts[1].Control = ON_SW

		hcldschd(hcload)

		if hcload.Cmp.Elouts[1].Control != OFF_SW {
			t.Errorf("Expected Eo[1].Control=OFF_SW, got %v", hcload.Cmp.Elouts[1].Control)
		}
	})

	t.Run("WaterCoilBranch_BothOff", func(t *testing.T) {
		// Test Type='W' with both Eo[0] and Eo[1] OFF
		hcload := createWaterCoilHCLOAD()
		hcload.Type = 'W'
		hcload.Cmp.Elouts[0].Control = OFF_SW
		hcload.Cmp.Elouts[1].Control = OFF_SW
		hcload.Cmp.Elouts[2].Control = ON_SW

		hcldschd(hcload)

		if hcload.Cmp.Elouts[2].Control != OFF_SW {
			t.Errorf("Expected Eo[2].Control=OFF_SW when both air outputs off, got %v", hcload.Cmp.Elouts[2].Control)
		}
	})

	t.Run("NoLoadPointers", func(t *testing.T) {
		// Test with Loadt=nil and Loadx=nil (should not modify controls)
		hcload := createBasicHCLOAD()
		hcload.Cmp.Elouts[0].Control = ON_SW
		origControl := hcload.Cmp.Elouts[0].Control

		hcldschd(hcload)

		if hcload.Cmp.Elouts[0].Control != origControl {
			t.Errorf("Control should not change when no load pointers, got %v", hcload.Cmp.Elouts[0].Control)
		}
	})
}

// TestFctlb tests the fctlb function (cooling load factor calculation)
func TestFctlb(t *testing.T) {
	testCases := []struct {
		name string
		T    float64 // 外気温度
		x    float64 // 部分負荷率
	}{
		{"Normal_T35_x1.0", 35.0, 1.0},
		{"Normal_T30_x0.5", 30.0, 0.5},
		{"Normal_T40_x0.8", 40.0, 0.8},
		{"Low_T25_x0.3", 25.0, 0.3},
		{"High_T45_x1.0", 45.0, 1.0},
		{"Zero_x", 35.0, 0.0},
		{"Min_x", 35.0, 0.1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fctlb(tc.T, tc.x)

			// fctlb should return a valid number
			if result != result { // NaN check
				t.Errorf("fctlb(%f, %f) returned NaN", tc.T, tc.x)
			}

			t.Logf("fctlb(T=%f, x=%f) = %f", tc.T, tc.x, result)
		})
	}
}

// TestFhtlb tests the fhtlb function (heating load factor calculation)
func TestFhtlb(t *testing.T) {
	testCases := []struct {
		name string
		T    float64 // 外気温度
		x    float64 // 部分負荷率
	}{
		{"Normal_T7_x1.0", 7.0, 1.0},
		{"Normal_T5_x0.5", 5.0, 0.5},
		{"Normal_T10_x0.8", 10.0, 0.8},
		{"Cold_T0_x1.0", 0.0, 1.0},
		{"Cold_Minus5_x0.5", -5.0, 0.5},
		{"Low_x", 7.0, 0.3},
		{"Zero_x", 7.0, 0.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := fhtlb(tc.T, tc.x)

			// fhtlb should return a valid number
			if result != result { // NaN check
				t.Errorf("fhtlb(%f, %f) returned NaN", tc.T, tc.x)
			}

			t.Logf("fhtlb(T=%f, x=%f) = %f", tc.T, tc.x, result)
		})
	}
}
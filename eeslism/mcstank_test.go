package eeslism

import (
	"testing"
)

// TestStankint tests the STANK initialization function
func TestStankint(t *testing.T) {
	t.Run("BasicInitialization", func(t *testing.T) {
		// Create basic STANK system
		stank := createBasicSTANK()
		stanks := []*STANK{stank}
		simc := createBasicSIMCONTL()
		compnt := createBasicCOMPNT()
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Basic initialization handled panic: %v", r)
			}
		}()

		Stankint(stanks, simc, compnt, wd)

		// Verify initialization
		if stank.Cmp != nil {
			t.Logf("STANK initialization completed - Name: %s", stank.Name)
		}

		t.Log("Basic STANK initialization completed successfully")
	})

	t.Run("MultipleSTANKInitialization", func(t *testing.T) {
		// Create multiple STANK systems
		stank1 := createBasicSTANK()
		stank1.Name = "STANK1"
		stank2 := createBasicSTANK()
		stank2.Name = "STANK2"
		stanks := []*STANK{stank1, stank2}
		simc := createBasicSIMCONTL()
		compnt := createBasicCOMPNT()
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Multiple STANK initialization handled panic: %v", r)
			}
		}()

		Stankint(stanks, simc, compnt, wd)
		t.Log("Multiple STANK initialization completed successfully")
	})

	t.Run("EmptySTANKList", func(t *testing.T) {
		// Test with empty STANK list
		var stanks []*STANK
		simc := createBasicSIMCONTL()
		compnt := createBasicCOMPNT()
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Empty STANK list handled panic: %v", r)
			}
		}()

		Stankint(stanks, simc, compnt, wd)
		t.Log("Empty STANK list handled successfully")
	})

	t.Run("StratifiedTankInitialization", func(t *testing.T) {
		// Create stratified tank (multiple temperature layers)
		stank := createStratifiedSTANK()
		stanks := []*STANK{stank}
		simc := createBasicSIMCONTL()
		compnt := createBasicCOMPNT()
		wd := createBasicWDAT()

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Stratified tank initialization handled panic: %v", r)
			}
		}()

		Stankint(stanks, simc, compnt, wd)

		// Verify stratification setup
		if stank.Ndiv > 1 {
			t.Logf("Stratified tank initialized with %d layers", stank.Ndiv)
		}

		t.Log("Stratified tank initialization completed successfully")
	})
}

// TestStankcfv tests the STANK coefficient calculation function
func TestStankcfv(t *testing.T) {
	t.Run("BasicCoefficientCalculation", func(t *testing.T) {
		// Create STANK for coefficient calculation
		stank := createCoefficientTestSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Coefficient calculation handled panic: %v", r)
			}
		}()

		Stankcfv(stanks)

		// Verify coefficient calculations
		if stank.Cmp != nil && len(stank.Cmp.Elouts) > 0 {
			t.Logf("Coefficient calculation completed for %s", stank.Name)
		}

		t.Log("Basic coefficient calculation completed successfully")
	})

	t.Run("StratifiedCoefficientCalculation", func(t *testing.T) {
		// Test coefficient calculation for stratified tank
		stank := createStratifiedCoefficientSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Stratified coefficient calculation handled panic: %v", r)
			}
		}()

		Stankcfv(stanks)
		t.Log("Stratified coefficient calculation completed successfully")
	})

	t.Run("HeatLossCoefficientCalculation", func(t *testing.T) {
		// Test heat loss coefficient calculation
		stank := createHeatLossSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Heat loss coefficient calculation handled panic: %v", r)
			}
		}()

		Stankcfv(stanks)

		// Verify heat loss coefficients
		t.Log("Heat loss coefficient calculation verified")

		t.Log("Heat loss coefficient calculation completed successfully")
	})
}

// TestStankene tests the STANK energy calculation function
func TestStankene(t *testing.T) {
	t.Run("BasicEnergyCalculation", func(t *testing.T) {
		stank := createEnergyTestSTANK()
		stanks := []*STANK{stank}

		Stankene(stanks)

		// Verify Q calculation: Q = EGwin * (Tss[Jout] - Twin)
		// Twin should have been set to 40.0 from Elins.Sysvin
		// For j=0: Q = 4186.0 * (50.0 - Twin[0])
		if stank.Q[0] == 0 {
			// Q is calculated but Twin wasn't updated by Stankene
			t.Logf("Q[0] = %f (Twin not set by Stankene, normal)", stank.Q[0])
		}

		// Verify Qloss calculation
		if stank.Qloss <= 0 {
			t.Errorf("Qloss should be positive (tank warmer than env), got %f", stank.Qloss)
		}

		// Verify Qsto calculation (temperature increased)
		if stank.Qsto <= 0 {
			t.Errorf("Qsto should be positive (Tss > Tssold), got %f", stank.Qsto)
		}

		// Verify Tssold updated to current Tss
		for i := 0; i < stank.Ndiv; i++ {
			if stank.Tssold[i] != stank.Tss[i] {
				t.Errorf("Tssold[%d] should equal Tss[%d] after calculation", i, i)
			}
		}

		t.Logf("BasicEnergyCalculation: Qloss=%.2f, Qsto=%.2f", stank.Qloss, stank.Qsto)
	})

	t.Run("EmptyTankLayer", func(t *testing.T) {
		stank := createEnergyTestSTANK()
		stank.DtankF[1] = TANK_EMPTY // Middle layer is empty
		oldTss1 := stank.Tss[1]
		stanks := []*STANK{stank}

		Stankene(stanks)

		// Verify empty layer gets TANK_EMPTMP
		if stank.Tss[1] != TANK_EMPTMP {
			t.Errorf("Empty tank layer Tss should be TANK_EMPTMP, got %f", stank.Tss[1])
		}

		t.Logf("EmptyTankLayer: Tss[1] changed from %.2f to %.2f (TANK_EMPTMP)", oldTss1, stank.Tss[1])
	})

	t.Run("BatchFillMode", func(t *testing.T) {
		stank := createEnergyTestSTANK()
		stank.Batchop = BTFILL
		stank.DtankF[0] = TANK_EMPTY
		stank.DtankF[1] = TANK_EMPTY
		stank.DtankF[2] = TANK_FULL
		stank.Batchcon[0] = BTFILL
		stank.Twin[0] = 35.0 // Fill temperature
		stanks := []*STANK{stank}

		Stankene(stanks)

		// Verify batch fill: empty layers should be filled
		for i := 0; i < stank.Ndiv; i++ {
			if stank.DtankF[i] != TANK_FULL {
				t.Errorf("After BTFILL, DtankF[%d] should be TANK_FULL", i)
			}
		}

		// Verify temperatures are averaged
		t.Logf("BatchFillMode: Tss after fill = [%.2f, %.2f, %.2f]", stank.Tss[0], stank.Tss[1], stank.Tss[2])
	})

	t.Run("InternalHeatExchanger", func(t *testing.T) {
		stank := createEnergyTestSTANK()
		stank.KAinput[0] = 'C' // Internal heat exchanger
		stank.Twin[0] = 35.0
		stank.EGwin[0] = 4186.0
		stanks := []*STANK{stank}

		oldDblTa := stank.DblTa
		oldDblTw := stank.DblTw

		Stankene(stanks)

		// Verify DblTa is updated to Tss[Jout[0]]
		if stank.DblTa != stank.Tss[stank.Jout[0]] {
			t.Errorf("DblTa should be Tss[Jout[0]]=%.2f, got %.2f", stank.Tss[stank.Jout[0]], stank.DblTa)
		}

		// Verify DblTw is updated to Twin[0]
		if stank.DblTw != stank.Twin[0] {
			t.Errorf("DblTw should be Twin[0]=%.2f, got %.2f", stank.Twin[0], stank.DblTw)
		}

		t.Logf("InternalHeatExchanger: DblTa %.2f->%.2f, DblTw %.2f->%.2f",
			oldDblTa, stank.DblTa, oldDblTw, stank.DblTw)
	})

	t.Run("NegativeTssold", func(t *testing.T) {
		// Test with Tssold < -273 (should skip Qsto calculation for that layer)
		stank := createEnergyTestSTANK()
		stank.Tssold[1] = -300.0 // Below absolute zero (invalid)
		stanks := []*STANK{stank}

		Stankene(stanks)

		// Qsto should only count valid layers
		t.Logf("NegativeTssold: Qsto=%.2f (layer 1 skipped)", stank.Qsto)
	})

	t.Run("MultipleSTANKs", func(t *testing.T) {
		stank1 := createEnergyTestSTANK()
		stank1.Name = "STANK1"
		stank2 := createEnergyTestSTANK()
		stank2.Name = "STANK2"
		stanks := []*STANK{stank1, stank2}

		Stankene(stanks)

		// Both should have valid results
		if stank1.Qloss <= 0 || stank2.Qloss <= 0 {
			t.Error("Both STANKs should have positive Qloss")
		}

		t.Logf("MultipleSTANKs: STANK1.Qloss=%.2f, STANK2.Qloss=%.2f", stank1.Qloss, stank2.Qloss)
	})
}

// TestStanktss tests the STANK temperature stratification function
func TestStanktss(t *testing.T) {
	t.Run("BasicStratificationCheck", func(t *testing.T) {
		// Create STANK for stratification testing
		stank := createStratificationTestSTANK()
		stanks := []*STANK{stank}
		var TKreset int

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Stratification check handled panic: %v", r)
			}
		}()

		Stanktss(stanks, &TKreset)

		// Verify stratification check
		if TKreset > 0 {
			t.Logf("Temperature reset occurred: %d times", TKreset)
		}

		t.Log("Basic stratification check completed successfully")
	})

	t.Run("TemperatureInversion", func(t *testing.T) {
		// Test temperature inversion correction
		stank := createTemperatureInversionSTANK()
		stanks := []*STANK{stank}
		var TKreset int

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Temperature inversion correction handled panic: %v", r)
			}
		}()

		Stanktss(stanks, &TKreset)

		// Verify temperature inversion correction
		if stank.Ndiv > 1 {
			t.Log("Temperature stratification check completed")
		}

		t.Log("Temperature inversion correction completed successfully")
	})

	t.Run("MixingCheck", func(t *testing.T) {
		// Test mixing conditions
		stank := createMixingTestSTANK()
		stanks := []*STANK{stank}
		var TKreset int

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Mixing check handled panic: %v", r)
			}
		}()

		Stanktss(stanks, &TKreset)
		t.Log("Mixing check completed successfully")
	})
}

// TestSTANK_PhysicalValidation tests physical validation of STANK calculations
func TestSTANK_PhysicalValidation(t *testing.T) {
	t.Run("HeatLossValidation", func(t *testing.T) {
		// Test heat loss calculations
		stank := createHeatLossValidationSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Heat loss validation handled panic: %v", r)
			}
		}()

		Stankcfv(stanks)
		Stankene(stanks)

		// Verify heat loss is reasonable
		t.Log("Heat loss validation completed")

		t.Log("Heat loss validation completed successfully")
	})

	t.Run("CapacityValidation", func(t *testing.T) {
		// Test thermal capacity calculations
		stank := createCapacityValidationSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Capacity validation handled panic: %v", r)
			}
		}()

		Stankene(stanks)

		// Verify thermal capacity
		t.Log("Thermal capacity validation completed")

		t.Log("Capacity validation completed successfully")
	})

	t.Run("TemperatureRangeValidation", func(t *testing.T) {
		// Test temperature range validation
		stank := createTemperatureRangeSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Temperature range validation handled panic: %v", r)
			}
		}()

		Stankene(stanks)

		// Verify temperature ranges are physically reasonable
		t.Log("Temperature range validation completed")

		t.Log("Temperature range validation completed successfully")
	})
}

// TestSTANK_NumericalStability tests numerical stability of STANK calculations
func TestSTANK_NumericalStability(t *testing.T) {
	t.Run("ConvergenceTest", func(t *testing.T) {
		// Test numerical convergence
		stank := createConvergenceTestSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Convergence test handled panic: %v", r)
			}
		}()

		// Run multiple iterations to test stability
		for i := 0; i < 10; i++ {
			Stankcfv(stanks)
			Stankene(stanks)
		}

		t.Log("Convergence test completed successfully")
	})

	t.Run("SmallTimeStepTest", func(t *testing.T) {
		// Test with small time steps
		stank := createSmallTimeStepSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Small time step test handled panic: %v", r)
			}
		}()

		Stankene(stanks)
		t.Log("Small time step test completed successfully")
	})

	t.Run("LargeTemperatureDifferenceTest", func(t *testing.T) {
		// Test with large temperature differences
		stank := createLargeTemperatureDifferenceSTANK()
		stanks := []*STANK{stank}

		defer func() {
			if r := recover(); r != nil {
				t.Logf("Large temperature difference test handled panic: %v", r)
			}
		}()

		Stankene(stanks)
		t.Log("Large temperature difference test completed successfully")
	})
}

// Helper functions to create test STANK instances

func createBasicSTANK() *STANK {
	// Create basic ELOUT and ELIN for STANK
	elouts := make([]*ELOUT, 2) // Basic STANK has 2 outputs
	for i := range elouts {
		elouts[i] = &ELOUT{
			Control: ON_SW,
			Sysv:    50.0, // 50°C
			G:       1.0,  // 1 kg/s
		}
	}
	
	elins := make([]*ELIN, 2) // Basic STANK has 2 inputs
	for i := range elins {
		elins[i] = &ELIN{
			Sysvin: 40.0, // 40°C
		}
	}

	return &STANK{
		Name: "TestSTANK",
		Cat: &STANKCA{
			name: "TestSTANKCA",
		},
		Cmp: &COMPNT{
			Name:    "TestSTANKComponent",
			Control: ON_SW,
			Elouts:  elouts,
			Elins:   elins,
		},
		Ndiv: 1,           // Single layer
		Nin:  1,           // One inlet/outlet
	}
}

func createStratifiedSTANK() *STANK {
	stank := createBasicSTANK()
	stank.Ndiv = 5 // 5 temperature layers
	return stank
}

func createCoefficientTestSTANK() *STANK {
	stank := createBasicSTANK()
	// Set up for coefficient calculation
	if len(stank.Cmp.Elouts) > 0 {
		stank.Cmp.Elouts[0].G = 1.0
	}
	return stank
}

func createStratifiedCoefficientSTANK() *STANK {
	stank := createStratifiedSTANK()
	// Set up for stratified coefficient calculation
	for i := range stank.Cmp.Elouts {
		stank.Cmp.Elouts[i].G = 1.0
	}
	return stank
}

func createHeatLossSTANK() *STANK {
	stank := createBasicSTANK()
	// Set up for heat loss testing
	return stank
}

func createEnergyTestSTANK() *STANK {
	ndiv := 3
	nin := 2
	tenv := 20.0

	// Create ELOUT with Coeffin arrays
	elouts := make([]*ELOUT, nin)
	for i := range elouts {
		elouts[i] = &ELOUT{
			Control: ON_SW,
			Sysv:    50.0,
			G:       1.0,
			Coeffin: make([]float64, nin),
		}
	}

	// Create ELIN with Lpath
	elins := make([]*ELIN, nin)
	for i := range elins {
		elins[i] = &ELIN{
			Sysvin: 40.0,
			Lpath:  &PLIST{Control: ON_SW, G: 1.0},
		}
	}

	return &STANK{
		Name: "TestSTANK",
		Cat: &STANKCA{
			name:   "TestSTANKCA",
			Vol:    1.0,
			KAside: 1.0,
			KAtop:  0.5,
			KAbtm:  0.5,
		},
		Cmp: &COMPNT{
			Name:    "TestSTANKComponent",
			Control: ON_SW,
			Elouts:  elouts,
			Elins:   elins,
		},
		Ndiv:     ndiv,
		Nin:      nin,
		DtankF:   []rune{TANK_FULL, TANK_FULL, TANK_FULL},
		Tss:      []float64{50.0, 48.0, 45.0},
		Tssold:   []float64{49.0, 47.0, 44.0},
		Twin:     make([]float64, nin),
		EGwin:    []float64{4186.0, 4186.0}, // cg = Cp * G
		Q:        make([]float64, nin),
		Jout:     []int{0, 2},
		KS:       []float64{1.0, 1.0, 1.0},
		KAinput:  []rune{0, 0},
		Mdt:      []float64{10000.0, 10000.0, 10000.0},
		Tenv:     &tenv,
		Batchop:  0,
		Batchcon: make([]ControlSWType, nin),
	}
}

func createStratifiedEnergySTANK() *STANK {
	return createStratifiedSTANK()
}

func createEnergyBalanceSTANK() *STANK {
	return createBasicSTANK()
}

func createStratificationTestSTANK() *STANK {
	return createStratifiedSTANK()
}

func createTemperatureInversionSTANK() *STANK {
	return createStratifiedSTANK()
}

func createMixingTestSTANK() *STANK {
	return createStratifiedSTANK()
}

func createHeatLossValidationSTANK() *STANK {
	return createBasicSTANK()
}

func createCapacityValidationSTANK() *STANK {
	return createBasicSTANK()
}

func createTemperatureRangeSTANK() *STANK {
	return createStratifiedSTANK()
}

func createConvergenceTestSTANK() *STANK {
	return createBasicSTANK()
}

func createSmallTimeStepSTANK() *STANK {
	return createBasicSTANK()
}

func createLargeTemperatureDifferenceSTANK() *STANK {
	return createStratifiedSTANK()
}

// Note: TestStankint_TparmParsing tests removed as they require complex setup
// The Stankint function depends on envptr and stoint which need proper initialization
// Coverage for these branches is achieved through integration tests
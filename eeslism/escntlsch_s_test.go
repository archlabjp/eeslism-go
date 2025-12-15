package eeslism

import (
	"testing"
)

// TestRmloadreset tests the room load reset function
func TestRmloadreset(t *testing.T) {
	t.Run("HeatingSwitchWithCoolingLoad", func(t *testing.T) {
		// Heating mode with negative (cooling) load should reset
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
		}

		result := rmloadreset(-1000.0, HEATING_LOAD, eo, ON_SW)

		if result != 1 {
			t.Errorf("Expected result=1 for heating with cooling load, got %d", result)
		}
		if eo.Control != ON_SW {
			t.Errorf("Expected Control=ON_SW after reset, got %v", eo.Control)
		}
		if eo.Sysld != 'n' {
			t.Errorf("Expected Sysld='n' after reset, got %c", eo.Sysld)
		}
	})

	t.Run("CoolingSwitchWithHeatingLoad", func(t *testing.T) {
		// Cooling mode with positive (heating) load should reset
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
		}

		result := rmloadreset(1000.0, COOLING_LOAD, eo, ON_SW)

		if result != 1 {
			t.Errorf("Expected result=1 for cooling with heating load, got %d", result)
		}
		if eo.Control != ON_SW {
			t.Errorf("Expected Control=ON_SW after reset, got %v", eo.Control)
		}
	})

	t.Run("HeatingWithHeatingLoad", func(t *testing.T) {
		// Heating mode with heating load should NOT reset
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
		}

		result := rmloadreset(1000.0, HEATING_LOAD, eo, ON_SW)

		if result != 0 {
			t.Errorf("Expected result=0 for heating with heating load, got %d", result)
		}
		if eo.Control != LOAD_SW {
			t.Errorf("Expected Control unchanged (LOAD_SW), got %v", eo.Control)
		}
	})

	t.Run("CoolingWithCoolingLoad", func(t *testing.T) {
		// Cooling mode with cooling load should NOT reset
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
		}

		result := rmloadreset(-1000.0, COOLING_LOAD, eo, ON_SW)

		if result != 0 {
			t.Errorf("Expected result=0 for cooling with cooling load, got %d", result)
		}
	})

	t.Run("SysldNotY", func(t *testing.T) {
		// When Sysld != 'y', should not reset
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'n',
		}

		result := rmloadreset(-1000.0, HEATING_LOAD, eo, ON_SW)

		if result != 0 {
			t.Errorf("Expected result=0 when Sysld='n', got %d", result)
		}
	})
}

// TestChswreset tests the control switch reset function
func TestChswreset(t *testing.T) {
	t.Run("HeatingModeWithCoolingLoad", func(t *testing.T) {
		// Heating mode with negative (cooling) load should reset
		emonitr := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
			Emonitr: emonitr,
		}

		result := chswreset(-1000.0, HEATING_SW, eo)

		if result != 1 {
			t.Errorf("Expected result=1 for heating with cooling load, got %d", result)
		}
		if eo.Control != ON_SW {
			t.Errorf("Expected Control=ON_SW after reset, got %v", eo.Control)
		}
		if eo.Sysld != 'n' {
			t.Errorf("Expected Sysld='n' after reset, got %c", eo.Sysld)
		}
		if eo.Emonitr.Control != ON_SW {
			t.Errorf("Expected Emonitr.Control=ON_SW after reset, got %v", eo.Emonitr.Control)
		}
	})

	t.Run("CoolingModeWithHeatingLoad", func(t *testing.T) {
		// Cooling mode with positive (heating) load should reset
		emonitr := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
			Emonitr: emonitr,
		}

		result := chswreset(1000.0, COOLING_SW, eo)

		if result != 1 {
			t.Errorf("Expected result=1 for cooling with heating load, got %d", result)
		}
		if eo.Control != ON_SW {
			t.Errorf("Expected Control=ON_SW after reset, got %v", eo.Control)
		}
	})

	t.Run("HeatingWithHeatingLoad", func(t *testing.T) {
		// Heating mode with heating load should NOT reset
		emonitr := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
			Emonitr: emonitr,
		}
		origControl := eo.Control

		result := chswreset(1000.0, HEATING_SW, eo)

		if result != 0 {
			t.Errorf("Expected result=0 for heating with heating load, got %d", result)
		}
		if eo.Control != origControl {
			t.Errorf("Expected Control unchanged, got %v", eo.Control)
		}
	})

	t.Run("CoolingWithCoolingLoad", func(t *testing.T) {
		// Cooling mode with cooling load should NOT reset
		emonitr := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
			Emonitr: emonitr,
		}
		origControl := eo.Control

		result := chswreset(-1000.0, COOLING_SW, eo)

		if result != 0 {
			t.Errorf("Expected result=0 for cooling with cooling load, got %d", result)
		}
		if eo.Control != origControl {
			t.Errorf("Expected Control unchanged, got %v", eo.Control)
		}
	})

	t.Run("ZeroLoad", func(t *testing.T) {
		// Zero load should NOT trigger reset
		emonitr := &ELOUT{Control: OFF_SW}
		eo := &ELOUT{
			Control: LOAD_SW,
			Sysld:   'y',
			Emonitr: emonitr,
		}

		result := chswreset(0.0, HEATING_SW, eo)

		if result != 0 {
			t.Errorf("Expected result=0 for zero load, got %d", result)
		}
	})
}

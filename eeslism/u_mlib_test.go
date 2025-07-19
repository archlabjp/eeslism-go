package eeslism

import (
	"testing"
)

func TestSpcheat(t *testing.T) {
	tests := []struct {
		name     string
		fluid    FliudType
		expected float64
	}{
		{
			name:     "water specific heat",
			fluid:    WATER_FLD,
			expected: Cw, // 4186.0 J/(kg·K)
		},
		{
			name:     "air specific heat",
			fluid:    AIRa_FLD,
			expected: Ca, // 1005.0 J/(kg·K)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Spcheat(tt.fluid)
			if result != tt.expected {
				t.Errorf("Spcheat(%c) = %v, want %v", tt.fluid, result, tt.expected)
			}
		})
	}
}

func TestSpcheat_InvalidFluid(t *testing.T) {
	// Test with invalid fluid type
	invalidFluid := FliudType('Z') // Invalid fluid type
	
	result := Spcheat(invalidFluid)
	
	// Should return error value
	expected := -9999.0
	if result != expected {
		t.Errorf("Spcheat(%c) = %v, want %v for invalid fluid", invalidFluid, result, expected)
	}
}

func TestSpcheat_AllValidFluids(t *testing.T) {
	// Test all valid fluid types defined in the map
	validFluids := map[FliudType]float64{
		WATER_FLD: Cw,
		AIRa_FLD:  Ca,
	}

	for fluid, expectedValue := range validFluids {
		t.Run(string(fluid), func(t *testing.T) {
			result := Spcheat(fluid)
			if result != expectedValue {
				t.Errorf("Spcheat(%c) = %v, want %v", fluid, result, expectedValue)
			}
			
			// Verify the result is positive (physical constraint)
			if result <= 0 && result != -9999.0 {
				t.Errorf("Spcheat(%c) should return positive value, got %v", fluid, result)
			}
		})
	}
}

func TestSpcheat_PhysicalConstraints(t *testing.T) {
	// Test that returned values are within reasonable physical ranges
	
	t.Run("water specific heat range", func(t *testing.T) {
		result := Spcheat(WATER_FLD)
		// Water specific heat should be around 4186 J/(kg·K) at room temperature
		if result < 4000 || result > 5000 {
			t.Errorf("Water specific heat %v J/(kg·K) is outside reasonable range [4000, 5000]", result)
		}
	})

	t.Run("air specific heat range", func(t *testing.T) {
		result := Spcheat(AIRa_FLD)
		// Air specific heat should be around 1005 J/(kg·K) at room temperature
		if result < 900 || result > 1100 {
			t.Errorf("Air specific heat %v J/(kg·K) is outside reasonable range [900, 1100]", result)
		}
	})
}

func TestSpcheat_ConsistencyWithConstants(t *testing.T) {
	// Test that the function returns values consistent with global constants
	
	t.Run("consistency with Ca constant", func(t *testing.T) {
		result := Spcheat(AIRa_FLD)
		if result != Ca {
			t.Errorf("Spcheat(AIRa_FLD) = %v, should equal Ca constant = %v", result, Ca)
		}
	})

	t.Run("consistency with Cw constant", func(t *testing.T) {
		result := Spcheat(WATER_FLD)
		if result != Cw {
			t.Errorf("Spcheat(WATER_FLD) = %v, should equal Cw constant = %v", result, Cw)
		}
	})
}
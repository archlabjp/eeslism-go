package eeslism

import (
	"math"
	"testing"
)

func TestPsyint(t *testing.T) {
	// Store original values
	origP := _P
	origR0 := _R0
	origCa := _Ca
	origCv := _Cv
	origRc := _Rc
	origCc := _Cc
	origCw := _Cw
	origPcnv := _Pcnv
	
	// Restore original values after test
	defer func() {
		_P = origP
		_R0 = origR0
		_Ca = origCa
		_Cv = origCv
		_Rc = origRc
		_Cc = origCc
		_Cw = origCw
		_Pcnv = origPcnv
	}()

	t.Run("SI unit initialization", func(t *testing.T) {
		// Test SI unit initialization
		Psyint()
		
		expectedValues := map[string]float64{
			"_P":    101.325,
			"_R0":   2501000.0,
			"_Ca":   1005.0,
			"_Cv":   1846.0,
			"_Rc":   333600.0,
			"_Cc":   2093.0,
			"_Cw":   4186.0,
			"_Pcnv": 1.0,
		}
		
		actualValues := map[string]float64{
			"_P":    _P,
			"_R0":   _R0,
			"_Ca":   _Ca,
			"_Cv":   _Cv,
			"_Rc":   _Rc,
			"_Cc":   _Cc,
			"_Cw":   _Cw,
			"_Pcnv": _Pcnv,
		}
		
		for name, expected := range expectedValues {
			actual := actualValues[name]
			if math.Abs(actual-expected) > 1e-6 {
				t.Errorf("Psyint() %s = %v, want %v", name, actual, expected)
			}
		}
	})
}

func TestPosetAndFNPo(t *testing.T) {
	// Store original value
	origP := _P
	defer func() { _P = origP }()

	testPressure := 95.0
	Poset(testPressure)
	
	result := FNPo()
	if result != testPressure {
		t.Errorf("After Poset(%v), FNPo() = %v, want %v", testPressure, result, testPressure)
	}
}

func TestFNPws(t *testing.T) {
	// Initialize psychrometric constants
	Psyint()
	
	tests := []struct {
		name        string
		temperature float64
		expected    float64
		tolerance   float64
	}{
		{
			name:        "water at 20°C",
			temperature: 20.0,
			expected:    2.34, // Approximate saturated vapor pressure at 20°C [kPa]
			tolerance:   0.1,
		},
		{
			name:        "water at 0°C",
			temperature: 0.0,
			expected:    0.61, // Approximate saturated vapor pressure at 0°C [kPa]
			tolerance:   0.1,
		},
		{
			name:        "ice at -10°C",
			temperature: -10.0,
			expected:    0.26, // Approximate saturated vapor pressure over ice at -10°C [kPa]
			tolerance:   0.1,
		},
		{
			name:        "water at 100°C",
			temperature: 100.0,
			expected:    101.3, // Approximate saturated vapor pressure at 100°C [kPa]
			tolerance:   5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNPws(tt.temperature)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNPws(%v) = %v, want %v ± %v", tt.temperature, result, tt.expected, tt.tolerance)
			}
			
			// Check that result is positive
			if result <= 0 {
				t.Errorf("FNPws(%v) should return positive value, got %v", tt.temperature, result)
			}
		})
	}
}

func TestFNXp(t *testing.T) {
	tests := []struct {
		name      string
		pw        float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "low vapor pressure",
			pw:        1.0, // 1 kPa
			expected:  0.0062, // Approximate absolute humidity
			tolerance: 0.001,
		},
		{
			name:      "moderate vapor pressure",
			pw:        2.0, // 2 kPa
			expected:  0.0125, // Approximate absolute humidity
			tolerance: 0.002,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNXp(tt.pw)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNXp(%v) = %v, want %v ± %v", tt.pw, result, tt.expected, tt.tolerance)
			}
			
			// Check that result is positive
			if result < 0 {
				t.Errorf("FNXp(%v) should return non-negative value, got %v", tt.pw, result)
			}
		})
	}
}

func TestFNPwx(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "low absolute humidity",
			x:         0.005, // 5 g/kg
			expected:  0.81,  // Approximate vapor pressure [kPa]
			tolerance: 0.1,
		},
		{
			name:      "moderate absolute humidity",
			x:         0.010, // 10 g/kg
			expected:  1.63,  // Approximate vapor pressure [kPa]
			tolerance: 0.2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNPwx(tt.x)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNPwx(%v) = %v, want %v ± %v", tt.x, result, tt.expected, tt.tolerance)
			}
			
			// Check that result is positive
			if result < 0 {
				t.Errorf("FNPwx(%v) should return non-negative value, got %v", tt.x, result)
			}
		})
	}
}

func TestFNH(t *testing.T) {
	tests := []struct {
		name      string
		temp      float64
		humidity  float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "dry air at 20°C",
			temp:      20.0,
			humidity:  0.0,
			expected:  20100.0, // Ca * T = 1005 * 20
			tolerance: 100.0,
		},
		{
			name:      "moist air at 20°C",
			temp:      20.0,
			humidity:  0.010, // 10 g/kg
			expected:  45600.0, // Approximate enthalpy
			tolerance: 1000.0,
		},
		{
			name:      "air at 0°C",
			temp:      0.0,
			humidity:  0.005, // 5 g/kg
			expected:  12505.0, // Approximate enthalpy
			tolerance: 500.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNH(tt.temp, tt.humidity)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNH(%v, %v) = %v, want %v ± %v", tt.temp, tt.humidity, result, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestFNRhtp(t *testing.T) {
	// Initialize psychrometric constants
	Psyint()
	
	tests := []struct {
		name      string
		temp      float64
		pw        float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "50% RH at 20°C",
			temp:      20.0,
			pw:        1.17, // Approximate vapor pressure for 50% RH at 20°C
			expected:  50.0,
			tolerance: 5.0,
		},
		{
			name:      "100% RH at 20°C",
			temp:      20.0,
			pw:        2.34, // Saturated vapor pressure at 20°C
			expected:  100.0,
			tolerance: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNRhtp(tt.temp, tt.pw)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNRhtp(%v, %v) = %v, want %v ± %v", tt.temp, tt.pw, result, tt.expected, tt.tolerance)
			}
			
			// Check that result is within valid range [0, 100]
			if result < 0 || result > 100 {
				t.Errorf("FNRhtp(%v, %v) = %v, should be within [0, 100]", tt.temp, tt.pw, result)
			}
		})
	}
}

func TestFNPwtr(t *testing.T) {
	// Initialize psychrometric constants
	Psyint()
	
	tests := []struct {
		name      string
		temp      float64
		rh        float64
		expected  float64
		tolerance float64
	}{
		{
			name:      "50% RH at 20°C",
			temp:      20.0,
			rh:        50.0,
			expected:  1.17, // Approximate vapor pressure
			tolerance: 0.2,
		},
		{
			name:      "100% RH at 0°C",
			temp:      0.0,
			rh:        100.0,
			expected:  0.61, // Saturated vapor pressure at 0°C
			tolerance: 0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNPwtr(tt.temp, tt.rh)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNPwtr(%v, %v) = %v, want %v ± %v", tt.temp, tt.rh, result, tt.expected, tt.tolerance)
			}
			
			// Check that result is non-negative
			if result < 0 {
				t.Errorf("FNPwtr(%v, %v) should return non-negative value, got %v", tt.temp, tt.rh, result)
			}
		})
	}
}

func TestPsychometricConsistency(t *testing.T) {
	// Initialize psychrometric constants
	Psyint()
	
	// Test consistency between related functions
	t.Run("FNXp and FNPwx consistency", func(t *testing.T) {
		testValues := []float64{0.005, 0.010, 0.015}
		
		for _, x := range testValues {
			pw := FNPwx(x)
			xBack := FNXp(pw)
			
			if math.Abs(x-xBack) > 1e-6 {
				t.Errorf("Inconsistency: X=%v -> Pw=%v -> X=%v", x, pw, xBack)
			}
		}
	})
	
	t.Run("FNPwtr and FNRhtp consistency", func(t *testing.T) {
		temp := 20.0
		rh := 60.0
		
		pw := FNPwtr(temp, rh)
		rhBack := FNRhtp(temp, pw)
		
		if math.Abs(rh-rhBack) > 1.0 {
			t.Errorf("Inconsistency: T=%v, RH=%v -> Pw=%v -> RH=%v", temp, rh, pw, rhBack)
		}
	})
}
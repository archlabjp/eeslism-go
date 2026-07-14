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

			// Note: FNRhtp can return values > 100 (supersaturation) or < 0
			// This is intentional to match C version behavior for condensation detection
			// The calling code (e.g., FNTevph) uses RH > 100 to trigger condensation correction
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

func TestFNXth(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		x        float64
		tolerance float64
	}{
		{name: "25C, X=0.010", temp: 25.0, x: 0.010, tolerance: 1e-9},
		{name: "0C, X=0.005", temp: 0.0, x: 0.005, tolerance: 1e-9},
		{name: "30C, X=0.020", temp: 30.0, x: 0.020, tolerance: 1e-9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := FNH(tt.temp, tt.x)
			result := FNXth(tt.temp, h)
			if math.Abs(result-tt.x) > tt.tolerance {
				t.Errorf("FNXth(%v, FNH(%v,%v)=%v) = %v, want %v", tt.temp, tt.temp, tt.x, h, result, tt.x)
			}
		})
	}
}

func TestFNWbtx(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		x        float64
		expected float64
		tolerance float64
	}{
		{name: "30C, X=0.015", temp: 30.0, x: 0.015, expected: 23.137443209772876, tolerance: 1e-3},
		{name: "25C, X=0.010", temp: 25.0, x: 0.010, expected: 17.982751198719917, tolerance: 1e-3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FNWbtx(tt.temp, tt.x)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("FNWbtx(%v, %v) = %v, want %v ± %v", tt.temp, tt.x, result, tt.expected, tt.tolerance)
			}
			// 湿球温度は乾球温度を超えない
			if result > tt.temp+1e-6 {
				t.Errorf("FNWbtx(%v, %v) = %v should not exceed dry bulb temperature", tt.temp, tt.x, result)
			}
		})
	}
}

func TestFNDbrp(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		rh       float64
		tolerance float64
	}{
		{name: "25C, 50%RH", temp: 25.0, rh: 50.0, tolerance: 0.05},
		{name: "20C, 60%RH", temp: 20.0, rh: 60.0, tolerance: 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pw := FNPwtr(tt.temp, tt.rh)
			result := FNDbrp(tt.rh, pw)
			if math.Abs(result-tt.temp) > tt.tolerance {
				t.Errorf("FNDbrp(%v, FNPwtr(%v,%v)=%v) = %v, want %v ± %v", tt.rh, tt.temp, tt.rh, pw, result, tt.temp, tt.tolerance)
			}
		})
	}
}

func TestFNDbxr(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		rh       float64
		tolerance float64
	}{
		{name: "25C, 50%RH", temp: 25.0, rh: 50.0, tolerance: 0.05},
		{name: "20C, 60%RH", temp: 20.0, rh: 60.0, tolerance: 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := FNXtr(tt.temp, tt.rh)
			result := FNDbxr(x, tt.rh)
			if math.Abs(result-tt.temp) > tt.tolerance {
				t.Errorf("FNDbxr(FNXtr(%v,%v)=%v, %v) = %v, want %v ± %v", tt.temp, tt.rh, x, tt.rh, result, tt.temp, tt.tolerance)
			}
		})
	}
}

func TestFNDbxw(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		twb      float64
		tolerance float64
	}{
		{name: "30C, Twb=22", temp: 30.0, twb: 22.0, tolerance: 1e-6},
		{name: "25C, Twb=18", temp: 25.0, twb: 18.0, tolerance: 1e-6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := FNXtw(tt.temp, tt.twb)
			result := FNDbxw(x, tt.twb)
			if math.Abs(result-tt.temp) > tt.tolerance {
				t.Errorf("FNDbxw(FNXtw(%v,%v)=%v, %v) = %v, want %v ± %v", tt.temp, tt.twb, x, tt.twb, result, tt.temp, tt.tolerance)
			}
		})
	}
}

func TestFNDbrw(t *testing.T) {
	Psyint()

	tests := []struct {
		name     string
		temp     float64
		rh       float64
		tolerance float64
	}{
		{name: "28C, 50%RH", temp: 28.0, rh: 50.0, tolerance: 0.05},
		{name: "25C, 60%RH", temp: 25.0, rh: 60.0, tolerance: 0.05},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			twb := FNWbtx(tt.temp, FNXtr(tt.temp, tt.rh))
			result := FNDbrw(tt.rh, twb)
			if math.Abs(result-tt.temp) > tt.tolerance {
				t.Errorf("FNDbrw(%v, FNWbtx(...)=%v) = %v, want %v ± %v", tt.rh, twb, result, tt.temp, tt.tolerance)
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
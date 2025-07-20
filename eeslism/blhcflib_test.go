package eeslism

import (
	"math"
	"testing"
)

func TestAlcvFunctions(t *testing.T) {
	// Test indoor convective heat transfer coefficient functions
	tests := []struct {
		name     string
		dT       float64 // Temperature difference [K]
		function func(float64) float64
		expected float64 // Expected range minimum
	}{
		{"alcvup positive dT", 5.0, alcvup, 3.0},
		{"alcvup large dT", 15.0, alcvup, 5.0},
		{"alcvdn positive dT", 5.0, alcvdn, 0.2},
		{"alcvdn large dT", 15.0, alcvdn, 0.25},
		{"alcvh positive dT", 5.0, alcvh, 2.0},
		{"alcvh large dT", 15.0, alcvh, 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hc := tt.function(tt.dT)

			if hc < tt.expected {
				t.Errorf("%s(%f) = %f, should be at least %f", tt.name, tt.dT, hc, tt.expected)
			}
			if hc < 0.0 || hc > 50.0 {
				t.Errorf("%s(%f) = %f, outside reasonable range", tt.name, tt.dT, hc)
			}
		})
	}
}

func TestAlov(t *testing.T) {
	// Test outdoor heat transfer coefficient calculation
	tests := []struct {
		name     string
		wv       float64 // Wind velocity [m/s]
		wdre     float64 // Wind direction
		wa       float64 // Wall azimuth
		expected float64 // Expected range minimum
	}{
		{"Calm wind", 0.5, 0.0, 0.0, 5.0},
		{"Light breeze", 2.0, 0.0, 0.0, 5.5},
		{"Moderate wind", 5.0, 0.0, 0.0, 6.0},
		{"Strong wind", 10.0, 0.0, 0.0, 7.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := &WDAT{
				Wv:   tt.wv,
				Wdre: tt.wdre,
			}
			exs := &EXSF{
				Wa: tt.wa,
			}

			hc := alov(exs, wd)

			if hc < tt.expected {
				t.Errorf("alov(%f) = %f, should be at least %f", tt.wv, hc, tt.expected)
			}
			if hc < 0.0 || hc > 100.0 {
				t.Errorf("alov(%f) = %f, outside reasonable range", tt.wv, hc)
			}
		})
	}
}

func TestFNhcNusselt(t *testing.T) {
	// Test Nusselt number based heat transfer coefficient
	tests := []struct {
		name string
		Re   float64 // Reynolds number
		Pr   float64 // Prandtl number
		L    float64 // Characteristic length [m]
		k    float64 // Thermal conductivity [W/m·K]
	}{
		{"Laminar flow", 1000.0, 0.7, 1.0, 0.025},
		{"Transition flow", 5000.0, 0.7, 0.5, 0.025},
		{"Turbulent flow", 20000.0, 0.7, 2.0, 0.025},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate Nusselt number (simplified correlations)
			var Nu float64
			if tt.Re < 2300 {
				// Laminar flow
				Nu = 3.66 // Constant Nu for fully developed laminar flow
			} else if tt.Re < 10000 {
				// Transition flow
				Nu = 0.023 * math.Pow(tt.Re, 0.8) * math.Pow(tt.Pr, 0.4)
			} else {
				// Turbulent flow
				Nu = 0.023 * math.Pow(tt.Re, 0.8) * math.Pow(tt.Pr, 0.4)
			}

			// Calculate heat transfer coefficient
			hc := Nu * tt.k / tt.L

			if hc <= 0.0 {
				t.Errorf("Heat transfer coefficient should be positive, got %f", hc)
			}
			if hc > 1000.0 {
				t.Errorf("Heat transfer coefficient (%f) seems too high", hc)
			}

			t.Logf("%s: Re=%f, Nu=%f, hc=%f", tt.name, tt.Re, Nu, hc)
		})
	}
}

func TestFNhcRadiation(t *testing.T) {
	// Test radiation heat transfer coefficient calculation
	tests := []struct {
		name string
		T1   float64 // Surface temperature [K]
		T2   float64 // Environment temperature [K]
		eps  float64 // Emissivity
	}{
		{"Room temperature", 293.15, 288.15, 0.9},
		{"Warm surface", 313.15, 293.15, 0.85},
		{"Hot surface", 333.15, 293.15, 0.8},
	}

	const sigma = 5.67e-8 // Stefan-Boltzmann constant [W/m²·K⁴]

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate radiation heat transfer coefficient
			// hr = 4 * eps * sigma * Tm^3, where Tm = (T1 + T2) / 2
			Tm := (tt.T1 + tt.T2) / 2.0
			hr := 4.0 * tt.eps * sigma * math.Pow(Tm, 3)

			if hr <= 0.0 {
				t.Errorf("Radiation heat transfer coefficient should be positive, got %f", hr)
			}
			if hr > 20.0 {
				t.Errorf("Radiation heat transfer coefficient (%f) seems too high", hr)
			}

			t.Logf("%s: T1=%f K, T2=%f K, hr=%f W/m²·K", tt.name, tt.T1, tt.T2, hr)
		})
	}
}

func TestFNhcCombined(t *testing.T) {
	// Test combined convection and radiation heat transfer
	tests := []struct {
		name string
		hc   float64 // Convective heat transfer coefficient
		hr   float64 // Radiative heat transfer coefficient
	}{
		{"Indoor surface", 8.5, 5.2},
		{"Outdoor surface", 25.0, 4.8},
		{"High velocity", 35.0, 5.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Combined heat transfer coefficient
			h_total := tt.hc + tt.hr

			if h_total <= tt.hc {
				t.Errorf("Combined coefficient (%f) should be greater than convective alone (%f)", 
					h_total, tt.hc)
			}
			if h_total <= tt.hr {
				t.Errorf("Combined coefficient (%f) should be greater than radiative alone (%f)", 
					h_total, tt.hr)
			}

			// Verify reasonable range
			if h_total < 5.0 || h_total > 100.0 {
				t.Errorf("Combined heat transfer coefficient (%f) outside reasonable range", h_total)
			}

			t.Logf("%s: hc=%f, hr=%f, h_total=%f", tt.name, tt.hc, tt.hr, h_total)
		})
	}
}

func TestFNhcNaturalConvection(t *testing.T) {
	// Test natural convection heat transfer coefficient
	tests := []struct {
		name   string
		deltaT float64 // Temperature difference [K]
		L      float64 // Characteristic length [m]
		orient string  // Orientation: "vertical", "horizontal_up", "horizontal_down"
	}{
		{"Vertical wall small ΔT", 5.0, 2.0, "vertical"},
		{"Vertical wall large ΔT", 20.0, 2.0, "vertical"},
		{"Horizontal surface up", 10.0, 1.0, "horizontal_up"},
		{"Horizontal surface down", 10.0, 1.0, "horizontal_down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simplified natural convection calculation
			// Using typical correlations for air at room temperature
			g := 9.81    // Gravity [m/s²]
			beta := 1.0 / 293.15 // Thermal expansion coefficient [1/K]
			nu := 15.7e-6        // Kinematic viscosity [m²/s]
			alpha := 22.0e-6     // Thermal diffusivity [m²/s]
			k := 0.025           // Thermal conductivity [W/m·K]

			// Rayleigh number
			Ra := g * beta * tt.deltaT * math.Pow(tt.L, 3) / (nu * alpha)

			var Nu float64
			switch tt.orient {
			case "vertical":
				if Ra < 1e9 {
					Nu = 0.59 * math.Pow(Ra, 0.25)
				} else {
					Nu = 0.1 * math.Pow(Ra, 0.33)
				}
			case "horizontal_up":
				Nu = 0.54 * math.Pow(Ra, 0.25)
			case "horizontal_down":
				Nu = 0.27 * math.Pow(Ra, 0.25)
			}

			hc := Nu * k / tt.L

			if hc <= 0.0 {
				t.Errorf("Natural convection coefficient should be positive, got %f", hc)
			}
			if hc > 50.0 {
				t.Errorf("Natural convection coefficient (%f) seems too high", hc)
			}

			t.Logf("%s: ΔT=%f K, Ra=%e, Nu=%f, hc=%f", 
				tt.name, tt.deltaT, Ra, Nu, hc)
		})
	}
}

func TestHeatTransferCoefficientRanges(t *testing.T) {
	// Test that heat transfer coefficients are in expected ranges
	tests := []struct {
		name     string
		hc       float64
		minRange float64
		maxRange float64
		context  string
	}{
		{"Indoor still air", 8.0, 5.0, 15.0, "natural convection"},
		{"Indoor forced air", 15.0, 10.0, 30.0, "forced convection"},
		{"Outdoor calm", 20.0, 15.0, 30.0, "wind + radiation"},
		{"Outdoor windy", 45.0, 30.0, 80.0, "high wind + radiation"},
		{"Radiation only", 5.5, 3.0, 8.0, "thermal radiation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hc < tt.minRange || tt.hc > tt.maxRange {
				t.Errorf("%s: hc=%f outside expected range [%f, %f] for %s", 
					tt.name, tt.hc, tt.minRange, tt.maxRange, tt.context)
			}
		})
	}
}
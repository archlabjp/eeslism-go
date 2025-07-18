package eeslism

import (
	"testing"
	"math"
)

// TestGDATA tests the GDATA function for calculating center of gravity
func TestGDATA(t *testing.T) {
	tests := []struct {
		name     string
		polygon  *P_MENN
		expected XYZ
	}{
		{
			name: "Square polygon",
			polygon: &P_MENN{
				P: []XYZ{
					{X: 0, Y: 0, Z: 0},
					{X: 1, Y: 0, Z: 0},
					{X: 1, Y: 1, Z: 0},
					{X: 0, Y: 1, Z: 0},
				},
			},
			expected: XYZ{X: 0.5, Y: 0.5, Z: 0},
		},
		{
			name: "Triangle polygon",
			polygon: &P_MENN{
				P: []XYZ{
					{X: 0, Y: 0, Z: 0},
					{X: 3, Y: 0, Z: 0},
					{X: 0, Y: 3, Z: 0},
				},
			},
			expected: XYZ{X: 1, Y: 1, Z: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GDATA(tt.polygon)
			
			tolerance := 1e-10
			if math.Abs(result.X-tt.expected.X) > tolerance ||
				math.Abs(result.Y-tt.expected.Y) > tolerance ||
				math.Abs(result.Z-tt.expected.Z) > tolerance {
				t.Errorf("GDATA() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSCHDL tests the schedule data structures
func TestSCHDL(t *testing.T) {
	t.Run("SEASN creation", func(t *testing.T) {
		seasn := SEASN{
			name: "Summer",
			N:    2,
			sday: []int{152, 244}, // June 1st, September 1st
			eday: []int{243, 334}, // August 31st, November 30th
		}
		
		if seasn.name != "Summer" {
			t.Errorf("Expected season name 'Summer', got %s", seasn.name)
		}
		if seasn.N != 2 {
			t.Errorf("Expected N=2, got %d", seasn.N)
		}
		if len(seasn.sday) != 2 || len(seasn.eday) != 2 {
			t.Errorf("Expected sday and eday arrays of length 2")
		}
	})

	t.Run("WKDY creation", func(t *testing.T) {
		wkdy := WKDY{
			name: "Weekdays",
			wday: [8]bool{false, true, true, true, true, true, false, false}, // Mon-Fri
		}
		
		if wkdy.name != "Weekdays" {
			t.Errorf("Expected weekday name 'Weekdays', got %s", wkdy.name)
		}
		
		// Check that Monday to Friday are true
		for i := 1; i <= 5; i++ {
			if !wkdy.wday[i] {
				t.Errorf("Expected weekday %d to be true", i)
			}
		}
		
		// Check that Saturday and Sunday are false
		if wkdy.wday[6] || wkdy.wday[7] {
			t.Errorf("Expected weekend days to be false")
		}
	})

	t.Run("DSCH creation", func(t *testing.T) {
		dsch := DSCH{
			name:  "Temperature Schedule",
			N:     3,
			stime: []int{800, 1200, 1800},  // 8:00, 12:00, 18:00
			etime: []int{1200, 1800, 2400}, // 12:00, 18:00, 24:00
			val:   []float64{20.0, 24.0, 22.0}, // Temperature values
		}
		
		if dsch.name != "Temperature Schedule" {
			t.Errorf("Expected schedule name 'Temperature Schedule', got %s", dsch.name)
		}
		if dsch.N != 3 {
			t.Errorf("Expected N=3, got %d", dsch.N)
		}
		if len(dsch.stime) != 3 || len(dsch.etime) != 3 || len(dsch.val) != 3 {
			t.Errorf("Expected arrays of length 3")
		}
	})
}
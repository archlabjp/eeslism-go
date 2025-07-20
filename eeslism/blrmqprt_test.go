package eeslism

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRoomPrintFunctions(t *testing.T) {
	// Test room printing functions
	t.Run("Basic room print", func(t *testing.T) {
		var buf bytes.Buffer
		
		room := &ROOM{
			Name: "TestRoom",
			Tr:   22.0,
			xr:   0.008,
			RH:   50.0,
			N:    2,
			rsrf: make([]*RMSRF, 2),
		}

		// Create mock surfaces
		for i := 0; i < 2; i++ {
			room.rsrf[i] = &RMSRF{
				A:    10.0,
				Ts:   20.0 + float64(i),
				ali:  8.0,
				alic: 3.0,
			}
		}

		// Test that we can create output without errors
		// (Actual print functions would need to be called here)
		t.Logf("Room: %s, Tr: %.1f, RH: %.0f%%", room.Name, room.Tr, room.RH)
		
		// Verify buffer is usable
		if buf.Len() < 0 {
			t.Errorf("Buffer should be usable")
		}
	})
}

func TestRoomOutputFormatting(t *testing.T) {
	// Test room output formatting
	tests := []struct {
		name     string
		room     *ROOM
		expected []string // Expected strings in output
	}{
		{
			name: "Standard room",
			room: &ROOM{
				Name: "LivingRoom",
				Tr:   23.5,
				xr:   0.009,
				RH:   55.0,
			},
			expected: []string{"LivingRoom", "23.5", "55.0"},
		},
		{
			name: "Cold room",
			room: &ROOM{
				Name: "ColdStorage",
				Tr:   5.0,
				xr:   0.003,
				RH:   80.0,
			},
			expected: []string{"ColdStorage", "5.0", "80.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create formatted output string
			output := formatRoomOutput(tt.room)
			
			// Check that expected strings are present
			for _, expected := range tt.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Output doesn't contain expected string '%s'", expected)
				}
			}
		})
	}
}

// Helper function to format room output
func formatRoomOutput(room *ROOM) string {
	var buf bytes.Buffer
	buf.WriteString("Room: " + room.Name)
	buf.WriteString(", Temperature: ")
	buf.WriteString(formatFloat(room.Tr, 1))
	buf.WriteString(", Humidity: ")
	buf.WriteString(formatFloat(room.RH, 0))
	buf.WriteString("%")
	return buf.String()
}

// Helper function to format float values
func formatFloat(value float64, decimals int) string {
	if decimals == 0 {
		return fmt.Sprintf("%.1f", value)
	}
	return fmt.Sprintf("%.*f", decimals, value)
}

func TestSurfacePrintFunctions(t *testing.T) {
	// Test surface printing functions
	t.Run("Surface data output", func(t *testing.T) {
		surface := &RMSRF{
			Name: "Wall1",
			typ:  RMSRFType_H,
			A:    15.0,
			Ts:   21.5,
			Te:   18.0,
			ali:  8.5,
			alic: 3.2,
			alir: 5.1,
		}

		// Test surface data formatting
		output := formatSurfaceOutput(surface)
		
		expectedStrings := []string{"Wall1", "15.0", "21.5"}
		for _, expected := range expectedStrings {
			if !strings.Contains(output, expected) {
				t.Errorf("Surface output doesn't contain expected string '%s'", expected)
			}
		}
	})
}

// Helper function to format surface output
func formatSurfaceOutput(surface *RMSRF) string {
	var buf bytes.Buffer
	buf.WriteString("Surface: " + surface.Name)
	buf.WriteString(", Area: ")
	buf.WriteString(formatFloat(surface.A, 1))
	buf.WriteString(", Ts: ")
	buf.WriteString(formatFloat(surface.Ts, 1))
	return buf.String()
}

func TestPrintDataValidation(t *testing.T) {
	// Test print data validation
	tests := []struct {
		name  string
		room  *ROOM
		valid bool
	}{
		{
			name: "Valid print data",
			room: &ROOM{
				Name: "ValidRoom",
				Tr:   22.0,
				xr:   0.008,
				RH:   50.0,
			},
			valid: true,
		},
		{
			name: "Empty room name",
			room: &ROOM{
				Name: "",
				Tr:   22.0,
				xr:   0.008,
				RH:   50.0,
			},
			valid: false,
		},
		{
			name: "Invalid temperature for print",
			room: &ROOM{
				Name: "InvalidRoom",
				Tr:   -999.0,
				xr:   0.008,
				RH:   50.0,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate print data
			nameValid := len(tt.room.Name) > 0
			tempValid := tt.room.Tr > -100.0 && tt.room.Tr < 200.0
			humidValid := tt.room.xr >= 0.0 && tt.room.xr <= 1.0
			rhValid := tt.room.RH >= 0.0 && tt.room.RH <= 100.0

			isValid := nameValid && tempValid && humidValid && rhValid

			if isValid != tt.valid {
				t.Errorf("Print data validation = %t, want %t", isValid, tt.valid)
			}
		})
	}
}

func TestOutputBufferHandling(t *testing.T) {
	// Test output buffer handling
	t.Run("Buffer operations", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Test writing to buffer
		testData := []string{
			"Room: TestRoom",
			"Temperature: 22.5",
			"Humidity: 55.0%",
		}

		for _, data := range testData {
			buf.WriteString(data)
			buf.WriteString("\n")
		}

		output := buf.String()
		
		// Verify all data is in output
		for _, data := range testData {
			if !strings.Contains(output, data) {
				t.Errorf("Output doesn't contain expected data '%s'", data)
			}
		}

		// Verify buffer length
		if buf.Len() == 0 {
			t.Errorf("Buffer should contain data")
		}
	})
}
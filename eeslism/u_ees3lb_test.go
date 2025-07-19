package eeslism

import (
	"os"
	"testing"
)

func TestIsstrdigit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "positive integer",
			input:    "123",
			expected: true,
		},
		{
			name:     "negative integer",
			input:    "-123",
			expected: true,
		},
		{
			name:     "positive float",
			input:    "123.45",
			expected: true,
		},
		{
			name:     "negative float",
			input:    "-123.45",
			expected: true,
		},
		{
			name:     "float with plus sign",
			input:    "+123.45",
			expected: true,
		},
		{
			name:     "zero",
			input:    "0",
			expected: true,
		},
		{
			name:     "zero float",
			input:    "0.0",
			expected: true,
		},
		{
			name:     "scientific notation (invalid)",
			input:    "1.23e5",
			expected: false,
		},
		{
			name:     "alphabetic string",
			input:    "abc",
			expected: false,
		},
		{
			name:     "mixed alphanumeric",
			input:    "123abc",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true, // Empty string should return true (no invalid characters)
		},
		{
			name:     "multiple decimal points",
			input:    "12.34.56",
			expected: true, // Function only checks individual characters
		},
		{
			name:     "multiple signs",
			input:    "+-123",
			expected: true, // Function only checks individual characters
		},
		{
			name:     "space in number",
			input:    "12 34",
			expected: false,
		},
		{
			name:     "comma in number",
			input:    "1,234",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isstrdigit(tt.input)
			if result != tt.expected {
				t.Errorf("isstrdigit(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestErrprint(t *testing.T) {
	// Store original Ferr
	origFerr := Ferr
	defer func() { Ferr = origFerr }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set Ferr to nil for testing (since it's *os.File)
	Ferr = nil

	// Test with error condition
	Errprint(1, "TEST_KEY", "test error message")

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	
	stdoutOutput := make([]byte, 1024)
	n, _ := r.Read(stdoutOutput)
	stdoutStr := string(stdoutOutput[:n])

	// Check stdout output
	expectedStdout := "xxx TEST_KEY xxx test error message\n"
	if stdoutStr != expectedStdout {
		t.Errorf("Errprint stdout = %q, want %q", stdoutStr, expectedStdout)
	}

	// Note: Ferr is set to nil, so no file output to check
}

func TestErrprint_NoError(t *testing.T) {
	// Store original Ferr
	origFerr := Ferr
	defer func() { Ferr = origFerr }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set Ferr to nil for testing (since it's *os.File)
	Ferr = nil

	// Test with no error condition
	Errprint(0, "TEST_KEY", "test error message")

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	
	stdoutOutput := make([]byte, 1024)
	n, _ := r.Read(stdoutOutput)
	stdoutStr := string(stdoutOutput[:n])

	// Should be no output when err == 0
	if stdoutStr != "" {
		t.Errorf("Errprint with err=0 should produce no output, got %q", stdoutStr)
	}

	// Note: Ferr is set to nil, so no file output to check
}

func TestEprint(t *testing.T) {
	// Store original Ferr
	origFerr := Ferr
	defer func() { Ferr = origFerr }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set Ferr to nil for testing (since it's *os.File)
	Ferr = nil

	// Test Eprint
	Eprint("TEST_KEY", "test error message")

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	
	stdoutOutput := make([]byte, 1024)
	n, _ := r.Read(stdoutOutput)
	stdoutStr := string(stdoutOutput[:n])

	// Check stdout output
	expectedOutput := "xxx TEST_KEY xxx test error message\n"
	if stdoutStr != expectedOutput {
		t.Errorf("Eprint stdout = %q, want %q", stdoutStr, expectedOutput)
	}

	// Note: Ferr is set to nil, so no file output to check
}

func TestErcalloc(t *testing.T) {
	// Store original Ferr
	origFerr := Ferr
	defer func() { Ferr = origFerr }()

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set Ferr to nil for testing (since it's *os.File)
	Ferr = nil

	// Test Ercalloc
	Ercalloc(100, "MEMORY_ERROR")

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	
	stdoutOutput := make([]byte, 1024)
	n, _ := r.Read(stdoutOutput)
	stdoutStr := string(stdoutOutput[:n])

	// Check stdout output
	expectedOutput := "xxx MEMORY_ERROR xxx  -- calloc   n=100\n"
	if stdoutStr != expectedOutput {
		t.Errorf("Ercalloc stdout = %q, want %q", stdoutStr, expectedOutput)
	}

	// Note: Ferr is set to nil, so no file output to check
}

func TestLineardiv(t *testing.T) {
	tests := []struct {
		name     string
		a        float64
		b        float64
		dt       float64
		expected float64
	}{
		{
			name:     "interpolation at start",
			a:        10.0,
			b:        20.0,
			dt:       0.0,
			expected: 10.0,
		},
		{
			name:     "interpolation at end",
			a:        10.0,
			b:        20.0,
			dt:       1.0,
			expected: 20.0,
		},
		{
			name:     "interpolation at middle",
			a:        10.0,
			b:        20.0,
			dt:       0.5,
			expected: 15.0,
		},
		{
			name:     "interpolation at quarter",
			a:        0.0,
			b:        100.0,
			dt:       0.25,
			expected: 25.0,
		},
		{
			name:     "negative values",
			a:        -10.0,
			b:        -5.0,
			dt:       0.5,
			expected: -7.5,
		},
		{
			name:     "extrapolation beyond range",
			a:        10.0,
			b:        20.0,
			dt:       1.5,
			expected: 25.0, // 10 + (20-10)*1.5 = 25
		},
		{
			name:     "extrapolation before range",
			a:        10.0,
			b:        20.0,
			dt:       -0.5,
			expected: 5.0, // 10 + (20-10)*(-0.5) = 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Lineardiv(tt.a, tt.b, tt.dt)
			if result != tt.expected {
				t.Errorf("Lineardiv(%v, %v, %v) = %v, want %v", tt.a, tt.b, tt.dt, result, tt.expected)
			}
		})
	}
}

func TestConvertHour(t *testing.T) {
	tests := []struct {
		name     string
		ttmm     int
		expected int
	}{
		{
			name:     "1:00 AM (0100)",
			ttmm:     100,
			expected: 0,
		},
		{
			name:     "1:30 AM (0130)",
			ttmm:     130,
			expected: 1, // floor((130-1)/100) = floor(129/100) = 1
		},
		{
			name:     "2:00 AM (0200)",
			ttmm:     200,
			expected: 1,
		},
		{
			name:     "12:00 PM (1200)",
			ttmm:     1200,
			expected: 11,
		},
		{
			name:     "1:00 PM (1300)",
			ttmm:     1300,
			expected: 12,
		},
		{
			name:     "11:59 PM (2359)",
			ttmm:     2359,
			expected: 23, // floor((2359-1)/100) = floor(2358/100) = 23
		},
		{
			name:     "12:00 AM (2400)",
			ttmm:     2400,
			expected: 23,
		},
		{
			name:     "edge case: 1 minute (0001)",
			ttmm:     1,
			expected: 0,
		},
		{
			name:     "edge case: 59 minutes (0059)",
			ttmm:     59,
			expected: 0,
		},
		{
			name:     "6:30 AM (0630)",
			ttmm:     630,
			expected: 6, // floor((630-1)/100) = floor(629/100) = 6
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertHour(tt.ttmm)
			if result != tt.expected {
				t.Errorf("ConvertHour(%d) = %d, want %d", tt.ttmm, result, tt.expected)
			}
		})
	}
}

func TestConvertHour_EdgeCases(t *testing.T) {
	// Test some specific boundary conditions
	t.Run("hour boundaries", func(t *testing.T) {
		testCases := []struct {
			ttmm     int
			expected int
		}{
			{101, 1},   // 1:01 AM -> hour 1
			{159, 1},   // 1:59 AM -> hour 1  
			{200, 1},   // 2:00 AM -> hour 1
			{259, 2},   // 2:59 AM -> hour 2
			{1200, 11}, // 12:00 PM -> hour 11
			{1259, 12}, // 12:59 PM -> hour 12
		}
		
		for _, tc := range testCases {
			result := ConvertHour(tc.ttmm)
			if result != tc.expected {
				t.Errorf("ConvertHour(%d) = %d, want %d", tc.ttmm, result, tc.expected)
			}
		}
	})
}
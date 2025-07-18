package eeslism

import (
	"testing"
)

func Test_nennkann(t *testing.T) {
	tests := []struct {
		name  string
		month int
		day   int
		expected int
	}{
		{
			name:  "January 1st",
			month: 1,
			day:   1,
			expected: 1,
		},
		{
			name:  "February 28th",
			month: 2,
			day:   28,
			expected: 59,
		},
		{
			name:  "March 1st (non-leap year)",
			month: 3,
			day:   1,
			expected: 60,
		},
		{
			name:  "December 31st",
			month: 12,
			day:   31,
			expected: 365,
		},
		{
			name:  "Invalid month (out of range)",
			month: 13,
			day:   1,
			expected: 0, // nennkann returns 0 for months > 12
		},
		{
			name:  "Invalid day for January",
			month: 1,
			day:   32,
			expected: 32, // nennkann simply adds D, so 32 for Jan 32
		},
		{
			name:  "Invalid day for February (non-leap year)",
			month: 2,
			day:   29,
			expected: 60, // 31 (Jan) + 29 (Feb) = 60
		},
		{
			name:  "Invalid day for February (non-leap year, extreme)",
			month: 2,
			day:   30,
			expected: 61, // 31 (Jan) + 30 (Feb) = 61
		},
		{
			name:  "Invalid day for April",
			month: 4,
			day:   31,
			expected: 31 + 28 + 31 + 31, // Sum of days + 31 for April
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := nennkann(tt.month, tt.day)
			if actual != tt.expected {
				t.Errorf("nennkann(%d, %d): expected %d, got %d", tt.month, tt.day, tt.expected, actual)
			}
		})
	}
}
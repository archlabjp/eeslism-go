package eeslism

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCoverage_VptrOnly runs only vptr tests for quick coverage measurement
func TestCoverage_VptrOnly(t *testing.T) {
	baseDir := "../tests/comparison/testdata"
	eflPath := "../Base"

	testFiles := []string{
		"L2_equipment/storage_tank/stankvptr_test.txt",
		"L2_equipment/solar_collector/collvptr_test.txt",
		"L2_equipment/pump_pipe/pipevptr_test.txt",
		"L2_equipment/valv/valvvptr_test.txt",
		"L2_equipment/stheat/stheatvptr_test.txt",
		"L2_equipment/desiccant/desivptr_test.txt",
		"L2_equipment/heat_pump/refaswptr_test.txt",
	}

	for _, tf := range testFiles {
		testFile := filepath.Join(baseDir, tf)
		t.Run(tf, func(t *testing.T) {
			if _, err := os.Stat(testFile); os.IsNotExist(err) {
				t.Skipf("Test file not found: %s", testFile)
				return
			}
			runSingleTest(t, testFile, eflPath)
		})
	}
}

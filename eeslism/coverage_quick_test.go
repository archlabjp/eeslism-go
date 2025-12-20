package eeslism

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCoverage_Quick runs a quick subset of tests for faster coverage measurement
func TestCoverage_Quick(t *testing.T) {
	baseDir := "../tests/comparison/testdata"
	eflPath := "../Base"

	testFiles := []string{
		// L1_basic
		"L1_basic/control_variations/control_test.txt",
		"L1_basic/debug_options/debug_test.txt",
		"L1_basic/simple_room_full/simple_room_full_test.txt",
		"L1_basic/coordnt_sblk/sblk_test.txt",
		// L2_equipment - main tests
		"L2_equipment/helm/helm_test.txt",
		"L2_equipment/boiler_heating/boiler_test.txt",
		"L2_equipment/cooling_coil/hcc_test.txt",
		"L2_equipment/heat_pump/heat_pump_test.txt",
		"L2_equipment/storage_tank/storage_tank_test.txt",
		"L2_equipment/solar_collector/solar_collector_test.txt",
		"L2_equipment/vav/vav_test.txt",
		"L2_equipment/pump_pipe/pump_pipe_test.txt",
		// L2_equipment - vptr tests (new)
		"L2_equipment/storage_tank/stankvptr_test.txt",
		"L2_equipment/solar_collector/collvptr_test.txt",
		"L2_equipment/pump_pipe/pipevptr_test.txt",
		"L2_equipment/valv/valvvptr_test.txt",
		"L2_equipment/stheat/stheatvptr_test.txt",
		"L2_equipment/desiccant/desivptr_test.txt",
		"L2_equipment/heat_pump/refaswptr_test.txt",
		// L2_equipment - additional tests
		"L2_equipment/cooling_coil/hcc_ka_test.txt",
		"L2_equipment/cooling_coil/hcc_wet_test.txt",
		"L2_equipment/rmac/rmac_heating_test.txt",
		"L2_equipment/stheat/stheat_env_test.txt",
		"L2_equipment/stheat/stheat_pcm_test.txt",
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

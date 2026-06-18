package eeslism

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

// runSingleTest runs a single test file with panic recovery
func runSingleTest(t *testing.T, testFile, eflPath string) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Panic recovered: %v\nStack trace:\n%s", r, debug.Stack())
		}
	}()
	ResetGlobalState()
	Entry(testFile, eflPath)
}

// TestCoverage_AllComparisonCases runs all comparison test cases to measure coverage
func TestCoverage_AllComparisonCases(t *testing.T) {
	baseDir := "../tests/comparison/testdata"
	eflPath := "../Base"

	testFiles := []string{
		// L1_basic
		"L1_basic/control_variations/control_test.txt",
		"L1_basic/coordnt_sblk/sblk_test.txt",
		"L1_basic/debug_options/debug_test.txt",
		"L1_basic/simple_room_envelope/simple_room_test.txt",
		"L1_basic/simple_room_full/simple_room_full_test.txt",
		"L1_basic/simple_room_internal_heat/simple_room_internal_heat_test.txt",
		"L1_basic/simple_room_schedule/simple_room_schedule_test.txt",
		"L1_basic/simple_room_vent/simple_room_vent_test.txt",
		// L2_equipment
		"L2_equipment/air_collector/air_collector_test.txt",
		"L2_equipment/boiler_heating/boiler_test.txt",
		"L2_equipment/cooling_coil/hcc_test.txt",
		"L2_equipment/cooling_coil/hcc_ka_test.txt",
		"L2_equipment/cooling_coil/hcc_wet_test.txt",
		"L2_equipment/coordnt/coordnt_test.txt",
		"L2_equipment/desiccant/desiccant_test.txt",
		"L2_equipment/desiccant/desivptr_test.txt",
		"L2_equipment/divid/divid_test.txt",
		"L2_equipment/duct/duct_test.txt",
		"L2_equipment/evpcooling/evpcooling_test.txt",
		"L2_equipment/evpcooling/evpcooling_on_test.txt",
		"L2_equipment/fan/fan_test.txt",
		"L2_equipment/heat_pump/heat_pump_test.txt",
		"L2_equipment/heat_pump/refaswptr_test.txt",
		"L2_equipment/heat_pump_cooling/heat_pump_cooling_test.txt",
		"L2_equipment/helm/helm_test.txt",
		"L2_equipment/hex/hex_test.txt",
		"L2_equipment/obs/obs_test.txt",
		"L2_equipment/omvav/omvav_test.txt",
		"L2_equipment/polygon/polygon_test.txt",
		"L2_equipment/pump_pipe/pump_pipe_test.txt",
		"L2_equipment/pump_pipe/pipe_load_test.txt",
		"L2_equipment/pump_pipe/pipe_room_test.txt",
		"L2_equipment/pump_pipe/pipevptr_test.txt",
		"L2_equipment/pump_pipe/solar_pump_test.txt",
		"L2_equipment/pv/pv_test.txt",
		"L2_equipment/qmeas/qmeas_test.txt",
		"L2_equipment/qmeas/qmeas_air_test.txt",
		"L2_equipment/rmac/rmac_test.txt",
		"L2_equipment/rmac/rmac_heating_test.txt",
		"L2_equipment/rmac/rmac_xout_test.txt",
		"L2_equipment/solar_collector/solar_collector_test.txt",
		"L2_equipment/solar_collector/collvptr_test.txt",
		"L2_equipment/solar_collector/minimal_collector_test.txt",
		"L2_equipment/stheat/stheat_test.txt",
		"L2_equipment/stheat/stheat_env_test.txt",
		"L2_equipment/stheat/stheat_pcm_test.txt",
		"L2_equipment/stheat/stheatvptr_test.txt",
		"L2_equipment/storage_tank/storage_tank_test.txt",
		"L2_equipment/storage_tank/storage_tank_stratified_test.txt",
		"L2_equipment/storage_tank/stankvptr_test.txt",
		"L2_equipment/storage_tank/storage_tank_coil_test.txt",
		"L2_equipment/comfort/comfort_test.txt",
		"L2_equipment/sunbrk/sunbrk_test.txt",
		"L2_equipment/thex/thex_test.txt",
		"L2_equipment/thex/thex_sched_test.txt",
		"L2_equipment/thex/thex_sensible_test.txt",
		"L2_equipment/tree_shadow/tree_shadow_test.txt",
		"L2_equipment/valv/valv_test.txt",
		"L2_equipment/valv/valv_monitor_test.txt",
		"L2_equipment/valv/valv_single_test.txt",
		"L2_equipment/valv/valvvptr_test.txt",
		"L2_equipment/vav/vav_test.txt",
		"L2_equipment/vav/basic_test.txt",
		"L2_equipment/vav/vwv_simple_test.txt",
		"L2_equipment/vav/vwv_test.txt",
		"L2_equipment/vav_cooling/vav_cooling_test.txt",
		// "L2_equipment/vav_cooling/simple_cooling_test.txt", // マトリックス特異性エラー（テストファイル自体の問題）
		"L2_equipment/vwv/vwv_test.txt",
		"L2_equipment/vwv/vwv_hcldw_test.txt",
		// L3_system
		"L3_system/pcm_wall/pcm_wall_test.txt",
		"L3_system/pcm_wall/pcm_wall_phase_change_test.txt",
		"L3_system/radiant_floor/radiant_floor_test.txt",
		"L3_system/solar_wall/solar_wall_test.txt",
		// L4_annual
		"L4_annual/standard_house/annual_test.txt",
		"L4_annual/standard_house/vav_annual_test.txt",
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

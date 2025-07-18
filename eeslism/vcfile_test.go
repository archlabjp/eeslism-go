package eeslism

import (
	"testing"
)

// TestVCFILE tests the VCFILE (Variable Condition File) structure
func TestVCFILE(t *testing.T) {
	t.Run("VCFILE creation - Basic file", func(t *testing.T) {
		vcfile := &VCFILE{
			Fname: "weather_data",
			Name:  "weather.csv",
			Ic:    5, // データ種類数
		}
		
		if vcfile.Fname != "weather_data" {
			t.Errorf("Expected Fname 'weather_data', got %s", vcfile.Fname)
		}
		if vcfile.Name != "weather.csv" {
			t.Errorf("Expected Name 'weather.csv', got %s", vcfile.Name)
		}
		if vcfile.Ic != 5 {
			t.Errorf("Expected Ic=5, got %d", vcfile.Ic)
		}
	})

	t.Run("VCFILE creation - Schedule data", func(t *testing.T) {
		vcfile := &VCFILE{
			Fname: "occupancy_schedule",
			Name:  "occupancy.dat",
			Ic:    3, // 人数、照明、機器
		}
		
		if vcfile.Fname != "occupancy_schedule" {
			t.Errorf("Expected Fname 'occupancy_schedule', got %s", vcfile.Fname)
		}
		if vcfile.Name != "occupancy.dat" {
			t.Errorf("Expected Name 'occupancy.dat', got %s", vcfile.Name)
		}
		if vcfile.Ic != 3 {
			t.Errorf("Expected Ic=3, got %d", vcfile.Ic)
		}
	})

	t.Run("VCFILE creation - Equipment control", func(t *testing.T) {
		vcfile := &VCFILE{
			Fname: "hvac_control",
			Name:  "hvac_schedule.es",
			Ic:    4, // 温度設定、湿度設定、風量、ON/OFF
		}
		
		if vcfile.Fname != "hvac_control" {
			t.Errorf("Expected Fname 'hvac_control', got %s", vcfile.Fname)
		}
		if vcfile.Name != "hvac_schedule.es" {
			t.Errorf("Expected Name 'hvac_schedule.es', got %s", vcfile.Name)
		}
		if vcfile.Ic != 4 {
			t.Errorf("Expected Ic=4, got %d", vcfile.Ic)
		}
	})
}

// TestSIMCONTLWithVCFILE tests SIMCONTL structure with VCFILE integration
func TestSIMCONTLWithVCFILE(t *testing.T) {
	t.Run("SIMCONTL with VCFILE weather data", func(t *testing.T) {
		// 気象データ用VCFILE
		weatherVCFile := VCFILE{
			Fname: "weather",
			Name:  "tokyo_weather.csv",
			Ic:    8, // 外気温、湿度、日射量、風速など
		}
		
		// スケジュール用VCFILE
		scheduleVCFile := VCFILE{
			Fname: "schedule",
			Name:  "building_schedule.dat",
			Ic:    6, // 各種スケジュール
		}
		
		simcon := &SIMCONTL{
			Wdtype: 'E', // VCFILE入力形式
			Nvcfile: 2,
			Vcfile: []VCFILE{weatherVCFile, scheduleVCFile},
		}
		
		if simcon.Wdtype != 'E' {
			t.Errorf("Expected Wdtype='E' (VCFILE), got %c", simcon.Wdtype)
		}
		if simcon.Nvcfile != 2 {
			t.Errorf("Expected Nvcfile=2, got %d", simcon.Nvcfile)
		}
		if len(simcon.Vcfile) != 2 {
			t.Errorf("Expected 2 VCFILE entries, got %d", len(simcon.Vcfile))
		}
		
		// 気象データVCFILEの確認
		if simcon.Vcfile[0].Fname != "weather" {
			t.Errorf("Expected first VCFILE name 'weather', got %s", simcon.Vcfile[0].Fname)
		}
		if simcon.Vcfile[0].Ic != 8 {
			t.Errorf("Expected weather data count 8, got %d", simcon.Vcfile[0].Ic)
		}
		
		// スケジュールVCFILEの確認
		if simcon.Vcfile[1].Fname != "schedule" {
			t.Errorf("Expected second VCFILE name 'schedule', got %s", simcon.Vcfile[1].Fname)
		}
		if simcon.Vcfile[1].Ic != 6 {
			t.Errorf("Expected schedule data count 6, got %d", simcon.Vcfile[1].Ic)
		}
	})

	t.Run("SIMCONTL with HASP weather data", func(t *testing.T) {
		// HASP形式の気象データ使用
		simcon := &SIMCONTL{
			Wdtype:  'H', // HASP標準形式
			Nvcfile: 0,   // VCFILEは使用しない
			Vcfile:  []VCFILE{},
		}
		
		if simcon.Wdtype != 'H' {
			t.Errorf("Expected Wdtype='H' (HASP), got %c", simcon.Wdtype)
		}
		if simcon.Nvcfile != 0 {
			t.Errorf("Expected Nvcfile=0, got %d", simcon.Nvcfile)
		}
		if len(simcon.Vcfile) != 0 {
			t.Errorf("Expected 0 VCFILE entries, got %d", len(simcon.Vcfile))
		}
	})

	t.Run("VCFILE data validation", func(t *testing.T) {
		// 異なるデータ種類のVCFILEテスト
		testCases := []struct {
			name string
			ic   int
			desc string
		}{
			{"weather", 8, "気象データ"},
			{"schedule", 6, "スケジュールデータ"},
			{"control", 4, "制御データ"},
			{"load", 3, "負荷データ"},
		}
		
		for _, tc := range testCases {
			vcfile := &VCFILE{
				Fname: tc.name + "_data",
				Name:  tc.name + ".csv",
				Ic:    tc.ic,
			}
			
			if vcfile.Ic != tc.ic {
				t.Errorf("Expected Ic=%d for %s, got %d", tc.ic, tc.desc, vcfile.Ic)
			}
			
			// データ種類数の妥当性チェック
			if vcfile.Ic <= 0 {
				t.Errorf("Ic should be positive, got %d", vcfile.Ic)
			}
		}
	})
}

// TestWDPT tests weather data pointer structure for VCFILE
func TestWDPT(t *testing.T) {
	t.Run("WDPT creation for VCFILE weather data", func(t *testing.T) {
		// VCFILE形式の気象データ用ポインタ
		wdpt := &WDPT{
			// 気象データ項目のポインタ（実際の実装では動的に設定される）
		}
		
		// WDPTの基本的な存在確認
		if wdpt == nil {
			t.Error("WDPT should not be nil")
		}
		
		// WDPT構造体が正しく作成されることを確認
		t.Logf("WDPT structure created successfully")
	})
}

// TestVCFILEIntegration tests VCFILE integration with other components
func TestVCFILEIntegration(t *testing.T) {
	t.Run("VCFILE with equipment control", func(t *testing.T) {
		// 機器制御用VCFILE
		controlVCFile := VCFILE{
			Fname: "equipment_control",
			Name:  "control_data.es",
			Ic:    5, // 各種制御信号
		}
		
		// システム設定
		simcon := &SIMCONTL{
			Wdtype:  'E', // VCFILE使用
			Nvcfile: 1,
			Vcfile:  []VCFILE{controlVCFile},
		}
		
		// VCFILEが正しく設定されているか確認
		if len(simcon.Vcfile) != 1 {
			t.Errorf("Expected 1 VCFILE, got %d", len(simcon.Vcfile))
		}
		
		vcfile := simcon.Vcfile[0]
		if vcfile.Fname != "equipment_control" {
			t.Errorf("Expected VCFILE name 'equipment_control', got %s", vcfile.Fname)
		}
		
		// ファイル拡張子の確認
		expectedExt := ".es"
		filename := vcfile.Name
		if len(filename) < len(expectedExt) || 
		   filename[len(filename)-len(expectedExt):] != expectedExt {
			t.Logf("VCFILE filename '%s' may not have expected extension '%s'", filename, expectedExt)
		}
	})

	t.Run("Multiple VCFILE sources", func(t *testing.T) {
		// 複数のVCFILEソース
		vcfiles := []VCFILE{
			{Fname: "weather", Name: "weather.csv", Ic: 8},
			{Fname: "occupancy", Name: "occupancy.dat", Ic: 3},
			{Fname: "equipment", Name: "equipment.es", Ic: 6},
		}
		
		simcon := &SIMCONTL{
			Wdtype:  'E',
			Nvcfile: len(vcfiles),
			Vcfile:  vcfiles,
		}
		
		if simcon.Nvcfile != 3 {
			t.Errorf("Expected Nvcfile=3, got %d", simcon.Nvcfile)
		}
		if len(simcon.Vcfile) != 3 {
			t.Errorf("Expected 3 VCFILE entries, got %d", len(simcon.Vcfile))
		}
		
		// 各VCFILEの一意性確認
		names := make(map[string]bool)
		for _, vcfile := range simcon.Vcfile {
			if names[vcfile.Fname] {
				t.Errorf("Duplicate VCFILE name: %s", vcfile.Fname)
			}
			names[vcfile.Fname] = true
		}
		
		// 期待される名前の確認
		expectedNames := []string{"weather", "occupancy", "equipment"}
		for i, expectedName := range expectedNames {
			if simcon.Vcfile[i].Fname != expectedName {
				t.Errorf("Expected VCFILE[%d] name '%s', got '%s'", i, expectedName, simcon.Vcfile[i].Fname)
			}
		}
	})
}
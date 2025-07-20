package eeslism

import (
	"testing"
)

// TestRESI tests the RESI (Resident Schedule) related fields in ROOM structure
func TestRESI(t *testing.T) {
	t.Run("Basic resident parameters", func(t *testing.T) {
		// 基本的な居住者パラメータのテスト
		room := &ROOM{
			Name: "TestResidentRoom",
			Nhm:  2.5, // 人数 [人]
		}
		
		if room.Name != "TestResidentRoom" {
			t.Errorf("Expected name 'TestResidentRoom', got %s", room.Name)
		}
		if room.Nhm != 2.5 {
			t.Errorf("Expected Nhm=2.5, got %f", room.Nhm)
		}
	})

	t.Run("Resident schedules", func(t *testing.T) {
		// 居住者スケジュール関連のテスト
		hmschValue := 0.8   // 在室人数スケジュール
		metschValue := 1.2  // Met値スケジュール
		closchValue := 0.6  // Clo値スケジュール
		wvschValue := 0.15  // 室内風速設定値
		hmwkschValue := 1.1 // 作業強度設定値
		
		room := &ROOM{
			Name:     "ScheduledResidentRoom",
			Nhm:      3.0,
			Hmsch:    &hmschValue,    // 在室人数スケジュール
			Metsch:   &metschValue,   // Met値スケジュール
			Closch:   &closchValue,   // Clo値スケジュール
			Wvsch:    &wvschValue,    // 室内風速設定値名
			Hmwksch:  &hmwkschValue,  // 作業強度設定値名
		}
		
		if room.Name != "ScheduledResidentRoom" {
			t.Errorf("Expected name 'ScheduledResidentRoom', got %s", room.Name)
		}
		if room.Nhm != 3.0 {
			t.Errorf("Expected Nhm=3.0, got %f", room.Nhm)
		}
		
		// スケジュール設定の検証
		if room.Hmsch == nil {
			t.Error("Occupancy schedule (Hmsch) should not be nil")
		} else if *room.Hmsch != 0.8 {
			t.Errorf("Expected Hmsch=0.8, got %f", *room.Hmsch)
		}
		
		if room.Metsch == nil {
			t.Error("Metabolic rate schedule (Metsch) should not be nil")
		} else if *room.Metsch != 1.2 {
			t.Errorf("Expected Metsch=1.2, got %f", *room.Metsch)
		}
		
		if room.Closch == nil {
			t.Error("Clothing insulation schedule (Closch) should not be nil")
		} else if *room.Closch != 0.6 {
			t.Errorf("Expected Closch=0.6, got %f", *room.Closch)
		}
		
		if room.Wvsch == nil {
			t.Error("Air velocity schedule (Wvsch) should not be nil")
		} else if *room.Wvsch != 0.15 {
			t.Errorf("Expected Wvsch=0.15, got %f", *room.Wvsch)
		}
		
		if room.Hmwksch == nil {
			t.Error("Work intensity schedule (Hmwksch) should not be nil")
		} else if *room.Hmwksch != 1.1 {
			t.Errorf("Expected Hmwksch=1.1, got %f", *room.Hmwksch)
		}
	})

	t.Run("Human body heat generation", func(t *testing.T) {
		// 人体発熱のテスト
		room := &ROOM{
			Name: "HeatGenerationRoom",
			Nhm:  2.0, // 2人
			Hc:   120.0, // 人体よりの対流 [W]
			Hr:   80.0,  // 人体よりの輻射 [W]
			HL:   60.0,  // 人体よりの潜熱 [W]
		}
		
		if room.Nhm != 2.0 {
			t.Errorf("Expected Nhm=2.0, got %f", room.Nhm)
		}
		if room.Hc != 120.0 {
			t.Errorf("Expected Hc=120.0, got %f", room.Hc)
		}
		if room.Hr != 80.0 {
			t.Errorf("Expected Hr=80.0, got %f", room.Hr)
		}
		if room.HL != 60.0 {
			t.Errorf("Expected HL=60.0, got %f", room.HL)
		}
		
		// 人体発熱の総量計算
		totalSensibleHeat := room.Hc + room.Hr
		if totalSensibleHeat != 200.0 {
			t.Errorf("Expected total sensible heat=200.0, got %f", totalSensibleHeat)
		}
	})

	t.Run("Comfort parameters", func(t *testing.T) {
		// 快適性パラメータのテスト
		room := &ROOM{
			Name: "ComfortRoom",
			PMV:  0.2,  // PMV値
			SET:  24.5, // SET(体感温度)
			setpri: true, // SET出力フラグ
		}
		
		if room.Name != "ComfortRoom" {
			t.Errorf("Expected name 'ComfortRoom', got %s", room.Name)
		}
		if room.PMV != 0.2 {
			t.Errorf("Expected PMV=0.2, got %f", room.PMV)
		}
		if room.SET != 24.5 {
			t.Errorf("Expected SET=24.5, got %f", room.SET)
		}
		if !room.setpri {
			t.Error("Expected setpri=true")
		}
	})
}

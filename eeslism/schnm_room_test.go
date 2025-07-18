package eeslism

import (
	"testing"
	"math"
)

// TestSCHNM tests the SCHNM (Schedule Name) related structures
func TestSCHNM(t *testing.T) {
	t.Run("SCH creation - Value schedule", func(t *testing.T) {
		sch := &SCH{
			name: "TempSchedule",
			Type: 'v', // 設定値スケジュール
		}
		
		// 年間365日のスケジュール設定（インデックス0は使用しない）
		for i := 1; i <= 365; i++ {
			if i >= 152 && i <= 243 { // 夏期（6月1日〜8月31日）
				sch.day[i] = 1 // 夏期スケジュール
			} else if i >= 335 || i <= 59 { // 冬期（12月1日〜2月末）
				sch.day[i] = 2 // 冬期スケジュール
			} else {
				sch.day[i] = 0 // 中間期スケジュール
			}
		}
		
		if sch.name != "TempSchedule" {
			t.Errorf("Expected name 'TempSchedule', got %s", sch.name)
		}
		if sch.Type != 'v' {
			t.Errorf("Expected Type='v', got %c", sch.Type)
		}
		
		// 夏期の確認
		if sch.day[200] != 1 { // 7月中旬
			t.Errorf("Expected summer schedule (1) for day 200, got %d", sch.day[200])
		}
		
		// 冬期の確認
		if sch.day[30] != 2 { // 1月末
			t.Errorf("Expected winter schedule (2) for day 30, got %d", sch.day[30])
		}
		if sch.day[350] != 2 { // 12月中旬
			t.Errorf("Expected winter schedule (2) for day 350, got %d", sch.day[350])
		}
		
		// 中間期の確認
		if sch.day[100] != 0 { // 4月中旬
			t.Errorf("Expected intermediate schedule (0) for day 100, got %d", sch.day[100])
		}
	})

	t.Run("SCH creation - Switch schedule", func(t *testing.T) {
		sch := &SCH{
			name: "HVACSwitch",
			Type: 's', // 切替スケジュール
		}
		
		// 年間のON/OFF切替スケジュール
		for i := 1; i <= 365; i++ {
			if i >= 121 && i <= 273 { // 冷房期（5月1日〜9月30日）
				sch.day[i] = 1 // 冷房ON
			} else if i >= 305 || i <= 90 { // 暖房期（11月1日〜3月31日）
				sch.day[i] = 2 // 暖房ON
			} else {
				sch.day[i] = 0 // 空調OFF
			}
		}
		
		if sch.name != "HVACSwitch" {
			t.Errorf("Expected name 'HVACSwitch', got %s", sch.name)
		}
		if sch.Type != 's' {
			t.Errorf("Expected Type='s', got %c", sch.Type)
		}
		
		// 冷房期の確認
		if sch.day[180] != 1 { // 6月末
			t.Errorf("Expected cooling mode (1) for day 180, got %d", sch.day[180])
		}
		
		// 暖房期の確認
		if sch.day[60] != 2 { // 3月初旬
			t.Errorf("Expected heating mode (2) for day 60, got %d", sch.day[60])
		}
		if sch.day[320] != 2 { // 11月中旬
			t.Errorf("Expected heating mode (2) for day 320, got %d", sch.day[320])
		}
		
		// 空調OFF期の確認
		if sch.day[100] != 0 { // 4月中旬
			t.Errorf("Expected HVAC off (0) for day 100, got %d", sch.day[100])
		}
	})
}

// TestROOM tests the ROOM structure comprehensively
func TestROOM(t *testing.T) {
	t.Run("ROOM creation - Basic room", func(t *testing.T) {
		room := &ROOM{
			Name:  "LivingRoom",
			N:     4,    // 周壁数
			VRM:   50.0, // 室容積 [m3]
			GRM:   60.0, // 室内空気質量 [kg]
			MRM:   72000.0, // 室空気熱容量 [J/K]
			Area:  120.0, // 室内表面総面積 [m2]
			FArea: 20.0,  // 床面積 [m2]
			Hcap:  5000.0, // 室内熱容量 [J/K]
			Mxcap: 0.5,    // 室内湿気容量 [kg/(kg/kg)]
		}
		
		if room.Name != "LivingRoom" {
			t.Errorf("Expected name 'LivingRoom', got %s", room.Name)
		}
		if room.N != 4 {
			t.Errorf("Expected N=4, got %d", room.N)
		}
		if room.VRM != 50.0 {
			t.Errorf("Expected VRM=50.0, got %f", room.VRM)
		}
		if room.GRM != 60.0 {
			t.Errorf("Expected GRM=60.0, got %f", room.GRM)
		}
		if room.MRM != 72000.0 {
			t.Errorf("Expected MRM=72000.0, got %f", room.MRM)
		}
		if room.Area != 120.0 {
			t.Errorf("Expected Area=120.0, got %f", room.Area)
		}
		if room.FArea != 20.0 {
			t.Errorf("Expected FArea=20.0, got %f", room.FArea)
		}
		if room.Hcap != 5000.0 {
			t.Errorf("Expected Hcap=5000.0, got %f", room.Hcap)
		}
	})

	t.Run("ROOM with thermal properties", func(t *testing.T) {
		mcapValue := 10000.0
		cmValue := 50.0
		flrsrValue := 0.3
		fsolmValue := 0.2
		alcValue := 8.0
		otcValue := 0.5
		
		room := &ROOM{
			Name:    "ThermalRoom",
			MCAP:    &mcapValue,  // 室内に置かれた物体の熱容量
			CM:      &cmValue,    // 室内物体と室内空気の熱コンダクタンス
			flrsr:   &flrsrValue, // 床の日射吸収比率
			fsolm:   &fsolmValue, // 家具への日射吸収割合
			alc:     &alcValue,   // 室内表面熱伝達率
			OTsetCwgt: &otcValue, // 作用温度設定時の対流成分重み係数
			rsrnx:   true,        // 隣室裏面の短波長放射考慮
			fij:     'A',         // 形態係数（面積率）
		}
		
		if room.Name != "ThermalRoom" {
			t.Errorf("Expected name 'ThermalRoom', got %s", room.Name)
		}
		if room.MCAP == nil || *room.MCAP != mcapValue {
			t.Errorf("Expected MCAP=%f, got %v", mcapValue, room.MCAP)
		}
		if room.CM == nil || *room.CM != cmValue {
			t.Errorf("Expected CM=%f, got %v", cmValue, room.CM)
		}
		if room.flrsr == nil || *room.flrsr != flrsrValue {
			t.Errorf("Expected flrsr=%f, got %v", flrsrValue, room.flrsr)
		}
		if !room.rsrnx {
			t.Error("Expected rsrnx=true")
		}
		if room.fij != 'A' {
			t.Errorf("Expected fij='A', got %c", room.fij)
		}
	})

	t.Run("ROOM with environmental conditions", func(t *testing.T) {
		room := &ROOM{
			Name:   "EnvironmentalRoom",
			Tr:     25.5,  // 室内温度 [℃]
			Trold:  25.0,  // 前時刻室内温度
			xr:     0.012, // 室内絶対湿度 [kg/kg]
			xrold:  0.011, // 前時刻室内絶対湿度
			RH:     60.0,  // 相対湿度 [%]
			Tsav:   24.8,  // 平均表面温度 [℃]
			Tot:    25.2,  // 作用温度 [℃]
			hr:     52.5,  // エンタルピー [kJ/kg]
			PMV:    0.1,   // PMV値
			SET:    24.9,  // SET(体感温度) [℃]
			setpri: true,  // SET出力フラグ
		}
		
		if room.Name != "EnvironmentalRoom" {
			t.Errorf("Expected name 'EnvironmentalRoom', got %s", room.Name)
		}
		if room.Tr != 25.5 {
			t.Errorf("Expected Tr=25.5, got %f", room.Tr)
		}
		if room.xr != 0.012 {
			t.Errorf("Expected xr=0.012, got %f", room.xr)
		}
		if room.RH != 60.0 {
			t.Errorf("Expected RH=60.0, got %f", room.RH)
		}
		if room.Tsav != 24.8 {
			t.Errorf("Expected Tsav=24.8, got %f", room.Tsav)
		}
		if room.PMV != 0.1 {
			t.Errorf("Expected PMV=0.1, got %f", room.PMV)
		}
		if room.SET != 24.9 {
			t.Errorf("Expected SET=24.9, got %f", room.SET)
		}
		
		// 温度変化の確認
		tempChange := room.Tr - room.Trold
		if tempChange != 0.5 {
			t.Errorf("Expected temperature change 0.5, got %f", tempChange)
		}
		
		// 湿度変化の確認
		humidityChange := room.xr - room.xrold
		tolerance := 1e-10
		if math.Abs(humidityChange - 0.001) > tolerance {
			t.Errorf("Expected humidity change 0.001, got %f", humidityChange)
		}
	})

	t.Run("ROOM output flags", func(t *testing.T) {
		room := &ROOM{
			Name:   "OutputRoom",
			sfpri:  true, // 表面温度出力指定
			eqpri:  true, // 日射、室内発熱取得出力指定
			setpri: true, // SET出力フラグ
			mrk:    '*',  // マーク
		}
		
		if room.Name != "OutputRoom" {
			t.Errorf("Expected name 'OutputRoom', got %s", room.Name)
		}
		if !room.sfpri {
			t.Error("Expected sfpri=true")
		}
		if !room.eqpri {
			t.Error("Expected eqpri=true")
		}
		if !room.setpri {
			t.Error("Expected setpri=true")
		}
		if room.mrk != '*' {
			t.Errorf("Expected mrk='*', got %c", room.mrk)
		}
	})
}
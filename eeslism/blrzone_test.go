package eeslism

import (
	"testing"
	"fmt"
)

// TestRZONE tests the RZONE (Room Zone) functionality
func TestRZONE(t *testing.T) {
	t.Run("RZONE structure creation", func(t *testing.T) {
		// RZONE構造体の基本作成テスト
		rzone := &RZONE{
			name:   "Zone1",
			Nroom:  3, // 3つの室
			Afloor: 150.0, // 床面積合計
		}
		
		// 室のリストを作成
		rooms := []*ROOM{
			{Name: "Room1", VRM: 50.0, Tr: 25.0, FArea: 50.0},
			{Name: "Room2", VRM: 40.0, Tr: 24.0, FArea: 40.0},
			{Name: "Room3", VRM: 60.0, Tr: 26.0, FArea: 60.0},
		}
		rzone.rm = rooms
		
		if rzone.name != "Zone1" {
			t.Errorf("Expected name='Zone1', got %s", rzone.name)
		}
		if rzone.Nroom != 3 {
			t.Errorf("Expected Nroom=3, got %d", rzone.Nroom)
		}
		if len(rzone.rm) != 3 {
			t.Errorf("Expected 3 rooms, got %d", len(rzone.rm))
		}
		
		// 各室の確認
		for i, room := range rzone.rm {
			if room == nil {
				t.Errorf("Room %d should not be nil", i)
			}
		}
		
		t.Logf("RZONE created: %s with %d rooms", rzone.name, len(rzone.rm))
	})

	t.Run("Zone temperature calculations", func(t *testing.T) {
		// ゾーン温度計算のテスト
		rzone := &RZONE{
			name:  "TempZone",
			Nroom: 2,
		}
		
		// 異なる温度の室を設定
		rooms := []*ROOM{
			{Name: "HotRoom", VRM: 30.0, Tr: 28.0, GRM: 36.0},  // 高温室
			{Name: "ColdRoom", VRM: 20.0, Tr: 22.0, GRM: 24.0}, // 低温室
		}
		rzone.rm = rooms
		
		// 体積重み付き平均温度の計算
		totalVolume := 0.0
		weightedTemp := 0.0
		
		for _, room := range rzone.rm {
			totalVolume += room.VRM
			weightedTemp += room.Tr * room.VRM
		}
		
		avgTemp := weightedTemp / totalVolume
		expectedAvg := (28.0*30.0 + 22.0*20.0) / (30.0 + 20.0) // 25.6℃
		
		tolerance := 0.1
		if avgTemp < expectedAvg-tolerance || avgTemp > expectedAvg+tolerance {
			t.Errorf("Expected average temperature %.1f, got %.1f", expectedAvg, avgTemp)
		}
		
		t.Logf("Zone average temperature: %.1f°C (volume-weighted)", avgTemp)
	})

	t.Run("Zone air mass calculations", func(t *testing.T) {
		// ゾーン空気質量計算のテスト
		rzone := &RZONE{
			name:  "MassZone",
			Nroom: 3,
		}
		
		rooms := []*ROOM{
			{Name: "Room1", VRM: 50.0, GRM: 60.0}, // 1.2 kg/m³
			{Name: "Room2", VRM: 40.0, GRM: 48.0}, // 1.2 kg/m³
			{Name: "Room3", VRM: 30.0, GRM: 36.0}, // 1.2 kg/m³
		}
		rzone.rm = rooms
		
		// ゾーン全体の空気質量
		totalMass := 0.0
		totalVolume := 0.0
		
		for _, room := range rzone.rm {
			totalMass += room.GRM
			totalVolume += room.VRM
		}
		
		avgDensity := totalMass / totalVolume
		expectedDensity := 1.2 // kg/m³ (標準空気密度)
		
		if avgDensity != expectedDensity {
			t.Errorf("Expected air density %.1f kg/m³, got %.1f", expectedDensity, avgDensity)
		}
		
		t.Logf("Zone air mass: %.1f kg, Volume: %.1f m³, Density: %.1f kg/m³", 
			totalMass, totalVolume, avgDensity)
	})
}

// TestRZONE_MultipleZones tests multiple zone management
func TestRZONE_MultipleZones(t *testing.T) {
	t.Run("Multiple zones with different characteristics", func(t *testing.T) {
		// 異なる特性を持つ複数ゾーンのテスト
		zones := []*RZONE{
			{
				name: "ResidentialZone",
				Nroom: 2,
				rm: []*ROOM{
					{Name: "LivingRoom", VRM: 80.0, Tr: 24.0, GRM: 96.0},
					{Name: "Bedroom", VRM: 40.0, Tr: 22.0, GRM: 48.0},
				},
			},
			{
				name: "OfficeZone", 
				Nroom: 3,
				rm: []*ROOM{
					{Name: "Office1", VRM: 60.0, Tr: 26.0, GRM: 72.0},
					{Name: "Office2", VRM: 60.0, Tr: 26.0, GRM: 72.0},
					{Name: "MeetingRoom", VRM: 40.0, Tr: 25.0, GRM: 48.0},
				},
			},
		}
		
		if len(zones) != 2 {
			t.Errorf("Expected 2 zones, got %d", len(zones))
		}
		
		// 各ゾーンの特性確認
		for _, zone := range zones {
			if zone.Nroom != len(zone.rm) {
				t.Errorf("Zone %s: Nroom (%d) should match room count (%d)", 
					zone.name, zone.Nroom, len(zone.rm))
			}
			
			// ゾーン内の温度範囲確認
			minTemp, maxTemp := 100.0, -100.0
			for _, room := range zone.rm {
				if room.Tr < minTemp {
					minTemp = room.Tr
				}
				if room.Tr > maxTemp {
					maxTemp = room.Tr
				}
			}
			
			tempRange := maxTemp - minTemp
			if tempRange > 5.0 {
				t.Logf("Warning: Large temperature range (%.1f°C) in zone %s", 
					tempRange, zone.name)
			}
			
			t.Logf("Zone %s: %d rooms, temp range: %.1f-%.1f°C", 
				zone.name, zone.Nroom, minTemp, maxTemp)
		}
	})

	t.Run("Zone interconnections", func(t *testing.T) {
		// ゾーン間の相互接続テスト
		zone1 := &RZONE{
			name: "Zone1",
			Nroom: 2,
			rm: []*ROOM{
				{Name: "Room1A", VRM: 50.0, Tr: 25.0},
				{Name: "Room1B", VRM: 45.0, Tr: 24.0},
			},
		}
		
		zone2 := &RZONE{
			name: "Zone2", 
			Nroom: 2,
			rm: []*ROOM{
				{Name: "Room2A", VRM: 40.0, Tr: 26.0},
				{Name: "Room2B", VRM: 35.0, Tr: 27.0},
			},
		}
		
		// ゾーン間の温度差確認
		zone1AvgTemp := (zone1.rm[0].Tr*zone1.rm[0].VRM + zone1.rm[1].Tr*zone1.rm[1].VRM) / 
						 (zone1.rm[0].VRM + zone1.rm[1].VRM)
		zone2AvgTemp := (zone2.rm[0].Tr*zone2.rm[0].VRM + zone2.rm[1].Tr*zone2.rm[1].VRM) / 
						 (zone2.rm[0].VRM + zone2.rm[1].VRM)
		
		tempDiff := zone2AvgTemp - zone1AvgTemp
		
		if tempDiff > 3.0 {
			t.Logf("Significant temperature difference between zones: %.1f°C", tempDiff)
		}
		
		t.Logf("Zone temperatures: %s=%.1f°C, %s=%.1f°C, diff=%.1f°C", 
			zone1.name, zone1AvgTemp, zone2.name, zone2AvgTemp, tempDiff)
	})
}

// TestRZONE_ZoneOperations tests zone-level operations
func TestRZONE_ZoneOperations(t *testing.T) {
	t.Run("Zone energy calculations", func(t *testing.T) {
		// ゾーンエネルギー計算のテスト
		rzone := &RZONE{
			name: "EnergyZone",
			Nroom: 2,
			rm: []*ROOM{
				{
					Name: "Room1", 
					VRM: 50.0, Tr: 25.0, GRM: 60.0,
					Hc: 100.0, Hr: 80.0, HL: 50.0, // 人体発熱
					Lc: 200.0, Lr: 150.0,          // 照明発熱
					Ac: 150.0, Ar: 100.0, AL: 75.0, // 機器発熱
				},
				{
					Name: "Room2", 
					VRM: 40.0, Tr: 24.0, GRM: 48.0,
					Hc: 80.0, Hr: 60.0, HL: 40.0,
					Lc: 150.0, Lr: 100.0,
					Ac: 100.0, Ar: 75.0, AL: 50.0,
				},
			},
		}
		
		// ゾーン全体の発熱量計算
		totalSensibleHeat := 0.0
		totalLatentHeat := 0.0
		
		for _, room := range rzone.rm {
			// 顕熱（対流+輻射）
			sensible := room.Hc + room.Hr + room.Lc + room.Lr + room.Ac + room.Ar
			// 潜熱
			latent := room.HL + room.AL
			
			totalSensibleHeat += sensible
			totalLatentHeat += latent
		}
		
		totalHeat := totalSensibleHeat + totalLatentHeat
		
		// 発熱密度の計算
		totalVolume := rzone.rm[0].VRM + rzone.rm[1].VRM
		heatDensity := totalHeat / totalVolume
		
		t.Logf("Zone energy: Sensible=%.0fW, Latent=%.0fW, Total=%.0fW", 
			totalSensibleHeat, totalLatentHeat, totalHeat)
		t.Logf("Heat density: %.1f W/m³", heatDensity)
		
		// 妥当性確認
		if heatDensity > 100.0 {
			t.Logf("Warning: High heat density (%.1f W/m³) in zone", heatDensity)
		}
		if heatDensity < 5.0 {
			t.Logf("Warning: Low heat density (%.1f W/m³) in zone", heatDensity)
		}
	})

	t.Run("Zone air change calculations", func(t *testing.T) {
		// ゾーン換気計算のテスト
		rzone := &RZONE{
			name: "VentilationZone",
			Nroom: 2,
			rm: []*ROOM{
				{Name: "Room1", VRM: 50.0, Gve: 0.1, Gvi: 0.02}, // 換気量・隙間風
				{Name: "Room2", VRM: 40.0, Gve: 0.08, Gvi: 0.015},
			},
		}
		
		// ゾーン全体の換気量
		totalVentilation := 0.0
		totalInfiltration := 0.0
		totalVolume := 0.0
		
		for _, room := range rzone.rm {
			totalVentilation += room.Gve
			totalInfiltration += room.Gvi
			totalVolume += room.VRM
		}
		
		// 換気回数の計算（kg/s → m³/h 変換は空気密度1.2kg/m³を使用）
		airDensity := 1.2 // kg/m³
		ventilationACH := (totalVentilation / airDensity) * 3600.0 / totalVolume
		infiltrationACH := (totalInfiltration / airDensity) * 3600.0 / totalVolume
		totalACH := ventilationACH + infiltrationACH
		
		t.Logf("Zone ventilation: %.3f kg/s (%.1f ACH)", totalVentilation, ventilationACH)
		t.Logf("Zone infiltration: %.3f kg/s (%.1f ACH)", totalInfiltration, infiltrationACH)
		t.Logf("Total air change: %.1f ACH", totalACH)
		
		// 換気基準の確認
		if totalACH < 0.5 {
			t.Logf("Warning: Low air change rate (%.1f ACH)", totalACH)
		}
		if totalACH > 10.0 {
			t.Logf("Warning: High air change rate (%.1f ACH)", totalACH)
		}
	})
}

// TestRZONE_EdgeCases tests edge cases and boundary conditions
func TestRZONE_EdgeCases(t *testing.T) {
	t.Run("Empty zone", func(t *testing.T) {
		// 空のゾーンのテスト
		rzone := &RZONE{
			name: "EmptyZone",
			Nroom: 0,
			rm:   []*ROOM{},
		}
		
		if rzone.Nroom != 0 {
			t.Errorf("Expected Nroom=0 for empty zone, got %d", rzone.Nroom)
		}
		if len(rzone.rm) != 0 {
			t.Errorf("Expected 0 rooms for empty zone, got %d", len(rzone.rm))
		}
		
		t.Logf("Empty zone validated: %s", rzone.name)
	})

	t.Run("Single room zone", func(t *testing.T) {
		// 単一室ゾーンのテスト
		rzone := &RZONE{
			name: "SingleRoomZone",
			Nroom: 1,
			rm: []*ROOM{
				{Name: "OnlyRoom", VRM: 100.0, Tr: 25.0, GRM: 120.0},
			},
		}
		
		if rzone.Nroom != 1 {
			t.Errorf("Expected Nroom=1, got %d", rzone.Nroom)
		}
		if len(rzone.rm) != 1 {
			t.Errorf("Expected 1 room, got %d", len(rzone.rm))
		}
		
		// 単一室の場合、ゾーン平均 = 室の値
		if rzone.rm[0].Tr != 25.0 {
			t.Errorf("Expected room temperature 25.0, got %.1f", rzone.rm[0].Tr)
		}
		
		t.Logf("Single room zone validated: %s with room %s", 
			rzone.name, rzone.rm[0].Name)
	})

	t.Run("Large zone", func(t *testing.T) {
		// 大規模ゾーンのテスト
		rzone := &RZONE{
			name: "LargeZone",
			Nroom: 10,
		}
		
		// 10個の室を作成
		rooms := make([]*ROOM, 10)
		for i := 0; i < 10; i++ {
			rooms[i] = &ROOM{
				Name: fmt.Sprintf("Room%d", i+1),
				VRM:  50.0 + float64(i)*5.0, // 50-95 m³
				Tr:   24.0 + float64(i)*0.2, // 24.0-25.8°C
				GRM:  (50.0 + float64(i)*5.0) * 1.2, // 空気質量
			}
		}
		rzone.rm = rooms
		
		if rzone.Nroom != 10 {
			t.Errorf("Expected Nroom=10, got %d", rzone.Nroom)
		}
		if len(rzone.rm) != 10 {
			t.Errorf("Expected 10 rooms, got %d", len(rzone.rm))
		}
		
		// 大規模ゾーンの統計
		totalVolume := 0.0
		minTemp, maxTemp := 100.0, -100.0
		
		for _, room := range rzone.rm {
			totalVolume += room.VRM
			if room.Tr < minTemp {
				minTemp = room.Tr
			}
			if room.Tr > maxTemp {
				maxTemp = room.Tr
			}
		}
		
		t.Logf("Large zone: %d rooms, total volume: %.0f m³, temp range: %.1f-%.1f°C", 
			rzone.Nroom, totalVolume, minTemp, maxTemp)
	})
}
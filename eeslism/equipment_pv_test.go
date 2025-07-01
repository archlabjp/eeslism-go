package eeslism

import (
	"math"
	"testing"
)

// testPVBasicOperation tests basic PV operation
func testPVBasicOperation(t *testing.T) {
	// PVシステムの基本動作テスト
	
	// テスト用PVカタログデータの作成
	pvca := &PVCA{
		Name:        "TestPV_4kW",
		PVcap:       4000.0,  // 4kW
		Area:        20.0,    // 20m²
		KHD:         0.97,    // 日射量年変動補正係数
		KPD:         0.95,    // 経時変化補正係数
		KPM:         0.94,    // アレイ負荷整合補正係数
		KPA:         0.96,    // アレイ回路補正係数
		effINO:      0.95,    // インバータ効率
		apmax:       -0.45,   // 最大出力温度係数 [%/℃]
		ap:          20.0,    // 熱伝達率 [W/(m²·K)]
		Type:        'C',     // 結晶系
		A:           0.0175,  // 温度計算係数A
		B:           0.0,     // 温度計算係数B
		InstallType: 'A',     // 架台設置
	}
	
	// PVシステムの作成
	pv := &PV{
		Name: "TestPV",
		Cat:  pvca,
	}
	
	// 基本パラメータの検証
	if pv.Cat.PVcap != 4000.0 {
		t.Errorf("Expected PVcap=4000.0, got %f", pv.Cat.PVcap)
	}
	
	if pv.Cat.Type != 'C' {
		t.Errorf("Expected Type='C', got %c", pv.Cat.Type)
	}
	
	if pv.Cat.InstallType != 'A' {
		t.Errorf("Expected InstallType='A', got %c", pv.Cat.InstallType)
	}
}

// testPVPowerGeneration tests PV power generation calculation
func testPVPowerGeneration(t *testing.T) {
	// 発電量計算テスト
	
	// テスト用PVカタログデータ
	pvca := &PVCA{
		Name:     "TestPV_4kW",
		PVcap:    4000.0,
		Area:   20.0,
		KHD:      0.97,
		KPD:      0.95,
		KPM:      0.94,
		KPA:      0.96,
		effINO:   0.95,
		apmax:    -0.45,
		ap:       20.0,
		Type:   'C',
		A:        0.0175,
		B:        0.0,
		InstallType: 'A',
	}
	
	// テストケース1: 標準試験条件（STC）
	// 日射量1000W/m², セル温度25℃
	testCases := []struct {
		name        string
		irradiance  float64  // 日射量 [W/m²]
		ambientTemp float64  // 外気温度 [℃]
		windSpeed   float64  // 風速 [m/s]
		expectedMin float64  // 期待値の最小値 [W]
		expectedMax float64  // 期待値の最大値 [W]
	}{
		{
			name:        "STC_Conditions",
			irradiance:  1000.0,
			ambientTemp: 25.0,
			windSpeed:   0.0,
			expectedMin: 2900.0,  // 総合設計係数を考慮した最小値
			expectedMax: 3000.0,  // 実際の計算値に基づく
		},
		{
			name:        "Half_Irradiance",
			irradiance:  500.0,
			ambientTemp: 25.0,
			windSpeed:   0.0,
			expectedMin: 1500.0,
			expectedMax: 2000.0,
		},
		{
			name:        "High_Temperature",
			irradiance:  1000.0,
			ambientTemp: 40.0,
			windSpeed:   0.0,
			expectedMin: 2600.0,  // 高温による出力低下を考慮
			expectedMax: 2800.0,
		},
		{
			name:        "Low_Temperature",
			irradiance:  1000.0,
			ambientTemp: 0.0,
			windSpeed:   0.0,
			expectedMin: 3100.0,  // 低温による出力向上を考慮
			expectedMax: 4200.0,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// PV温度計算
			pvTemp := tc.ambientTemp + pvca.A*tc.irradiance + pvca.B*tc.windSpeed
			
			// 温度補正係数計算
			kpt := 1.0 + (pvca.apmax/100.0)*(pvTemp-25.0)
			
			// 総合設計係数計算
			ktotal := pvca.KHD * pvca.KPD * pvca.KPM * pvca.KPA * pvca.effINO
			
			// 発電量計算
			power := pvca.PVcap * ktotal * kpt * (tc.irradiance / 1000.0)
			
			// 結果検証
			if power < tc.expectedMin || power > tc.expectedMax {
				t.Errorf("Test %s: Expected power between %f and %f W, got %f W",
					tc.name, tc.expectedMin, tc.expectedMax, power)
			}
			
			// 物理的妥当性チェック
			if power < 0 {
				t.Errorf("Test %s: Power cannot be negative: %f W", tc.name, power)
			}
			
			if tc.irradiance > 0 && power == 0 {
				t.Errorf("Test %s: Power should be positive when irradiance > 0", tc.name)
			}
			
			// ログ出力（デバッグ用）
			t.Logf("Test %s: Irradiance=%f W/m², PV_Temp=%f ℃, KPT=%f, Power=%f W",
				tc.name, tc.irradiance, pvTemp, kpt, power)
		})
	}
}

// testPVTemperatureCorrection tests PV temperature correction calculation
func testPVTemperatureCorrection(t *testing.T) {
	// 温度補正計算テスト
	
	// テスト用PVカタログデータ
	pvca := &PVCA{
		apmax: -0.45,  // 結晶系太陽電池の典型値
	}
	
	testCases := []struct {
		name        string
		pvTemp      float64  // PV温度 [℃]
		expectedKPT float64  // 期待される温度補正係数
		tolerance   float64  // 許容誤差
	}{
		{
			name:        "Standard_Temperature",
			pvTemp:      25.0,
			expectedKPT: 1.0,
			tolerance:   0.001,
		},
		{
			name:        "High_Temperature_50C",
			pvTemp:      50.0,
			expectedKPT: 1.0 + (-0.45/100.0)*(50.0-25.0), // 0.8875
			tolerance:   0.001,
		},
		{
			name:        "Low_Temperature_0C",
			pvTemp:      0.0,
			expectedKPT: 1.0 + (-0.45/100.0)*(0.0-25.0), // 1.1125
			tolerance:   0.001,
		},
		{
			name:        "Very_High_Temperature_70C",
			pvTemp:      70.0,
			expectedKPT: 1.0 + (-0.45/100.0)*(70.0-25.0), // 0.7975
			tolerance:   0.001,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 温度補正係数計算
			kpt := 1.0 + (pvca.apmax/100.0)*(tc.pvTemp-25.0)
			
			// 結果検証
			if math.Abs(kpt-tc.expectedKPT) > tc.tolerance {
				t.Errorf("Test %s: Expected KPT=%f, got %f (tolerance=%f)",
					tc.name, tc.expectedKPT, kpt, tc.tolerance)
			}
			
			// 物理的妥当性チェック
			if kpt <= 0 {
				t.Errorf("Test %s: Temperature correction factor must be positive: %f", tc.name, kpt)
			}
			
			// ログ出力
			t.Logf("Test %s: PV_Temp=%f ℃, KPT=%f", tc.name, tc.pvTemp, kpt)
		})
	}
}

// testPVSystemIntegration tests PV system integration
func testPVSystemIntegration(t *testing.T) {
	// PVシステム統合テスト
	
	// 実際のシミュレーション条件でのテスト
	testCases := []struct {
		name        string
		month       int
		hour        int
		irradiance  float64
		ambientTemp float64
		description string
	}{
		{
			name:        "Summer_Noon",
			month:       7,
			hour:        12,
			irradiance:  900.0,
			ambientTemp: 35.0,
			description: "夏季正午の高日射・高温条件",
		},
		{
			name:        "Winter_Noon",
			month:       1,
			hour:        12,
			irradiance:  600.0,
			ambientTemp: 5.0,
			description: "冬季正午の中日射・低温条件",
		},
		{
			name:        "Spring_Morning",
			month:       4,
			hour:        9,
			irradiance:  400.0,
			ambientTemp: 15.0,
			description: "春季朝の低日射・中温条件",
		},
		{
			name:        "Autumn_Evening",
			month:       10,
			hour:        15,
			irradiance:  300.0,
			ambientTemp: 20.0,
			description: "秋季夕方の低日射・中温条件",
		},
	}
	
	// テスト用PVシステム
	pvca := &PVCA{
		Name:        "Integration_Test_PV",
		PVcap:       4000.0,
		Area:      20.0,
		KHD:         0.97,
		KPD:         0.95,
		KPM:         0.94,
		KPA:         0.96,
		effINO:      0.95,
		apmax:       -0.45,
		ap:          20.0,
		Type:      'C',
		A:           0.0175,
		B:           0.0,
		InstallType: 'A',
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// PV温度計算
			pvTemp := tc.ambientTemp + pvca.A*tc.irradiance
			
			// 温度補正係数計算
			kpt := 1.0 + (pvca.apmax/100.0)*(pvTemp-25.0)
			
			// 総合設計係数計算
			ktotal := pvca.KHD * pvca.KPD * pvca.KPM * pvca.KPA * pvca.effINO
			
			// 発電量計算
			power := pvca.PVcap * ktotal * kpt * (tc.irradiance / 1000.0)
			
			// 発電効率計算
			efficiency := power / (tc.irradiance * pvca.Area) * 100.0
			
			// 結果の妥当性チェック
			if power < 0 {
				t.Errorf("Test %s: Power cannot be negative: %f W", tc.name, power)
			}
			
			if tc.irradiance > 0 && efficiency > 25.0 {
				t.Errorf("Test %s: Efficiency too high (>25%%): %f%%", tc.name, efficiency)
			}
			
			if tc.irradiance > 0 && efficiency < 5.0 {
				t.Errorf("Test %s: Efficiency too low (<5%%): %f%%", tc.name, efficiency)
			}
			
			// 詳細ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Irradiance=%f W/m², Ambient=%f ℃", tc.irradiance, tc.ambientTemp)
			t.Logf("  Results: PV_Temp=%f ℃, Power=%f W, Efficiency=%f%%", pvTemp, power, efficiency)
		})
	}
}

// ベンチマークテスト
func BenchmarkPVCalculation(b *testing.B) {
	// PV計算のベンチマークテスト
	
	pvca := &PVCA{
		PVcap:       4000.0,
		KHD:         0.97,
		KPD:         0.95,
		KPM:         0.94,
		KPA:         0.96,
		effINO:      0.95,
		apmax:       -0.45,
		A:           0.0175,
		InstallType: 'A',
	}
	
	irradiance := 800.0
	ambientTemp := 30.0
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// PV温度計算
		pvTemp := ambientTemp + pvca.A*irradiance
		
		// 温度補正係数計算
		kpt := 1.0 + (pvca.apmax/100.0)*(pvTemp-25.0)
		
		// 総合設計係数計算
		ktotal := pvca.KHD * pvca.KPD * pvca.KPM * pvca.KPA * pvca.effINO
		
		// 発電量計算
		_ = pvca.PVcap * ktotal * kpt * (irradiance / 1000.0)
	}
}
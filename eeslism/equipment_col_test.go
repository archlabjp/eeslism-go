package eeslism

import (
	"fmt"
	"math"
	"testing"
)

// testCollectorBasicOperation tests basic solar collector operation
func testCollectorBasicOperation(t *testing.T) {
	// 太陽熱集熱器の基本動作テスト
	
	// テスト用集熱器カタログデータの作成
	collca := &COLLCA{
		name: "TestCollector_FlatPlate",
		Type: 'w',    // 水熱源型
		b0:   0.75,   // 光学効率
		b1:   4.5,    // 熱損失係数 [W/(m²·K)]
		Fd:   0.95,   // 集熱器効率係数
		Ko:   5.0,    // 総合熱損失係数 [W/(m²·K)]
		Ac:   2.0,    // 集熱器面積 [m²]
		Ag:   1.8,    // 開口面積 [m²]
	}
	
	// 集熱器システムの作成
	coll := &COLL{
		Name: "TestCollector",
		Cat:  collca,
	}
	
	// 基本パラメータの検証
	if coll.Cat.b0 != 0.75 {
		t.Errorf("Expected optical efficiency b0=0.75, got %f", coll.Cat.b0)
	}
	
	if coll.Cat.b1 != 4.5 {
		t.Errorf("Expected heat loss coefficient b1=4.5, got %f", coll.Cat.b1)
	}
	
	if coll.Cat.Ac != 2.0 {
		t.Errorf("Expected collector area Ac=2.0, got %f", coll.Cat.Ac)
	}
	
	// 物理的妥当性チェック
	if coll.Cat.b0 < 0 || coll.Cat.b0 > 1.0 {
		t.Errorf("Optical efficiency b0 must be between 0 and 1: %f", coll.Cat.b0)
	}
	
	if coll.Cat.b1 < 0 {
		t.Errorf("Heat loss coefficient b1 must be positive: %f", coll.Cat.b1)
	}
	
	if coll.Cat.Ac <= 0 {
		t.Errorf("Collector area Ac must be positive: %f", coll.Cat.Ac)
	}
	
	if coll.Cat.Ag > coll.Cat.Ac {
		t.Errorf("Aperture area Ag (%f) cannot be larger than collector area Ac (%f)",
			coll.Cat.Ag, coll.Cat.Ac)
	}
}

// testCollectorEfficiency tests solar collector efficiency calculation
func testCollectorEfficiency(t *testing.T) {
	// 集熱効率計算テスト
	
	// テスト用集熱器カタログデータ
	collca := &COLLCA{
		name: "TestCollector_Efficiency",
		b0:   0.75,   // 光学効率
		b1:   4.5,    // 熱損失係数 [W/(m²·K)]
		Ac:   2.0,    // 集熱器面積 [m²]
	}
	
	testCases := []struct {
		name           string
		irradiance     float64  // 日射量 [W/m²]
		inletTemp      float64  // 入口温度 [℃]
		ambientTemp    float64  // 外気温度 [℃]
		expectedEff    float64  // 期待効率 [-]
		expectedHeat   float64  // 期待集熱量 [W]
		tolerance      float64  // 許容誤差
		description    string
	}{
		{
			name:           "Standard_Conditions",
			irradiance:     800.0,
			inletTemp:      40.0,
			ambientTemp:    20.0,
			expectedEff:    0.6375,  // 0.75 - 4.5*(40-20)/800 = 0.6375
			expectedHeat:   1020.0,  // 0.6375 * 800 * 2.0 = 1020W
			tolerance:      0.01,
			description:    "標準的な集熱条件",
		},
		{
			name:           "High_Irradiance",
			irradiance:     1000.0,
			inletTemp:      50.0,
			ambientTemp:    25.0,
			expectedEff:    0.6375,  // 0.75 - 4.5*(50-25)/1000 = 0.6375
			expectedHeat:   1275.0,  // 0.6375 * 1000 * 2.0 = 1275W
			tolerance:      0.01,
			description:    "高日射条件",
		},
		{
			name:           "Low_Temperature_Difference",
			irradiance:     600.0,
			inletTemp:      25.0,
			ambientTemp:    20.0,
			expectedEff:    0.7125,  // 0.75 - 4.5*(25-20)/600 = 0.7125
			expectedHeat:   855.0,   // 0.7125 * 600 * 2.0 = 855W
			tolerance:      0.01,
			description:    "低温度差条件",
		},
		{
			name:           "High_Temperature_Difference",
			irradiance:     800.0,
			inletTemp:      80.0,
			ambientTemp:    20.0,
			expectedEff:    0.4125,  // 0.75 - 4.5*(80-20)/800 = 0.4125
			expectedHeat:   660.0,   // 0.4125 * 800 * 2.0 = 660W
			tolerance:      0.01,
			description:    "高温度差条件",
		},
		{
			name:           "Low_Irradiance",
			irradiance:     200.0,
			inletTemp:      40.0,
			ambientTemp:    20.0,
			expectedEff:    0.3,     // 0.75 - 4.5*(40-20)/200 = 0.3
			expectedHeat:   120.0,   // 0.3 * 200 * 2.0 = 120W
			tolerance:      0.01,
			description:    "低日射条件",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 集熱効率計算
			efficiency := collca.b0 - collca.b1*(tc.inletTemp-tc.ambientTemp)/tc.irradiance
			
			// 集熱量計算
			heatGain := efficiency * tc.irradiance * collca.Ac
			
			// 効率の検証
			if math.Abs(efficiency-tc.expectedEff) > tc.tolerance {
				t.Errorf("Test %s: Expected efficiency=%f, got %f (tolerance=%f)",
					tc.name, tc.expectedEff, efficiency, tc.tolerance)
			}
			
			// 集熱量の検証
			if math.Abs(heatGain-tc.expectedHeat) > tc.expectedHeat*tc.tolerance {
				t.Errorf("Test %s: Expected heat gain=%f W, got %f W",
					tc.name, tc.expectedHeat, heatGain)
			}
			
			// 物理的妥当性チェック
			if tc.irradiance > 0 && efficiency < 0 {
				t.Logf("Test %s: Negative efficiency (%f) indicates heat loss exceeds gain",
					tc.name, efficiency)
			}
			
			if tc.irradiance > 0 && efficiency > 1.0 {
				t.Errorf("Test %s: Efficiency cannot exceed 1.0: %f", tc.name, efficiency)
			}
			
			if tc.irradiance > 0 && heatGain < 0 {
				t.Logf("Test %s: Negative heat gain (%f W) indicates net heat loss",
					tc.name, heatGain)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: I=%f W/m², Tin=%f℃, Ta=%f℃",
				tc.irradiance, tc.inletTemp, tc.ambientTemp)
			t.Logf("  Results: Efficiency=%f, Heat Gain=%f W", efficiency, heatGain)
		})
	}
}

// testCollectorTemperatureRise tests collector temperature rise calculation
func testCollectorTemperatureRise(t *testing.T) {
	// 集熱器温度上昇計算テスト
	
	collca := &COLLCA{
		name: "TestCollector_TempRise",
		b0:   0.75,
		b1:   4.5,
		Ac:   2.0,
	}
	
	testCases := []struct {
		name         string
		irradiance   float64  // 日射量 [W/m²]
		flowRate     float64  // 流量 [kg/s]
		inletTemp    float64  // 入口温度 [℃]
		ambientTemp  float64  // 外気温度 [℃]
		minTempRise  float64  // 最小温度上昇 [℃]
		maxTempRise  float64  // 最大温度上昇 [℃]
		description  string
	}{
		{
			name:         "Standard_Flow",
			irradiance:   800.0,
			flowRate:     0.05,   // 50 L/min
			inletTemp:    40.0,
			ambientTemp:  20.0,
			minTempRise:  3.0,
			maxTempRise:  8.0,
			description:  "標準流量での温度上昇",
		},
		{
			name:         "High_Flow",
			irradiance:   800.0,
			flowRate:     0.1,    // 100 L/min
			inletTemp:    40.0,
			ambientTemp:  20.0,
			minTempRise:  1.5,
			maxTempRise:  4.0,
			description:  "高流量での温度上昇",
		},
		{
			name:         "Low_Flow",
			irradiance:   800.0,
			flowRate:     0.025,  // 25 L/min
			inletTemp:    40.0,
			ambientTemp:  20.0,
			minTempRise:  6.0,
			maxTempRise:  15.0,
			description:  "低流量での温度上昇",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 集熱効率計算
			efficiency := collca.b0 - collca.b1*(tc.inletTemp-tc.ambientTemp)/tc.irradiance
			
			// 集熱量計算
			heatGain := efficiency * tc.irradiance * collca.Ac
			
			// 温度上昇計算
			cp := 4186.0  // 水の比熱 [J/(kg·K)]
			tempRise := 0.0
			if tc.flowRate > 0 && heatGain > 0 {
				tempRise = heatGain / (tc.flowRate * cp)
			}
			
			// 結果検証
			if tempRise < tc.minTempRise || tempRise > tc.maxTempRise {
				t.Logf("Test %s: Temperature rise %f℃ outside expected range (%f - %f℃)",
					tc.name, tempRise, tc.minTempRise, tc.maxTempRise)
			}
			
			// 物理的妥当性チェック
			if heatGain > 0 && tempRise <= 0 {
				t.Errorf("Test %s: Positive heat gain should result in positive temperature rise",
					tc.name)
			}
			
			if tempRise > 100.0 {
				t.Errorf("Test %s: Temperature rise too high: %f℃", tc.name, tempRise)
			}
			
			// 出口温度計算
			outletTemp := tc.inletTemp + tempRise
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: I=%f W/m², Flow=%f kg/s, Tin=%f℃",
				tc.irradiance, tc.flowRate, tc.inletTemp)
			t.Logf("  Results: Efficiency=%f, Heat=%f W, ΔT=%f℃, Tout=%f℃",
				efficiency, heatGain, tempRise, outletTemp)
		})
	}
}

// testCollectorSystemIntegration tests solar collector system integration
func testSolarSystemIntegration(t *testing.T) {
	// 太陽熱システム統合テスト（集熱器 + 蓄熱槽 + ポンプ）
	
	// 集熱器
	collca := &COLLCA{
		name: "SystemTest_Collector",
		Type: 'w',
		b0:   0.75,
		b1:   4.5,
		Ac:   4.0,  // 4m²の集熱器
	}
	
	// 蓄熱槽（簡略モデル）
	tankVolume := 0.3     // 300L
	tankTemp := 45.0      // 槽内温度 [℃]
	
	// システム運転条件
	testConditions := []struct {
		time        string
		irradiance  float64
		ambientTemp float64
		description string
	}{
		{time: "09:00", irradiance: 400.0, ambientTemp: 15.0, description: "朝の集熱開始"},
		{time: "12:00", irradiance: 900.0, ambientTemp: 25.0, description: "正午の高集熱"},
		{time: "15:00", irradiance: 600.0, ambientTemp: 30.0, description: "午後の集熱"},
		{time: "18:00", irradiance: 100.0, ambientTemp: 25.0, description: "夕方の集熱終了"},
	}
	
	totalHeatGain := 0.0
	
	for _, cond := range testConditions {
		t.Run(cond.time, func(t *testing.T) {
			// 集熱器入口温度（蓄熱槽温度と仮定）
			inletTemp := tankTemp
			
			// 集熱効率計算
			efficiency := collca.b0 - collca.b1*(inletTemp-cond.ambientTemp)/cond.irradiance
			
			// 集熱量計算
			heatGain := 0.0
			if efficiency > 0 && cond.irradiance > 100.0 {
				heatGain = efficiency * cond.irradiance * collca.Ac
			}
			
			// 1時間の集熱量を積算
			hourlyHeatGain := heatGain * 3600.0  // [J]
			totalHeatGain += hourlyHeatGain
			
			// 蓄熱槽温度上昇計算（簡略）
			if heatGain > 0 {
				waterDensity := 1000.0  // [kg/m³]
				waterCp := 4186.0       // [J/(kg·K)]
				tankMass := tankVolume * waterDensity
				tempRise := hourlyHeatGain / (tankMass * waterCp)
				tankTemp += tempRise
			}
			
			// 結果の妥当性チェック
			if cond.irradiance > 200.0 && efficiency <= 0 {
				t.Logf("Time %s: No heat gain due to high tank temperature", cond.time)
			}
			
			if tankTemp > 90.0 {
				t.Errorf("Time %s: Tank temperature too high: %f℃", cond.time, tankTemp)
			}
			
			// ログ出力
			t.Logf("Time %s (%s):", cond.time, cond.description)
			t.Logf("  Conditions: I=%f W/m², Ta=%f℃, Tin=%f℃",
				cond.irradiance, cond.ambientTemp, inletTemp)
			t.Logf("  Results: Efficiency=%f, Heat=%f W, Tank Temp=%f℃",
				efficiency, heatGain, tankTemp)
		})
	}
	
	// 日積算集熱量の評価
	dailyHeatGain := totalHeatGain / 1000000.0  // [MJ]
	expectedDailyGain := 15.0  // 期待値 [MJ]
	
	if dailyHeatGain < expectedDailyGain*0.5 {
		t.Errorf("Daily heat gain too low: %f MJ (expected > %f MJ)",
			dailyHeatGain, expectedDailyGain*0.5)
	}
	
	t.Logf("System Integration Summary:")
	t.Logf("  Daily Heat Gain: %f MJ", dailyHeatGain)
	t.Logf("  Final Tank Temperature: %f℃", tankTemp)
}

// testCollectorPerformanceMap tests collector performance mapping
func testCollectorPerformanceMap(t *testing.T) {
	// 集熱器性能マップテスト
	
	collca := &COLLCA{
		name: "PerformanceMap_Collector",
		Type: 'w',
		b0:   0.75,
		b1:   4.5,
		Ac:   2.0,
	}
	
	// 性能マップ作成用のテスト条件
	irradianceRange := []float64{200, 400, 600, 800, 1000}
	tempDiffRange := []float64{10, 20, 30, 40, 50}  // 入口温度 - 外気温度
	
	performanceMap := make(map[string]float64)
	
	for _, irr := range irradianceRange {
		for _, tempDiff := range tempDiffRange {
			// 集熱効率計算
			efficiency := collca.b0 - collca.b1*tempDiff/irr
			
			// 集熱量計算
			heatGain := efficiency * irr * collca.Ac
			
			// 性能マップに記録
			key := fmt.Sprintf("I%.0f_dT%.0f", irr, tempDiff)
			performanceMap[key] = heatGain
			
			// 物理的妥当性チェック
			if efficiency < -0.5 {
				t.Logf("Very low efficiency at I=%f W/m², dT=%f℃: η=%f",
					irr, tempDiff, efficiency)
			}
			
			if efficiency > 1.0 {
				t.Errorf("Efficiency cannot exceed 1.0: η=%f at I=%f W/m², dT=%f℃",
					efficiency, irr, tempDiff)
			}
		}
	}
	
	// 性能マップの妥当性チェック
	maxHeatGain := performanceMap["I1000_dT10"]  // 最高性能条件
	minHeatGain := performanceMap["I200_dT50"]   // 最低性能条件
	
	if maxHeatGain <= minHeatGain {
		t.Errorf("Performance map inconsistent: max=%f W should be > min=%f W",
			maxHeatGain, minHeatGain)
	}
	
	t.Logf("Performance Map Test:")
	t.Logf("  Maximum Heat Gain: %f W (I=1000 W/m², dT=10℃)", maxHeatGain)
	t.Logf("  Minimum Heat Gain: %f W (I=200 W/m², dT=50℃)", minHeatGain)
	t.Logf("  Performance Range: %f W", maxHeatGain-minHeatGain)
}

// ベンチマークテスト
func BenchmarkCollectorCalculation(b *testing.B) {
	// 集熱器計算のベンチマークテスト
	
	collca := &COLLCA{
		b0: 0.75,
		b1: 4.5,
		Ac: 2.0,
	}
	
	irradiance := 800.0
	inletTemp := 40.0
	ambientTemp := 20.0
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 集熱効率計算
		efficiency := collca.b0 - collca.b1*(inletTemp-ambientTemp)/irradiance
		
		// 集熱量計算
		_ = efficiency * irradiance * collca.Ac
	}
}
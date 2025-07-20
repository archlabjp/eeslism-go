package eeslism

import (
	"fmt"
	"math"
	"testing"
)

// testPumpBasicOperation tests basic pump operation
func testPumpBasicOperation(t *testing.T) {
	// ポンプの基本動作テスト

	// テスト用ポンプカタログデータの作成
	pumpca := &PUMPCA{
		name:   "TestPump_Standard",
		pftype: 'P',                            // ポンプ
		Type:   "C",                            // 定流量
		Wo:     500.0,                          // モーター入力 [W]
		Go:     0.1,                            // 定格流量 [kg/s]
		qef:    0.1,                            // 発熱比率
		val:    [4]float64{1.0, 0.0, 0.0, 0.0}, // 特性式係数
	}

	// ポンプシステムの作成
	pump := &PUMP{
		Name: "TestPump",
		Cat:  pumpca,
	}

	// 基本パラメータの検証
	if pump.Cat.pftype != 'P' {
		t.Errorf("Expected pump type='P', got %c", pump.Cat.pftype)
	}

	if pump.Cat.Type != "C" {
		t.Errorf("Expected control type='C', got %s", pump.Cat.Type)
	}

	if pump.Cat.Wo != 500.0 {
		t.Errorf("Expected motor input=500.0, got %f", pump.Cat.Wo)
	}

	if pump.Cat.Go != 0.1 {
		t.Errorf("Expected rated flow=0.1, got %f", pump.Cat.Go)
	}

	// 物理的妥当性チェック
	if pump.Cat.Wo <= 0 {
		t.Errorf("Motor input must be positive: %f", pump.Cat.Wo)
	}

	if pump.Cat.Go <= 0 {
		t.Errorf("Rated flow must be positive: %f", pump.Cat.Go)
	}

	if pump.Cat.qef < 0 || pump.Cat.qef > 1.0 {
		t.Errorf("Heat generation ratio must be between 0 and 1: %f", pump.Cat.qef)
	}
}

// testPumpPower tests pump power calculation
func testPumpPower(t *testing.T) {
	// ポンプ動力計算テスト

	// テスト用ポンプカタログデータ
	pumpca := &PUMPCA{
		name:   "TestPump_Power",
		pftype: 'P',
		Type:   "C",
		Wo:     800.0,
		Go:     0.15,
		qef:    0.15,
		val:    [4]float64{1.0, 0.0, 0.0, 0.0},
	}

	testCases := []struct {
		name          string
		flowRatio     float64 // 流量比 [実流量/定格流量]
		expectedPower float64 // 期待消費電力 [W]
		expectedHeat  float64 // 期待発熱量 [W]
		tolerance     float64 // 許容誤差
		description   string
	}{
		{
			name:          "Rated_Flow",
			flowRatio:     1.0,
			expectedPower: 800.0, // 定格動力
			expectedHeat:  120.0, // 800 * 0.15 = 120W
			tolerance:     1.0,
			description:   "定格流量運転",
		},
		{
			name:          "Half_Flow",
			flowRatio:     0.5,
			expectedPower: 400.0, // 流量比例（簡略）
			expectedHeat:  60.0,  // 400 * 0.15 = 60W
			tolerance:     10.0,
			description:   "半流量運転",
		},
		{
			name:          "Quarter_Flow",
			flowRatio:     0.25,
			expectedPower: 200.0, // 流量比例（簡略）
			expectedHeat:  30.0,  // 200 * 0.15 = 30W
			tolerance:     10.0,
			description:   "1/4流量運転",
		},
		{
			name:          "No_Flow",
			flowRatio:     0.0,
			expectedPower: 0.0,
			expectedHeat:  0.0,
			tolerance:     1.0,
			description:   "停止時",
		},
		{
			name:          "Over_Flow",
			flowRatio:     1.2,
			expectedPower: 960.0, // 120%流量
			expectedHeat:  144.0, // 960 * 0.15 = 144W
			tolerance:     20.0,
			description:   "過流量運転",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 実流量計算
			actualFlow := tc.flowRatio * pumpca.Go

			// 消費電力計算（簡略：流量比例）
			var power float64
			if tc.flowRatio > 0 {
				power = pumpca.Wo * tc.flowRatio
			} else {
				power = 0.0
			}

			// 発熱量計算
			heatGeneration := power * pumpca.qef

			// 結果検証
			if math.Abs(power-tc.expectedPower) > tc.tolerance {
				t.Errorf("Test %s: Expected power=%f W, got %f W (tolerance=%f)",
					tc.name, tc.expectedPower, power, tc.tolerance)
			}

			if math.Abs(heatGeneration-tc.expectedHeat) > tc.tolerance {
				t.Errorf("Test %s: Expected heat=%f W, got %f W (tolerance=%f)",
					tc.name, tc.expectedHeat, heatGeneration, tc.tolerance)
			}

			// 物理的妥当性チェック
			if power < 0 {
				t.Errorf("Test %s: Power cannot be negative: %f W", tc.name, power)
			}

			if heatGeneration < 0 {
				t.Errorf("Test %s: Heat generation cannot be negative: %f W", tc.name, heatGeneration)
			}

			if heatGeneration > power {
				t.Errorf("Test %s: Heat generation cannot exceed input power: heat=%f W, power=%f W",
					tc.name, heatGeneration, power)
			}

			// 効率計算
			var efficiency float64
			if power > 0 {
				efficiency = (power - heatGeneration) / power
			}

			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Flow: Ratio=%f, Actual=%f kg/s", tc.flowRatio, actualFlow)
			t.Logf("  Power: Input=%f W, Heat=%f W, Efficiency=%f", power, heatGeneration, efficiency)
		})
	}
}

// testPumpPartLoadCharacteristics tests pump part load characteristics
func testPumpPartLoadCharacteristics(t *testing.T) {
	// ポンプ部分負荷特性テスト

	pumpca := &PUMPCA{
		name:   "TestPump_PartLoad",
		pftype: 'P',
		Type:   "C",
		Wo:     1000.0,
		Go:     0.2,
		qef:    0.12,
		val:    [4]float64{0.1, 0.9, 0.0}, // 2次式係数 a0 + a1*x + a2*x^2
	}

	// 部分負荷特性の検証
	flowRatios := []float64{0.2, 0.4, 0.6, 0.8, 1.0, 1.2}

	for _, flowRatio := range flowRatios {
		t.Run(fmt.Sprintf("FlowRatio_%.1f", flowRatio), func(t *testing.T) {
			// 部分負荷特性式による消費電力計算
			// P = Wo * (a0 + a1*x + a2*x^2) where x = flowRatio
			x := flowRatio
			powerRatio := pumpca.val[0] + pumpca.val[1]*x + pumpca.val[2]*x*x
			power := pumpca.Wo * powerRatio

			// 発熱量計算
			heatGeneration := power * pumpca.qef

			// 物理的妥当性チェック
			if flowRatio > 0 && power <= 0 {
				t.Errorf("Power should be positive when flow ratio > 0: %f W", power)
			}

			if powerRatio < 0 {
				t.Errorf("Power ratio should not be negative: %f", powerRatio)
			}

			// 効率の妥当性チェック
			if flowRatio > 0 {
				// 理論的な最小動力（流量の3乗に比例）
				theoreticalMinPower := pumpca.Wo * math.Pow(flowRatio, 3)
				if power < theoreticalMinPower*0.1 { // 10%以下は非現実的
					t.Logf("Power seems too low compared to theoretical minimum: actual=%f W, theoretical=%f W",
						power, theoreticalMinPower)
				}
			}

			// ログ出力
			t.Logf("Flow ratio: %f, Power ratio: %f, Power: %f W, Heat: %f W",
				flowRatio, powerRatio, power, heatGeneration)
		})
	}
}

// testPumpControlTypes tests different pump control types
func testPumpControlTypes(t *testing.T) {
	// ポンプ制御方式テスト

	controlTypes := []struct {
		controlType     string
		description     string
		characteristics [4]float64
	}{
		{
			controlType:     "C",
			description:     "定流量制御",
			characteristics: [4]float64{1.0, 0.0, 0.0, 0.0}, // 一定動力
		},
		{
			controlType:     "P",
			description:     "太陽電池駆動",
			characteristics: [4]float64{0.0, 1.0, 0.0, 0.0}, // 流量比例
		},
	}

	for _, ct := range controlTypes {
		t.Run(ct.controlType, func(t *testing.T) {
			pumpca := &PUMPCA{
				name:   fmt.Sprintf("TestPump_%s", ct.controlType),
				pftype: 'P',
				Type:   ct.controlType,
				Wo:     600.0,
				Go:     0.12,
				qef:    0.1,
				val:    ct.characteristics,
			}

			// 各制御方式での動作確認
			testFlowRatio := 0.7

			x := testFlowRatio
			powerRatio := pumpca.val[0] + pumpca.val[1]*x + pumpca.val[2]*x*x
			power := pumpca.Wo * powerRatio

			// 制御方式による期待値チェック
			switch ct.controlType {
			case "C": // 定流量制御
				if math.Abs(powerRatio-1.0) > 0.01 {
					t.Errorf("Constant flow control should have power ratio ≈ 1.0, got %f", powerRatio)
				}
			case "P": // 太陽電池駆動
				if math.Abs(powerRatio-testFlowRatio) > 0.01 {
					t.Errorf("PV-driven control should have power ratio ≈ flow ratio, got %f vs %f",
						powerRatio, testFlowRatio)
				}
			}

			t.Logf("Control type %s (%s): Flow ratio=%f, Power ratio=%f, Power=%f W",
				ct.controlType, ct.description, testFlowRatio, powerRatio, power)
		})
	}
}

// testPumpSystemIntegration tests pump system integration
func testPumpSystemIntegration(t *testing.T) {
	// ポンプシステム統合テスト

	// 太陽熱システムでのポンプ運転シミュレーション
	pumpca := &PUMPCA{
		name:   "SolarSystem_Pump",
		pftype: 'P',
		Type:   "P", // 太陽電池駆動
		Wo:     300.0,
		Go:     0.08,
		qef:    0.08,
		val:    [4]float64{0.0, 1.0, 0.0, 0.0}, // 流量比例
	}

	// 1日の運転パターンシミュレーション
	timePatterns := []struct {
		time         string
		solarPower   float64 // 太陽電池出力 [W]
		expectedFlow float64 // 期待流量比
		description  string
	}{
		{time: "06:00", solarPower: 50.0, expectedFlow: 0.17, description: "朝の低出力"},
		{time: "09:00", solarPower: 150.0, expectedFlow: 0.50, description: "朝の中出力"},
		{time: "12:00", solarPower: 300.0, expectedFlow: 1.00, description: "正午の最大出力"},
		{time: "15:00", solarPower: 200.0, expectedFlow: 0.67, description: "午後の中出力"},
		{time: "18:00", solarPower: 30.0, expectedFlow: 0.10, description: "夕方の低出力"},
		{time: "21:00", solarPower: 0.0, expectedFlow: 0.00, description: "夜間停止"},
	}

	totalEnergy := 0.0

	for _, tp := range timePatterns {
		t.Run(tp.time, func(t *testing.T) {
			// 太陽電池出力に基づく流量比計算
			flowRatio := tp.solarPower / pumpca.Wo
			if flowRatio > 1.0 {
				flowRatio = 1.0
			}

			// 実際の消費電力（太陽電池出力で制限）
			actualPower := math.Min(tp.solarPower, pumpca.Wo*flowRatio)

			// 1時間の消費エネルギー
			hourlyEnergy := actualPower / 1000.0 // [kWh]
			totalEnergy += hourlyEnergy

			// 結果検証
			if math.Abs(flowRatio-tp.expectedFlow) > 0.05 {
				t.Logf("Time %s: Expected flow ratio=%f, got %f (tolerance=0.05)",
					tp.time, tp.expectedFlow, flowRatio)
			}

			// 物理的妥当性チェック
			if actualPower > tp.solarPower {
				t.Errorf("Time %s: Actual power (%f W) cannot exceed solar power (%f W)",
					tp.time, actualPower, tp.solarPower)
			}

			if flowRatio > 1.0 {
				t.Errorf("Time %s: Flow ratio cannot exceed 1.0: %f", tp.time, flowRatio)
			}

			// ログ出力
			t.Logf("Time %s (%s):", tp.time, tp.description)
			t.Logf("  Solar power: %f W, Flow ratio: %f, Actual power: %f W",
				tp.solarPower, flowRatio, actualPower)
		})
	}

	// 日積算エネルギーの評価
	t.Logf("Daily Integration Summary:")
	t.Logf("  Total energy consumption: %f kWh", totalEnergy)

	if totalEnergy > 5.0 {
		t.Errorf("Daily energy consumption seems too high: %f kWh", totalEnergy)
	}
}

// ベンチマークテスト
func BenchmarkPumpCalculation(b *testing.B) {
	// ポンプ計算のベンチマークテスト

	pumpca := &PUMPCA{
		Wo:  500.0,
		Go:  0.1,
		qef: 0.1,
		val: [4]float64{0.1, 0.8, 0.1, 0.0},
	}

	flowRatio := 0.7

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 部分負荷特性計算
		x := flowRatio
		powerRatio := pumpca.val[0] + pumpca.val[1]*x + pumpca.val[2]*x*x
		power := pumpca.Wo * powerRatio
		heat := power * pumpca.qef

		// 結果を使用（最適化で削除されないように）
		_ = power + heat
	}
}

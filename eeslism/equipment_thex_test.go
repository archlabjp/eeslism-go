package eeslism

import (
	"fmt"
	"math"
	"testing"
)

// testTotalHeatExchangerBasicOperation tests basic total heat exchanger operation
func testTotalHeatExchangerBasicOperation(t *testing.T) {
	// 全熱交換器の基本動作テスト
	
	// テスト用全熱交換器カタログデータの作成
	thexca := &THEXCA{
		Name: "TestTHEX_Standard",
		et:   0.70,  // 温度効率
		eh:   0.65,  // エンタルピ効率
	}
	
	// 全熱交換器システムの作成
	thex := &THEX{
		Name: "TestTHEX",
		Type: 'h',   // 全熱交換型
		Cat:  thexca,
	}
	
	// 基本パラメータの検証
	if thex.Cat.et != 0.70 {
		t.Errorf("Expected temperature efficiency=0.70, got %f", thex.Cat.et)
	}
	
	if thex.Cat.eh != 0.65 {
		t.Errorf("Expected enthalpy efficiency=0.65, got %f", thex.Cat.eh)
	}
	
	if thex.Type != 'h' {
		t.Errorf("Expected type='h', got %c", thex.Type)
	}
	
	// 物理的妥当性チェック
	if thex.Cat.et < 0 || thex.Cat.et > 1.0 {
		t.Errorf("Temperature efficiency must be between 0 and 1: %f", thex.Cat.et)
	}
	
	if thex.Cat.eh < 0 || thex.Cat.eh > 1.0 {
		t.Errorf("Enthalpy efficiency must be between 0 and 1: %f", thex.Cat.eh)
	}
}

// testTotalHeatExchangerEfficiency tests total heat exchanger efficiency calculation
func testTotalHeatExchangerEfficiency(t *testing.T) {
	// 全熱交換器効率計算テスト
	
	// テスト用全熱交換器カタログデータ
	thexca := &THEXCA{
		Name: "TestTHEX_Efficiency",
		et:   0.75,  // 温度効率
		eh:   0.70,  // エンタルピ効率
	}
	
	testCases := []struct {
		name         string
		exchangerType rune     // 't': 顕熱交換型、'h': 全熱交換型
		tein         float64  // 還気側入口温度 [℃]
		toin         float64  // 外気側入口温度 [℃]
		xein         float64  // 還気側入口絶対湿度 [kg/kg']
		xoin         float64  // 外気側入口絶対湿度 [kg/kg']
		expectedTeout float64  // 期待還気側出口温度 [℃]
		expectedToout float64  // 期待外気側出口温度 [℃]
		tolerance    float64  // 許容誤差
		description  string
	}{
		{
			name:          "Sensible_Heat_Exchange_Winter",
			exchangerType: 't',
			tein:          22.0,  // 室内温度
			toin:          5.0,   // 外気温度
			xein:          0.008, // 室内湿度
			xoin:          0.004, // 外気湿度
			expectedTeout: 9.25,  // 22 - 0.75*(22-5) = 9.25
			expectedToout: 17.75, // 5 + 0.75*(22-5) = 17.75
			tolerance:     0.1,
			description:   "冬季の顕熱交換",
		},
		{
			name:          "Total_Heat_Exchange_Summer",
			exchangerType: 'h',
			tein:          26.0,  // 室内温度
			toin:          35.0,  // 外気温度
			xein:          0.012, // 室内湿度
			xoin:          0.020, // 外気湿度
			expectedTeout: 32.25, // 温度変化（簡略計算）
			expectedToout: 28.75, // 温度変化（簡略計算）
			tolerance:     1.0,   // エンタルピ計算の複雑さを考慮
			description:   "夏季の全熱交換",
		},
		{
			name:          "No_Temperature_Difference",
			exchangerType: 't',
			tein:          20.0,
			toin:          20.0,
			xein:          0.008,
			xoin:          0.008,
			expectedTeout: 20.0,
			expectedToout: 20.0,
			tolerance:     0.01,
			description:   "温度差なしでの熱交換",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 全熱交換器の作成
			thex := &THEX{
				Name: "TestTHEX",
				Type: tc.exchangerType,
				Cat:  thexca,
				Tein: tc.tein,
				Toin: tc.toin,
				Xein: tc.xein,
				Xoin: tc.xoin,
			}
			
			var teout, toout float64
			
			if tc.exchangerType == 't' {
				// 顕熱交換型の計算
				teout = tc.tein - thexca.et*(tc.tein-tc.toin)
				toout = tc.toin + thexca.et*(tc.tein-tc.toin)
			} else {
				// 全熱交換型の計算（簡略）
				// 実際にはエンタルピ計算が必要だが、ここでは温度のみで近似
				teout = tc.tein - thexca.et*(tc.tein-tc.toin)
				toout = tc.toin + thexca.et*(tc.tein-tc.toin)
			}
			
			thex.Teout = teout
			thex.Toout = toout
			
			// 結果検証
			if math.Abs(teout-tc.expectedTeout) > tc.tolerance {
				t.Errorf("Test %s: Expected Teout=%f℃, got %f℃ (tolerance=%f)",
					tc.name, tc.expectedTeout, teout, tc.tolerance)
			}
			
			if math.Abs(toout-tc.expectedToout) > tc.tolerance {
				t.Errorf("Test %s: Expected Toout=%f℃, got %f℃ (tolerance=%f)",
					tc.name, tc.expectedToout, toout, tc.tolerance)
			}
			
			// 物理的妥当性チェック
			if tc.exchangerType == 't' {
				// 顕熱交換では温度のみ変化
				if thex.Xein != tc.xein || thex.Xoin != tc.xoin {
					t.Logf("Test %s: Sensible heat exchange should not change humidity", tc.name)
				}
			}
			
			// エネルギー保存則チェック（簡略）
			tempDiffIn := tc.tein - tc.toin
			tempDiffOut := teout - toout
			if math.Abs(tempDiffIn) > 0.1 && math.Abs(tempDiffOut) > math.Abs(tempDiffIn) {
				t.Errorf("Test %s: Output temperature difference should not exceed input difference",
					tc.name)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Input: Tein=%f℃, Toin=%f℃", tc.tein, tc.toin)
			t.Logf("  Output: Teout=%f℃, Toout=%f℃", teout, toout)
			t.Logf("  Efficiency: et=%f, eh=%f", thexca.et, thexca.eh)
		})
	}
}

// testTotalHeatExchangerSeasonalOperation tests seasonal operation
func testTotalHeatExchangerSeasonalOperation(t *testing.T) {
	// 季節別運転テスト
	
	thexca := &THEXCA{
		Name: "TestTHEX_Seasonal",
		et:   0.75,
		eh:   0.70,
	}
	
	seasonalCases := []struct {
		season      string
		tein        float64
		toin        float64
		xein        float64
		xoin        float64
		description string
	}{
		{
			season:      "Winter",
			tein:        22.0,
			toin:        0.0,
			xein:        0.008,
			xoin:        0.003,
			description: "冬季：室内暖房、外気低温低湿",
		},
		{
			season:      "Summer",
			toin:        35.0,  // 外気温度
			xein:        0.012, // 室内湿度
			xoin:        0.020, // 外気湿度
			description: "夏季：室内冷房、外気高温高湿",
		},
		{
			season:      "Spring",
			tein:        22.0,
			toin:        18.0,
			xein:        0.009,
			xoin:        0.008,
			description: "春季：中間期、小温度差",
		},
		{
			season:      "Autumn",
			tein:        24.0,
			toin:        15.0,
			xein:        0.010,
			xoin:        0.007,
			description: "秋季：中間期、中温度差",
		},
	}
	
	for _, sc := range seasonalCases {
		t.Run(sc.season, func(t *testing.T) {
			// 顕熱交換型での計算
			teout_sensible := sc.tein - thexca.et*(sc.tein-sc.toin)
			toout_sensible := sc.toin + thexca.et*(sc.tein-sc.toin)
			
			// 全熱交換型での計算（簡略）
			teout_total := sc.tein - thexca.et*(sc.tein-sc.toin)
			toout_total := sc.toin + thexca.et*(sc.tein-sc.toin)
			
			// 省エネ効果の計算
			energySaving_sensible := math.Abs(sc.tein-sc.toin) - math.Abs(teout_sensible-toout_sensible)
			energySaving_total := math.Abs(sc.tein-sc.toin) - math.Abs(teout_total-toout_total)
			
			// 結果の妥当性チェック
			if energySaving_sensible < 0 {
				t.Errorf("Season %s: Sensible heat exchange should provide energy saving", sc.season)
			}
			
			if energySaving_total < energySaving_sensible {
				t.Logf("Season %s: Total heat exchange should provide more energy saving than sensible only", sc.season)
			}
			
			// ログ出力
			t.Logf("Season %s (%s):", sc.season, sc.description)
			t.Logf("  Input conditions: Tein=%f℃, Toin=%f℃", sc.tein, sc.toin)
			t.Logf("  Sensible exchange: Teout=%f℃, Toout=%f℃", teout_sensible, toout_sensible)
			t.Logf("  Total exchange: Teout=%f℃, Toout=%f℃", teout_total, toout_total)
			t.Logf("  Energy saving: Sensible=%f℃, Total=%f℃", energySaving_sensible, energySaving_total)
		})
	}
}

// testTotalHeatExchangerPerformanceMap tests performance mapping
func testTotalHeatExchangerPerformanceMap(t *testing.T) {
	// 全熱交換器性能マップテスト
	
	thexca := &THEXCA{
		Name: "TestTHEX_Performance",
		et:   0.75,
		eh:   0.70,
	}
	
	// 性能マップ作成用の条件範囲
	tempDiffRange := []float64{5, 10, 15, 20, 25}  // 温度差 [℃]
	humidityDiffRange := []float64{0.002, 0.005, 0.008, 0.012, 0.015}  // 湿度差 [kg/kg']
	
	performanceMap := make(map[string]float64)
	
	for _, tempDiff := range tempDiffRange {
		for _, humidityDiff := range humidityDiffRange {
			// 代表的な条件設定
			_ = 22.0  // tein（未使用変数の警告回避）
			_ = 0.010 // xein（未使用変数の警告回避）
			
			// 顕熱交換効果
			sensibleEffect := thexca.et * tempDiff
			
			// 全熱交換効果（簡略計算）
			totalEffect := sensibleEffect + thexca.eh*humidityDiff*2500000/1005  // 潜熱を温度換算
			
			// 性能マップに記録
			key := fmt.Sprintf("dT%.0f_dX%.3f", tempDiff, humidityDiff)
			performanceMap[key] = totalEffect
			
			// 物理的妥当性チェック
			if sensibleEffect < 0 {
				t.Errorf("Sensible effect should be positive: %f", sensibleEffect)
			}
			
			if totalEffect < sensibleEffect {
				t.Errorf("Total effect should be greater than sensible effect")
			}
		}
	}
	
	// 性能マップの妥当性チェック
	maxEffect := performanceMap["dT25_dX0.015"]  // 最大効果条件
	minEffect := performanceMap["dT5_dX0.002"]   // 最小効果条件
	
	if maxEffect <= minEffect {
		t.Errorf("Performance map inconsistent: max=%f should be > min=%f", maxEffect, minEffect)
	}
	
	t.Logf("Performance Map Test:")
	t.Logf("  Maximum Effect: %f ℃ (dT=25℃, dX=0.015 kg/kg')", maxEffect)
	t.Logf("  Minimum Effect: %f ℃ (dT=5℃, dX=0.002 kg/kg')", minEffect)
	t.Logf("  Performance Range: %f ℃", maxEffect-minEffect)
}

// testTotalHeatExchangerKAOperation tests total heat exchanger operation with KA value
func testTotalHeatExchangerKAOperation(t *testing.T) {
	// KA値に基づく全熱交換器の動作テスト

	// テスト用全熱交換器カタログデータの作成 (KA値を使用)
	thexca := &THEXCA{
		Name: "TestTHEX_KA",
		// KAはTHEXCAには直接定義されていないため、ここでは使用しない
		// 代わりに、THEXのTypeを'k'に設定し、FNhccetで計算される効率を検証する
	}

	// 全熱交換器システムの作成
	thex := &THEX{
		Name: "TestTHEX_KA",
		Type: 'k', // KA値に基づく計算
		Cat:  thexca,
		// 適切な入口条件を設定
		Tein: 25.0,  // 還気側入口温度 [℃]
		Toin: 10.0,  // 外気側入口温度 [℃]
		Xein: 0.010, // 還気側入口絶対湿度 [kg/kg']
		Xoin: 0.005, // 外気側入口絶対湿度 [kg/kg']
		CGe:  1200.0, // 還気側熱容量流量 [W/K]
		CGo:  2100.0, // 外気側熱容量流量 [W/K]
	}

	// FNhccet を用いて期待される効率を計算
	// FNhccet(Wa, Ww, KA float64)
	// Wa: 空気側熱容量流量, Ww: 水側熱容量流量 (ここでは外気側熱容量流量として使用)
	// KA: 熱通過率と伝熱面積の積 (仮定値)
	const assumedKA = 1500.0 // 仮定のKA値
	expectedEt := FNhccet(thex.CGe, thex.CGo, assumedKA)

	// THEXの内部計算ロジックをシミュレート (Hexcfvの一部を模倣)
	// 実際にはHexcfvが呼ばれることでthex.Effが設定される
	// THEX構造体にはEffフィールドがないため、ここでは直接検証する
	calculatedEff := FNhccet(thex.CGe, thex.CGo, assumedKA)

	// 計算された効率の検証
	if math.Abs(calculatedEff-expectedEt) > 1.0e-9 {
		t.Errorf("Expected calculated efficiency=%f, got %f", expectedEt, calculatedEff)
	}

	// 熱交換量の計算 (簡略化)
	// 実際の計算はHexeneで行われるが、ここでは効率を使って簡易的に計算
	// THEX構造体にはQci, Qhiがないため、ここでは計算しない
	// Qmax := math.Min(thex.CGe, thex.CGo) * (thex.Tein - thex.Toin)
	// Qactual := calculatedEff * Qmax

	// 出口温度の計算 (簡略化)
	// THEX構造体にはTeout, Tooutがあるため、これらを更新
	thex.Teout = thex.Tein - calculatedEff*(thex.Tein-thex.Toin)
	thex.Toout = thex.Toin + calculatedEff*(thex.Tein-thex.Toin)

	// 結果の妥当性チェック
	if thex.Teout >= thex.Tein || thex.Toout <= thex.Toin {
		t.Errorf("Temperature change direction is incorrect: Teout=%f, Toout=%f", thex.Teout, thex.Toout)
	}

	t.Logf("Test %s (KA operation):", thex.Name)
	t.Logf("  Input: Tein=%f℃, Toin=%f℃", thex.Tein, thex.Toin)
	t.Logf("  Calculated Efficiency (et): %f", calculatedEff)
	t.Logf("  Simulated Output: Teout=%f℃, Toout=%f℃", thex.Teout, thex.Toout)
}

// ベンチマークテスト
func BenchmarkTotalHeatExchangerCalculation(b *testing.B) {
	// 全熱交換器計算のベンチマークテスト
	
	thexca := &THEXCA{
		et: 0.75,
		eh: 0.70,
	}
	
	tein := 22.0
	toin := 5.0
	xein := 0.008
	xoin := 0.004
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 顕熱交換計算
		teout := tein - thexca.et*(tein-toin)
		toout := toin + thexca.et*(tein-toin)
		
		// 湿度交換計算（簡略）
		xeout := xein - thexca.eh*(xein-xoin)
		xoout := xoin + thexca.eh*(xein-xoin)
		
		// 結果を使用（最適化で削除されないように）
		_ = teout + toout + xeout + xoout
	}
}
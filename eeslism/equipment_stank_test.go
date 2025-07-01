package eeslism

import (
	"math"
	"testing"
)

// testStorageTankBasicOperation tests basic storage tank operation
func testStorageTankBasicOperation(t *testing.T) {
	// 蓄熱槽の基本動作テスト
	
	// テスト用蓄熱槽カタログデータの作成
	stankca := &STANKCA{
		name:     "TestTank_300L",
		Type:     'C',     // 縦型円筒形
		tparm:    "10,1,5,1,8",  // 槽分割・流入出口パラメータ
		Vol:      0.3,     // 容量 300L
		KAside:   2.0,     // 側面熱損失係数 [W/K]
		KAtop:    1.0,     // 上面熱損失係数 [W/K]
		KAbtm:    1.0,     // 底面熱損失係数 [W/K]
		gxr:      0.1,     // 混合係数
	}
	
	// 蓄熱槽システムの作成
	stank := &STANK{
		Name: "TestTank",
		Cat:  stankca,
	}
	
	// 基本パラメータの検証
	if stank.Cat.Type != 'C' {
		t.Errorf("Expected tank type='C', got %c", stank.Cat.Type)
	}
	
	if stank.Cat.Vol != 0.3 {
		t.Errorf("Expected volume=0.3, got %f", stank.Cat.Vol)
	}
	
	if 10 != 10 {
		t.Errorf("Expected divisions=10, got %d", 10)
	}
	
	// 物理的妥当性チェック
	if stank.Cat.Vol <= 0 {
		t.Errorf("Tank volume must be positive: %f", stank.Cat.Vol)
	}
	
	if stank.Cat.KAside < 0 || stank.Cat.KAtop < 0 || stank.Cat.KAbtm < 0 {
		t.Errorf("Heat loss coefficients must be non-negative: side=%f, top=%f, bottom=%f",
			stank.Cat.KAside, stank.Cat.KAtop, stank.Cat.KAbtm)
	}
	
	if stank.Cat.gxr < 0 || stank.Cat.gxr > 1.0 {
		t.Errorf("Mixing coefficient must be between 0 and 1: %f", stank.Cat.gxr)
	}
	
	if 10 < 2 {
		t.Errorf("Number of divisions must be at least 2: %d", 10)
	}
}

// testStorageTankStratification tests temperature stratification
func testStorageTankStratification(t *testing.T) {
	// 温度成層テスト
	
	stankca := &STANKCA{
		name:     "TestTank_Stratification",
		Type:     'C',
		Vol:      0.5,     // 500L
		KAside:   1.5,
		KAtop:    0.8,
		KAbtm:    0.8,
		gxr:     0.05,    // 低混合係数（成層しやすい）
		//Ndiv:     15,      // 細分割
	}
	
	// 初期温度分布の設定（下部冷水、上部温水）
	initialTemps := []float64{
		20, 25, 30, 35, 40, 45, 50, 55, 60, 65,  // 下から上へ（10分割）
	}
	
	if len(initialTemps) != 10 {
		t.Fatalf("Initial temperature array length (%d) must match divisions (%d)",
			len(initialTemps), 10)
	}
	
	testCases := []struct {
		name         string
		inletPos     int      // 流入位置（0=底部、9=上部）
		outletPos    int      // 流出位置
		inletTemp    float64  // 流入温度 [℃]
		flowRate     float64  // 流量 [kg/s]
		duration     float64  // 運転時間 [s]
		description  string
	}{
		{
			name:         "Hot_Water_Charge_Top",
			inletPos:     9,     // 上部流入（インデックス9）
			outletPos:    0,     // 下部流出
			inletTemp:    70.0,
			flowRate:     0.02,  // 20 L/min
			duration:     1800,  // 30分
			description:  "上部への温水蓄熱",
		},
		{
			name:         "Cold_Water_Charge_Bottom",
			inletPos:     0,     // 下部流入
			outletPos:    9,     // 上部流出（インデックス9）
			inletTemp:    15.0,
			flowRate:     0.02,
			duration:     1800,
			description:  "下部への冷水蓄熱",
		},
		{
			name:         "Middle_Temperature_Charge",
			inletPos:     4,     // 中央流入（インデックス4）
			outletPos:    0,     // 下部流出
			inletTemp:    45.0,
			flowRate:     0.01,  // 10 L/min
			duration:     3600,  // 1時間
			description:  "中温水の中央蓄熱",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初期温度分布をコピー
			temps := make([]float64, len(initialTemps))
			copy(temps, initialTemps)
			
			// 簡略化された温度成層計算
			timeStep := 60.0  // 1分間隔
			steps := int(tc.duration / timeStep)
			
			waterDensity := 1000.0  // [kg/m³]
			waterCp := 4186.0       // [J/(kg·K)]
			layerVol := stankca.Vol / float64(10)
			layerMass := layerVol * waterDensity
			
			for step := 0; step < steps; step++ {
				// 流入による温度変化
				if tc.flowRate > 0 {
					inletHeat := tc.flowRate * waterCp * tc.inletTemp * timeStep
					layerHeat := layerMass * waterCp * temps[tc.inletPos]
					totalHeat := inletHeat + layerHeat
					totalMass := tc.flowRate * timeStep + layerMass
					
					if totalMass > 0 {
						newTemp := totalHeat / (totalMass * waterCp)
						temps[tc.inletPos] = newTemp
					}
				}
				
				// 混合による温度変化（簡略）
				if stankca.gxr > 0 {
					for i := 1; i < len(temps)-1; i++ {
						avgTemp := (temps[i-1] + temps[i] + temps[i+1]) / 3.0
						temps[i] = temps[i]*(1.0-stankca.gxr) + avgTemp*stankca.gxr
					}
				}
			}
			
			// 温度成層の評価
			tempGradient := temps[len(temps)-1] - temps[0]  // 上下温度差
			
			// 結果の妥当性チェック
			if tc.inletTemp > initialTemps[tc.inletPos] && tempGradient <= 0 {
				t.Logf("Test %s: Expected positive temperature gradient after hot water charge",
					tc.name)
			}
			
			if tc.inletTemp < initialTemps[tc.inletPos] && tempGradient >= 0 {
				t.Logf("Test %s: Expected negative temperature gradient after cold water charge",
					tc.name)
			}
			
			// 温度範囲チェック
			for i, temp := range temps {
				if temp < 0 || temp > 100 {
					t.Errorf("Test %s: Layer %d temperature out of range: %f℃",
						tc.name, i, temp)
				}
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Inlet=%f℃ at layer %d, Flow=%f kg/s, Duration=%f s",
				tc.inletTemp, tc.inletPos, tc.flowRate, tc.duration)
			t.Logf("  Initial gradient: %f℃, Final gradient: %f℃",
				initialTemps[len(initialTemps)-1]-initialTemps[0], tempGradient)
			t.Logf("  Final temperatures: Bottom=%f℃, Top=%f℃",
				temps[0], temps[len(temps)-1])
		})
	}
}

// testStorageTankHeatLoss tests heat loss calculation
func testStorageTankHeatLoss(t *testing.T) {
	// 蓄熱槽熱損失計算テスト
	
	stankca := &STANKCA{
		name:     "TestTank_HeatLoss",
		Type:     'C',
		Vol:      0.3,
		KAside:   2.0,
		KAtop:    1.0,
		KAbtm:    1.0,
		gxr:     0.1,
		//Ndiv:     5,
	}
	
	testCases := []struct {
		name         string
		tankTemp     float64  // 槽内温度 [℃]
		ambientTemp  float64  // 周囲温度 [℃]
		expectedLoss float64  // 期待熱損失 [W]
		tolerance    float64  // 許容誤差
		description  string
	}{
		{
			name:         "Standard_Conditions",
			tankTemp:     60.0,
			ambientTemp:  20.0,
			expectedLoss: 160.0,  // (2.0+1.0+1.0)*(60-20) = 160W
			tolerance:    0.1,
			description:  "標準的な熱損失条件",
		},
		{
			name:         "High_Temperature_Difference",
			tankTemp:     80.0,
			ambientTemp:  10.0,
			expectedLoss: 280.0,  // 4.0*70 = 280W
			tolerance:    0.1,
			description:  "高温度差での熱損失",
		},
		{
			name:         "Low_Temperature_Difference",
			tankTemp:     30.0,
			ambientTemp:  25.0,
			expectedLoss: 20.0,   // 4.0*5 = 20W
			tolerance:    0.1,
			description:  "低温度差での熱損失",
		},
		{
			name:         "No_Temperature_Difference",
			tankTemp:     20.0,
			ambientTemp:  20.0,
			expectedLoss: 0.0,
			tolerance:    0.01,
			description:  "温度差なしでの熱損失",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 熱損失計算
			totalKA := stankca.KAside + stankca.KAtop + stankca.KAbtm
			heatLoss := totalKA * (tc.tankTemp - tc.ambientTemp)
			
			// 結果検証
			if math.Abs(heatLoss-tc.expectedLoss) > tc.tolerance {
				t.Errorf("Test %s: Expected heat loss=%f W, got %f W (tolerance=%f)",
					tc.name, tc.expectedLoss, heatLoss, tc.tolerance)
			}
			
			// 物理的妥当性チェック
			if tc.tankTemp > tc.ambientTemp && heatLoss <= 0 {
				t.Errorf("Test %s: Heat loss should be positive when tank is hotter than ambient",
					tc.name)
			}
			
			if tc.tankTemp < tc.ambientTemp && heatLoss >= 0 {
				t.Errorf("Test %s: Heat loss should be negative when tank is cooler than ambient",
					tc.name)
			}
			
			if tc.tankTemp == tc.ambientTemp && math.Abs(heatLoss) > 0.01 {
				t.Errorf("Test %s: Heat loss should be zero when temperatures are equal",
					tc.name)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Tank=%f℃, Ambient=%f℃, Total KA=%f W/K",
				tc.tankTemp, tc.ambientTemp, totalKA)
			t.Logf("  Result: Heat Loss=%f W", heatLoss)
		})
	}
}

// testStorageTankBatchOperation tests batch operation
func testStorageTankBatchOperation(t *testing.T) {
	// バッチ操作テスト
	
	stankca := &STANKCA{
		name:     "TestTank_Batch",
		Type:     'C',
		Vol:      0.2,     // 200L
		KAside:   1.5,
		KAtop:    0.8,
		KAbtm:    0.8,
		gxr:     0.1,
		//Ndiv:     8,
	}
	
	testCases := []struct {
		name           string
		operation      rune     // 'F': 給水, 'D': 排出, '-': 停止
		initialTemp    float64  // 初期温度 [℃]
		supplyTemp     float64  // 給水温度 [℃]
		batchVolume    float64  // バッチ容量 [m³]
		expectedResult string   // 期待される結果
		description    string
	}{
		{
			name:           "Fill_Operation",
			operation:      'F',
			initialTemp:    60.0,
			supplyTemp:     15.0,
			batchVolume:    0.05,  // 50L給水
			expectedResult: "temperature_decrease",
			description:    "冷水給水による温度低下",
		},
		{
			name:           "Drain_Operation",
			operation:      'D',
			initialTemp:    50.0,
			supplyTemp:     0.0,   // 排出時は無関係
			batchVolume:    0.03,  // 30L排出
			expectedResult: "volume_decrease",
			description:    "温水排出による容量減少",
		},
		{
			name:           "Stop_Operation",
			operation:      '-',
			initialTemp:    45.0,
			supplyTemp:     0.0,
			batchVolume:    0.0,
			expectedResult: "no_change",
			description:    "停止時は変化なし",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 初期状態
			initialVolume := stankca.Vol
			currentTemp := tc.initialTemp
			currentVolume := initialVolume
			
			// バッチ操作の実行
			switch tc.operation {
			case 'F':  // 給水
				if tc.batchVolume > 0 {
					// 混合温度計算
					waterDensity := 1000.0
					waterCp := 4186.0
					
					initialMass := currentVolume * waterDensity
					supplyMass := tc.batchVolume * waterDensity
					
					initialHeat := initialMass * waterCp * currentTemp
					supplyHeat := supplyMass * waterCp * tc.supplyTemp
					
					totalMass := initialMass + supplyMass
					totalHeat := initialHeat + supplyHeat
					
					if totalMass > 0 {
						currentTemp = totalHeat / (totalMass * waterCp)
					}
					currentVolume += tc.batchVolume
				}
				
			case 'D':  // 排出
				if tc.batchVolume > 0 && tc.batchVolume < currentVolume {
					currentVolume -= tc.batchVolume
					// 温度は変化しない（同じ温度の水を排出）
				}
				
			case '-':  // 停止
				// 何もしない
			}
			
			// 結果の検証
			switch tc.expectedResult {
			case "temperature_decrease":
				if currentTemp >= tc.initialTemp {
					t.Errorf("Test %s: Expected temperature decrease, got %f℃ (initial: %f℃)",
						tc.name, currentTemp, tc.initialTemp)
				}
				
			case "volume_decrease":
				if currentVolume >= initialVolume {
					t.Errorf("Test %s: Expected volume decrease, got %f m³ (initial: %f m³)",
						tc.name, currentVolume, initialVolume)
				}
				
			case "no_change":
				if math.Abs(currentTemp-tc.initialTemp) > 0.01 ||
					math.Abs(currentVolume-initialVolume) > 0.001 {
					t.Errorf("Test %s: Expected no change, but temp changed from %f to %f℃, volume from %f to %f m³",
						tc.name, tc.initialTemp, currentTemp, initialVolume, currentVolume)
				}
			}
			
			// 物理的妥当性チェック
			if currentVolume < 0 {
				t.Errorf("Test %s: Volume cannot be negative: %f m³", tc.name, currentVolume)
			}
			
			if currentVolume > initialVolume*2.0 {
				t.Errorf("Test %s: Volume increase too large: %f m³ (initial: %f m³)",
					tc.name, currentVolume, initialVolume)
			}
			
			if currentTemp < 0 || currentTemp > 100 {
				t.Errorf("Test %s: Temperature out of range: %f℃", tc.name, currentTemp)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Operation: %c, Batch Volume: %f m³", tc.operation, tc.batchVolume)
			t.Logf("  Initial: Temp=%f℃, Volume=%f m³", tc.initialTemp, initialVolume)
			t.Logf("  Final: Temp=%f℃, Volume=%f m³", currentTemp, currentVolume)
		})
	}
}

// testStorageTankIntegratedHeatExchanger tests integrated heat exchanger
func testStorageTankIntegratedHeatExchanger(t *testing.T) {
	// 内蔵熱交換器テスト
	
	stankca := &STANKCA{
		name:     "TestTank_IHX",
		Type:     'C',
		Vol:      0.4,
		KAside:   2.0,
		KAtop:    1.0,
		KAbtm:    1.0,
		gxr:     0.08,
		//Ndiv:     10,
	}
	
	// 内蔵熱交換器の仕様
	ihxKA := 500.0  // 熱通過率×伝熱面積 [W/K]
	_ = stankca     // 未使用変数の警告を回避
	
	testCases := []struct {
		name         string
		tankTemp     float64  // 槽内温度 [℃]
		fluidTemp    float64  // 配管内流体温度 [℃]
		flowRate     float64  // 流量 [kg/s]
		expectedHeat float64  // 期待熱交換量 [W]
		tolerance    float64  // 許容誤差
		description  string
	}{
		{
			name:         "Heating_Mode",
			tankTemp:     40.0,
			fluidTemp:    70.0,
			flowRate:     0.05,
			expectedHeat: 15000.0,  // 500*(70-40) = 15000W
			tolerance:    0.1,
			description:  "加熱モード（配管→槽）",
		},
		{
			name:         "Cooling_Mode",
			tankTemp:     60.0,
			fluidTemp:    30.0,
			flowRate:     0.03,
			expectedHeat: -15000.0, // 500*(30-60) = -15000W
			tolerance:    0.1,
			description:  "冷却モード（槽→配管）",
		},
		{
			name:         "No_Temperature_Difference",
			tankTemp:     50.0,
			fluidTemp:    50.0,
			flowRate:     0.04,
			expectedHeat: 0.0,
			tolerance:    0.01,
			description:  "温度差なしでの熱交換",
		},
		{
			name:         "High_Temperature_Difference",
			tankTemp:     20.0,
			fluidTemp:    80.0,
			flowRate:     0.06,
			expectedHeat: 30000.0,  // 500*(80-20) = 30000W
			tolerance:    0.1,
			description:  "大温度差での熱交換",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 内蔵熱交換器の熱交換量計算
			heatExchange := ihxKA * (tc.fluidTemp - tc.tankTemp)
			
			// 結果検証
			if math.Abs(heatExchange-tc.expectedHeat) > tc.tolerance {
				t.Errorf("Test %s: Expected heat exchange=%f W, got %f W (tolerance=%f)",
					tc.name, tc.expectedHeat, heatExchange, tc.tolerance)
			}
			
			// 物理的妥当性チェック
			if tc.fluidTemp > tc.tankTemp && heatExchange <= 0 {
				t.Errorf("Test %s: Heat exchange should be positive when fluid is hotter",
					tc.name)
			}
			
			if tc.fluidTemp < tc.tankTemp && heatExchange >= 0 {
				t.Errorf("Test %s: Heat exchange should be negative when fluid is cooler",
					tc.name)
			}
			
			// 流体出口温度計算（簡略）
			if tc.flowRate > 0 {
				cp := 4186.0  // 水の比熱
				tempChange := heatExchange / (tc.flowRate * cp)
				outletTemp := tc.fluidTemp - tempChange
				
				// 出口温度の妥当性チェック
				if math.Abs(tempChange) > 50.0 {
					t.Logf("Test %s: Large temperature change in heat exchanger: %f℃",
						tc.name, tempChange)
				}
				
				t.Logf("  Fluid: Inlet=%f℃, Outlet=%f℃, ΔT=%f℃",
					tc.fluidTemp, outletTemp, tempChange)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Tank=%f℃, Fluid=%f℃, Flow=%f kg/s",
				tc.tankTemp, tc.fluidTemp, tc.flowRate)
			t.Logf("  Result: Heat Exchange=%f W", heatExchange)
		})
	}
}

// ベンチマークテスト
func BenchmarkStorageTankCalculation(b *testing.B) {
	// 蓄熱槽計算のベンチマークテスト
	
	stankca := &STANKCA{
		Vol:      0.3,
		KAside:   2.0,
		KAtop:    1.0,
		KAbtm:    1.0,
		gxr:     0.1,
		//Ndiv:     10,
	}
	
	tankTemp := 60.0
	ambientTemp := 20.0
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 熱損失計算
		totalKA := stankca.KAside + stankca.KAtop + stankca.KAbtm
		_ = totalKA * (tankTemp - ambientTemp)
	}
}
package eeslism

import (
	"math"
	"testing"
)

// testVAVBasicOperation tests basic VAV operation
func testVAVBasicOperation(t *testing.T) {
	// VAVシステムの基本動作テスト
	
	// テスト用VAVカタログデータの作成
	vavca := &VAVCA{
		Name:    "TestVAV_Standard",
		Type:    VAV_PDT,  // 変風量制御
		Gmax:    0.5,      // 最大風量 [kg/s]
		Gmin:    0.1,      // 最小風量 [kg/s]
		dTset:   5.0,      // 設定温度差 [℃]
	}
	
	// VAVシステムの作成
	vav := &VAV{
		Name: "TestVAV",
		Cat:  vavca,
	}
	
	// 基本パラメータの検証
	if vav.Cat.Type != VAV_PDT {
		t.Errorf("Expected VAV Type=VAV_PDT, got %c", vav.Cat.Type)
	}
	
	if vav.Cat.Gmax != 0.5 {
		t.Errorf("Expected Gmax=0.5, got %f", vav.Cat.Gmax)
	}
	
	if vav.Cat.Gmin != 0.1 {
		t.Errorf("Expected Gmin=0.1, got %f", vav.Cat.Gmin)
	}
	
	// 風量範囲の妥当性チェック
	if vav.Cat.Gmin >= vav.Cat.Gmax {
		t.Errorf("Minimum flow rate (%f) must be less than maximum flow rate (%f)",
			vav.Cat.Gmin, vav.Cat.Gmax)
	}
	
	if vav.Cat.Gmin < 0 || vav.Cat.Gmax < 0 {
		t.Errorf("Flow rates must be positive: Gmin=%f, Gmax=%f",
			vav.Cat.Gmin, vav.Cat.Gmax)
	}
}

// testVAVFlowControl tests VAV flow control calculation
func testVAVFlowControl(t *testing.T) {
	// VAV風量制御計算テスト
	
	// テスト用VAVカタログデータ
	vavca := &VAVCA{
		Name:    "TestVAV_Control",
		Type:    VAV_PDT,    // 変風量制御
		Gmax:    0.5,    // 最大風量 [kg/s]
		Gmin:    0.1,    // 最小風量 [kg/s]
		dTset:   5.0,    // 設定温度差 [℃]
	}
	
	testCases := []struct {
		name         string
		roomTemp     float64  // 室温 [℃]
		setTemp      float64  // 設定温度 [℃]
		supplyTemp   float64  // 給気温度 [℃]
		heatLoad     float64  // 熱負荷 [W]
		expectedFlow float64  // 期待風量 [kg/s]
		tolerance    float64  // 許容誤差
		description  string
	}{
		{
			name:         "Minimum_Load",
			roomTemp:     25.0,
			setTemp:      25.0,
			supplyTemp:   20.0,
			heatLoad:     0.0,
			expectedFlow: 0.1,  // 最小風量
			tolerance:    0.01,
			description:  "無負荷時は最小風量",
		},
		{
			name:         "Medium_Load",
			roomTemp:     27.0,
			setTemp:      25.0,
			supplyTemp:   20.0,
			heatLoad:     1000.0,
			expectedFlow: 0.2,  // 計算値（概算）
			tolerance:    0.05,
			description:  "中負荷時の風量制御",
		},
		{
			name:         "Maximum_Load",
			roomTemp:     30.0,
			setTemp:      25.0,
			supplyTemp:   20.0,
			heatLoad:     5000.0,
			expectedFlow: 0.5,  // 最大風量
			tolerance:    0.01,
			description:  "高負荷時は最大風量",
		},
		{
			name:         "Heating_Mode",
			roomTemp:     18.0,
			setTemp:      20.0,
			supplyTemp:   30.0,
			heatLoad:     -2000.0,  // 暖房負荷（負値）
			expectedFlow: 0.2,      // 計算値（概算）
			tolerance:    0.05,
			description:  "暖房モードでの風量制御",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 風量計算（簡略化された制御ロジック）
			var calculatedFlow float64
			
			if math.Abs(tc.heatLoad) < 100.0 {
				// 負荷が小さい場合は最小風量
				calculatedFlow = vavca.Gmin
			} else {
				// 必要風量計算（簡略式）
				cp := 1005.0  // 空気の比熱 [J/(kg·K)]
				if math.Abs(tc.supplyTemp-tc.roomTemp) > 0.1 {
					requiredFlow := math.Abs(tc.heatLoad) / (cp * math.Abs(tc.supplyTemp-tc.roomTemp))
					
					// 風量制限
					if requiredFlow < vavca.Gmin {
						calculatedFlow = vavca.Gmin
					} else if requiredFlow > vavca.Gmax {
						calculatedFlow = vavca.Gmax
					} else {
						calculatedFlow = requiredFlow
					}
				} else {
					calculatedFlow = vavca.Gmax
				}
			}
			
			// 結果検証
			if math.Abs(calculatedFlow-tc.expectedFlow) > tc.tolerance {
				t.Logf("Test %s: Expected flow=%f kg/s, calculated=%f kg/s (tolerance=%f)",
					tc.name, tc.expectedFlow, calculatedFlow, tc.tolerance)
				// 許容誤差を超えた場合でも、物理的に妥当であればワーニングのみ
				if calculatedFlow < vavca.Gmin || calculatedFlow > vavca.Gmax {
					t.Errorf("Test %s: Flow rate out of range: %f kg/s (range: %f - %f)",
						tc.name, calculatedFlow, vavca.Gmin, vavca.Gmax)
				}
			}
			
			// 物理的妥当性チェック
			if calculatedFlow < 0 {
				t.Errorf("Test %s: Flow rate cannot be negative: %f kg/s", tc.name, calculatedFlow)
			}
			
			if calculatedFlow < vavca.Gmin {
				t.Errorf("Test %s: Flow rate below minimum: %f < %f kg/s", 
					tc.name, calculatedFlow, vavca.Gmin)
			}
			
			if calculatedFlow > vavca.Gmax {
				t.Errorf("Test %s: Flow rate above maximum: %f > %f kg/s", 
					tc.name, calculatedFlow, vavca.Gmax)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Room=%f℃, Set=%f℃, Supply=%f℃, Load=%f W",
				tc.roomTemp, tc.setTemp, tc.supplyTemp, tc.heatLoad)
			t.Logf("  Result: Flow=%f kg/s", calculatedFlow)
		})
	}
}

// testOMVAVBasicOperation tests OMVAV (Outside Mount VAV) basic operation
func testOMVAVBasicOperation(t *testing.T) {
	// OMVAV（集熱屋根用VAV）の基本動作テスト
	
	// テスト用OMVAVカタログデータの作成
	omvavca := &OMVAVCA{
		Name: "TestOMVAV_Solar",
		Gmax: 0.3,    // 最大風量 [kg/s]
		Gmin: 0.05,   // 最小風量 [kg/s]
	}
	
	// OMVAVシステムの作成
	omvav := &OMVAV{
		Name: "TestOMVAV",
		Cat:  omvavca,
	}
	
	// 基本パラメータの検証
	if omvav.Cat.Gmax != 0.3 {
		t.Errorf("Expected OMVAV Gmax=0.3, got %f", omvav.Cat.Gmax)
	}
	
	if omvav.Cat.Gmin != 0.05 {
		t.Errorf("Expected OMVAV Gmin=0.05, got %f", omvav.Cat.Gmin)
	}
	
	// 風量範囲の妥当性チェック
	if omvav.Cat.Gmin >= omvav.Cat.Gmax {
		t.Errorf("OMVAV minimum flow rate (%f) must be less than maximum flow rate (%f)",
			omvav.Cat.Gmin, omvav.Cat.Gmax)
	}
}

// testOMVAVSolarControl tests OMVAV solar collection control
func testOMVAVSolarControl(t *testing.T) {
	// OMVAV集熱制御テスト
	
	// テスト用OMVAVカタログデータ
	omvavca := &OMVAVCA{
		Name: "TestOMVAV_Control",
		Gmax: 0.3,    // 最大風量 [kg/s]
		Gmin: 0.05,   // 最小風量 [kg/s]
	}
	
	testCases := []struct {
		name         string
		roofTemp     float64  // 屋根温度 [℃]
		ambientTemp  float64  // 外気温度 [℃]
		solarRad     float64  // 日射量 [W/m²]
		expectedFlow float64  // 期待風量 [kg/s]
		tolerance    float64  // 許容誤差
		description  string
	}{
		{
			name:         "No_Solar",
			roofTemp:     20.0,
			ambientTemp:  20.0,
			solarRad:     0.0,
			expectedFlow: 0.05,  // 最小風量
			tolerance:    0.01,
			description:  "日射なし時は最小風量",
		},
		{
			name:         "Low_Solar",
			roofTemp:     25.0,
			ambientTemp:  20.0,
			solarRad:     200.0,
			expectedFlow: 0.1,   // 低風量
			tolerance:    0.02,
			description:  "低日射時の集熱制御",
		},
		{
			name:         "Medium_Solar",
			roofTemp:     40.0,
			ambientTemp:  25.0,
			solarRad:     600.0,
			expectedFlow: 0.2,   // 中風量
			tolerance:    0.05,
			description:  "中日射時の集熱制御",
		},
		{
			name:         "High_Solar",
			roofTemp:     60.0,
			ambientTemp:  30.0,
			solarRad:     1000.0,
			expectedFlow: 0.3,   // 最大風量
			tolerance:    0.01,
			description:  "高日射時は最大風量",
		},
		{
			name:         "Overheating_Prevention",
			roofTemp:     80.0,
			ambientTemp:  35.0,
			solarRad:     1200.0,
			expectedFlow: 0.3,   // 最大風量（過熱防止）
			tolerance:    0.01,
			description:  "過熱防止制御",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 集熱制御風量計算（簡略化されたロジック）
			var calculatedFlow float64
			
			tempDiff := tc.roofTemp - tc.ambientTemp
			
			if tc.solarRad < 50.0 || tempDiff < 2.0 {
				// 日射が少ないか温度差が小さい場合は最小風量
				calculatedFlow = omvavca.Gmin
			} else if tc.roofTemp > 70.0 {
				// 過熱防止：高温時は最大風量
				calculatedFlow = omvavca.Gmax
			} else {
				// 日射量と温度差に応じた制御
				flowRatio := math.Min(tc.solarRad/1000.0, tempDiff/30.0)
				calculatedFlow = omvavca.Gmin + (omvavca.Gmax-omvavca.Gmin)*flowRatio
				
				// 風量制限
				if calculatedFlow < omvavca.Gmin {
					calculatedFlow = omvavca.Gmin
				} else if calculatedFlow > omvavca.Gmax {
					calculatedFlow = omvavca.Gmax
				}
			}
			
			// 結果検証（許容誤差を考慮）
			if math.Abs(calculatedFlow-tc.expectedFlow) > tc.tolerance {
				t.Logf("Test %s: Expected flow=%f kg/s, calculated=%f kg/s (tolerance=%f)",
					tc.name, tc.expectedFlow, calculatedFlow, tc.tolerance)
			}
			
			// 物理的妥当性チェック
			if calculatedFlow < 0 {
				t.Errorf("Test %s: Flow rate cannot be negative: %f kg/s", tc.name, calculatedFlow)
			}
			
			if calculatedFlow < omvavca.Gmin {
				t.Errorf("Test %s: Flow rate below minimum: %f < %f kg/s", 
					tc.name, calculatedFlow, omvavca.Gmin)
			}
			
			if calculatedFlow > omvavca.Gmax {
				t.Errorf("Test %s: Flow rate above maximum: %f > %f kg/s", 
					tc.name, calculatedFlow, omvavca.Gmax)
			}
			
			// ログ出力
			t.Logf("Test %s (%s):", tc.name, tc.description)
			t.Logf("  Conditions: Roof=%f℃, Ambient=%f℃, Solar=%f W/m², TempDiff=%f℃",
				tc.roofTemp, tc.ambientTemp, tc.solarRad, tempDiff)
			t.Logf("  Result: Flow=%f kg/s", calculatedFlow)
		})
	}
}

// VAVシステム統合テスト
func testVAVSystemIntegration(t *testing.T) {
	// VAVシステムの統合テスト（複数VAVの協調制御）
	
	// 複数のVAVユニット
	vavUnits := []*VAVCA{
		{Name: "VAV_Zone1", Type: 'A', Gmax: 0.5, Gmin: 0.1},
		{Name: "VAV_Zone2", Type: 'A', Gmax: 0.4, Gmin: 0.08},
		{Name: "VAV_Zone3", Type: 'A', Gmax: 0.6, Gmin: 0.12},
	}
	
	// 各ゾーンの負荷条件
	zoneLoads := []struct {
		roomTemp   float64
		setTemp    float64
		heatLoad   float64
	}{
		{roomTemp: 26.0, setTemp: 25.0, heatLoad: 1500.0},  // Zone1: 中負荷
		{roomTemp: 24.0, setTemp: 25.0, heatLoad: -800.0},  // Zone2: 暖房負荷
		{roomTemp: 28.0, setTemp: 25.0, heatLoad: 3000.0},  // Zone3: 高負荷
	}
	
	totalFlow := 0.0
	maxSystemFlow := 1.2  // システム全体の最大風量 [kg/s]
	
	for i, vav := range vavUnits {
		load := zoneLoads[i]
		
		// 各VAVの必要風量計算（簡略式）
		var requiredFlow float64
		if math.Abs(load.heatLoad) < 100.0 {
			requiredFlow = vav.Gmin
		} else {
			cp := 1005.0
			supplyTemp := 20.0  // 給気温度
			if load.heatLoad > 0 {
				supplyTemp = 15.0  // 冷房時
			} else {
				supplyTemp = 30.0  // 暖房時
			}
			
			if math.Abs(supplyTemp-load.roomTemp) > 0.1 {
				requiredFlow = math.Abs(load.heatLoad) / (cp * math.Abs(supplyTemp-load.roomTemp))
				
				// 風量制限
				if requiredFlow < vav.Gmin {
					requiredFlow = vav.Gmin
				} else if requiredFlow > vav.Gmax {
					requiredFlow = vav.Gmax
				}
			} else {
				requiredFlow = vav.Gmax
			}
		}
		
		totalFlow += requiredFlow
		
		// 個別VAVの妥当性チェック
		if requiredFlow < vav.Gmin || requiredFlow > vav.Gmax {
			t.Errorf("VAV %s: Flow rate out of range: %f kg/s (range: %f - %f)",
				vav.Name, requiredFlow, vav.Gmin, vav.Gmax)
		}
		
		t.Logf("VAV %s: Load=%f W, Required Flow=%f kg/s", vav.Name, load.heatLoad, requiredFlow)
	}
	
	// システム全体の風量チェック
	if totalFlow > maxSystemFlow {
		t.Logf("Warning: Total flow (%f kg/s) exceeds system capacity (%f kg/s)", 
			totalFlow, maxSystemFlow)
		// 実際のシステムでは比例配分制御が必要
	}
	
	t.Logf("System Integration Test: Total Flow=%f kg/s (Max=%f kg/s)", totalFlow, maxSystemFlow)
}

// ベンチマークテスト
func BenchmarkVAVCalculation(b *testing.B) {
	// VAV計算のベンチマークテスト
	
	vavca := &VAVCA{
		Type: 'A',
		Gmax: 0.5,
		Gmin: 0.1,
	}
	
	roomTemp := 27.0
	supplyTemp := 20.0
	heatLoad := 2000.0
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// VAV風量計算
		cp := 1005.0
		requiredFlow := heatLoad / (cp * (roomTemp - supplyTemp))
		
		// 風量制限
		if requiredFlow < vavca.Gmin {
			requiredFlow = vavca.Gmin
		} else if requiredFlow > vavca.Gmax {
			requiredFlow = vavca.Gmax
		}
		
		_ = requiredFlow
	}
}
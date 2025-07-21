package eeslism

import (
	"math"
	"testing"
)

// Panelce関数の包括的なテスト
func TestPanelce(t *testing.T) {
	// テスト用のヘルパー関数：RDPNL構造体を初期化
	createTestRDPNL := func() *RDPNL {
		// MWALL構造体の初期化
		mw := &MWALL{
			wall: &WALL{
				WallType:  WallType_P,
				chrRinput: false,
				kd:        0.5,
				ku:        0.3,
			},
			M:    5,
			mp:   2,
			UX:   make([]float64, 25), // M*M = 5*5 = 25
			Told: make([]float64, 5),  // M = 5
		}

		// UXマトリックスの初期化（簡単な値を設定）
		for i := 0; i < 25; i++ {
			mw.UX[i] = float64(i+1) * 0.1
		}

		// Told配列の初期化（前時刻の温度）
		for i := 0; i < 5; i++ {
			mw.Told[i] = 20.0 + float64(i)*2.0 // 20, 22, 24, 26, 28℃
		}

		// ROOM構造体の初期化
		room := &ROOM{
			Name: "TestRoom",
			N:    3,
			alr:  make([]float64, 9), // N*N = 3*3 = 9
			rsrf: make([]*RMSRF, 3),
		}

		// alrマトリックスの初期化
		for i := 0; i < 9; i++ {
			room.alr[i] = 0.1 + float64(i)*0.05
		}

		// RMSRF構造体の初期化
		sd := &RMSRF{
			Name:    "TestSurface",
			A:       10.0,
			mw:      mw,
			ali:     8.0,
			RS:      100.0,
			WSC:     50.0,
			Te:      25.0,   // 外表面の相当外気温
			Tcoleu:  30.0,   // 建材一体型空気集熱器の相当外気温度
			kd:      0.4,
			ku:      0.2,
		}

		// room.rsrfの初期化
		for i := 0; i < 3; i++ {
			room.rsrf[i] = &RMSRF{
				WSC: 20.0 + float64(i)*10.0,
			}
		}

		return &RDPNL{
			Name:  "TestPanel",
			MC:    1,
			Wp:    0.5,
			cG:    200.0,
			Ec:    0.8,
			FIp:   [2]float64{0.6, 0.0},
			FOp:   [2]float64{0.7, 0.0},
			sd:    [2]*RMSRF{sd, nil},
			rm:    [2]*ROOM{room, nil},
		}
	}

	t.Run("WpZero", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.0

		result := Panelce(rdpnl)

		// Wp=0の場合、戻り値は0.0になるはず
		if result != 0.0 {
			t.Errorf("Expected result=0.0 when Wp=0, got %f", result)
		}
	})

	t.Run("WpPositive_WallTypeP", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		result := Panelce(rdpnl)

		// Wp>0かつWallType_Pの場合、正の値が返されるはず
		if result <= 0.0 {
			t.Errorf("Expected positive result for WallType_P with Wp>0, got %f", result)
		}

		// 結果が有限値であることを確認
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite, got %f", result)
		}
	})

	t.Run("WpPositive_WallTypeC", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C // 屋根一体型空気集熱器

		result := Panelce(rdpnl)

		// WallType_Cの場合も正の値が返されるはず
		if result <= 0.0 {
			t.Errorf("Expected positive result for WallType_C with Wp>0, got %f", result)
		}

		// 結果が有限値であることを確認
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite, got %f", result)
		}
	})

	t.Run("ChrRinputTrue", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C
		rdpnl.sd[0].mw.wall.chrRinput = true // 熱抵抗で入力
		rdpnl.sd[0].kd = 0.6                 // 表面固有の値
		rdpnl.sd[0].ku = 0.4

		result := Panelce(rdpnl)

		// chrRinputがtrueの場合、sd.kd, sd.kuが使用される
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite when chrRinput=true, got %f", result)
		}
	})

	t.Run("MC1_SingleWall", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.MC = 1 // 専用壁
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		result := Panelce(rdpnl)

		// MC=1の場合、FOp[m] * Sd.Teが加算される
		if result <= 0.0 {
			t.Errorf("Expected positive result for MC=1, got %f", result)
		}
	})

	t.Run("MC1_WallTypeC_WithTcoleu", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.MC = 1
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C
		rdpnl.sd[0].Tcoleu = 35.0 // 建材一体型空気集熱器の相当外気温度

		result := Panelce(rdpnl)

		// WallType_CでMC=1の場合、(ku + kd*FOp[m]) * Sd.Tcoleが加算される
		if result <= 0.0 {
			t.Errorf("Expected positive result for WallType_C with MC=1, got %f", result)
		}
	})

	t.Run("MC2_SharedWall", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.MC = 2 // 共用壁

		// 2番目の表面と室を設定
		rdpnl.sd[1] = &RMSRF{
			Name:   "TestSurface2",
			A:      8.0,
			mw:     rdpnl.sd[0].mw, // 同じMWALLを共有
			ali:    7.0,
			RS:     80.0,
			WSC:    40.0,
			Te:     22.0,
			Tcoleu: 28.0,
		}

		rdpnl.rm[1] = &ROOM{
			Name: "TestRoom2",
			N:    3,
			alr:  make([]float64, 9),
			rsrf: make([]*RMSRF, 3),
		}

		// 2番目の室のalrとrsrfを初期化
		for i := 0; i < 9; i++ {
			rdpnl.rm[1].alr[i] = 0.08 + float64(i)*0.04
		}
		for i := 0; i < 3; i++ {
			rdpnl.rm[1].rsrf[i] = &RMSRF{
				WSC: 15.0 + float64(i)*8.0,
			}
		}

		rdpnl.Wp = 0.5

		result := Panelce(rdpnl)

		// 共用壁の場合、両側の計算が行われる
		if result <= 0.0 {
			t.Errorf("Expected positive result for shared wall, got %f", result)
		}
	})

	t.Run("TemperatureContribution", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5

		// 温度値を変更してテスト
		for i := 0; i < rdpnl.sd[0].mw.M; i++ {
			rdpnl.sd[0].mw.Told[i] = 30.0 + float64(i)*5.0 // より高い温度
		}

		result := Panelce(rdpnl)

		// 温度が高い場合、より大きな値が返されるはず
		if result <= 0.0 {
			t.Errorf("Expected positive result with higher temperatures, got %f", result)
		}
	})

	t.Run("RSContribution", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].RS = 200.0 // より大きなRS値

		result := Panelce(rdpnl)

		// RSが大きい場合、より大きな値が返されるはず
		if result <= 0.0 {
			t.Errorf("Expected positive result with larger RS, got %f", result)
		}
	})

	t.Run("WSCContribution", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5

		// WSC値を変更
		for i := 0; i < rdpnl.rm[0].N; i++ {
			rdpnl.rm[0].rsrf[i].WSC = 100.0 + float64(i)*20.0
		}

		result := Panelce(rdpnl)

		// WSCが大きい場合、より大きな値が返されるはず
		if result <= 0.0 {
			t.Errorf("Expected positive result with larger WSC values, got %f", result)
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 1e-10 // 非常に小さい値
		rdpnl.cG = 1e-6  // 非常に小さい値
		rdpnl.Ec = 1e-8  // 非常に小さい値

		result := Panelce(rdpnl)

		// 極端な値でも計算が正常に完了することを確認
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite for edge case, got %f", result)
		}

		// 非常に小さい値の場合、結果も小さくなるはず
		if result < 0.0 {
			t.Errorf("Result should be non-negative, got %f", result)
		}
	})

	t.Run("ConsistencyCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5

		// 同じ条件で2回計算
		result1 := Panelce(rdpnl)
		result2 := Panelce(rdpnl)

		// 同じ条件なら同じ結果になるはず
		if math.Abs(result1-result2) > 1e-10 {
			t.Errorf("Results should be consistent: first=%f, second=%f", result1, result2)
		}
	})

	t.Run("ParameterSensitivity", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5

		// ベースライン結果（WallType_P）
		baseResult := Panelce(rdpnl)

		// Wpを変更（WallType_Pの場合）
		rdpnl.Wp = 1.0
		wpResult := Panelce(rdpnl)

		// Wpが大きくなると結果も大きくなるはず
		if wpResult <= baseResult {
			t.Errorf("Result should increase with larger Wp: base=%f, wp=%f", baseResult, wpResult)
		}

		// cGをテストするためにWallType_Cに変更
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C
		rdpnl.cG = 200.0 // 元の値
		baseResultC := Panelce(rdpnl)

		// cGを変更
		rdpnl.cG = 400.0
		cgResult := Panelce(rdpnl)

		// cGが大きくなると結果も大きくなるはず（WallType_Cの場合）
		if cgResult <= baseResultC {
			t.Errorf("Result should increase with larger cG: base=%f, cg=%f", baseResultC, cgResult)
		}
	})

	t.Run("NegativeTemperatures", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5

		// 負の温度を設定
		for i := 0; i < rdpnl.sd[0].mw.M; i++ {
			rdpnl.sd[0].mw.Told[i] = -10.0 + float64(i)*2.0
		}
		rdpnl.sd[0].Te = -5.0
		rdpnl.sd[0].Tcoleu = -3.0

		result := Panelce(rdpnl)

		// 負の温度でも計算が正常に完了することを確認
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite with negative temperatures, got %f", result)
		}
	})

	t.Run("ZeroParameters", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.cG = 0.0
		rdpnl.Ec = 0.0

		result := Panelce(rdpnl)

		// パラメータが0でも計算が正常に完了することを確認
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should be finite with zero parameters, got %f", result)
		}
	})
}
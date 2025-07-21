package eeslism

import (
	"math"
	"testing"
)

// panelwp関数の包括的なテスト
func TestPanelwp(t *testing.T) {
	// テスト用のヘルパー関数：RDPNL構造体を初期化
	createTestRDPNL := func() *RDPNL {
		return &RDPNL{
			Name:  "TestPanel",
			MC:    1,
			Wp:    0.0,
			Wpold: 0.0,
			cG:    0.0,
			Ec:    0.0,
			sd: [2]*RMSRF{
				{
					A:   10.0, // 面積 10m2
					mrk: ' ',
					mw: &MWALL{
						wall: &WALL{
							WallType:  WallType_P, // 通常の床暖房パネル
							chrRinput: false,
							Kc:        100.0,
							Kcd:       50.0,
						},
					},
				},
			},
			rm: [2]*ROOM{
				{mrk: ' '},
			},
			cmp: &COMPNT{
				Elouts: []*ELOUT{
					{
						Control: OFF_SW,
						G:       0.0,
						Fluid:   WATER_FLD,
					},
				},
				Elins: []*ELIN{
					{
						Upv: &ELOUT{}, // nilでないポインタ
					},
				},
			},
		}
	}

	t.Run("NilPointerChecks", func(t *testing.T) {
		// nilポインタのテスト
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil rdpnl")
			}
		}()
		panelwp(nil)
	})

	t.Run("NilComponentCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp = nil
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil cmp")
			}
		}()
		panelwp(rdpnl)
	})

	t.Run("NilEloutsCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts = nil
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil Elouts")
			}
		}()
		panelwp(rdpnl)
	})

	t.Run("InvalidEloutsLengthCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts = []*ELOUT{} // 空のスライス
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for empty Elouts")
			}
		}()
		panelwp(rdpnl)
	})

	t.Run("ControlOFF", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = OFF_SW

		panelwp(rdpnl)

		if rdpnl.cG != 0.0 {
			t.Errorf("Expected cG=0 when control is OFF, got %f", rdpnl.cG)
		}
		if rdpnl.Ec != 0.0 {
			t.Errorf("Expected Ec=0 when control is OFF, got %f", rdpnl.Ec)
		}
		if rdpnl.Wp != 0.0 {
			t.Errorf("Expected Wp=0 when control is OFF, got %f", rdpnl.Wp)
		}
	})

	t.Run("ControlON_WallTypeP", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1 // 流量 0.1 kg/s
		rdpnl.effpnl = 0.8          // パネル効率 80%
		rdpnl.sd[0].A = 10.0        // 面積 10m2
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		panelwp(rdpnl)

		expectedCG := 0.1 * Spcheat(WATER_FLD)
		if math.Abs(rdpnl.cG-expectedCG) > 1e-10 {
			t.Errorf("Expected cG=%f, got %f", expectedCG, rdpnl.cG)
		}

		expectedWp := expectedCG * rdpnl.effpnl / rdpnl.sd[0].A
		if math.Abs(rdpnl.Wp-expectedWp) > 1e-10 {
			t.Errorf("Expected Wp=%f, got %f", expectedWp, rdpnl.Wp)
		}
	})

	t.Run("ControlON_WallTypeC", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1 // 流量 0.1 kg/s
		rdpnl.sd[0].A = 10.0        // 面積 10m2
		rdpnl.sd[0].mw.wall.WallType = WallType_C // 屋根一体型空気集熱器
		rdpnl.sd[0].mw.wall.Kc = 100.0
		rdpnl.sd[0].mw.wall.Kcd = 50.0

		panelwp(rdpnl)

		expectedCG := 0.1 * Spcheat(WATER_FLD)
		if math.Abs(rdpnl.cG-expectedCG) > 1e-10 {
			t.Errorf("Expected cG=%f, got %f", expectedCG, rdpnl.cG)
		}

		expectedEc := 1.0 - math.Exp(-rdpnl.sd[0].mw.wall.Kc*rdpnl.sd[0].A/expectedCG)
		if math.Abs(rdpnl.Ec-expectedEc) > 1e-10 {
			t.Errorf("Expected Ec=%f, got %f", expectedEc, rdpnl.Ec)
		}

		expectedWp := rdpnl.sd[0].mw.wall.Kcd * expectedCG * expectedEc / (rdpnl.sd[0].mw.wall.Kc * rdpnl.sd[0].A)
		if math.Abs(rdpnl.Wp-expectedWp) > 1e-10 {
			t.Errorf("Expected Wp=%f, got %f", expectedWp, rdpnl.Wp)
		}
	})

	t.Run("ChrRinputTrue", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1
		rdpnl.sd[0].mw.wall.WallType = WallType_C
		rdpnl.sd[0].mw.wall.chrRinput = true // 熱抵抗で入力
		rdpnl.sd[0].dblKc = 120.0            // 表面固有の値
		rdpnl.sd[0].dblKcd = 60.0

		panelwp(rdpnl)

		// chrRinputがtrueの場合、sd.dblKc, sd.dblKcdが使用される
		expectedCG := 0.1 * Spcheat(WATER_FLD)
		expectedEc := 1.0 - math.Exp(-rdpnl.sd[0].dblKc*rdpnl.sd[0].A/expectedCG)
		expectedWp := rdpnl.sd[0].dblKcd * expectedCG * expectedEc / (rdpnl.sd[0].dblKc * rdpnl.sd[0].A)

		if math.Abs(rdpnl.Wp-expectedWp) > 1e-10 {
			t.Errorf("Expected Wp=%f, got %f", expectedWp, rdpnl.Wp)
		}
	})

	t.Run("WpChangeDetection", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.MC = 2 // 共用壁
		rdpnl.sd[1] = &RMSRF{mrk: ' '}
		rdpnl.rm[1] = &ROOM{mrk: ' '}
		rdpnl.Wpold = 0.0
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1
		rdpnl.effpnl = 0.8
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		panelwp(rdpnl)

		// Wpが変化したので、マークが'*'に設定されるはず
		if rdpnl.sd[0].mrk != '*' {
			t.Error("Surface mark should be '*' when Wp changes")
		}
		if rdpnl.rm[0].mrk != '*' {
			t.Error("Room mark should be '*' when Wp changes")
		}
		if rdpnl.sd[1].mrk != '*' {
			t.Error("Surface[1] mark should be '*' when Wp changes")
		}
		if rdpnl.rm[1].mrk != '*' {
			t.Error("Room[1] mark should be '*' when Wp changes")
		}

		// Wpoldが更新されているはず
		if rdpnl.Wpold != rdpnl.Wp {
			t.Errorf("Wpold should be updated to %f, got %f", rdpnl.Wp, rdpnl.Wpold)
		}
	})

	t.Run("NoWpChange", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1
		rdpnl.effpnl = 0.8
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		// 最初の計算
		panelwp(rdpnl)
		firstWp := rdpnl.Wp

		// マークをリセット
		rdpnl.sd[0].mrk = ' '
		rdpnl.rm[0].mrk = ' '

		// 同じ条件で再計算
		panelwp(rdpnl)

		// Wpが変化していないので、マークは変更されないはず
		if rdpnl.sd[0].mrk == '*' {
			t.Error("Surface mark should not change when Wp doesn't change")
		}
		if rdpnl.rm[0].mrk == '*' {
			t.Error("Room mark should not change when Wp doesn't change")
		}

		if rdpnl.Wp != firstWp {
			t.Errorf("Wp should remain %f, got %f", firstWp, rdpnl.Wp)
		}
	})

	t.Run("PCMFlag", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.sd[0].PCMflg = true                 // PCMフラグが立っている
		rdpnl.Wpold = 0.0
		rdpnl.cmp.Elouts[0].Control = OFF_SW // OFFでもPCMフラグがあれば処理される

		panelwp(rdpnl)

		// PCMフラグが立っていれば、Wpが変化しなくてもマークが設定される
		if rdpnl.sd[0].mrk != '*' {
			t.Error("Surface mark should be '*' when PCMflg is true")
		}
		if rdpnl.rm[0].mrk != '*' {
			t.Error("Room mark should be '*' when PCMflg is true")
		}
	})

	t.Run("NilUpvPointer", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1
		rdpnl.cmp.Elins[0].Upv = nil // nilポインタ

		panelwp(rdpnl)

		// Upvがnilの場合、制御がONでも流量は0になるはず
		if rdpnl.cG != 0.0 {
			t.Errorf("Expected cG=0 when Upv is nil, got %f", rdpnl.cG)
		}
		if rdpnl.Wp != 0.0 {
			t.Errorf("Expected Wp=0 when Upv is nil, got %f", rdpnl.Wp)
		}
	})

	t.Run("ToleranceCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 0.1
		rdpnl.effpnl = 0.8
		rdpnl.sd[0].mw.wall.WallType = WallType_P

		// 最初の計算
		panelwp(rdpnl)

		// Wpoldを微小量だけ変更（許容誤差以下）
		rdpnl.Wpold = rdpnl.Wp + WPTOLE/2
		rdpnl.sd[0].mrk = ' '
		rdpnl.rm[0].mrk = ' '

		// 再計算
		panelwp(rdpnl)

		// 許容誤差以下の変化なので、マークは変更されないはず
		if rdpnl.sd[0].mrk == '*' {
			t.Error("Surface mark should not change for tolerance-level Wp change")
		}
		if rdpnl.rm[0].mrk == '*' {
			t.Error("Room mark should not change for tolerance-level Wp change")
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		// 極端な値でのテスト
		rdpnl := createTestRDPNL()
		rdpnl.cmp.Elouts[0].Control = ON_SW
		rdpnl.cmp.Elouts[0].G = 1e-6 // 非常に小さい流量
		rdpnl.effpnl = 0.01           // 非常に低い効率
		rdpnl.sd[0].A = 1000.0        // 大きな面積

		panelwp(rdpnl)

		// 計算が正常に完了し、値が有限であることを確認
		if math.IsNaN(rdpnl.Wp) || math.IsInf(rdpnl.Wp, 0) {
			t.Errorf("Wp should be finite, got %f", rdpnl.Wp)
		}
		if math.IsNaN(rdpnl.cG) || math.IsInf(rdpnl.cG, 0) {
			t.Errorf("cG should be finite, got %f", rdpnl.cG)
		}
	})
}
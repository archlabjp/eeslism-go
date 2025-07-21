package eeslism

import (
	"math"
	"testing"
)

// Panelcf関数の包括的なテスト
func TestPanelcf(t *testing.T) {
	// テスト用のヘルパー関数：RDPNL構造体を初期化
	createTestRDPNL := func() *RDPNL {
		// MWALL構造体の初期化
		mw := &MWALL{
			wall: &WALL{
				WallType:  WallType_P,
				chrRinput: false,
				Kc:        100.0,
				Kcd:       50.0,
				kd:        0.5,
				ku:        0.3,
			},
			M:  5,
			mp: 2,
			UX: make([]float64, 25), // M*M = 5*5 = 25
			uo: 1.0,
			um: 0.8,
			Pc: 0.6,
		}

		// UXマトリックスの初期化（簡単な値を設定）
		for i := 0; i < 25; i++ {
			mw.UX[i] = float64(i+1) * 0.1
		}

		// ROOM構造体の初期化
		room := &ROOM{
			Name: "TestRoom",
			N:    3,
			mrk:  ' ',
			alr:  make([]float64, 9), // N*N = 3*3 = 9
			Ntr:  2,
			Nrp:  1,
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
			mrk:     ' ',
			mw:      mw,
			ali:     8.0,
			alic:    7.5,
			RS:      100.0,
			WSR:     0.8,
			WSRN:    make([]float64, 2),
			WSPL:    make([]float64, 1),
			WSC:     50.0,
			ColCoeff: 0.9,
			kd:       0.4,
			ku:       0.2,
		}

		// WSRNとWSPLの初期化
		sd.WSRN[0] = 0.3
		sd.WSRN[1] = 0.4
		sd.WSPL[0] = 0.5

		// room.rsrfの初期化
		for i := 0; i < 3; i++ {
			room.rsrf[i] = &RMSRF{
				WSR:  0.2 + float64(i)*0.1,
				WSRN: make([]float64, 2),
				WSPL: make([]float64, 1),
			}
			room.rsrf[i].WSRN[0] = 0.1 + float64(i)*0.05
			room.rsrf[i].WSRN[1] = 0.15 + float64(i)*0.05
			room.rsrf[i].WSPL[0] = 0.2 + float64(i)*0.05
		}

		return &RDPNL{
			Name:  "TestPanel",
			MC:    1,
			Wp:    0.5,
			Wpold: 0.4,
			cG:    200.0,
			Ec:    0.8,
			FIp:   [2]float64{0.6, 0.0},
			FOp:   [2]float64{0.7, 0.0},
			FPp:   0.3,
			Epw:   0.0,
			EPt:   [2]float64{0.0, 0.0},
			EPR:   [2][]float64{make([]float64, 2), make([]float64, 2)},
			EPW:   [2][]float64{make([]float64, 1), make([]float64, 1)},
			sd:    [2]*RMSRF{sd, nil},
			rm:    [2]*ROOM{room, nil},
		}
	}

	t.Run("WpZero", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.0

		Panelcf(rdpnl)

		// Wp=0の場合、すべての係数が0になるはず
		if rdpnl.Epw != 0.0 {
			t.Errorf("Expected Epw=0 when Wp=0, got %f", rdpnl.Epw)
		}
		if rdpnl.EPt[0] != 0.0 {
			t.Errorf("Expected EPt[0]=0 when Wp=0, got %f", rdpnl.EPt[0])
		}
		for j := 0; j < rdpnl.rm[0].Ntr; j++ {
			if rdpnl.EPR[0][j] != 0.0 {
				t.Errorf("Expected EPR[0][%d]=0 when Wp=0, got %f", j, rdpnl.EPR[0][j])
			}
		}
		for j := 0; j < rdpnl.rm[0].Nrp; j++ {
			if rdpnl.EPW[0][j] != 0.0 {
				t.Errorf("Expected EPW[0][%d]=0 when Wp=0, got %f", j, rdpnl.EPW[0][j])
			}
		}
	})

	t.Run("WpPositive_WallTypeP", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_P
		rdpnl.sd[0].mrk = '*' // 係数行列再作成フラグ

		Panelcf(rdpnl)

		// Wp>0かつWallType_Pの場合、係数が計算されるはず
		if rdpnl.FIp[0] == 0.0 {
			t.Error("FIp[0] should be calculated when Wp>0 and WallType_P")
		}
		if rdpnl.FOp[0] == 0.0 {
			t.Error("FOp[0] should be calculated when Wp>0 and WallType_P")
		}
		if rdpnl.FPp == 0.0 {
			t.Error("FPp should be calculated when Wp>0 and WallType_P")
		}
		if rdpnl.EPt[0] == 0.0 {
			t.Error("EPt[0] should be calculated when Wp>0 and WallType_P")
		}
		if rdpnl.Epw == 0.0 {
			t.Error("Epw should be calculated when Wp>0 and WallType_P")
		}
	})

	t.Run("WpPositive_WallTypeC", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C // 屋根一体型空気集熱器
		rdpnl.sd[0].mrk = '*'

		Panelcf(rdpnl)

		// WallType_Cの場合、異なる計算式が使用されるはず
		if rdpnl.FOp[0] == 0.0 {
			t.Error("FOp[0] should be calculated when Wp>0 and WallType_C")
		}
		if rdpnl.EPt[0] == 0.0 {
			t.Error("EPt[0] should be calculated when Wp>0 and WallType_C")
		}
	})

	t.Run("ChrRinputTrue", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mw.wall.WallType = WallType_C
		rdpnl.sd[0].mw.wall.chrRinput = true // 熱抵抗で入力
		rdpnl.sd[0].kd = 0.6                 // 表面固有の値
		rdpnl.sd[0].ku = 0.4
		rdpnl.sd[0].mrk = '*'

		Panelcf(rdpnl)

		// chrRinputがtrueの場合、sd.kdが使用されるはず
		// 計算が正常に完了することを確認
		if math.IsNaN(rdpnl.EPt[0]) || math.IsInf(rdpnl.EPt[0], 0) {
			t.Errorf("EPt[0] should be finite when chrRinput=true, got %f", rdpnl.EPt[0])
		}
	})

	t.Run("MC2_SharedWall", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.MC = 2 // 共用壁
		rdpnl.Wp = 0.5

		// 2番目の表面と室を設定
		rdpnl.sd[1] = &RMSRF{
			Name: "TestSurface2",
			A:    8.0,
			mrk:  '*',
			mw:   rdpnl.sd[0].mw, // 同じMWALLを共有
			ali:  7.0,
			alic: 6.5,
			RS:   80.0,
			WSR:  0.7,
		}

		rdpnl.rm[1] = &ROOM{
			Name: "TestRoom2",
			N:    3,
			mrk:  '*',
			alr:  make([]float64, 9),
			Ntr:  2,
			Nrp:  1,
			rsrf: make([]*RMSRF, 3),
		}

		// 2番目の室のalrとrsrfを初期化
		for i := 0; i < 9; i++ {
			rdpnl.rm[1].alr[i] = 0.08 + float64(i)*0.04
		}
		for i := 0; i < 3; i++ {
			rdpnl.rm[1].rsrf[i] = &RMSRF{
				WSR:  0.15 + float64(i)*0.08,
				WSRN: make([]float64, 2),
				WSPL: make([]float64, 1),
			}
			rdpnl.rm[1].rsrf[i].WSRN[0] = 0.08 + float64(i)*0.04
			rdpnl.rm[1].rsrf[i].WSRN[1] = 0.12 + float64(i)*0.04
			rdpnl.rm[1].rsrf[i].WSPL[0] = 0.16 + float64(i)*0.04
		}

		rdpnl.sd[0].mrk = '*'

		Panelcf(rdpnl)

		// 共用壁の場合、両側の係数が計算されるはず
		if rdpnl.FIp[1] == 0.0 {
			t.Error("FIp[1] should be calculated for shared wall")
		}
		if rdpnl.FOp[1] == 0.0 {
			t.Error("FOp[1] should be calculated for shared wall")
		}
		if rdpnl.EPt[1] == 0.0 {
			t.Error("EPt[1] should be calculated for shared wall")
		}
	})

	t.Run("MarkNotSet", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mrk = ' ' // マークが設定されていない
		rdpnl.sd[0].PCMflg = false

		// 初期値を設定
		initialEPt := rdpnl.EPt[0]
		initialEpw := rdpnl.Epw

		Panelcf(rdpnl)

		// マークが設定されていない場合、係数の再計算は行われないはず
		if rdpnl.EPt[0] != initialEPt {
			t.Error("EPt[0] should not change when mark is not set")
		}
		if rdpnl.Epw != initialEpw {
			t.Error("Epw should not change when mark is not set")
		}
	})

	t.Run("PCMFlag", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mrk = ' '     // マークは設定されていない
		rdpnl.sd[0].PCMflg = true // PCMフラグが立っている

		Panelcf(rdpnl)

		// PCMフラグが立っている場合、マークに関係なく計算されるはず
		if rdpnl.EPt[0] == 0.0 {
			t.Error("EPt[0] should be calculated when PCMflg is true")
		}
	})

	t.Run("EPRCalculation", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mrk = '*'

		Panelcf(rdpnl)

		// EPR配列の計算確認
		for j := 0; j < rdpnl.rm[0].Ntr; j++ {
			if math.IsNaN(rdpnl.EPR[0][j]) || math.IsInf(rdpnl.EPR[0][j], 0) {
				t.Errorf("EPR[0][%d] should be finite, got %f", j, rdpnl.EPR[0][j])
			}
		}
	})

	t.Run("EPWCalculation", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mrk = '*'

		Panelcf(rdpnl)

		// EPW配列の計算確認
		for j := 0; j < rdpnl.rm[0].Nrp; j++ {
			if math.IsNaN(rdpnl.EPW[0][j]) || math.IsInf(rdpnl.EPW[0][j], 0) {
				t.Errorf("EPW[0][%d] should be finite, got %f", j, rdpnl.EPW[0][j])
			}
		}
	})

	t.Run("EdgeCases", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 1e-10 // 非常に小さい値
		rdpnl.sd[0].mrk = '*'
		rdpnl.sd[0].A = 1e-6 // 非常に小さい面積

		Panelcf(rdpnl)

		// 極端な値でも計算が正常に完了することを確認
		if math.IsNaN(rdpnl.EPt[0]) || math.IsInf(rdpnl.EPt[0], 0) {
			t.Errorf("EPt[0] should be finite for edge case, got %f", rdpnl.EPt[0])
		}
		if math.IsNaN(rdpnl.Epw) || math.IsInf(rdpnl.Epw, 0) {
			t.Errorf("Epw should be finite for edge case, got %f", rdpnl.Epw)
		}
	})

	t.Run("ConsistencyCheck", func(t *testing.T) {
		rdpnl := createTestRDPNL()
		rdpnl.Wp = 0.5
		rdpnl.sd[0].mrk = '*'

		// 同じ条件で2回計算
		Panelcf(rdpnl)
		firstEPt := rdpnl.EPt[0]
		firstEpw := rdpnl.Epw

		rdpnl.sd[0].mrk = '*' // マークを再設定
		Panelcf(rdpnl)

		// 同じ条件なら同じ結果になるはず
		if math.Abs(rdpnl.EPt[0]-firstEPt) > 1e-10 {
			t.Errorf("EPt[0] should be consistent: first=%f, second=%f", firstEPt, rdpnl.EPt[0])
		}
		if math.Abs(rdpnl.Epw-firstEpw) > 1e-10 {
			t.Errorf("Epw should be consistent: first=%f, second=%f", firstEpw, rdpnl.Epw)
		}
	})
}
package eeslism

import (
	"bytes"
	"testing"
)

func TestOPIhor(t *testing.T) {
	// テスト用のダミーデータを作成
	// P_MENN, WDAT, bekt, XYZ 構造体はeeslism/MODEL.goとeeslism/wthrd.goからコピーまたは参照
	// 実際のテストでは、これらの構造体のフィールドに適切な値を設定する必要があります。

	// ダミーのP_MENN構造体
	mp := make([]*P_MENN, 1)
	mp[0] = &P_MENN{
		opname:  "TestPanel",
		faia:    0.5,
		faig:    0.5,
		refg:    0.2,
		sum:     0.1,
		faiwall: [500]float64{0: 0.1}, // ダミー値
		P:       []XYZ{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}, {X: 1, Y: 1, Z: 0}, {X: 0, Y: 1, Z: 0}},
		e:       XYZ{X: 0, Y: 0, Z: 1}, // 上向きの法線ベクトル
	}

	// ダミーのWDAT構造体
	Wd := &WDAT{
		T:    20.0,  // 気温
		Idn:  800.0, // 法線面直達日射
		Isky: 200.0, // 水平面天空日射
		RN:   -50.0, // 夜間輻射
		Rsky: 300.0, // 大気放射量
		Sh:   0.8,   // 太陽光線の方向余弦 (ns) - 日中
		Sw:   0.0,
		Ss:   0.0,
	}

	// ダミーのlp (P_MENNのスライス)
	lp := make([]*P_MENN, 1)
	lp[0] = &P_MENN{
		opname:  "TestLP",
		faia:    0.5,
		faig:    0.5,
		refg:    0.2,
		sum:     0.1,
		faiwall: [500]float64{0: 0.1}, // ダミー値
		P:       []XYZ{{X: 0, Y: 0, Z: 0}, {X: 1, Y: 0, Z: 0}, {X: 1, Y: 1, Z: 0}, {X: 0, Y: 1, Z: 0}},
		e:       XYZ{X: 0, Y: 0, Z: 1}, // 上向きの法線ベクトル
	}

	// ダミーのullp, ulmp (bektのスライス)
	ullp := make([]*bekt, 1)
	ulmp := make([]*bekt, 1)
	ullp[0] = &bekt{}
	ulmp[0] = &bekt{}

	// ダミーのgp (XYZの2次元スライス)
	gp := make([][]XYZ, 1)
	gp[0] = []XYZ{{X: 0.5, Y: 0.5, Z: 0}, {X: INAN, Y: INAN, Z: INAN}} // 終端マーカー

	// テストケース
	tests := []struct {
		name       string
		wd         *WDAT
		monten     int
		dayprn     bool
		expectedIw float64 // 期待されるIwの値 (簡略化のため、ここでは計算せず0.0)
		expectedRn float64 // 期待されるrnの値 (簡略化のため、ここでは計算せず0.0)
	}{
		{
			name:       "Daytime_Monten_Enabled_Dayprn_Enabled",
			wd:         Wd,
			monten:     1,
			dayprn:     true,
			expectedIw: 0.0, // 実際の計算結果に基づいて更新
			expectedRn: 0.0, // 実際の計算結果に基づいて更新
		},
		{
			name:       "Daytime_Monten_Disabled_Dayprn_Enabled",
			wd:         Wd,
			monten:     0,
			dayprn:     true,
			expectedIw: 0.0, // 実際の計算結果に基づいて更新
			expectedRn: 0.0, // 実際の計算結果に基づいて更新
		},
		{
			name:       "Nighttime_Monten_Enabled_Dayprn_Enabled",
			wd:         &WDAT{Sh: -0.1, T: 10.0, RN: -30.0, Rsky: 200.0}, // 夜間
			monten:     1,
			dayprn:     true,
			expectedIw: 0.0,
			expectedRn: 0.0,
		},
		{
			name:       "Nighttime_Monten_Disabled_Dayprn_Enabled",
			wd:         &WDAT{Sh: -0.1, T: 10.0, RN: -30.0, Rsky: 200.0}, // 夜間
			monten:     0,
			dayprn:     true,
			expectedIw: 0.0,
			expectedRn: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 出力バッファ
			var bufFp bytes.Buffer
			var bufFp1 bytes.Buffer

			// dayprnの値を設定 (グローバル変数または引数として渡す必要がある)
			// 現状、dayprnはグローバル変数として扱われているため、テスト前に設定
			dayprn = tt.dayprn

			// OPIhor関数を呼び出し
			OPIhor(&bufFp, &bufFp1, len(lp), len(mp), mp, lp, tt.wd, ullp, ulmp, gp, 1, tt.monten)

			// 結果の検証
			// ここでは、出力がエラーなく行われたことと、主要な計算結果が期待通りかを確認します。
			// 厳密な数値比較は、OPIhorの内部計算ロジックが複雑なため、別途詳細なテストケースが必要になります。

			// 日中のテストケースの場合、Idre, Idf, Iwが計算されていることを確認
			if tt.wd.Sh > 0 { // 日中
				if mp[0].Idre == 0.0 && mp[0].Idf == 0.0 && mp[0].Iw == 0.0 {
					t.Errorf("Daytime calculation failed: Idre, Idf, Iw are all zero")
				}
			} else { // 夜間
				if mp[0].Idre != 0.0 || mp[0].Idf != 0.0 || mp[0].Iw != 0.0 {
					t.Errorf("Nighttime calculation failed: Idre, Idf, Iw should be zero")
				}
			}

			// rnとReffが計算されていることを確認
			if mp[0].rn == 0.0 && mp[0].Reff == 0.0 {
				t.Errorf("rn and Reff are both zero")
			}

			// dayprnがtrueの場合、出力バッファに内容があることを確認
			if tt.dayprn {
				if bufFp.Len() == 0 && tt.wd.Sh > 0 { // 日中の場合
					t.Errorf("fp buffer is empty when dayprn is true and daytime")
				}
				if bufFp1.Len() == 0 {
					t.Errorf("fp1 buffer is empty when dayprn is true")
				}
			} else {
				if bufFp.Len() != 0 {
					t.Errorf("fp buffer is not empty when dayprn is false")
				}
				if bufFp1.Len() != 0 {
					t.Errorf("fp1 buffer is not empty when dayprn is false")
				}
			}

			// 厳密な数値比較の例 (必要に応じて追加)
			// if math.Abs(mp[0].Iw - tt.expectedIw) > 0.001 {
			//     t.Errorf("Iw mismatch. Expected %f, Got %f", tt.expectedIw, mp[0].Iw)
			// }
		})
	}
}

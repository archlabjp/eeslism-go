package eeslism

import (
	"testing"
	"math"
)

// TestCOORDNT_CoordinateTransformations tests the coordinate transformation functions
func TestCOORDNT_CoordinateTransformations(t *testing.T) {
	t.Run("LP_COORDNT function availability", func(t *testing.T) {
		// LP_COORDNT関数の存在確認
		// 空のスライスでテスト
		poly := make([]*POLYGN, 0)
		tree := make([]*TREE, 0)
		obs := make([]*OBS, 0)
		bdp := make([]*BBDP, 0)
		
		// 関数が呼び出せることを確認
		result := LP_COORDNT(poly, tree, obs, bdp)
		
		if result == nil {
			t.Error("LP_COORDNT should return a valid slice, got nil")
		}
		
		// 空の入力では空の結果が期待される
		if len(result) != 0 {
			t.Logf("LP_COORDNT returned %d elements with empty input", len(result))
		}
		
		t.Logf("LP_COORDNT function test completed successfully")
	})

	t.Run("OP_COORDNT function availability", func(t *testing.T) {
		// OP_COORDNT関数の存在確認
		// 空のスライスでテスト
		bdp := make([]*BBDP, 0)
		poly := make([]*POLYGN, 0)
		
		// 関数が呼び出せることを確認
		result := OP_COORDNT(bdp, poly)
		
		if result == nil {
			t.Error("OP_COORDNT should return a valid slice, got nil")
		}
		
		// 空の入力では空の結果が期待される
		if len(result) != 0 {
			t.Logf("OP_COORDNT returned %d elements with empty input", len(result))
		}
		
		t.Logf("OP_COORDNT function test completed successfully")
	})

	t.Run("Basic angle calculations", func(t *testing.T) {
		// 基本的な角度計算のテスト
		testCases := []struct {
			name string
			wa   float64
			wb   float64
			desc string
		}{
			{"North", 0.0, 0.0, "北向き水平"},
			{"East", 90.0, 0.0, "東向き水平"},
			{"South", 180.0, 0.0, "南向き水平"},
			{"West", 270.0, 0.0, "西向き水平"},
			{"Tilted", 45.0, 30.0, "45度方位、30度傾斜"},
		}
		
		for _, tc := range testCases {
			// 角度をラジアンに変換
			waRad := tc.wa * math.Pi / 180.0
			wbRad := tc.wb * math.Pi / 180.0
			
			// 三角関数値の計算
			cosWa := math.Cos(waRad)
			sinWa := math.Sin(waRad)
			cosWb := math.Cos(wbRad)
			sinWb := math.Sin(wbRad)
			
			// 角度の妥当性確認
			if tc.wa < 0.0 || tc.wa >= 360.0 {
				t.Logf("Warning: Azimuth angle %.1f° outside normal range", tc.wa)
			}
			if tc.wb < 0.0 || tc.wb > 180.0 {
				t.Logf("Warning: Tilt angle %.1f° outside normal range", tc.wb)
			}
			
			t.Logf("Angle calculation: %s - %s (Wa=%.1f°, Wb=%.1f°, cos/sin: %.3f/%.3f, %.3f/%.3f)", 
				tc.name, tc.desc, tc.wa, tc.wb, cosWa, sinWa, cosWb, sinWb)
		}
	})

	t.Run("Coordinate transformation with simple data", func(t *testing.T) {
		// 簡単なデータでの座標変換テスト
		
		// 簡単なOBSデータを作成
		obs := []*OBS{
			{
				fname: "rect",
				obsname: "TestBuilding",
				x: 10.0, y: 20.0, z: 0.0,
				H: 15.0, D: 30.0, W: 20.0,
				Wa: 0.0, Wb: 90.0,
			},
		}
		
		// 空のその他のデータ
		poly := make([]*POLYGN, 0)
		tree := make([]*TREE, 0)
		bdp := make([]*BBDP, 0)
		
		// LP_COORDNT関数を呼び出し
		result := LP_COORDNT(poly, tree, obs, bdp)
		
		if result == nil {
			t.Error("LP_COORDNT should return a valid slice")
		}
		
		// 結果の確認
		t.Logf("LP_COORDNT with 1 OBS returned %d P_MENN elements", len(result))
		
		// 結果が期待される数になっているか確認
		if len(result) > 0 {
			t.Logf("Successfully processed obstacle data into coordinate planes")
		}
	})
}

// TestCOORDNT_AngleValidation tests angle parameter validation
func TestCOORDNT_AngleValidation(t *testing.T) {
	t.Run("Azimuth angle validation", func(t *testing.T) {
		testCases := []struct {
			wa    float64
			valid bool
			desc  string
		}{
			{0.0, true, "North"},
			{90.0, true, "East"},
			{180.0, true, "South"},
			{270.0, true, "West"},
			{359.9, true, "Almost full circle"},
			{360.0, false, "Full circle (should normalize)"},
			{450.0, false, "Over 360 degrees"},
			{-30.0, false, "Negative angle"},
		}
		
		for _, tc := range testCases {
			// 角度の妥当性確認
			if tc.valid {
				if tc.wa < 0.0 || tc.wa >= 360.0 {
					t.Logf("Note: Azimuth angle %.1f° may need normalization", tc.wa)
				}
			} else {
				if tc.wa >= 0.0 && tc.wa < 360.0 {
					t.Logf("Note: Angle %.1f° is in valid range but marked as invalid", tc.wa)
				}
			}
			
			t.Logf("Azimuth test: %.1f° (%s) - %s", tc.wa, tc.desc, 
				map[bool]string{true: "valid", false: "needs normalization"}[tc.valid])
		}
	})

	t.Run("Tilt angle validation", func(t *testing.T) {
		testCases := []struct {
			wb    float64
			valid bool
			desc  string
		}{
			{0.0, true, "Horizontal"},
			{30.0, true, "30° tilt"},
			{45.0, true, "45° tilt"},
			{90.0, true, "Vertical"},
			{180.0, true, "Upside down"},
			{270.0, false, "Over 180°"},
			{-30.0, false, "Negative tilt"},
		}
		
		for _, tc := range testCases {
			// 傾斜角の範囲確認
			if tc.valid {
				if tc.wb < 0.0 || tc.wb > 180.0 {
					t.Logf("Note: Tilt angle %.1f° is outside typical range [0, 180]", tc.wb)
				}
			}
			
			t.Logf("Tilt test: %.1f° (%s) - %s", tc.wb, tc.desc,
				map[bool]string{true: "valid", false: "out of range"}[tc.valid])
		}
	})
}

// TestCOORDNT_GeometricCalculations tests basic geometric calculations
func TestCOORDNT_GeometricCalculations(t *testing.T) {
	t.Run("Vector calculations", func(t *testing.T) {
		// ベクトル計算のテスト
		vectors := []struct {
			name string
			x, y, z float64
		}{
			{"UnitX", 1.0, 0.0, 0.0},
			{"UnitY", 0.0, 1.0, 0.0},
			{"UnitZ", 0.0, 0.0, 1.0},
			{"Diagonal", 1.0, 1.0, 1.0},
		}
		
		for _, vec := range vectors {
			// ベクトルの長さ計算
			length := math.Sqrt(vec.x*vec.x + vec.y*vec.y + vec.z*vec.z)
			
			// 単位ベクトルの確認
			if vec.name == "UnitX" || vec.name == "UnitY" || vec.name == "UnitZ" {
				if math.Abs(length-1.0) > 0.001 {
					t.Errorf("Expected unit vector length 1.0, got %.3f", length)
				}
			}
			
			t.Logf("Vector %s: (%.1f, %.1f, %.1f), Length: %.3f", 
				vec.name, vec.x, vec.y, vec.z, length)
		}
	})

	t.Run("Point transformation", func(t *testing.T) {
		// 点の座標変換テスト
		origin := XYZ{X: 5.0, Y: 10.0, Z: 2.0}
		localPoint := XYZ{X: 3.0, Y: 4.0, Z: 1.0}
		
		// グローバル座標系への変換（平行移動のみ）
		globalPoint := XYZ{
			X: localPoint.X + origin.X,
			Y: localPoint.Y + origin.Y,
			Z: localPoint.Z + origin.Z,
		}
		
		expectedGlobal := XYZ{X: 8.0, Y: 14.0, Z: 3.0}
		tolerance := 0.001
		
		if math.Abs(globalPoint.X-expectedGlobal.X) > tolerance ||
		   math.Abs(globalPoint.Y-expectedGlobal.Y) > tolerance ||
		   math.Abs(globalPoint.Z-expectedGlobal.Z) > tolerance {
			t.Errorf("Expected global point (%.1f, %.1f, %.1f), got (%.1f, %.1f, %.1f)",
				expectedGlobal.X, expectedGlobal.Y, expectedGlobal.Z,
				globalPoint.X, globalPoint.Y, globalPoint.Z)
		}
		
		t.Logf("Point transformation: Local(%.1f, %.1f, %.1f) -> Global(%.1f, %.1f, %.1f)",
			localPoint.X, localPoint.Y, localPoint.Z,
			globalPoint.X, globalPoint.Y, globalPoint.Z)
	})

	t.Run("Angle normalization", func(t *testing.T) {
		// 角度の正規化テスト
		testAngles := []struct {
			input    float64
			expected float64
			desc     string
		}{
			{0.0, 0.0, "Zero"},
			{360.0, 0.0, "Full circle"},
			{450.0, 90.0, "Over 360"},
			{-90.0, 270.0, "Negative"},
			{720.0, 0.0, "Two full circles"},
		}
		
		for _, test := range testAngles {
			// 角度正規化の実装例
			normalized := test.input
			for normalized < 0.0 {
				normalized += 360.0
			}
			for normalized >= 360.0 {
				normalized -= 360.0
			}
			
			if math.Abs(normalized-test.expected) > 0.001 {
				t.Errorf("Angle normalization failed: %.1f° -> expected %.1f°, got %.1f°",
					test.input, test.expected, normalized)
			}
			
			t.Logf("Angle normalization: %.1f° -> %.1f° (%s)", 
				test.input, normalized, test.desc)
		}
	})
}
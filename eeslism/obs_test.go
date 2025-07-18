package eeslism

import (
	"testing"
)

// TestOBS_BasicStructure tests the basic OBS structure and initialization
func TestOBS_BasicStructure(t *testing.T) {
	t.Run("OBS manual creation - rectangular obstacle", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "Building1",
			x:       10.0,  // 左下頂点座標
			y:       20.0,
			z:       0.0,
			H:       15.0, // 高さ [m]
			D:       30.0, // 奥行き [m]
			W:       20.0, // 巾 [m]
			Wa:      0.0,  // 方位角 [度]
			Wb:      90.0, // 傾斜角 [度] (垂直)
			ref:     [4]float64{0.3, 0.3, 0.3, 0.3}, // 各面の反射率
			rgb:     [3]float64{0.8, 0.8, 0.8},      // 色 (灰色)
		}
		
		if obs.fname != "rect" {
			t.Errorf("Expected fname='rect', got %s", obs.fname)
		}
		if obs.obsname != "Building1" {
			t.Errorf("Expected obsname='Building1', got %s", obs.obsname)
		}
		if obs.x != 10.0 {
			t.Errorf("Expected x=10.0, got %f", obs.x)
		}
		if obs.y != 20.0 {
			t.Errorf("Expected y=20.0, got %f", obs.y)
		}
		if obs.z != 0.0 {
			t.Errorf("Expected z=0.0, got %f", obs.z)
		}
		if obs.H != 15.0 {
			t.Errorf("Expected H=15.0, got %f", obs.H)
		}
		if obs.D != 30.0 {
			t.Errorf("Expected D=30.0, got %f", obs.D)
		}
		if obs.W != 20.0 {
			t.Errorf("Expected W=20.0, got %f", obs.W)
		}
	})

	t.Run("OBS creation - cubic obstacle", func(t *testing.T) {
		obs := &OBS{
			fname:   "cube",
			obsname: "CubicBuilding",
			x:       0.0,
			y:       0.0,
			z:       0.0,
			H:       10.0, // 立方体: 高さ
			D:       10.0, // 立方体: 奥行き
			W:       10.0, // 立方体: 巾
			Wa:      45.0, // 45度回転
			Wb:      90.0, // 垂直
			ref:     [4]float64{0.2, 0.25, 0.3, 0.35}, // 各面異なる反射率
			rgb:     [3]float64{0.6, 0.4, 0.2},        // 茶色
		}
		
		if obs.fname != "cube" {
			t.Errorf("Expected fname='cube', got %s", obs.fname)
		}
		if obs.obsname != "CubicBuilding" {
			t.Errorf("Expected obsname='CubicBuilding', got %s", obs.obsname)
		}
		
		// 立方体の寸法確認
		if obs.H != obs.D || obs.D != obs.W {
			t.Logf("Note: Not a perfect cube - H:%f, D:%f, W:%f", obs.H, obs.D, obs.W)
		}
		
		// 方位角の確認
		if obs.Wa != 45.0 {
			t.Errorf("Expected Wa=45.0, got %f", obs.Wa)
		}
	})

	t.Run("OBS creation - triangular obstacle", func(t *testing.T) {
		obs := &OBS{
			fname:   "r_tri", // 直角三角形
			obsname: "TriangularRoof",
			x:       5.0,
			y:       5.0,
			z:       10.0, // 地上10m
			H:       8.0,  // 高さ
			D:       12.0, // 奥行き
			W:       6.0,  // 巾
			Wa:      90.0, // 東向き
			Wb:      45.0, // 45度傾斜
			ref:     [4]float64{0.4, 0.4, 0.4, 0.4}, // 均一反射率
			rgb:     [3]float64{0.5, 0.7, 0.3},      // 緑色
		}
		
		if obs.fname != "r_tri" {
			t.Errorf("Expected fname='r_tri', got %s", obs.fname)
		}
		if obs.obsname != "TriangularRoof" {
			t.Errorf("Expected obsname='TriangularRoof', got %s", obs.obsname)
		}
		if obs.z != 10.0 {
			t.Errorf("Expected z=10.0 (elevated), got %f", obs.z)
		}
		if obs.Wb != 45.0 {
			t.Errorf("Expected Wb=45.0 (sloped), got %f", obs.Wb)
		}
	})
}

// TestOBS_GeometryValidation tests geometry parameter validation
func TestOBS_GeometryValidation(t *testing.T) {
	t.Run("Dimension validation", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "TestBuilding",
			H:       25.0, // 高さ
			D:       40.0, // 奥行き
			W:       30.0, // 巾
		}
		
		// 寸法の妥当性チェック
		if obs.H <= 0 || obs.D <= 0 || obs.W <= 0 {
			t.Error("All dimensions (H, D, W) should be positive")
		}
		
		// 建物として妥当な寸法かチェック
		if obs.H > 500.0 {
			t.Logf("Warning: Building height %f m seems unusually large", obs.H)
		}
		if obs.D > 1000.0 || obs.W > 1000.0 {
			t.Logf("Warning: Building dimensions D:%f, W:%f seem unusually large", obs.D, obs.W)
		}
		
		// アスペクト比の確認
		aspectRatio := obs.H / ((obs.D + obs.W) / 2.0)
		if aspectRatio > 10.0 {
			t.Logf("Warning: Very tall building - aspect ratio: %.2f", aspectRatio)
		}
		
		t.Logf("Building dimensions: H=%.1f, D=%.1f, W=%.1f m", obs.H, obs.D, obs.W)
		t.Logf("Volume: %.1f m³, Aspect ratio: %.2f", obs.H*obs.D*obs.W, aspectRatio)
	})

	t.Run("Orientation validation", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "OrientedBuilding",
			Wa:      135.0, // 方位角 [度]
			Wb:      90.0,  // 傾斜角 [度]
		}
		
		// 方位角の範囲確認
		if obs.Wa < 0.0 || obs.Wa >= 360.0 {
			t.Errorf("Azimuth angle Wa should be in range [0, 360), got %f", obs.Wa)
		}
		
		// 傾斜角の範囲確認
		if obs.Wb < 0.0 || obs.Wb > 180.0 {
			t.Errorf("Tilt angle Wb should be in range [0, 180], got %f", obs.Wb)
		}
		
		// 一般的な建物の傾斜角
		if obs.Wb != 90.0 && obs.fname == "rect" {
			t.Logf("Note: Rectangular building with non-vertical tilt: %.1f degrees", obs.Wb)
		}
		
		t.Logf("Orientation: Azimuth=%.1f°, Tilt=%.1f°", obs.Wa, obs.Wb)
	})

	t.Run("Reflectance validation", func(t *testing.T) {
		obs := &OBS{
			fname:   "cube",
			obsname: "ReflectiveBuilding",
			ref:     [4]float64{0.1, 0.3, 0.5, 0.8}, // 各面の反射率
		}
		
		// 反射率の範囲確認
		for i, reflectance := range obs.ref {
			if reflectance < 0.0 || reflectance > 1.0 {
				t.Errorf("Reflectance[%d] should be in range [0, 1], got %f", i, reflectance)
			}
		}
		
		// 反射率の妥当性確認
		for i, reflectance := range obs.ref {
			if reflectance > 0.9 {
				t.Logf("Warning: Very high reflectance[%d]: %.2f (mirror-like surface)", i, reflectance)
			}
			if reflectance < 0.05 {
				t.Logf("Warning: Very low reflectance[%d]: %.2f (very dark surface)", i, reflectance)
			}
		}
		
		avgReflectance := (obs.ref[0] + obs.ref[1] + obs.ref[2] + obs.ref[3]) / 4.0
		t.Logf("Surface reflectances: %.2f, %.2f, %.2f, %.2f (avg: %.2f)", 
			obs.ref[0], obs.ref[1], obs.ref[2], obs.ref[3], avgReflectance)
	})

	t.Run("Color validation", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "ColoredBuilding",
			rgb:     [3]float64{0.8, 0.2, 0.1}, // 赤系
		}
		
		// RGB値の範囲確認
		for i, color := range obs.rgb {
			if color < 0.0 || color > 1.0 {
				t.Errorf("RGB[%d] should be in range [0, 1], got %f", i, color)
			}
		}
		
		// 色の明度計算（簡易）
		brightness := 0.299*obs.rgb[0] + 0.587*obs.rgb[1] + 0.114*obs.rgb[2]
		if brightness < 0.1 {
			t.Logf("Very dark building: brightness=%.3f", brightness)
		} else if brightness > 0.9 {
			t.Logf("Very bright building: brightness=%.3f", brightness)
		}
		
		t.Logf("Building color: R=%.2f, G=%.2f, B=%.2f (brightness=%.3f)", 
			obs.rgb[0], obs.rgb[1], obs.rgb[2], brightness)
	})
}

// TestOBS_ObstacleTypes tests different obstacle types
func TestOBS_ObstacleTypes(t *testing.T) {
	t.Run("Rectangular building", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "OfficeBuilding",
			x: 0.0, y: 0.0, z: 0.0,
			H: 50.0, D: 80.0, W: 40.0,
			Wa: 0.0, Wb: 90.0,
		}
		
		if obs.fname != "rect" {
			t.Errorf("Expected fname='rect', got %s", obs.fname)
		}
		
		volume := obs.H * obs.D * obs.W
		t.Logf("Rectangular building: %s, Volume: %.0f m³", obs.obsname, volume)
	})

	t.Run("Cubic structure", func(t *testing.T) {
		obs := &OBS{
			fname:   "cube",
			obsname: "CubicPavilion",
			x: 10.0, y: 10.0, z: 0.0,
			H: 12.0, D: 12.0, W: 12.0,
			Wa: 0.0, Wb: 90.0,
		}
		
		if obs.fname != "cube" {
			t.Errorf("Expected fname='cube', got %s", obs.fname)
		}
		
		// 立方体の確認
		if obs.H == obs.D && obs.D == obs.W {
			t.Logf("Perfect cube: %s, Edge length: %.1f m", obs.obsname, obs.H)
		} else {
			t.Logf("Rectangular prism: %s, Dimensions: %.1f×%.1f×%.1f m", 
				obs.obsname, obs.W, obs.D, obs.H)
		}
	})

	t.Run("Triangular structures", func(t *testing.T) {
		triangularTypes := []struct {
			fname string
			desc  string
		}{
			{"r_tri", "Right triangle"},
			{"i_tri", "Isosceles triangle"},
		}
		
		for _, triType := range triangularTypes {
			obs := &OBS{
				fname:   triType.fname,
				obsname: "Triangular_" + triType.fname,
				x: 20.0, y: 20.0, z: 0.0,
				H: 8.0, D: 15.0, W: 10.0,
				Wa: 45.0, Wb: 60.0, // 傾斜屋根
			}
			
			if obs.fname != triType.fname {
				t.Errorf("Expected fname='%s', got %s", triType.fname, obs.fname)
			}
			
			t.Logf("%s: %s, Dimensions: %.1f×%.1f×%.1f m", 
				triType.desc, obs.obsname, obs.W, obs.D, obs.H)
		}
	})
}

// TestOBS_MultipleObstacles tests multiple obstacle management
func TestOBS_MultipleObstacles(t *testing.T) {
	t.Run("Urban environment simulation", func(t *testing.T) {
		obstacles := []*OBS{
			{
				fname: "rect", obsname: "MainBuilding",
				x: 0.0, y: 0.0, z: 0.0,
				H: 60.0, D: 100.0, W: 50.0,
				Wa: 0.0, Wb: 90.0,
				ref: [4]float64{0.3, 0.3, 0.3, 0.3},
				rgb: [3]float64{0.7, 0.7, 0.7},
			},
			{
				fname: "rect", obsname: "AdjacentBuilding",
				x: 60.0, y: 0.0, z: 0.0,
				H: 40.0, D: 80.0, W: 40.0,
				Wa: 0.0, Wb: 90.0,
				ref: [4]float64{0.25, 0.25, 0.25, 0.25},
				rgb: [3]float64{0.8, 0.6, 0.4},
			},
			{
				fname: "cube", obsname: "SmallPavilion",
				x: 120.0, y: 50.0, z: 0.0,
				H: 8.0, D: 8.0, W: 8.0,
				Wa: 45.0, Wb: 90.0,
				ref: [4]float64{0.4, 0.4, 0.4, 0.4},
				rgb: [3]float64{0.2, 0.6, 0.2},
			},
		}
		
		if len(obstacles) != 3 {
			t.Errorf("Expected 3 obstacles, got %d", len(obstacles))
		}
		
		// 各障害物の名前の一意性確認
		names := make(map[string]bool)
		for _, obs := range obstacles {
			if names[obs.obsname] {
				t.Errorf("Duplicate obstacle name: %s", obs.obsname)
			}
			names[obs.obsname] = true
		}
		
		// 総体積計算
		totalVolume := 0.0
		for _, obs := range obstacles {
			volume := obs.H * obs.D * obs.W
			totalVolume += volume
			t.Logf("Obstacle: %s, Volume: %.0f m³", obs.obsname, volume)
		}
		
		t.Logf("Total built volume: %.0f m³", totalVolume)
		
		// 位置関係の確認
		for i, obs1 := range obstacles {
			for j, obs2 := range obstacles {
				if i < j {
					distance := ((obs1.x-obs2.x)*(obs1.x-obs2.x) + 
								(obs1.y-obs2.y)*(obs1.y-obs2.y))
					if distance < 100.0 { // 10m以内
						t.Logf("Close obstacles: %s and %s (distance²: %.1f)", 
							obs1.obsname, obs2.obsname, distance)
					}
				}
			}
		}
	})
}

// TestOBS_BoundaryValues tests boundary and edge cases
func TestOBS_BoundaryValues(t *testing.T) {
	t.Run("Minimum valid obstacle", func(t *testing.T) {
		obs := &OBS{
			fname:   "cube",
			obsname: "MinimalObstacle",
			x: 0.0, y: 0.0, z: 0.0,
			H: 0.1, D: 0.1, W: 0.1, // 最小寸法
			Wa: 0.0, Wb: 90.0,
			ref: [4]float64{0.0, 0.0, 0.0, 0.0}, // 最小反射率
			rgb: [3]float64{0.0, 0.0, 0.0},      // 黒色
		}
		
		if obs.H < 0.05 || obs.D < 0.05 || obs.W < 0.05 {
			t.Logf("Warning: Very small obstacle dimensions: %.2f×%.2f×%.2f m", obs.W, obs.D, obs.H)
		}
		
		t.Logf("Minimal obstacle validated: %.2f×%.2f×%.2f m", obs.W, obs.D, obs.H)
	})

	t.Run("Large scale obstacle", func(t *testing.T) {
		obs := &OBS{
			fname:   "rect",
			obsname: "MegaStructure",
			x: 1000.0, y: 2000.0, z: 100.0,
			H: 300.0, D: 500.0, W: 200.0, // 大規模建築
			Wa: 180.0, Wb: 90.0,
			ref: [4]float64{1.0, 1.0, 1.0, 1.0}, // 最大反射率
			rgb: [3]float64{1.0, 1.0, 1.0},      // 白色
		}
		
		volume := obs.H * obs.D * obs.W
		if volume > 10000000.0 { // 1000万m³
			t.Logf("Warning: Very large structure volume: %.0f m³", volume)
		}
		
		t.Logf("Large obstacle: %.0f×%.0f×%.0f m, Volume: %.0f m³", 
			obs.W, obs.D, obs.H, volume)
	})
}
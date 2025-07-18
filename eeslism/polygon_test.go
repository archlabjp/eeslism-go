package eeslism

import (
	"testing"
	"math"
)

// TestPOLYGN_BasicStructure tests the basic POLYGN structure and initialization
func TestPOLYGN_BasicStructure(t *testing.T) {
	t.Run("POLYGN creation - RMP (Room Panel)", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "WallPanel1",
			wallname: "ExteriorWall",
			polyd:    4, // 四角形
			ref:      0.3,
			refg:     0.2,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},   // 左下
				{X: 4.0, Y: 0.0, Z: 0.0},   // 右下
				{X: 4.0, Y: 0.0, Z: 3.0},   // 右上
				{X: 0.0, Y: 0.0, Z: 3.0},   // 左上
			},
			grpx: 1.0,
			rgb:  [3]float64{0.8, 0.7, 0.6}, // ベージュ色
		}
		
		if polygon.polyknd != "RMP" {
			t.Errorf("Expected polyknd='RMP', got %s", polygon.polyknd)
		}
		if polygon.polyname != "WallPanel1" {
			t.Errorf("Expected polyname='WallPanel1', got %s", polygon.polyname)
		}
		if polygon.wallname != "ExteriorWall" {
			t.Errorf("Expected wallname='ExteriorWall', got %s", polygon.wallname)
		}
		if polygon.polyd != 4 {
			t.Errorf("Expected polyd=4, got %d", polygon.polyd)
		}
		if len(polygon.P) != 4 {
			t.Errorf("Expected 4 vertices, got %d", len(polygon.P))
		}
		if polygon.ref != 0.3 {
			t.Errorf("Expected ref=0.3, got %f", polygon.ref)
		}
		if polygon.refg != 0.2 {
			t.Errorf("Expected refg=0.2, got %f", polygon.refg)
		}
	})

	t.Run("POLYGN creation - OBS (Obstacle)", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "OBS",
			polyname: "BuildingFacade",
			wallname: "",
			polyd:    6, // 六角形
			ref:      0.25,
			refg:     0.15,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 2.0, Y: 0.0, Z: 0.0},
				{X: 3.0, Y: 1.5, Z: 0.0},
				{X: 2.0, Y: 3.0, Z: 0.0},
				{X: 0.0, Y: 3.0, Z: 0.0},
				{X: -1.0, Y: 1.5, Z: 0.0},
			},
			grpx: 2.0,
			rgb:  [3]float64{0.5, 0.5, 0.7}, // 青灰色
		}
		
		if polygon.polyknd != "OBS" {
			t.Errorf("Expected polyknd='OBS', got %s", polygon.polyknd)
		}
		if polygon.polyname != "BuildingFacade" {
			t.Errorf("Expected polyname='BuildingFacade', got %s", polygon.polyname)
		}
		if polygon.polyd != 6 {
			t.Errorf("Expected polyd=6, got %d", polygon.polyd)
		}
		if len(polygon.P) != 6 {
			t.Errorf("Expected 6 vertices, got %d", len(polygon.P))
		}
		if polygon.grpx != 2.0 {
			t.Errorf("Expected grpx=2.0, got %f", polygon.grpx)
		}
	})

	t.Run("POLYGN creation - Triangle", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "TriangularPanel",
			wallname: "RoofPanel",
			polyd:    3, // 三角形
			ref:      0.4,
			refg:     0.1,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 5.0, Y: 0.0, Z: 0.0},
				{X: 2.5, Y: 0.0, Z: 4.0},
			},
			grpx: 1.5,
			rgb:  [3]float64{0.6, 0.3, 0.2}, // 茶色
		}
		
		if polygon.polyknd != "RMP" {
			t.Errorf("Expected polyknd='RMP', got %s", polygon.polyknd)
		}
		if polygon.polyd != 3 {
			t.Errorf("Expected polyd=3 (triangle), got %d", polygon.polyd)
		}
		if len(polygon.P) != 3 {
			t.Errorf("Expected 3 vertices for triangle, got %d", len(polygon.P))
		}
	})
}

// TestPOLYGN_GeometryValidation tests geometry parameter validation
func TestPOLYGN_GeometryValidation(t *testing.T) {
	t.Run("Vertex count validation", func(t *testing.T) {
		testCases := []struct {
			name     string
			polyd    int
			vertices int
			valid    bool
		}{
			{"Triangle", 3, 3, true},
			{"Quadrilateral", 4, 4, true},
			{"Pentagon", 5, 5, true},
			{"Hexagon", 6, 6, true},
			{"Mismatch", 4, 3, false},
		}
		
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				vertices := make([]XYZ, tc.vertices)
				for i := range vertices {
					angle := float64(i) * 2.0 * math.Pi / float64(tc.vertices)
					vertices[i] = XYZ{
						X: math.Cos(angle),
						Y: math.Sin(angle),
						Z: 0.0,
					}
				}
				
				polygon := &POLYGN{
					polyknd:  "RMP",
					polyname: tc.name + "Test",
					polyd:    tc.polyd,
					P:        vertices,
				}
				
				// 頂点数の整合性チェック
				if tc.valid {
					if polygon.polyd != len(polygon.P) {
						t.Errorf("Vertex count mismatch: polyd=%d, actual vertices=%d", 
							polygon.polyd, len(polygon.P))
					}
				} else {
					if polygon.polyd == len(polygon.P) {
						t.Logf("Note: This case should have mismatched vertex count")
					}
				}
				
				// 多角形の妥当性
				if polygon.polyd < 3 {
					t.Errorf("Polygon should have at least 3 vertices, got %d", polygon.polyd)
				}
				if polygon.polyd > 6 {
					t.Logf("Warning: Complex polygon with %d vertices", polygon.polyd)
				}
			})
		}
	})

	t.Run("Reflectance validation", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "ReflectanceTest",
			polyd:    4,
			ref:      0.45, // 表面反射率
			refg:     0.25, // 前面地面反射率
		}
		
		// 反射率の範囲確認
		if polygon.ref < 0.0 || polygon.ref > 1.0 {
			t.Errorf("Surface reflectance should be in range [0, 1], got %f", polygon.ref)
		}
		if polygon.refg < 0.0 || polygon.refg > 1.0 {
			t.Errorf("Ground reflectance should be in range [0, 1], got %f", polygon.refg)
		}
		
		// 反射率の妥当性
		if polygon.ref > 0.9 {
			t.Logf("Warning: Very high surface reflectance: %.2f", polygon.ref)
		}
		if polygon.refg > 0.5 {
			t.Logf("Warning: High ground reflectance: %.2f", polygon.refg)
		}
		
		t.Logf("Reflectances: Surface=%.2f, Ground=%.2f", polygon.ref, polygon.refg)
	})

	t.Run("Color validation", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "OBS",
			polyname: "ColorTest",
			rgb:      [3]float64{0.2, 0.8, 0.3}, // 緑色
		}
		
		// RGB値の範囲確認
		for i, color := range polygon.rgb {
			if color < 0.0 || color > 1.0 {
				t.Errorf("RGB[%d] should be in range [0, 1], got %f", i, color)
			}
		}
		
		// 色の明度計算
		brightness := 0.299*polygon.rgb[0] + 0.587*polygon.rgb[1] + 0.114*polygon.rgb[2]
		
		t.Logf("Polygon color: R=%.2f, G=%.2f, B=%.2f (brightness=%.3f)", 
			polygon.rgb[0], polygon.rgb[1], polygon.rgb[2], brightness)
	})

	t.Run("Ground point distance validation", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "GroundPointTest",
			grpx:     3.5, // 前面地面代表点までの距離
		}
		
		if polygon.grpx <= 0.0 {
			t.Error("Ground point distance (grpx) should be positive")
		}
		if polygon.grpx > 100.0 {
			t.Logf("Warning: Very large ground point distance: %.1f m", polygon.grpx)
		}
		
		t.Logf("Ground point distance: %.1f m", polygon.grpx)
	})
}

// TestPOLYGN_PolygonTypes tests different polygon types
func TestPOLYGN_PolygonTypes(t *testing.T) {
	t.Run("RMP type polygons", func(t *testing.T) {
		rmpPolygons := []POLYGN{
			{
				polyknd: "RMP", polyname: "WallPanel", wallname: "ExteriorWall",
				polyd: 4, ref: 0.3, refg: 0.2,
			},
			{
				polyknd: "RMP", polyname: "RoofPanel", wallname: "Roof",
				polyd: 3, ref: 0.4, refg: 0.15,
			},
		}
		
		for _, polygon := range rmpPolygons {
			if polygon.polyknd != "RMP" {
				t.Errorf("Expected RMP type, got %s", polygon.polyknd)
			}
			if polygon.wallname == "" {
				t.Error("RMP polygon should have wallname")
			}
			
			t.Logf("RMP: %s (wall: %s, %d vertices)", 
				polygon.polyname, polygon.wallname, polygon.polyd)
		}
	})

	t.Run("OBS type polygons", func(t *testing.T) {
		obsPolygons := []POLYGN{
			{
				polyknd: "OBS", polyname: "BuildingA", wallname: "",
				polyd: 4, ref: 0.25, refg: 0.2,
			},
			{
				polyknd: "OBS", polyname: "ComplexBuilding", wallname: "",
				polyd: 6, ref: 0.3, refg: 0.18,
			},
		}
		
		for _, polygon := range obsPolygons {
			if polygon.polyknd != "OBS" {
				t.Errorf("Expected OBS type, got %s", polygon.polyknd)
			}
			// OBSタイプはwallnameが空でも良い
			
			t.Logf("OBS: %s (%d vertices)", polygon.polyname, polygon.polyd)
		}
	})
}

// TestPOLYGN_GeometricCalculations tests basic geometric calculations
func TestPOLYGN_GeometricCalculations(t *testing.T) {
	t.Run("Rectangle area calculation", func(t *testing.T) {
		// 4m × 3m の長方形
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "Rectangle",
			polyd:    4,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 4.0, Y: 0.0, Z: 0.0},
				{X: 4.0, Y: 0.0, Z: 3.0},
				{X: 0.0, Y: 0.0, Z: 3.0},
			},
		}
		
		// 簡易面積計算（長方形の場合）
		width := polygon.P[1].X - polygon.P[0].X
		height := polygon.P[2].Z - polygon.P[1].Z
		area := width * height
		
		expectedArea := 12.0 // 4m × 3m
		if math.Abs(area-expectedArea) > 0.001 {
			t.Errorf("Expected area %.1f, got %.1f", expectedArea, area)
		}
		
		t.Logf("Rectangle: %.1f × %.1f m, Area: %.1f m²", width, height, area)
	})

	t.Run("Triangle area calculation", func(t *testing.T) {
		// 直角三角形 (底辺4m、高さ3m)
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "Triangle",
			polyd:    3,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 4.0, Y: 0.0, Z: 0.0},
				{X: 0.0, Y: 0.0, Z: 3.0},
			},
		}
		
		// 直角三角形の面積計算
		base := polygon.P[1].X - polygon.P[0].X
		height := polygon.P[2].Z - polygon.P[0].Z
		area := 0.5 * base * height
		
		expectedArea := 6.0 // 0.5 × 4m × 3m
		if math.Abs(area-expectedArea) > 0.001 {
			t.Errorf("Expected area %.1f, got %.1f", expectedArea, area)
		}
		
		t.Logf("Triangle: base=%.1f m, height=%.1f m, Area: %.1f m²", base, height, area)
	})

	t.Run("Polygon centroid calculation", func(t *testing.T) {
		// 正方形の重心計算
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "Square",
			polyd:    4,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 2.0, Y: 0.0, Z: 0.0},
				{X: 2.0, Y: 0.0, Z: 2.0},
				{X: 0.0, Y: 0.0, Z: 2.0},
			},
		}
		
		// 重心計算（単純平均）
		var centroid XYZ
		for _, vertex := range polygon.P {
			centroid.X += vertex.X
			centroid.Y += vertex.Y
			centroid.Z += vertex.Z
		}
		centroid.X /= float64(len(polygon.P))
		centroid.Y /= float64(len(polygon.P))
		centroid.Z /= float64(len(polygon.P))
		
		expectedCentroid := XYZ{X: 1.0, Y: 0.0, Z: 1.0}
		tolerance := 0.001
		
		if math.Abs(centroid.X-expectedCentroid.X) > tolerance ||
		   math.Abs(centroid.Y-expectedCentroid.Y) > tolerance ||
		   math.Abs(centroid.Z-expectedCentroid.Z) > tolerance {
			t.Errorf("Expected centroid (%.1f, %.1f, %.1f), got (%.1f, %.1f, %.1f)",
				expectedCentroid.X, expectedCentroid.Y, expectedCentroid.Z,
				centroid.X, centroid.Y, centroid.Z)
		}
		
		t.Logf("Square centroid: (%.1f, %.1f, %.1f)", centroid.X, centroid.Y, centroid.Z)
	})
}

// TestPOLYGN_ComplexPolygons tests complex polygon configurations
func TestPOLYGN_ComplexPolygons(t *testing.T) {
	t.Run("Pentagon building facade", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "OBS",
			polyname: "PentagonBuilding",
			polyd:    5,
			ref:      0.35,
			refg:     0.2,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 3.0, Y: 0.0, Z: 0.0},
				{X: 4.5, Y: 0.0, Z: 2.0},
				{X: 1.5, Y: 0.0, Z: 4.0},
				{X: -1.5, Y: 0.0, Z: 2.0},
			},
			grpx: 2.5,
			rgb:  [3]float64{0.7, 0.6, 0.5},
		}
		
		if polygon.polyd != 5 {
			t.Errorf("Expected pentagon (5 vertices), got %d", polygon.polyd)
		}
		if len(polygon.P) != 5 {
			t.Errorf("Expected 5 vertices, got %d", len(polygon.P))
		}
		
		// 頂点の妥当性確認
		for i, vertex := range polygon.P {
			if vertex.Y != 0.0 {
				t.Errorf("Vertex %d should be on Y=0 plane, got Y=%f", i, vertex.Y)
			}
		}
		
		t.Logf("Pentagon building with %d vertices validated", len(polygon.P))
	})

	t.Run("Hexagonal skylight", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "HexagonalSkylight",
			wallname: "Skylight",
			polyd:    6,
			ref:      0.1, // 低反射率（透明に近い）
			refg:     0.2,
			grpx:     1.0,
			rgb:      [3]float64{0.8, 0.9, 1.0}, // 薄青色
		}
		
		// 正六角形の頂点を生成
		radius := 2.0
		for i := 0; i < 6; i++ {
			angle := float64(i) * math.Pi / 3.0 // 60度ずつ
			vertex := XYZ{
				X: radius * math.Cos(angle),
				Y: radius * math.Sin(angle),
				Z: 5.0, // 高さ5mに設置
			}
			polygon.P = append(polygon.P, vertex)
		}
		
		if polygon.polyd != 6 {
			t.Errorf("Expected hexagon (6 vertices), got %d", polygon.polyd)
		}
		if len(polygon.P) != 6 {
			t.Errorf("Expected 6 vertices, got %d", len(polygon.P))
		}
		
		// スカイライトとしての妥当性
		if polygon.ref > 0.3 {
			t.Logf("Warning: High reflectance for skylight: %.2f", polygon.ref)
		}
		
		t.Logf("Hexagonal skylight: radius=%.1f m, height=%.1f m", radius, polygon.P[0].Z)
	})
}

// TestPOLYGN_BoundaryValues tests boundary and edge cases
func TestPOLYGN_BoundaryValues(t *testing.T) {
	t.Run("Minimum triangle", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "RMP",
			polyname: "MinTriangle",
			polyd:    3,
			ref:      0.0, // 最小反射率
			refg:     0.0,
			P: []XYZ{
				{X: 0.0, Y: 0.0, Z: 0.0},
				{X: 0.1, Y: 0.0, Z: 0.0},
				{X: 0.0, Y: 0.0, Z: 0.1},
			},
			grpx: 0.1,
			rgb:  [3]float64{0.0, 0.0, 0.0}, // 黒色
		}
		
		if polygon.polyd != 3 {
			t.Errorf("Expected triangle (3 vertices), got %d", polygon.polyd)
		}
		if len(polygon.P) != 3 {
			t.Errorf("Expected 3 vertices, got %d", len(polygon.P))
		}
		
		// 最小三角形の面積
		area := 0.5 * 0.1 * 0.1 // 0.005 m²
		if area < 0.001 {
			t.Logf("Very small triangle area: %.6f m²", area)
		}
		
		t.Logf("Minimum triangle validated: area=%.6f m²", area)
	})

	t.Run("Maximum complexity polygon", func(t *testing.T) {
		polygon := &POLYGN{
			polyknd:  "OBS",
			polyname: "ComplexBuilding",
			polyd:    6, // 最大複雑度
			ref:      1.0, // 最大反射率
			refg:     1.0,
			grpx:     100.0, // 大きな距離
			rgb:      [3]float64{1.0, 1.0, 1.0}, // 白色
		}
		
		// 大きな六角形を生成
		radius := 50.0
		for i := 0; i < 6; i++ {
			angle := float64(i) * math.Pi / 3.0
			vertex := XYZ{
				X: radius * math.Cos(angle),
				Y: radius * math.Sin(angle),
				Z: 100.0, // 高層建築
			}
			polygon.P = append(polygon.P, vertex)
		}
		
		if polygon.ref != 1.0 || polygon.refg != 1.0 {
			t.Errorf("Expected maximum reflectance (1.0), got ref=%f, refg=%f", 
				polygon.ref, polygon.refg)
		}
		
		if len(polygon.P) != 6 {
			t.Errorf("Expected 6 vertices, got %d", len(polygon.P))
		}
		
		t.Logf("Maximum complexity polygon: radius=%.1f m, height=%.1f m", 
			radius, polygon.P[0].Z)
	})
}
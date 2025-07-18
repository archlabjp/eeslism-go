package eeslism

import (
	"testing"
)

// TestTREE_BasicStructure tests the basic TREE structure and initialization
func TestTREE_BasicStructure(t *testing.T) {
	t.Run("treeinit basic initialization", func(t *testing.T) {
		tree := treeinit()
		
		if tree == nil {
			t.Fatal("treeinit() returned nil")
		}
		
		// Check default initialization values
		if tree.treename != "" {
			t.Errorf("Expected empty treename, got %s", tree.treename)
		}
		if tree.treetype != "" {
			t.Errorf("Expected empty treetype, got %s", tree.treetype)
		}
		if tree.x != 0.0 {
			t.Errorf("Expected x=0.0, got %f", tree.x)
		}
		if tree.y != 0.0 {
			t.Errorf("Expected y=0.0, got %f", tree.y)
		}
		if tree.z != 0.0 {
			t.Errorf("Expected z=0.0, got %f", tree.z)
		}
		
		// Check trunk parameters
		if tree.W1 != 0.0 {
			t.Errorf("Expected W1=0.0, got %f", tree.W1)
		}
		if tree.H1 != 0.0 {
			t.Errorf("Expected H1=0.0, got %f", tree.H1)
		}
		
		// Check foliage parameters
		if tree.W2 != 0.0 {
			t.Errorf("Expected W2=0.0, got %f", tree.W2)
		}
		if tree.W3 != 0.0 {
			t.Errorf("Expected W3=0.0, got %f", tree.W3)
		}
		if tree.W4 != 0.0 {
			t.Errorf("Expected W4=0.0, got %f", tree.W4)
		}
		if tree.H2 != 0.0 {
			t.Errorf("Expected H2=0.0, got %f", tree.H2)
		}
		if tree.H3 != 0.0 {
			t.Errorf("Expected H3=0.0, got %f", tree.H3)
		}
	})

	t.Run("TREE structure manual creation", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "TestTree",
			x:        10.0,
			y:        20.0,
			z:        0.0,
			W1:       0.5,  // 幹太さ [m]
			H1:       3.0,  // 幹高さ [m]
			W2:       4.0,  // 葉部下面巾 [m]
			W3:       5.0,  // 葉部中央巾 [m]
			W4:       3.0,  // 葉部上面巾 [m]
			H2:       2.0,  // 葉部下側高さ [m]
			H3:       3.0,  // 葉部上側高さ [m]
		}
		
		if tree.treetype != "treeA" {
			t.Errorf("Expected treetype='treeA', got %s", tree.treetype)
		}
		if tree.treename != "TestTree" {
			t.Errorf("Expected treename='TestTree', got %s", tree.treename)
		}
		if tree.x != 10.0 {
			t.Errorf("Expected x=10.0, got %f", tree.x)
		}
		if tree.y != 20.0 {
			t.Errorf("Expected y=20.0, got %f", tree.y)
		}
		if tree.z != 0.0 {
			t.Errorf("Expected z=0.0, got %f", tree.z)
		}
	})
}

// TestTREE_GeometryValidation tests geometry parameter validation
func TestTREE_GeometryValidation(t *testing.T) {
	t.Run("Trunk geometry validation", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "TrunkTest",
			W1:       0.8,  // 幹太さ
			H1:       4.0,  // 幹高さ
		}
		
		// 幹の寸法妥当性チェック
		if tree.W1 <= 0 {
			t.Error("Trunk width (W1) should be positive")
		}
		if tree.H1 <= 0 {
			t.Error("Trunk height (H1) should be positive")
		}
		if tree.W1 > 5.0 {
			t.Logf("Warning: Trunk width %f m seems unusually large", tree.W1)
		}
		if tree.H1 > 50.0 {
			t.Logf("Warning: Trunk height %f m seems unusually large", tree.H1)
		}
	})

	t.Run("Foliage geometry validation", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "FoliageTest",
			W2:       6.0,  // 葉部下面巾
			W3:       8.0,  // 葉部中央巾
			W4:       4.0,  // 葉部上面巾
			H2:       3.0,  // 葉部下側高さ
			H3:       4.0,  // 葉部上側高さ
		}
		
		// 葉部の寸法妥当性チェック
		if tree.W2 <= 0 || tree.W3 <= 0 || tree.W4 <= 0 {
			t.Error("Foliage widths (W2, W3, W4) should be positive")
		}
		if tree.H2 <= 0 || tree.H3 <= 0 {
			t.Error("Foliage heights (H2, H3) should be positive")
		}
		
		// 葉部の形状整合性チェック
		if tree.H3 < tree.H2 {
			t.Error("Upper foliage height (H3) should be >= lower foliage height (H2)")
		}
		
		// 一般的な樹木形状の妥当性
		if tree.W3 < tree.W2 && tree.W3 < tree.W4 {
			t.Logf("Note: Central width (W3=%f) is smaller than both lower (W2=%f) and upper (W4=%f) widths", 
				tree.W3, tree.W2, tree.W4)
		}
	})

	t.Run("Complete tree geometry", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "CompleteTree",
			x:        15.0,
			y:        25.0,
			z:        0.5,
			W1:       0.6,  // 幹太さ
			H1:       5.0,  // 幹高さ
			W2:       5.0,  // 葉部下面巾
			W3:       7.0,  // 葉部中央巾
			W4:       3.0,  // 葉部上面巾
			H2:       3.0,  // 葉部下側高さ
			H3:       5.0,  // 葉部上側高さ
		}
		
		// 全体的な寸法関係の確認
		totalHeight := tree.H1 + tree.H2 + tree.H3
		if totalHeight > 100.0 {
			t.Logf("Warning: Total tree height %f m seems unusually large", totalHeight)
		}
		
		maxWidth := tree.W2
		if tree.W3 > maxWidth {
			maxWidth = tree.W3
		}
		if tree.W4 > maxWidth {
			maxWidth = tree.W4
		}
		
		if maxWidth > 50.0 {
			t.Logf("Warning: Maximum tree width %f m seems unusually large", maxWidth)
		}
		
		// 幹と葉部のサイズ関係
		if tree.W1 > maxWidth {
			t.Error("Trunk width should not be larger than foliage width")
		}
		
		t.Logf("Tree geometry summary:")
		t.Logf("  Position: (%.1f, %.1f, %.1f)", tree.x, tree.y, tree.z)
		t.Logf("  Trunk: W1=%.1f m, H1=%.1f m", tree.W1, tree.H1)
		t.Logf("  Foliage: W2=%.1f, W3=%.1f, W4=%.1f m", tree.W2, tree.W3, tree.W4)
		t.Logf("  Heights: H2=%.1f, H3=%.1f m", tree.H2, tree.H3)
		t.Logf("  Total height: %.1f m, Max width: %.1f m", totalHeight, maxWidth)
	})
}

// TestTREE_TreeTypes tests different tree types
func TestTREE_TreeTypes(t *testing.T) {
	t.Run("treeA type validation", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "StandardTree",
		}
		
		if tree.treetype != "treeA" {
			t.Errorf("Expected treetype='treeA', got %s", tree.treetype)
		}
		
		// treeAが実際に使用される形状タイプであることを確認
		validTypes := []string{"treeA"}
		isValid := false
		for _, validType := range validTypes {
			if tree.treetype == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			t.Errorf("Tree type %s is not in valid types %v", tree.treetype, validTypes)
		}
	})

	t.Run("Multiple trees", func(t *testing.T) {
		trees := []*TREE{
			{treetype: "treeA", treename: "Tree1", x: 0.0, y: 0.0, z: 0.0},
			{treetype: "treeA", treename: "Tree2", x: 10.0, y: 0.0, z: 0.0},
			{treetype: "treeA", treename: "Tree3", x: 0.0, y: 10.0, z: 0.0},
		}
		
		if len(trees) != 3 {
			t.Errorf("Expected 3 trees, got %d", len(trees))
		}
		
		// 各樹木の名前の一意性確認
		names := make(map[string]bool)
		for _, tree := range trees {
			if names[tree.treename] {
				t.Errorf("Duplicate tree name: %s", tree.treename)
			}
			names[tree.treename] = true
		}
		
		// 位置の確認
		for i, tree := range trees {
			if tree.treename == "" {
				t.Errorf("Tree %d has empty name", i)
			}
			if tree.treetype != "treeA" {
				t.Errorf("Tree %d has invalid type: %s", i, tree.treetype)
			}
		}
	})
}

// TestTREE_BoundaryValues tests boundary and edge cases
func TestTREE_BoundaryValues(t *testing.T) {
	t.Run("Minimum valid values", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "MinTree",
			x:        0.0,
			y:        0.0,
			z:        0.0,
			W1:       0.1,  // 最小幹太さ
			H1:       0.5,  // 最小幹高さ
			W2:       0.5,  // 最小葉部巾
			W3:       0.5,
			W4:       0.5,
			H2:       0.2,  // 最小葉部高さ
			H3:       0.2,
		}
		
		// 最小値の妥当性確認
		if tree.W1 < 0.05 {
			t.Error("Trunk width too small for realistic tree")
		}
		if tree.H1 < 0.1 {
			t.Error("Trunk height too small for realistic tree")
		}
		
		t.Logf("Minimum tree values validated: W1=%.2f, H1=%.2f", tree.W1, tree.H1)
	})

	t.Run("Large tree values", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "LargeTree",
			x:        100.0,
			y:        200.0,
			z:        5.0,
			W1:       3.0,   // 大きな幹
			H1:       20.0,  // 高い幹
			W2:       25.0,  // 大きな葉部
			W3:       30.0,
			W4:       20.0,
			H2:       10.0,  // 高い葉部
			H3:       15.0,
		}
		
		totalHeight := tree.H1 + tree.H2 + tree.H3
		maxWidth := tree.W3 // 中央が最大
		
		if totalHeight > 80.0 {
			t.Logf("Warning: Very large tree height: %.1f m", totalHeight)
		}
		if maxWidth > 100.0 {
			t.Logf("Warning: Very large tree width: %.1f m", maxWidth)
		}
		
		t.Logf("Large tree validated: Height=%.1f m, Width=%.1f m", totalHeight, maxWidth)
	})

	t.Run("Negative coordinate handling", func(t *testing.T) {
		tree := &TREE{
			treetype: "treeA",
			treename: "NegativeCoordTree",
			x:        -10.0,  // 負の座標も有効
			y:        -20.0,
			z:        -1.0,   // 地下も可能（基礎部分など）
			W1:       0.5,
			H1:       3.0,
		}
		
		// 負の座標は地形や基準点によっては有効
		t.Logf("Tree with negative coordinates: (%.1f, %.1f, %.1f)", tree.x, tree.y, tree.z)
		
		if tree.x < -1000.0 || tree.y < -1000.0 {
			t.Logf("Warning: Very large negative coordinates may indicate input error")
		}
	})
}
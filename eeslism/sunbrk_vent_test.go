package eeslism

import (
	"testing"
)

// TestSNBK tests the SNBK (SUNBRK) structure for shading devices
func TestSNBK(t *testing.T) {
	t.Run("SNBK creation - General overhang", func(t *testing.T) {
		snbk := &SNBK{
			Name: "TestOverhang",
			Type: 1, // 一般の庇(H)
			Ksi:  0, // 反転なし
			W:    2.0, // 開口部の高さ [m]
			H:    1.5, // 開口部の幅 [m]
			D:    0.8, // 庇の付け根から先端までの長さ [m]
			W1:   0.2, // 開口部の左端から壁の左端までの距離 [m]
			W2:   0.2, // 開口部の右端から壁の右端までの距離 [m]
			H1:   0.3, // 開口部の上端から壁の上端までの距離 [m]
			H2:   0.5, // 地面から開口部の下端までの高さ [m]
		}
		
		if snbk.Name != "TestOverhang" {
			t.Errorf("Expected name 'TestOverhang', got %s", snbk.Name)
		}
		if snbk.Type != 1 {
			t.Errorf("Expected Type=1 (general overhang), got %d", snbk.Type)
		}
		if snbk.Ksi != 0 {
			t.Errorf("Expected Ksi=0 (no inversion), got %d", snbk.Ksi)
		}
		if snbk.W != 2.0 {
			t.Errorf("Expected W=2.0, got %f", snbk.W)
		}
		if snbk.H != 1.5 {
			t.Errorf("Expected H=1.5, got %f", snbk.H)
		}
		if snbk.D != 0.8 {
			t.Errorf("Expected D=0.8, got %f", snbk.D)
		}
	})

	t.Run("SNBK creation - Side wall", func(t *testing.T) {
		snbk := &SNBK{
			Name: "TestSideWall",
			Type: 2, // 袖壁(HL)
			Ksi:  1, // 反転あり
			W:    1.8,
			H:    2.1,
			D:    0.5,
			W1:   0.0,
			W2:   0.4,
			H1:   0.2,
			H2:   0.3,
		}
		
		if snbk.Name != "TestSideWall" {
			t.Errorf("Expected name 'TestSideWall', got %s", snbk.Name)
		}
		if snbk.Type != 2 {
			t.Errorf("Expected Type=2 (side wall), got %d", snbk.Type)
		}
		if snbk.Ksi != 1 {
			t.Errorf("Expected Ksi=1 (with inversion), got %d", snbk.Ksi)
		}
	})
}

// TestSunblk tests the sunblk structure for attached shading devices
func TestSunblk(t *testing.T) {
	t.Run("sunblk creation", func(t *testing.T) {
		sb := &sunblk{
			sbfname: "HISASI",
			snbname: "TestHisasi",
			rgb:     [3]float64{0.5, 0.5, 0.5}, // Gray color
			x:       1.0,
			y:       2.0,
			D:       0.8,
			W:       2.0,
			H:       0.3,
			h:       2.5,
			WA:      180.0, // South facing
			ref:     0.3,   // Reflectance
		}
		
		if sb.sbfname != "HISASI" {
			t.Errorf("Expected sbfname 'HISASI', got %s", sb.sbfname)
		}
		if sb.snbname != "TestHisasi" {
			t.Errorf("Expected snbname 'TestHisasi', got %s", sb.snbname)
		}
		if sb.rgb[0] != 0.5 || sb.rgb[1] != 0.5 || sb.rgb[2] != 0.5 {
			t.Errorf("Expected rgb [0.5, 0.5, 0.5], got %v", sb.rgb)
		}
		if sb.ref != 0.3 {
			t.Errorf("Expected ref=0.3, got %f", sb.ref)
		}
		if sb.WA != 180.0 {
			t.Errorf("Expected WA=180.0, got %f", sb.WA)
		}
	})
}

// TestACHIR tests the ACHIR structure for inter-room ventilation
func TestACHIR(t *testing.T) {
	t.Run("ACHIR creation", func(t *testing.T) {
		room := &ROOM{Name: "TestRoom"}
		
		achir := &ACHIR{
			rm:   1,
			sch:  0,
			room: room,
			Gvr:  0.05, // 室間相互換気量 [kg/s]
		}
		
		if achir.rm != 1 {
			t.Errorf("Expected rm=1, got %d", achir.rm)
		}
		if achir.sch != 0 {
			t.Errorf("Expected sch=0, got %d", achir.sch)
		}
		if achir.room != room {
			t.Error("Room reference not set correctly")
		}
		if achir.Gvr != 0.05 {
			t.Errorf("Expected Gvr=0.05, got %f", achir.Gvr)
		}
	})
}

// TestRoomVentilation tests ventilation-related fields in ROOM structure
func TestRoomVentilation(t *testing.T) {
	t.Run("Room ventilation parameters", func(t *testing.T) {
		room := &ROOM{
			Name: "TestVentRoom",
			Gve:  0.1,  // 換気量 [kg/s]
			Gvi:  0.02, // 隙間風量 [kg/s]
			Nachr: 2,   // 室間相互換気の数
		}
		
		// 室間相互換気の設定
		room.achr = []*ACHIR{
			{rm: 2, sch: 1, room: &ROOM{Name: "AdjacentRoom1"}, Gvr: 0.03},
			{rm: 3, sch: 1, room: &ROOM{Name: "AdjacentRoom2"}, Gvr: 0.025},
		}
		
		if room.Name != "TestVentRoom" {
			t.Errorf("Expected name 'TestVentRoom', got %s", room.Name)
		}
		if room.Gve != 0.1 {
			t.Errorf("Expected Gve=0.1, got %f", room.Gve)
		}
		if room.Gvi != 0.02 {
			t.Errorf("Expected Gvi=0.02, got %f", room.Gvi)
		}
		if room.Nachr != 2 {
			t.Errorf("Expected Nachr=2, got %d", room.Nachr)
		}
		if len(room.achr) != 2 {
			t.Errorf("Expected 2 ACHIR entries, got %d", len(room.achr))
		}
		
		// 室間相互換気の検証
		if room.achr[0].Gvr != 0.03 {
			t.Errorf("Expected first ACHIR Gvr=0.03, got %f", room.achr[0].Gvr)
		}
		if room.achr[1].Gvr != 0.025 {
			t.Errorf("Expected second ACHIR Gvr=0.025, got %f", room.achr[1].Gvr)
		}
	})
}

// TestTRNX tests the TRNX structure for adjacent rooms
func TestTRNX(t *testing.T) {
	t.Run("TRNX creation", func(t *testing.T) {
		nextroom := &ROOM{Name: "AdjacentRoom"}
		sd := &RMSRF{Name: "SharedSurface"}
		
		trnx := &TRNX{
			nextroom: nextroom,
			sd:       sd,
		}
		
		if trnx.nextroom != nextroom {
			t.Error("Next room reference not set correctly")
		}
		if trnx.sd != sd {
			t.Error("Surface reference not set correctly")
		}
		if trnx.nextroom.Name != "AdjacentRoom" {
			t.Errorf("Expected adjacent room name 'AdjacentRoom', got %s", trnx.nextroom.Name)
		}
		if trnx.sd.Name != "SharedSurface" {
			t.Errorf("Expected surface name 'SharedSurface', got %s", trnx.sd.Name)
		}
	})
}
package eeslism

import (
	"testing"
)

// TestEXSF tests the EXSF (External Surface) structure
func TestEXSF(t *testing.T) {
	t.Run("EXSF creation - Basic external surface", func(t *testing.T) {
		aloValue := 25.0
		exsf := &EXSF{
			Name:   "South",
			Typ:    EXSFType_S,    // 一般外表面
			Wa:     0.0,           // 方位角 [度] (南向き)
			Wb:     90.0,          // 傾斜角 [度] (垂直)
			Rg:     0.2,           // 前面地物の日射反射率 [-]
			Alo:    &aloValue,     // 外表面総合熱伝達率 [W/m2K]
			Z:      0.0,           // 地中深さ [m] (地上)
			Erdff:  0.36e-6,       // 土の熱拡散率 [m2/s]
		}
		
		if exsf.Name != "South" {
			t.Errorf("Expected name 'South', got %s", exsf.Name)
		}
		if exsf.Typ != EXSFType_S {
			t.Errorf("Expected Typ=EXSFType_S, got %c", exsf.Typ)
		}
		if exsf.Wa != 0.0 {
			t.Errorf("Expected Wa=0.0, got %f", exsf.Wa)
		}
		if exsf.Wb != 90.0 {
			t.Errorf("Expected Wb=90.0, got %f", exsf.Wb)
		}
		if exsf.Rg != 0.2 {
			t.Errorf("Expected Rg=0.2, got %f", exsf.Rg)
		}
		if exsf.Alo == nil || *exsf.Alo != 25.0 {
			t.Errorf("Expected Alo=25.0, got %v", exsf.Alo)
		}
	})

	t.Run("EXSF creation - Horizontal surface", func(t *testing.T) {
		aloValue := 23.3
		exsf := &EXSF{
			Name:   "Hor",
			Typ:    EXSFType_S,    // 一般外表面
			Wa:     0.0,           // 方位角は水平面では意味なし
			Wb:     0.0,           // 傾斜角 [度] (水平)
			Rg:     0.15,          // 前面地物の日射反射率 [-]
			Alo:    &aloValue,     // 外表面総合熱伝達率 [W/m2K]
			Z:      0.0,
			Erdff:  0.36e-6,
		}
		
		if exsf.Name != "Hor" {
			t.Errorf("Expected name 'Hor', got %s", exsf.Name)
		}
		if exsf.Typ != EXSFType_S {
			t.Errorf("Expected Typ=EXSFType_S, got %c", exsf.Typ)
		}
		if exsf.Wb != 0.0 {
			t.Errorf("Expected Wb=0.0 (horizontal), got %f", exsf.Wb)
		}
		if exsf.Alo == nil || *exsf.Alo != 23.3 {
			t.Errorf("Expected Alo=23.3, got %v", exsf.Alo)
		}
	})

	t.Run("EXSF orientation angles", func(t *testing.T) {
		// 各方位のテスト
		orientations := []struct {
			name string
			wa   float64
			desc string
		}{
			{"South", 0.0, "南向き"},
			{"West", 90.0, "西向き"},
			{"North", 180.0, "北向き"},
			{"East", 270.0, "東向き"},
		}
		
		for _, orient := range orientations {
			aloValue := 25.0
			exsf := &EXSF{
				Name: orient.name,
				Typ:  EXSFType_S,
				Wa:   orient.wa,
				Wb:   90.0, // 垂直面
				Rg:   0.2,
				Alo:  &aloValue,
			}
			
			if exsf.Name != orient.name {
				t.Errorf("Expected name '%s', got %s", orient.name, exsf.Name)
			}
			if exsf.Wa != orient.wa {
				t.Errorf("Expected Wa=%f for %s (%s), got %f", orient.wa, orient.name, orient.desc, exsf.Wa)
			}
			
			// 方位角の範囲チェック
			if exsf.Wa < 0.0 || exsf.Wa >= 360.0 {
				t.Errorf("Azimuth angle out of range [0, 360): %f", exsf.Wa)
			}
		}
	})
}

// TestEXSFS tests the EXSFS (External Surfaces) collection structure
func TestEXSFS(t *testing.T) {
	t.Run("EXSFS creation and management", func(t *testing.T) {
		aloValue1 := 25.0
		aloValue2 := 23.3
		exsfs := &EXSFS{
			Exs: []*EXSF{
				{Name: "South", Typ: EXSFType_S, Wa: 0.0, Wb: 90.0, Rg: 0.2, Alo: &aloValue1},
				{Name: "West", Typ: EXSFType_S, Wa: 90.0, Wb: 90.0, Rg: 0.2, Alo: &aloValue1},
				{Name: "North", Typ: EXSFType_S, Wa: 180.0, Wb: 90.0, Rg: 0.2, Alo: &aloValue1},
				{Name: "Hor", Typ: EXSFType_S, Wa: 0.0, Wb: 0.0, Rg: 0.15, Alo: &aloValue2},
			},
		}
		
		if len(exsfs.Exs) != 4 {
			t.Errorf("Expected 4 external surfaces, got %d", len(exsfs.Exs))
		}
		
		// 各外表面の検証
		expectedNames := []string{"South", "West", "North", "Hor"}
		for i, expectedName := range expectedNames {
			if exsfs.Exs[i].Name != expectedName {
				t.Errorf("Expected surface %d name '%s', got %s", i, expectedName, exsfs.Exs[i].Name)
			}
		}
		
		// 水平面の検証
		horSurface := exsfs.Exs[3]
		if horSurface.Name != "Hor" {
			t.Errorf("Expected horizontal surface name 'Hor', got %s", horSurface.Name)
		}
		if horSurface.Wb != 0.0 {
			t.Errorf("Expected horizontal surface Wb=0.0, got %f", horSurface.Wb)
		}
	})
}

// TestPCM tests the PCM (Phase Change Material) structure
func TestPCM(t *testing.T) {
	t.Run("PCM creation - Basic paraffin wax", func(t *testing.T) {
		pcm := &PCM{
			Name:     "ParaffinWax28",
			Spctype:  'm', // モデルで設定
			Condtype: 'm', // モデルで設定
			Ql:       180000.0, // 潜熱量 [J/m3]
			Condl:    0.15,     // 液相の熱伝導率 [W/mK]
			Conds:    0.20,     // 固相の熱伝導率 [W/mK]
			Crol:     800000.0, // 液相の容積比熱 [J/m3K]
			Cros:     900000.0, // 固相の容積比熱 [J/m3K]
			Ts:       26.0,     // 固体から融解が始まる温度 [℃]
			Tl:       30.0,     // 液体から凝固が始まる温度 [℃]
			Tp:       28.0,     // 見かけの比熱のピーク温度 [℃]
			Iterate:  false,    // 収束計算なし
		}
		
		if pcm.Name != "ParaffinWax28" {
			t.Errorf("Expected name 'ParaffinWax28', got %s", pcm.Name)
		}
		if pcm.Spctype != 'm' {
			t.Errorf("Expected Spctype='m', got %c", pcm.Spctype)
		}
		if pcm.Condtype != 'm' {
			t.Errorf("Expected Condtype='m', got %c", pcm.Condtype)
		}
		if pcm.Ql != 180000.0 {
			t.Errorf("Expected Ql=180000.0, got %f", pcm.Ql)
		}
		if pcm.Ts != 26.0 {
			t.Errorf("Expected Ts=26.0, got %f", pcm.Ts)
		}
		if pcm.Tl != 30.0 {
			t.Errorf("Expected Tl=30.0, got %f", pcm.Tl)
		}
		if pcm.Tp != 28.0 {
			t.Errorf("Expected Tp=28.0, got %f", pcm.Tp)
		}
		
		// 温度関係の妥当性チェック
		if pcm.Ts > pcm.Tp || pcm.Tp > pcm.Tl {
			t.Errorf("Temperature relationship should be Ts <= Tp <= Tl, got Ts=%f, Tp=%f, Tl=%f", 
				pcm.Ts, pcm.Tp, pcm.Tl)
		}
	})

	t.Run("PCM with table data", func(t *testing.T) {
		// CHARTABLEの設定例
		chartable := CHARTABLE{
			filename:    "pcm_data.csv",
			PCMchar:     'E', // エンタルピー
			tabletype:   'h', // 見かけの比熱
			minTemp:     20.0,
			maxTemp:     35.0,
			itablerow:   15,
			minTempChng: 0.1,
		}
		
		pcm := &PCM{
			Name:     "TablePCM",
			Spctype:  't', // テーブル形式
			Condtype: 'm', // モデルで設定
			Chartable: [2]CHARTABLE{chartable, {}}, // エンタルピーテーブルのみ
		}
		
		if pcm.Spctype != 't' {
			t.Errorf("Expected Spctype='t', got %c", pcm.Spctype)
		}
		if pcm.Chartable[0].PCMchar != 'E' {
			t.Errorf("Expected PCMchar='E', got %c", pcm.Chartable[0].PCMchar)
		}
		if pcm.Chartable[0].tabletype != 'h' {
			t.Errorf("Expected tabletype='h', got %c", pcm.Chartable[0].tabletype)
		}
	})
}

// TestPCMSTATE tests the PCMSTATE structure for PCM state values
func TestPCMSTATE(t *testing.T) {
	t.Run("PCMSTATE creation", func(t *testing.T) {
		pcmName := "TestPCM"
		pcmstate := &PCMSTATE{
			Name:         &pcmName,
			TempPCMNodeL: 25.5, // PCM温度（左側節点）[℃]
			TempPCMNodeR: 26.2, // PCM温度（右側節点）[℃]
			TempPCMave:   25.85, // PCM温度（平均温度）[℃]
			CapmL:        2500.0, // PCM見かけの比熱（左側）[J/kgK]
			CapmR:        2600.0, // PCM見かけの比熱（右側）[J/kgK]
			LamdaL:       0.16,   // PCM熱伝導率（左側）[W/mK]
			LamdaR:       0.17,   // PCM熱伝導率（右側）[W/mK]
		}
		
		if pcmstate.Name == nil || *pcmstate.Name != "TestPCM" {
			t.Errorf("Expected name 'TestPCM', got %v", pcmstate.Name)
		}
		if pcmstate.TempPCMNodeL != 25.5 {
			t.Errorf("Expected TempPCMNodeL=25.5, got %f", pcmstate.TempPCMNodeL)
		}
		if pcmstate.TempPCMNodeR != 26.2 {
			t.Errorf("Expected TempPCMNodeR=26.2, got %f", pcmstate.TempPCMNodeR)
		}
		
		// 平均温度の妥当性チェック
		expectedAvg := (pcmstate.TempPCMNodeL + pcmstate.TempPCMNodeR) / 2.0
		if pcmstate.TempPCMave != expectedAvg {
			t.Errorf("Average temperature inconsistent: expected %f, got %f", 
				expectedAvg, pcmstate.TempPCMave)
		}
	})
}
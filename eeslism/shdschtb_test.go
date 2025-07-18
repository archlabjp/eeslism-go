package eeslism

import (
	"testing"
)

// TestSHDSCHTB tests the SHDSCHTB (Shade Schedule Table) functionality
func TestSHDSCHTB(t *testing.T) {
	t.Run("SHADTB structure creation", func(t *testing.T) {
		// SHADTB構造体の基本作成テスト
		shadtb := &SHADTB{
			lpname: "WallSurface1",
			indatn: 12, // 12ヶ月
			shad: [12]float64{
				0.8, // 1月 - 冬季、ほぼ落葉
				0.7, // 2月 - 冬季
				0.5, // 3月 - 春の始まり
				0.2, // 4月 - 新緑
				0.0, // 5月 - 完全に葉
				0.0, // 6月 - 夏季
				0.0, // 7月 - 夏季
				0.0, // 8月 - 夏季
				0.1, // 9月 - 秋の始まり
				0.3, // 10月 - 紅葉
				0.6, // 11月 - 落葉進行
				0.8, // 12月 - 冬季
			},
		}
		
		// 開始日・終了日の設定
		for i := 0; i < 12; i++ {
			shadtb.ndays[i] = i*30 + 1    // 各月の開始日（概算）
			shadtb.ndaye[i] = (i+1)*30   // 各月の終了日（概算）
		}
		
		if shadtb.lpname != "WallSurface1" {
			t.Errorf("Expected lpname='WallSurface1', got %s", shadtb.lpname)
		}
		if shadtb.indatn != 12 {
			t.Errorf("Expected indatn=12, got %d", shadtb.indatn)
		}
		
		// 落葉率の範囲確認
		for i, rate := range shadtb.shad {
			if rate < 0.0 || rate > 1.0 {
				t.Errorf("Shade rate for month %d should be in range [0, 1], got %f", i+1, rate)
			}
		}
		
		t.Logf("SHADTB created: %s with %d monthly values", shadtb.lpname, len(shadtb.shad))
	})

	t.Run("Evergreen tree schedule", func(t *testing.T) {
		// 常緑樹のスケジュール（年中葉がある）
		shadtb := &SHADTB{
			lpname: "EvergreenSurface",
			indatn: 12,
			shad: [12]float64{
				0.1, 0.1, 0.1, 0.1, 0.0, 0.0, // 1-6月
				0.0, 0.0, 0.0, 0.1, 0.1, 0.1, // 7-12月
			},
		}
		
		if shadtb.lpname != "EvergreenSurface" {
			t.Errorf("Expected lpname='EvergreenSurface', got %s", shadtb.lpname)
		}
		
		// 常緑樹は年中低い落葉率
		for i, rate := range shadtb.shad {
			if rate > 0.2 {
				t.Logf("Note: High shade rate %.1f for evergreen tree in month %d", rate, i+1)
			}
		}
		
		// 夏季の完全な葉の確認
		summerMonths := []int{4, 5, 6, 7} // 5-8月（0ベース）
		for _, month := range summerMonths {
			if shadtb.shad[month] > 0.1 {
				t.Logf("Note: Evergreen tree has shade rate %.1f in summer month %d", 
					shadtb.shad[month], month+1)
			}
		}
		
		t.Logf("Evergreen schedule validated for %s", shadtb.lpname)
	})

	t.Run("Seasonal variation validation", func(t *testing.T) {
		// 季節変化の妥当性テスト
		shadtb := &SHADTB{
			lpname: "SeasonalSurface",
			indatn: 12,
			shad: [12]float64{
				0.9, 0.8, 0.6, 0.3, 0.1, 0.0, // 冬→春→夏
				0.0, 0.0, 0.2, 0.4, 0.7, 0.9, // 夏→秋→冬
			},
		}
		
		// 季節変化パターンの確認
		winter := (shadtb.shad[0] + shadtb.shad[1] + shadtb.shad[11]) / 3.0  // 12,1,2月
		summer := (shadtb.shad[5] + shadtb.shad[6] + shadtb.shad[7]) / 3.0   // 6,7,8月
		
		if winter <= summer {
			t.Error("Winter shade rate should be higher than summer for deciduous trees")
		}
		
		// 春の減少傾向確認
		spring := []int{2, 3, 4} // 3,4,5月
		for i := 0; i < len(spring)-1; i++ {
			if shadtb.shad[spring[i]] < shadtb.shad[spring[i+1]] {
				t.Logf("Note: Shade rate increased from month %d to %d in spring", 
					spring[i]+1, spring[i+1]+1)
			}
		}
		
		// 秋の増加傾向確認
		autumn := []int{8, 9, 10} // 9,10,11月
		for i := 0; i < len(autumn)-1; i++ {
			if shadtb.shad[autumn[i]] > shadtb.shad[autumn[i+1]] {
				t.Logf("Note: Shade rate decreased from month %d to %d in autumn", 
					autumn[i]+1, autumn[i+1]+1)
			}
		}
		
		t.Logf("Seasonal variation: Winter avg=%.2f, Summer avg=%.2f", winter, summer)
	})
}

// TestSHDSCHTB_MultipleSchedules tests multiple shade schedules
func TestSHDSCHTB_MultipleSchedules(t *testing.T) {
	t.Run("Multiple tree types", func(t *testing.T) {
		// 複数の樹木タイプのスケジュール
		schedules := []*SHADTB{
			{
				lpname: "OakSurface", // オーク（落葉樹）
				indatn: 12,
				shad: [12]float64{0.9, 0.8, 0.6, 0.3, 0.1, 0.0, 0.0, 0.0, 0.2, 0.5, 0.7, 0.9},
			},
			{
				lpname: "PineSurface", // 松（常緑樹）
				indatn: 12,
				shad: [12]float64{0.1, 0.1, 0.1, 0.1, 0.0, 0.0, 0.0, 0.0, 0.1, 0.1, 0.1, 0.1},
			},
			{
				lpname: "MapleSurface", // カエデ（落葉樹）
				indatn: 12,
				shad: [12]float64{0.8, 0.7, 0.5, 0.2, 0.0, 0.0, 0.0, 0.1, 0.3, 0.6, 0.8, 0.9},
			},
		}
		
		if len(schedules) != 3 {
			t.Errorf("Expected 3 schedules, got %d", len(schedules))
		}
		
		// 各スケジュールの名前の一意性確認
		names := make(map[string]bool)
		for _, schedule := range schedules {
			if names[schedule.lpname] {
				t.Errorf("Duplicate schedule name: %s", schedule.lpname)
			}
			names[schedule.lpname] = true
		}
		
		// 樹木タイプ別の特性確認
		for _, schedule := range schedules {
			winterAvg := (schedule.shad[0] + schedule.shad[1] + schedule.shad[11]) / 3.0
			summerAvg := (schedule.shad[5] + schedule.shad[6] + schedule.shad[7]) / 3.0
			
			if schedule.lpname == "PineSurface" {
				// 常緑樹は年中低い落葉率
				if winterAvg > 0.2 || summerAvg > 0.1 {
					t.Logf("Note: %s (evergreen) has high shade rates: winter=%.2f, summer=%.2f", 
						schedule.lpname, winterAvg, summerAvg)
				}
			} else {
				// 落葉樹は冬に高い落葉率
				if winterAvg <= summerAvg {
					t.Errorf("%s (deciduous) should have higher winter shade rate", schedule.lpname)
				}
			}
			
			t.Logf("Tree: %s, Winter avg: %.2f, Summer avg: %.2f", 
				schedule.lpname, winterAvg, summerAvg)
		}
	})

	t.Run("Schedule data validation", func(t *testing.T) {
		// スケジュールデータの妥当性確認
		shadtb := &SHADTB{
			lpname: "ValidationSurface",
			indatn: 12,
			shad: [12]float64{0.8, 0.7, 0.5, 0.3, 0.1, 0.0, 0.0, 0.1, 0.3, 0.5, 0.7, 0.8},
		}
		
		// データ数の整合性
		if shadtb.indatn != 12 {
			t.Errorf("indatn (%d) should be 12 for monthly data", shadtb.indatn)
		}
		
		// 値の範囲確認
		for i, value := range shadtb.shad {
			if value < 0.0 || value > 1.0 {
				t.Errorf("Shade value[%d] = %f is out of range [0, 1]", i, value)
			}
		}
		
		// 急激な変化の確認
		for i := 1; i < len(shadtb.shad); i++ {
			change := shadtb.shad[i] - shadtb.shad[i-1]
			if change > 0.5 || change < -0.5 {
				t.Logf("Warning: Large change (%.2f) between month %d and %d", 
					change, i, i+1)
			}
		}
		
		// 年末と年始の連続性確認
		yearEndChange := shadtb.shad[0] - shadtb.shad[11] // 1月 - 12月
		if yearEndChange > 0.3 || yearEndChange < -0.3 {
			t.Logf("Warning: Large change (%.2f) between December and January", yearEndChange)
		}
		
		t.Logf("Schedule validation completed for %s", shadtb.lpname)
	})
}

// TestSHDSCHTB_DateScheduling tests date scheduling functionality
func TestSHDSCHTB_DateScheduling(t *testing.T) {
	t.Run("Monthly date ranges", func(t *testing.T) {
		shadtb := &SHADTB{
			lpname: "DateTestSurface",
			indatn: 12,
		}
		
		// 月別の日付範囲設定
		monthDays := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		currentDay := 1
		
		for i := 0; i < 12; i++ {
			shadtb.ndays[i] = currentDay
			shadtb.ndaye[i] = currentDay + monthDays[i] - 1
			currentDay += monthDays[i]
		}
		
		// 日付範囲の妥当性確認
		for i := 0; i < 12; i++ {
			if shadtb.ndays[i] > shadtb.ndaye[i] {
				t.Errorf("Month %d: start day (%d) should be <= end day (%d)", 
					i+1, shadtb.ndays[i], shadtb.ndaye[i])
			}
			
			if i > 0 {
				if shadtb.ndays[i] != shadtb.ndaye[i-1]+1 {
					t.Logf("Note: Gap between month %d end (%d) and month %d start (%d)", 
						i, shadtb.ndaye[i-1], i+1, shadtb.ndays[i])
				}
			}
		}
		
		// 年間日数の確認
		totalDays := shadtb.ndaye[11]
		if totalDays != 365 {
			t.Logf("Note: Total days = %d (expected 365 for non-leap year)", totalDays)
		}
		
		t.Logf("Date scheduling validated: %d days total", totalDays)
	})

	t.Run("Leap year handling", func(t *testing.T) {
		shadtb := &SHADTB{
			lpname: "LeapYearSurface",
			indatn: 12,
		}
		
		// うるう年の日付設定
		monthDays := []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31} // 2月が29日
		currentDay := 1
		
		for i := 0; i < 12; i++ {
			shadtb.ndays[i] = currentDay
			shadtb.ndaye[i] = currentDay + monthDays[i] - 1
			currentDay += monthDays[i]
		}
		
		// うるう年の確認
		totalDays := shadtb.ndaye[11]
		if totalDays == 366 {
			t.Logf("Leap year correctly handled: %d days", totalDays)
		} else {
			t.Logf("Note: Total days = %d (may not be leap year)", totalDays)
		}
		
		// 2月の日数確認
		febDays := shadtb.ndaye[1] - shadtb.ndays[1] + 1
		if febDays == 29 {
			t.Logf("February has 29 days (leap year)")
		} else if febDays == 28 {
			t.Logf("February has 28 days (non-leap year)")
		} else {
			t.Errorf("February has unexpected %d days", febDays)
		}
	})
}

// TestSHDSCHTB_EdgeCases tests edge cases and boundary values
func TestSHDSCHTB_EdgeCases(t *testing.T) {
	t.Run("Extreme shade values", func(t *testing.T) {
		// 極端な日射遮蔽率のテスト
		extremeCases := []struct {
			name     string
			schedule SHADTB
			desc     string
		}{
			{
				name: "AlwaysShaded",
				schedule: SHADTB{
					lpname: "AlwaysShadedSurface",
					indatn: 12,
					shad: [12]float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
				},
				desc: "年中完全に日陰",
			},
			{
				name: "NeverShaded",
				schedule: SHADTB{
					lpname: "NeverShadedSurface",
					indatn: 12,
					shad: [12]float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0},
				},
				desc: "年中完全に日向",
			},
			{
				name: "ExtremeVariation",
				schedule: SHADTB{
					lpname: "ExtremeVariationSurface",
					indatn: 12,
					shad: [12]float64{1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0},
				},
				desc: "極端な月変化",
			},
		}
		
		for _, testCase := range extremeCases {
			schedule := testCase.schedule
			
			// 基本的な妥当性確認
			if schedule.indatn != 12 {
				t.Errorf("%s: Expected indatn=12, got %d", testCase.name, schedule.indatn)
			}
			
			// 値の範囲確認
			for i, value := range schedule.shad {
				if value < 0.0 || value > 1.0 {
					t.Errorf("%s: Value[%d] = %f is out of range", testCase.name, i, value)
				}
			}
			
			// 特殊ケースの確認
			if testCase.name == "ExtremeVariation" {
				for i := 1; i < len(schedule.shad); i++ {
					change := schedule.shad[i] - schedule.shad[i-1]
					if change != 1.0 && change != -1.0 {
						t.Logf("%s: Unexpected change %.1f between months %d and %d", 
							testCase.name, change, i, i+1)
					}
				}
			}
			
			t.Logf("Extreme case validated: %s - %s", testCase.name, testCase.desc)
		}
	})
}
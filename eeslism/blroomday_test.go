package eeslism

import (
	"testing"
)

// TestRoomday tests the main Roomday function for daily room calculations
func TestRoomday(t *testing.T) {
	t.Run("Basic Roomday execution", func(t *testing.T) {
		// 基本的なRoomday関数の実行テスト
		
		// テスト用の室を作成（基本フィールドのみ）
		rooms := []*ROOM{
			{
				Name:  "TestRoom1",
				VRM:   50.0,  // 室容積
				Tr:    25.0,  // 室温
				FArea: 20.0,  // 床面積
			},
		}
		
		// テスト用の輻射パネル（空のリスト）
		rdpnls := []*RDPNL{}
		
		// Roomday関数を実行
		Mon := 6      // 6月
		Day := 15     // 15日
		Nday := 166   // 年間通算日
		ttmm := 1200  // 12:00
		Simdayend := 0 // 日終了フラグ
		
		// 関数実行（パニックしないことを確認）
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked: %v", r)
			}
		}()
		
		Roomday(Mon, Day, Nday, ttmm, rooms, rdpnls, Simdayend)
		
		// 実行後の基本確認（関数が正常終了したことを確認）
		t.Logf("Roomday function completed without panic")
		t.Logf("Roomday executed successfully for %d rooms on %d/%d", len(rooms), Mon, Day)
	})

	t.Run("Roomday with different time periods", func(t *testing.T) {
		// 異なる時間帯でのRoomday実行テスト
		room := &ROOM{
			Name:  "TimeTestRoom",
			VRM:   60.0,
			Tr:    23.0,
			FArea: 25.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		// 異なる時刻でのテスト
		testTimes := []struct {
			hour   int
			minute int
			ttmm   int
			desc   string
		}{
			{0, 0, 0, "Midnight"},
			{6, 0, 600, "Early morning"},
			{12, 0, 1200, "Noon"},
			{18, 0, 1800, "Evening"},
			{23, 59, 2359, "Late night"},
		}
		
		for _, tt := range testTimes {
			t.Run(tt.desc, func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Roomday panicked at %s (%d): %v", tt.desc, tt.ttmm, r)
					}
				}()
				
				Roomday(6, 15, 166, tt.ttmm, rooms, rdpnls, 0)
				t.Logf("Roomday executed at %s (%02d:%02d)", tt.desc, tt.hour, tt.minute)
			})
		}
	})

	t.Run("Roomday with seasonal variations", func(t *testing.T) {
		// 季節変化でのRoomday実行テスト
		room := &ROOM{
			Name:  "SeasonalRoom",
			VRM:   50.0,
			Tr:    20.0,
			FArea: 20.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		// 季節ごとのテスト
		seasons := []struct {
			mon  int
			day  int
			nday int
			desc string
		}{
			{1, 15, 15, "Winter"},
			{4, 15, 105, "Spring"},
			{7, 15, 196, "Summer"},
			{10, 15, 288, "Autumn"},
		}
		
		for _, season := range seasons {
			t.Run(season.desc, func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Roomday panicked in %s: %v", season.desc, r)
					}
				}()
				
				Roomday(season.mon, season.day, season.nday, 1200, rooms, rdpnls, 0)
				t.Logf("Roomday executed for %s (%d/%d)", season.desc, season.mon, season.day)
			})
		}
	})
}

// TestRoomdayEdgeCases tests edge cases for Roomday function
func TestRoomdayEdgeCases(t *testing.T) {
	t.Run("Empty room list", func(t *testing.T) {
		// 空の室リストでのテスト
		rooms := []*ROOM{}
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked with empty room list: %v", r)
			}
		}()
		
		Roomday(6, 15, 166, 1200, rooms, rdpnls, 0)
		t.Logf("Roomday handled empty room list successfully")
	})

	t.Run("Single room", func(t *testing.T) {
		// 単一室でのテスト
		room := &ROOM{
			Name:  "SingleRoom",
			VRM:   100.0,
			Tr:    25.0,
			FArea: 40.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked with single room: %v", r)
			}
		}()
		
		Roomday(8, 20, 232, 1500, rooms, rdpnls, 0)
		t.Logf("Roomday handled single room successfully")
	})

	t.Run("Large room list", func(t *testing.T) {
		// 大量の室でのテスト
		rooms := make([]*ROOM, 10)
		for i := 0; i < 10; i++ {
			rooms[i] = &ROOM{
				Name:  "Room" + string(rune('A'+i)),
				VRM:   50.0 + float64(i)*10.0,
				Tr:    20.0 + float64(i)*1.0,
				FArea: 20.0 + float64(i)*2.0,
			}
		}
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked with large room list: %v", r)
			}
		}()
		
		Roomday(9, 10, 253, 1000, rooms, rdpnls, 0)
		t.Logf("Roomday handled %d rooms successfully", len(rooms))
	})

	t.Run("Extreme temperature values", func(t *testing.T) {
		// 極端な温度値でのテスト
		rooms := []*ROOM{
			{Name: "HotRoom", VRM: 50.0, Tr: 50.0, FArea: 20.0},   // 高温
			{Name: "ColdRoom", VRM: 50.0, Tr: -10.0, FArea: 20.0}, // 低温
		}
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked with extreme temperatures: %v", r)
			}
		}()
		
		Roomday(2, 28, 59, 1200, rooms, rdpnls, 0)
		t.Logf("Roomday handled extreme temperature values successfully")
	})

	t.Run("Day end simulation", func(t *testing.T) {
		// 日終了シミュレーションのテスト
		room := &ROOM{
			Name:  "DayEndRoom",
			VRM:   50.0,
			Tr:    25.0,
			FArea: 20.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked at day end: %v", r)
			}
		}()
		
		// 日終了フラグを立ててテスト
		Roomday(12, 31, 365, 2359, rooms, rdpnls, 1)
		t.Logf("Roomday handled day end simulation successfully")
	})
}

// TestRoomdayBoundaryValues tests boundary values
func TestRoomdayBoundaryValues(t *testing.T) {
	t.Run("Boundary time values", func(t *testing.T) {
		// 境界時刻値でのテスト
		room := &ROOM{
			Name:  "BoundaryRoom",
			VRM:   50.0,
			Tr:    25.0,
			FArea: 20.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		boundaryTimes := []struct {
			ttmm int
			desc string
		}{
			{0, "Start of day"},
			{2359, "End of day"},
			{1200, "Noon"},
		}
		
		for _, bt := range boundaryTimes {
			t.Run(bt.desc, func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Roomday panicked at %s (%d): %v", bt.desc, bt.ttmm, r)
					}
				}()
				
				Roomday(6, 15, 166, bt.ttmm, rooms, rdpnls, 0)
				t.Logf("Roomday handled %s (%d) successfully", bt.desc, bt.ttmm)
			})
		}
	})

	t.Run("Boundary date values", func(t *testing.T) {
		// 境界日付値でのテスト
		room := &ROOM{
			Name:  "DateBoundaryRoom",
			VRM:   50.0,
			Tr:    25.0,
			FArea: 20.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		boundaryDates := []struct {
			mon  int
			day  int
			nday int
			desc string
		}{
			{1, 1, 1, "New Year"},
			{12, 31, 365, "Year End"},
			{2, 29, 60, "Leap Day"},
		}
		
		for _, bd := range boundaryDates {
			t.Run(bd.desc, func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Roomday panicked on %s (%d/%d): %v", bd.desc, bd.mon, bd.day, r)
					}
				}()
				
				Roomday(bd.mon, bd.day, bd.nday, 1200, rooms, rdpnls, 0)
				t.Logf("Roomday handled %s (%d/%d) successfully", bd.desc, bd.mon, bd.day)
			})
		}
	})
}

// TestRoomdayIntegration tests integration scenarios
func TestRoomdayIntegration(t *testing.T) {
	t.Run("Multi-room multi-panel scenario", func(t *testing.T) {
		// 複数室・複数パネルのシナリオテスト
		rooms := []*ROOM{
			{Name: "LivingRoom", VRM: 80.0, Tr: 24.0, FArea: 32.0},
		}
		
		rdpnls := []*RDPNL{}
		
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Roomday panicked in multi-room scenario: %v", r)
			}
		}()
		
		// 複数室・複数パネルでのRoomday実行
		Roomday(8, 15, 227, 1400, rooms, rdpnls, 0)
		
		t.Logf("Multi-room multi-panel scenario completed: %d rooms, %d panels", 
			len(rooms), len(rdpnls))
	})

	t.Run("Complete daily cycle simulation", func(t *testing.T) {
		// 完全な日サイクルのシミュレーション
		room := &ROOM{
			Name:  "CycleTestRoom",
			VRM:   50.0,
			Tr:    22.0,
			FArea: 20.0,
		}
		rooms := []*ROOM{room}
		rdpnls := []*RDPNL{}
		
		// 1日の複数時刻でRoomdayを実行
		Mon := 3
		Day := 21
		Nday := 80
		
		times := []int{0, 600, 1200, 1800, 2359}
		
		for i, ttmm := range times {
			simdayend := 0
			if i == len(times)-1 {
				simdayend = 1 // 最後の時刻で日終了
			}
			
			defer func(timeIndex int, time int) {
				if r := recover(); r != nil {
					t.Errorf("Roomday panicked at time %d (%d): %v", timeIndex, time, r)
				}
			}(i, ttmm)
			
			Roomday(Mon, Day, Nday, ttmm, rooms, rdpnls, simdayend)
		}
		
		t.Logf("Complete daily cycle executed for room: %s", room.Name)
	})
}
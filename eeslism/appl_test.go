package eeslism

import (
	"testing"
)

// TestRoomAppliances tests appliance-related fields in ROOM structure
func TestAppl(t *testing.T) {
	t.Run("Room lighting parameters", func(t *testing.T) {
		// 照明スケジュール値のダミー
		lightschValue := 0.8
		
		room := &ROOM{
			Name:     "TestLightRoom",
			Ltyp:     'x',    // 照明器具形式
			Light:    500.0,  // 照明器具容量 [W]
			Lightsch: &lightschValue, // 照明スケジュール
		}
		
		if room.Name != "TestLightRoom" {
			t.Errorf("Expected name 'TestLightRoom', got %s", room.Name)
		}
		if room.Ltyp != 'x' {
			t.Errorf("Expected Ltyp='x', got %c", room.Ltyp)
		}
		if room.Light != 500.0 {
			t.Errorf("Expected Light=500.0, got %f", room.Light)
		}
		if room.Lightsch == nil {
			t.Error("Lighting schedule (Lightsch) should not be nil")
		}
		if *room.Lightsch != 0.8 {
			t.Errorf("Expected Lightsch=0.8, got %f", *room.Lightsch)
		}
	})

	t.Run("Room appliance parameters", func(t *testing.T) {
		// スケジュール値のダミー
		asschValue := 0.7
		alschValue := 0.6
		
		room := &ROOM{
			Name:  "TestApplianceRoom",
			Apsc:  300.0, // 機器対流放熱容量 [W]
			Apsr:  200.0, // 機器輻射放熱容量 [W]
			Apl:   100.0, // 機器潜熱放熱容量 [W]
			Assch: &asschValue, // 機器顕熱スケジュール
			Alsch: &alschValue, // 機器潜熱スケジュール
		}
		
		if room.Name != "TestApplianceRoom" {
			t.Errorf("Expected name 'TestApplianceRoom', got %s", room.Name)
		}
		if room.Apsc != 300.0 {
			t.Errorf("Expected Apsc=300.0, got %f", room.Apsc)
		}
		if room.Apsr != 200.0 {
			t.Errorf("Expected Apsr=200.0, got %f", room.Apsr)
		}
		if room.Apl != 100.0 {
			t.Errorf("Expected Apl=100.0, got %f", room.Apl)
		}
		if room.Assch == nil {
			t.Error("Appliance sensible heat schedule (Assch) should not be nil")
		}
		if room.Alsch == nil {
			t.Error("Appliance latent heat schedule (Alsch) should not be nil")
		}
		if *room.Assch != 0.7 {
			t.Errorf("Expected Assch=0.7, got %f", *room.Assch)
		}
		if *room.Alsch != 0.6 {
			t.Errorf("Expected Alsch=0.6, got %f", *room.Alsch)
		}
	})
}

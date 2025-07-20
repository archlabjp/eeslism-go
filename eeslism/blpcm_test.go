package eeslism

import (
	"testing"
)

func TestPCMdata(t *testing.T) {
	t.Run("valid PCM data", func(t *testing.T) {
		input := `ParaffinWax28 Condl=0.15 Conds=0.20 Ql=180000 Ts=26 Tl=30 Tp=28 ;
				  PCM_HighPerf Condl=0.18 Conds=0.22 Ql=200000 Ts=24 Tl=26 Tp=25 -iterate ;
*`
		fi := NewEeTokens(input)

		var pcm []*PCM
		var pcmiterate rune

		PCMdata(fi, "test", &pcm, &pcmiterate)

		if len(pcm) != 2 {
			t.Fatalf("expected 2 PCM entries, got %d", len(pcm))
		}

		// Check first PCM entry
		pcm1 := pcm[0]
		if pcm1.Name != "ParaffinWax28" {
			t.Errorf("pcm[0].Name = %s, want ParaffinWax28", pcm1.Name)
		}
		if pcm1.Condl != 0.15 {
			t.Errorf("pcm[0].Condl = %f, want 0.15", pcm1.Condl)
		}
		if pcm1.Conds != 0.20 {
			t.Errorf("pcm[0].Conds = %f, want 0.20", pcm1.Conds)
		}
		if pcm1.Ql != 180000 {
			t.Errorf("pcm[0].Ql = %f, want 180000", pcm1.Ql)
		}
		if pcm1.Ts != 26 {
			t.Errorf("pcm[0].Ts = %f, want 26", pcm1.Ts)
		}
		if pcm1.Tl != 30 {
			t.Errorf("pcm[0].Tl = %f, want 30", pcm1.Tl)
		}
		if pcm1.Tp != 28 {
			t.Errorf("pcm[0].Tp = %f, want 28", pcm1.Tp)
		}
		if pcm1.Iterate != false {
			t.Errorf("pcm[0].Iterate = %t, want false", pcm1.Iterate)
		}

		// Check second PCM entry
		pcm2 := pcm[1]
		if pcm2.Name != "PCM_HighPerf" {
			t.Errorf("pcm[1].Name = %s, want PCM_HighPerf", pcm2.Name)
		}
		if pcm2.Iterate != true {
			t.Errorf("pcm[1].Iterate = %t, want true", pcm2.Iterate)
		}
	})
}

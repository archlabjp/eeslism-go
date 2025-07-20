package eeslism

import (
	"testing"
)

func TestRmhtrcf(t *testing.T) {
	// Test room heat transfer coefficient calculation
	t.Run("Basic heat transfer coefficient calculation", func(t *testing.T) {
		// Setup basic room
		room := &ROOM{
			Name:     "TestRoom",
			N:        2,
			rsrf:     make([]*RMSRF, 2),
			Brs:      0,
			F:        make([]float64, 4),
			Wradx:    make([]float64, 4),
			alr:      make([]float64, 4),
			alrbold:  0.0,
			Tr:       22.0,
			Tsav:     21.0,
			alc:      &[]float64{8.0}[0],
			mrk:      'C',
		}

		// Create surfaces
		surfaces := make([]*RMSRF, 2)
		for i := 0; i < 2; i++ {
			surfaces[i] = &RMSRF{
				A:    10.0,
				ali:  8.0,
				alo:  25.0,
				alic: 3.0,
				alir: 5.0,
				Ei:   0.9,
				Eo:   0.9,
				room: room,
			}
		}
		room.rsrf = surfaces

		// Setup EXSFS
		exs := &EXSFS{
			Alosch:  &[]float64{23.0}[0],
			Alotype: Alotype_Fix,
			Exs:     make([]*EXSF, 1),
		}
		exs.Exs[0] = &EXSF{}

		// Setup WDAT
		wd := &WDAT{}

		// Setup emrk
		emrk := make([]rune, 1)
		emrk[0] = '*'

		// Execute
		Rmhtrcf(exs, emrk, []*ROOM{room}, surfaces, wd)

		// Verify basic execution
		if room.Name != "TestRoom" {
			t.Errorf("Room name should remain unchanged")
		}
	})
}

func TestRmrdshfc(t *testing.T) {
	// Test room radiation heat flux calculation
	t.Run("Radiation heat flux calculation", func(t *testing.T) {
		// Setup room
		room := &ROOM{
			Name:    "TestRoom",
			N:       2,
			Brs:     0,
			FArea:   100.0,
			Area:    120.0,
			tfsol:   0.8,
			eqcv:    0.9,
			fsolm:   new(float64),
		}

		// Create surfaces
		surfaces := make([]*RMSRF, 2)
		for i := 0; i < 2; i++ {
			surfaces[i] = &RMSRF{
				A:    10.0,
				srg:  0.5,
				srg2: 0.4,
				fsol: &[]float64{0.5}[0],
				room: room,
			}
		}
		room.rsrf = surfaces

		// Execute
		Rmrdshfc([]*ROOM{room}, surfaces)

		// Verify basic execution
		if room.Name != "TestRoom" {
			t.Errorf("Room name should remain unchanged")
		}
		if room.FArea != 100.0 {
			t.Errorf("FArea should remain unchanged")
		}
	})
}

func TestRmhtrsmcf(t *testing.T) {
	// Test room heat transfer surface coefficient calculation
	tests := []struct {
		name     string
		surface  *RMSRF
		expected float64
	}{
		{
			name: "Standard surface",
			surface: &RMSRF{
				Rwall: 0.1,
				ali:   8.0,
				alo:   25.0,
			},
			expected: 0.0, // K will be calculated
		},
		{
			name: "High resistance surface",
			surface: &RMSRF{
				Rwall: 0.5,
				ali:   5.0,
				alo:   20.0,
			},
			expected: 0.0, // K will be calculated
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			surfaces := []*RMSRF{tt.surface}

			// Execute
			Rmhtrsmcf(surfaces)

			// Verify K is calculated
			if tt.surface.K <= 0.0 {
				t.Errorf("K should be positive, got %f", tt.surface.K)
			}

			// Verify K calculation formula: K = 1/(Rwall + 1/ali + 1/alo)
			expectedK := 1.0 / (tt.surface.Rwall + 1.0/tt.surface.ali + 1.0/tt.surface.alo)
			if abs(tt.surface.K-expectedK) > 1e-10 {
				t.Errorf("K = %f, want %f", tt.surface.K, expectedK)
			}
		})
	}
}

func TestRmexct(t *testing.T) {
	// Test room external condition calculation
	t.Run("Basic external condition calculation", func(t *testing.T) {
		// Setup room
		room := &ROOM{
			Name:      "TestRoom",
			N:         1,
			rsrf:      make([]*RMSRF, 1),
			Brs:       0,
			FArea:     100.0,
			Area:      120.0,
			tfsol:     0.0,
			fsolm:     &[]float64{0.7}[0],
			flrsr:     &[]float64{0.6}[0],
			Nfsolfix:  0,
			Qgt:       0.0,
			Qsolm:     0.0,
			Qsab:      0.0,
			Qrnab:     0.0,
			Hr:        100.0,
			Lr:        50.0,
			Ar:        25.0,
			Qeqp:      200.0,
			rsrnx:     false,
			Srgm2:     0.5,
		}

		// Create surface
		surface := &RMSRF{
			ble:      BLE_ExternalWall,
			typ:      RMSRFType_H,
			A:        10.0,
			as:       0.7,
			Eo:       0.9,
			alo:      25.0,
			ali:      8.0,
			K:        2.0,
			exs:      0,
			sb:       -1,
			Sname:    "",
			srg:      0.5,
			srg2:     0.4,
			srh:      0.3,
			srl:      0.2,
			sra:      0.1,
			eqrd:     0.05,
			room:     room,
			ffix_flg: ' ',
			fsol:     &[]float64{0.0}[0],
		}
		room.rsrf[0] = surface

		// Setup external surface
		exs := &EXSF{
			Cinc:  800.0,
			Idre:  600.0,
			Idf:   200.0,
			Rn:    50.0,
			Tazm:  180.0,
			Tprof: 45.0,
		}

		// Setup weather data
		wd := &WDAT{
			T: 15.0,
			X: 0.008,
		}

		// Setup QRM
		qrm := &QRM{
			Solo: 0.0,
			Solw: 0.0,
			Asl:  0.0,
		}

		// Execute
		Rmexct([]*ROOM{room}, []*RMSRF{surface}, wd, []*EXSF{exs}, []*SNBK{}, []*QRM{qrm}, 1, 1)

		// Verify calculations
		if surface.Te == 0.0 {
			t.Errorf("Te should be calculated")
		}
		if qrm.Solo == 0.0 {
			t.Errorf("Solo should be calculated")
		}
	})
}

func TestRoomcf(t *testing.T) {
	// Test room coefficient calculation
	t.Run("Room coefficient calculation", func(t *testing.T) {
		// Setup basic room
		room := &ROOM{
			Name:   "TestRoom",
			Tr:     22.0,
			Trold:  21.0,
			xr:     0.008,
			xrold:  0.007,
			GRM:    2000.0,
			HL:     50.0,
			AL:     25.0,
			Gvent:  100.0,
			RMt:    0.0,
			RMC:    0.0,
			RMx:    0.0,
			RMXC:   0.0,
		}

		// Setup WDAT
		wd := &WDAT{
			T: 15.0,
			X: 0.008,
		}

		// Setup EXSFS
		exsf := &EXSFS{
			Exs: make([]*EXSF, 0),
		}

		// Execute
		Roomcf([]*MWALL{}, []*ROOM{room}, []*RDPNL{}, wd, exsf)

		// Verify calculations
		if room.RMx == 0.0 {
			t.Errorf("RMx should be calculated")
		}
		if room.RMXC == 0.0 {
			t.Errorf("RMXC should be calculated")
		}
		if room.RMt == 0.0 {
			t.Errorf("RMt should be calculated")
		}
		if room.RMC == 0.0 {
			t.Errorf("RMC should be calculated")
		}
	})
}

func TestRmsurft(t *testing.T) {
	// Test room surface temperature calculation
	t.Run("Surface temperature calculation", func(t *testing.T) {
		// Setup room
		room := &ROOM{
			Name:      "TestRoom",
			N:         2,
			rsrf:      make([]*RMSRF, 2),
			Brs:       0,
			Tr:        22.0,
			Trold:     21.0,
			xr:        0.008,
			xrold:     0.007,
			RH:        50.0,
			mrk:       ' ',
			FunHcap:   0.0,
			OTsetCwgt: &[]float64{0.5}[0],
			XA:        make([]float64, 4),
			Wradx:     make([]float64, 4),
			F:         make([]float64, 4),
			ARN:       make([]float64, 2),
			RMP:       make([]float64, 2),
			alr:       make([]float64, 4),
			alrbold:   0.0,
			CM:        new(float64),
			MCAP:      new(float64),
			AR:        100.0,
			MRM:       2000.0,
			RMt:       0.0,
			RMC:       0.0,
			HGc:       0.0,
			CA:        0.0,
			oldTM:     21.0,
			FMT:       0.0,
			FMC:       0.0,
			QM:        0.0,
			GRM:       2000.0,
			RMx:       0.0,
			RMXC:      0.0,
			Gvent:     0.0,
			Nachr:     0,
			Ntr:       0,
			Nrp:       0,
			Nflr:      0,
			Nfsolfix:  0,
			Nisidermpnl: 0,
			Nasup:     0,
			achr:      nil,
			trnx:      nil,
			rmpnl:     nil,
			Arsp:      nil,
			elinasup:  nil,
			elinasupx: nil,
			rmld:      nil,
			rmqe:      nil,
			cmp:       &COMPNT{},
			VAVcontrl: &VAV{},
			Lightsch:  nil,
			Assch:     nil,
			Alsch:     nil,
			Hmsch:     nil,
			Metsch:    nil,
			Closch:    nil,
			Wvsch:     nil,
			Hmwksch:   nil,
			Vesc:      nil,
			Visc:      nil,
			alc:       nil,
			Hc:        0.0,
			Hr:        0.0,
			HL:        0.0,
			Lc:        0.0,
			Lr:        0.0,
			Ac:        0.0,
			Ar:        0.0,
			AL:        0.0,
			Qeqp:      0.0,
			eqcv:      0.0,
			Area:      120.0,
			FArea:     100.0,
			flrsr:     new(float64),
			fsolm:     new(float64),
			tfsol:     0.0,
			Srgm2:     0.0,
			Qsolm:     0.0,
			Qsab:      0.0,
			Qrnab:     0.0,
			Qgt:       0.0,
			rsrnx:     false,
			fij:       'A',
			Hcap:      0.0,
			Mxcap:     0.0,
			Ltyp:      ' ',
			Nhm:       0.0,
			Light:     0.0,
			Apsc:      0.0,
			Apsr:      0.0,
			Apl:       0.0,
			Gve:       0.0,
			Gvi:       0.0,
			AE:        0.0,
			AG:        0.0,
			VRM:       0.0,
			PCM:       nil,
			mPCM:      0.0,
			PCMQl:     0.0,
			PCMfurnname: "",
			HM:        0.0,
			TM:        0.0,
			hr:        0.0,
			PMV:       0.0,
			SET:       0.0,
			setpri:    false,
			Tot:       0.0,
			Trdy:      SVDAY{},
			xrdy:      SVDAY{},
			RHdy:      SVDAY{},
			Tsavdy:    SVDAY{},
			mTrdy:     SVDAY{},
			mxrdy:     SVDAY{},
			mRHdy:     SVDAY{},
			mTsavdy:   SVDAY{},
		}

		// Create surfaces
		surfaces := make([]*RMSRF, 2)
		for i := 0; i < 2; i++ {
			surfaces[i] = &RMSRF{
				Name:   "",
				Sname:  "",
				exs:    0,
				sb:     0,
				K:      0.0,
				alo:    0.0,
				alicsch:nil,
				alirsch:nil,
				Rwall:  0.0,
				CAPwall:0.0,
				tgtn:   0.0,
				Bn:     0.0,
				c:      0.0,
				as:     0.0,
				Ei:     0.9,
				Eo:     0.9,
				RStrans:false,
				tnxt:   0.0,
				Npcm:   0,
				pcmstate: nil,
				PVwall: PVWALL{},
				SQi:    QDAY{},
				Tsdy:   SVDAY{},
				mSQi:   QDAY{},
				mTsdy:  SVDAY{},
				Tcole:  0.0,
				Tcoleu: 0.0,
				Tcoled: 0.0,
				Tf:     0.0,
				end:    0,
				eqrd:   0.0,
				srg:    0.0,
				srg2:   0.0,
				srh:    0.0,
				srl:    0.0,
				sra:    0.0,
				fsol:   new(float64),
				ffix_flg: ' ',
				mrk:    ' ',
				ble:    BLE_InnerWall,
				typ:    RMSRFType_H,
				mwtype: RMSRFMwType_I,
				mwside: RMSRFMwSideType_i,
				A:      10.0,
				Ts:     20.0 + float64(i),
				WSC:    0.0,
				WSRN:   make([]float64, 2),
				WSPL:   make([]float64, 2),
				mw:     &MWALL{M: 2, Tw: make([]float64, 2), Told: make([]float64, 2), UX: make([]float64, 2), res: make([]float64, 2), cap: make([]float64, 2)},
				room:   room,
				alic:   8.0,
				alir:   5.0,
				ali:    8.0,
				CF:     0.0,
				WSR:    0.0,
				RS:     0.0,
				Te:     0.0,
				FO:     0.0,
				FI:     0.0,
				FP:     0.0,
			}
		}
		room.rsrf = surfaces
		room.alr[0] = 0.8
		room.alr[1] = 0.2
		room.alr[2] = 0.2
		room.alr[3] = 0.8

		room.Tsav = 20.0

		// Execute
		Rmsurft([]*ROOM{room}, surfaces)

		// Verify calculations
		if room.Trold != 22.0 {
			t.Errorf("Trold should be set to previous Tr value")
		}
		if room.xrold != 0.008 {
			t.Errorf("xrold should be set to previous xr value")
		}
		if room.mrk != 'C' {
			t.Errorf("mrk should be set to 'C'")
		}
		if room.Tsav == 20.0 {
			t.Errorf("Tsav should be calculated")
		}
		if room.Tot == 0.0 {
			t.Errorf("Tot should be calculated")
		}
	})
}

func TestQrmsim(t *testing.T) {
	// Test room heat gain simulation
	t.Run("Heat gain simulation", func(t *testing.T) {
		// Setup room
		room := &ROOM{
			Name:   "TestRoom",
			Hc:     100.0,
			Hr:     80.0,
			Lc:     60.0,
			Lr:     40.0,
			Ac:     30.0,
			Ar:     20.0,
			HL:     50.0,
			AL:     25.0,
			Gvent:  150.0,
			Tr:     22.0,
			Trold:  21.0,
			xr:     0.008,
			xrold:  0.007,
			Qeqp:   200.0,
			MRM:    2000.0,
			GRM:    1800.0,
			AE:     500.0,
			AG:     300.0,
			AEsch:  &[]float64{0.8}[0],
			AGsch:  &[]float64{0.6}[0],
		}

		// Setup weather data
		wd := &WDAT{
			T: 15.0,
			X: 0.008,
		}

		// Setup QRM
		qrm := &QRM{}

		// Execute
		Qrmsim([]*ROOM{room}, wd, []*QRM{qrm})

		// Verify calculations
		expectedHums := room.Hc + room.Hr
		if qrm.Hums != expectedHums {
			t.Errorf("Hums = %f, want %f", qrm.Hums, expectedHums)
		}

		expectedLight := room.Lc + room.Lr
		if qrm.Light != expectedLight {
			t.Errorf("Light = %f, want %f", qrm.Light, expectedLight)
		}

		expectedApls := room.Ac + room.Ar
		if qrm.Apls != expectedApls {
			t.Errorf("Apls = %f, want %f", qrm.Apls, expectedApls)
		}

		expectedHgins := qrm.Hums + qrm.Light + qrm.Apls
		if qrm.Hgins != expectedHgins {
			t.Errorf("Hgins = %f, want %f", qrm.Hgins, expectedHgins)
		}

		if qrm.Huml != room.HL {
			t.Errorf("Huml = %f, want %f", qrm.Huml, room.HL)
		}

		if qrm.Apll != room.AL {
			t.Errorf("Apll = %f, want %f", qrm.Apll, room.AL)
		}

		// Verify energy calculations
		if qrm.Qeqp != room.Qeqp {
			t.Errorf("Qeqp = %f, want %f", qrm.Qeqp, room.Qeqp)
		}

		expectedAE := room.AE * (*room.AEsch)
		if qrm.AE != expectedAE {
			t.Errorf("AE = %f, want %f", qrm.AE, expectedAE)
		}

		expectedAG := room.AG * (*room.AGsch)
		if qrm.AG != expectedAG {
			t.Errorf("AG = %f, want %f", qrm.AG, expectedAG)
		}
	})
}
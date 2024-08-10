package eeslism

type bekt struct {
	ps [][]float64
}

/*--付設障害物--*/
type sunblk struct {
	sbfname    string /*HISASI or BARUKONI or SODEKABE or MADOHIYOKE*/
	snbname    string
	rgb        [3]float64 /*--色--*/
	x, y       float64
	D, W, H, h float64
	WA         float64
	ref        float64 /*--反射率--*/
}

/*---窓---*/
type MADO struct {
	winname string     /*--名前--*/
	xr, yr  float64    /*--左下頂点座標--*/
	Ww, Wh  float64    /*--巾、高さ--*/
	ref     float64    /*--反射率--*/
	grpx    float64    /*--前面地面の代表点までの距離 初期値=1---*/
	rgb     [3]float64 /*--色--*/
}

/*---RMP---*/
type RRMP struct {
	rmpname  string /*--RMP名--*/
	wallname string
	sumWD    int        /*--窓の数--*/
	ref      float64    /*--反射率--*/
	xb0, yb0 float64    /*--左下頂点座標--*/
	Rw, Rh   float64    /*--巾、高さ--*/
	grpx     float64    /*--前面地面の代表点までの距離 初期値=1---*/
	rgb      [3]float64 /*--色--*/
	WD       []*MADO
}

/*---BDP---*/
type BBDP struct {
	bdpname         string    /*--BDP名--*/
	sumRMP, sumsblk int       /*--RMPの数、日よけの数--*/
	x0, y0, z0      float64   /*--左下頂点座標--*/
	Wa, Wb          float64   /*--方位角、傾斜角--*/
	exh, exw        float64   /*--巾、高さ--*/
	RMP             []*RRMP   /*RMP*/
	SBLK            []*sunblk /*SBLK*/

	// Satoh修正（2018/1/23）
	exsfname string
}

/*---OBS 外部障害物---*/
type OBS struct {
	fname   string     /*--rect or cube or r_tri or i_tri--*/
	obsname string     /*--名前--*/
	x, y, z float64    /*--左下頂点座標--*/
	H, D, W float64    /*--巾、奥行き、高さ--*/
	Wa      float64    /*--方位角--*/
	Wb      float64    /*--傾斜角--*/
	ref     [4]float64 /*--反射率--*/
	rgb     [3]float64 /*--色--*/
}

/*---樹木---*/
type TREE struct {
	treetype       string  /*--樹木の形Ａ，Ｂ，Ｃ--*/
	treename       string  /*--名前--*/
	x, y, z        float64 /*--幹部下面の中心座標--*/
	W1, W2, W3, W4 float64 /*--W1=幹太さ,W2=葉部下面巾,W3=葉部中央巾,W4=葉部上面巾--*/
	H1, H2, H3     float64 /*--H1=幹高さ,H2=葉部下側高さ,H3=葉部上側高さ--*/
}

/*---日射遮蔽率---*/
type SHADTB struct {
	lpname       string      /*--対象ＬＰ名--*/
	indatn       int         /*--入力データの数--*/
	ndays, ndaye [12]int     /*--スケジュール開始日と終了日--*/
	shad         [12]float64 /*--日射遮蔽率--*/
}

/*--多角形の頂点座標--*/
type XYZ struct {
	X, Y, Z float64
}

/*--OPW:受照窓面--*/
type WD_MENN struct {
	opwname string     /*--名前--*/
	P       []XYZ      /*--頂点--*/
	ref     float64    /*--反射率--*/
	grpx    float64    /*--前面地面の代表点までの距離 初期値=1---*/
	sumw    float64    /*--窓面の影面積の割合--*/
	rgb     [3]float64 /*--色R,G,B--*/
	polyd   int        /*--何角形か？--*/
}

/*--OP（受照面）,LP（被受照面）,MP(OP+OPW)--*/
type P_MENN struct {
	opname              string       /*--名前--*/
	rgb                 [3]float64   /*--色--*/
	wd, exs             int          /*--窓の数、方位番号--*/
	grpx                float64      /*--前面地面の代表点までの距離 初期値=1---*/
	faia                float64      /*--天空に対する形態係数--*/
	faig                float64      /*--地面に対する形態係数--*/
	faiwall             [200]float64 /*--外部障害物に対する形態係数--*/
	grpfaia             float64      /*--前面地面代表点から見た天空の形態係数--*/
	sum                 float64      /*--壁面の影面積--*/
	ref, refg           float64      /*--反射率、前面地面の反射率--*/
	wa                  float64      /*--面の方位角--*/
	wb                  float64      /*--面の傾斜角--*/
	Ihor, Idre, Idf, Iw float64      /*--日射量--*/
	Reff, rn            float64      /*--大気放射量、夜間放射量--*/
	Te, Teg             float64      /*--面の表面温度、前面地面の表面温度--*/
	shad                [366]float64 /*--面の日射遮蔽率--*/
	alo, as, Eo         float64      /*--外表面総合熱伝達率、日射吸収率、放射率--*/
	Nopw                int
	opw                 []WD_MENN
	polyd               int   /*--何角形か--*/
	P                   []XYZ /*--頂点座標(配列長はpolydに一致する)--*/
	e                   XYZ   /*--法線ベクトル--*/
	G                   XYZ   /*--中心点--*/
	grp                 XYZ   /*--前面地面代表点--*/
	sbflg               int   /*--付設障害物フラグ　付設障害物の場合：１、その他：０--*/
	wlflg               int   /*--外表面の種類 窓：1 壁：0 --*/
}

/*--LP(ポリゴン)直接入力用--*/
type POLYGN struct {
	polyknd   string     /*--ポリゴン種類(RMP OBS WD)--*/
	polyname  string     /*--名前--*/
	wallname  string     /*--壁名称--*/
	polyd     int        /*--何角形か? 3,4,5,6--*/
	ref, refg float64    /*--反射率、前面地面の反射率--*/
	P         []XYZ      /*--頂点--*/
	grpx      float64    /*--前面地面の代表点までの距離 初期値=1---*/
	rgb       [3]float64 /*--色--*/
}

/*---Sdstr 影面積のストア 110413 higuchi add---*/
type SHADSTR struct {
	sdsum []float64 /*--影面積--*/
}

/*--- 110413 higuchi end ----*/

type NOPLPMP struct {
	Nop int // 外部障害物の数
	Nlp int // 外部障害物の数
	Nmp int // 外部障害物の数(合計)
}

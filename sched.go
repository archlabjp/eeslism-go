package main

// 季節設定
type SEASN struct {
	name       string
	N          int
	sday, eday []int
	end        int // 要素数(インデックス0のみ)
}

// 曜日設定
type WKDY struct {
	name string
	wday [15]int
	end  int // 要素数(インデックス0のみ)
}

// 一日の設定量スケジュ－ル
type DSCH struct {
	name         string
	N            int
	stime, etime []int
	val          []float64
	end          int // 要素数(インデックス0のみ)
}

// 一日の切り替えスケジュ－ル
type DSCW struct {
	name         string
	dcode        [10]rune
	N            int
	stime, etime []int
	Nmod         int
	mode         []rune
	end          int // 要素数(インデックス0のみ)
}

type SCH struct /*スケジュ－ル*/
{
	name string
	Type rune
	day  [366]int
	end  int // 要素数(インデックス0のみ)
}

type SCHDL struct {
	Nsch  int // `Sch`の要素数
	Nscw  int // `Scw`の要素数
	Seasn []SEASN
	Wkdy  []WKDY
	Dsch  []DSCH
	Dscw  []DSCW
	Sch   []SCH
	Scw   []SCH
	Val   []float64 // `Sch`の要素数と同数
	Isw   []rune    // `Scw`の要素数と同数
}

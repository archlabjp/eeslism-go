package eeslism

import (
	"fmt"
	"math"
	"os"
	"unicode"
)

// 文字列s が数値かどうかの判定
func isstrdigit(s string) bool {
	for i := 0; i < len(s); i++ {
		if !unicode.IsDigit(rune(s[i])) {
			if s[i] != '.' && s[i] != '-' && s[i] != '+' {
				return false
			}
		}
	}
	return true
}

/* 入力データエラーの出力 */
func Errprint(err int, key string, s string) {
	if err != 0 {
		fmt.Printf(ERRFMTA, key, s)
		if Ferr != nil {
			fmt.Fprintf(Ferr, ERRFMTA, key, s)
		}
	}
}

func Eprint(key string, s string) {
	fmt.Printf(ERRFMTA, key, s)
	if Ferr != nil {
		fmt.Fprintf(Ferr, ERRFMTA, key, s)
	}
}

/* データの記憶域確保時のエラー出力 */
func Ercalloc(n int, errkey string) {
	s := fmt.Sprintf(" -- calloc   n=%d", n)
	Eprint(errkey, s)
}

func Preexit() {
	var NSTOP int
	fmt.Printf("Press Hit Return Key .......\n")
	if NSTOP == 0 {
		var buf [1]byte
		os.Stdin.Read(buf[:])
	} else {
		os.Exit(1)
	}
}

func Lineardiv(A, B, dt float64) float64 {
	return A + (B-A)*dt
}

// ttmmから1時間間隔の時刻へ変換する関数
// 0:01～1:00を1時（ここでは配列番号として0～23にしている）とする
func ConvertHour(ttmm int) int {
	tt := int(math.Floor(float64(ttmm-1) / 100.))
	return tt
}

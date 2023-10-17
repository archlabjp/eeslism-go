package eeslism

import (
	"fmt"
	"strconv"
	"strings"
)

// システム要素周囲条件（温度など）のポインター
func envptr(s string, Simc *SIMCONTL, Compnt []COMPNT, Wd *WDAT, Exsf *EXSFS) *float64 {
	var err error
	var vptr VPTR
	var dmy []MPATH
	var val *float64

	if isStrDigit(s) {
		// 固定値へのポインタを作成
		num, err2 := readFloat(s)
		if err2 != nil {
			panic(err2)
		}
		val = new(float64)
		*val = num
	} else {
		vptr, _, err = kynameptr(s, Simc, Compnt, 0, dmy, Wd, Exsf)
		if err == nil && vptr.Type == VAL_CTYPE {
			val = vptr.Ptr.(*float64)
		} else {
			fmt.Println("<*envptr>", s)
		}
	}

	if val == nil {
		fmt.Printf("xxxx  %s\n", s)
	}

	return val
}

func roomptr(s string, Compnt []COMPNT) *ROOM {
	var rm *ROOM

	for i := range Compnt {
		if s != "" && Compnt[i].Name != "" && strings.Compare(s, Compnt[i].Name) == 0 {
			rm, _ = Compnt[i].Eqp.(*ROOM)
			break
		}
	}

	return rm
}

func isStrDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

/*********** Satoh Create  2001/5/3 ********************/
func hccptr(c byte, s string, Compnt []COMPNT, m *rune) interface{} {
	var i int
	var h interface{}

	h = nil

	for i = range Compnt {
		if s != "" && s == Compnt[i].Name {
			if c == 'c' && Compnt[i].Eqptype == HCCOIL_TYPE {
				h = Compnt[i].Eqp.(*HCC)
				*m = 'c'
				return h
			} else if c == 'h' && Compnt[i].Eqptype == HCLOADW_TYPE {
				h = Compnt[i].Eqp.(*HCLOAD)
				*m = 'h'
				return h
			}
		}
	}

	return h
}

/*********** Satoh Create  2003/5/17 ********************/
/* 放射パネルの検索 */

func rdpnlptr(s string, Compnt []COMPNT) *RDPNL {
	var i int
	var h *RDPNL

	h = nil

	for i = range Compnt {
		if s == Compnt[i].Name {
			if Compnt[i].Eqptype == RDPANEL_TYPE {
				h = Compnt[i].Eqp.(*RDPNL)
				return h
			}
		}
	}

	return h
}

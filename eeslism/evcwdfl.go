package eeslism

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

/*
wdflinit (Weather Data File Initialization)

この関数は、気象データファイルから地域情報（緯度、経度、標準子午線、地盤温度など）を読み込み、
対応する構造体（`LOCAT`）に格納します。
また、気象データファイルから抽出された各気象要素（温度、湿度、日射量など）へのポインターを設定します。

建築環境工学的な観点:
- **地域情報の重要性**: 建物のエネルギーシミュレーションでは、
  その建物の位置する地域の気象条件が、
  熱負荷やエネルギー消費量に大きく影響します。
  この関数は、緯度（`loc.Lat`）、経度（`loc.Lon`）、
  標準子午線（`loc.Ls`）などの地理的情報を読み込み、
  太陽位置計算や日射量計算の基礎とします。
- **地盤温度の考慮**: `loc.Tgrav`（地盤温度）や`loc.DTgr`（地盤温度の時定数）は、
  地盤からの熱伝達をモデル化する際に用いられます。
  これは、地下室や基礎からの熱損失・熱取得を評価する上で重要です。
- **気象要素へのポインター設定**: 気象データファイルから読み込まれた各気象要素（`wp.Ta`：温度、`wp.Xa`：絶対湿度、
  `wp.Idn`：法線面直達日射量、`wp.Isky`：水平面天空日射量、`wp.Ihor`：水平面全天日射量、
  `wp.Cc`：雲量、`wp.Wdre`：風向、`wp.Wv`：風速、`wp.Rh`：相対湿度、`wp.Rn`：夜間放射量）へのポインターを設定します。
  これにより、シミュレーション中にこれらの気象データにアクセスし、
  熱負荷計算や機器の運転制御に利用することが可能になります。
- **給水温度の考慮**: `loc.Twsup`は、
  給水温度の月別データであり、
  給湯負荷計算などに用いられます。

この関数は、建物のエネルギーシミュレーションにおいて、
外部環境条件を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func wdflinit(Simc *SIMCONTL, Estl *ESTL, Tlist []TLIST) {
	var wp WDPT
	var s, ss, Err string
	var dt float64
	var id, N, i, m int

	if s = Estl.Wdloc; s == "" {
		return
	}

	Err = fmt.Sprintf(ERRFMT, "(wdflinit)")
	Simc.Loc = NewLOCAT()
	loc := Simc.Loc

	m = -1

	for _, field := range strings.Fields(s) {
		N = len(field)
		if st := strings.Index(field, "="); st >= 0 {
			ss = field[:st]
			dt, _ = strconv.ParseFloat(field[st+1:], 64)
			switch ss {
			case "Lat":
				loc.Lat = dt
			case "Lon":
				loc.Lon = dt
			case "Ls":
				loc.Ls = dt
			case "Tgrav":
				loc.Tgrav = dt
			case "DTgr":
				loc.DTgr = dt
			case "daymx":
				loc.Daymxert = int(dt)
			default:
				id = 1
			}
		} else {
			if field == "-" || m >= 0 {
				switch field {
				case "-Twsup":
					m = 0
				default:
					loc.Twsup[m], _ = strconv.ParseFloat(field, 64)
					m++
					if m > 11 {
						m = -1
					}
				}
			} else {
				loc.Name = field
			}
		}
		s = s[N:]
		for len(s) > 0 && unicode.IsSpace(rune(s[0])) {
			s = s[1:]
		}
	}

	if id != 0 {
		fmt.Printf("%s %s\n", Err, ss)
	}

	wp.Ta = nil
	wp.Xa = nil
	wp.Rn = nil
	wp.Ihor = nil
	wp.Rh = nil
	wp.Cc = nil
	wp.Wv = nil
	wp.Wdre = nil

	for i = 0; i < Estl.Ndata; i++ {
		t := &Tlist[i]
		if t.Name == "Wd" {
			s = t.Id
			val := t.Fval
			switch s {
			case "T": // 温度
				wp.Ta = val
			case "x": // 絶対湿度
				wp.Xa = val
			case "Idn": // 法線面直達日射量
				wp.Idn = val
			case "Isky": // 水平面天空日射量
				wp.Isky = val
			case "Ihor": // 水平面全天日射量
				wp.Ihor = val
			case "CC": // 雲量
				wp.Cc = val
			case "Wdre": // 風向
				wp.Wdre = val
			case "Wv": // 風速
				wp.Wv = val
			case "RH": // 相対湿度
				wp.Rh = val
			case "RN": // 夜間放射量
				wp.Rn = val
			}
		}
	}

	Simc.Wdpt = wp
}

/*
Wdflinput (Weather Data File Input)

この関数は、気象データファイルから読み込まれた各気象要素（温度、湿度、日射量など）を、
現在の時刻の気象データとして`WDAT`構造体（`Wd`）に格納します。
また、気象要素間の相互関係に基づいて、
不足しているデータを補完したり、夜間放射量を計算したりします。

建築環境工学的な観点:
- **気象データの更新**: シミュレーションの各時間ステップで、
  現在の時刻の気象データを`Wd`構造体に格納します。
  これにより、建物が受ける外部環境条件をリアルタイムでモデル化できます。
- **データ補完と整合性**: 気象データファイルによっては、
  全ての気象要素が揃っていない場合があります。
  この関数は、
  - `if Wd.X > 0.0 && Wd.RH < 0.0`: 絶対湿度（`Wd.X`）が与えられ、相対湿度（`Wd.RH`）が不足している場合に、
    `FNRhtx`関数を呼び出して相対湿度を計算します。
  - `else if Wd.X < 0.0 && Wd.RH > 0.0`: 相対湿度が与えられ、絶対湿度が不足している場合に、
    `FNXtr`関数を呼び出して絶対湿度を計算します。
  これにより、気象データの整合性を確保し、
  熱負荷計算の精度を向上させます。
- **夜間放射量の計算**: `Wd.RN`（夜間放射量）が与えられていない場合、
  雲量（`Wd.CC`）や空気の絶対湿度（`Wd.X`）に基づいて夜間放射量を計算します。
  夜間放射は、特に夜間の熱損失に影響を与える重要な要素です。
- **全天日射量の計算**: `wp.Ihor == nil` の場合、
  法線面直達日射量（`Wd.Idn`）と水平面天空日射量（`Wd.Isky`）から水平面全天日射量（`Wd.Ihor`）を計算します。

この関数は、建物のエネルギーシミュレーションにおいて、
外部環境条件を正確にモデル化し、
熱負荷計算、エネルギー消費量予測、
および省エネルギー対策の検討を行うための重要な役割を果たします。
*/
func Wdflinput(wp *WDPT, Wd *WDAT) {
	var Br float64

	Wd.T = wp.Ta[0]
	Wd.Idn = wp.Idn[0]
	Wd.Isky = wp.Isky[0]

	if wp.Ihor == nil {
		Wd.Ihor = Wd.Idn*Wd.Sh + Wd.Isky
	}

	if wp.Xa != nil {
		Wd.X = wp.Xa[0]
	} else {
		Wd.X = -999.0
	}

	if wp.Rh != nil {
		Wd.RH = wp.Rh[0]
	} else {
		Wd.RH = -999.0
	}

	if wp.Cc != nil {
		Wd.CC = wp.Cc[0]
	} else {
		Wd.CC = -999.0
	}

	if wp.Rn != nil {
		Wd.RN = wp.Rn[0]
	} else {
		Wd.RN = -999.0
	}

	if wp.Wv != nil {
		Wd.Wv = wp.Wv[0]
	} else {
		Wd.Wv = -999.0
	}

	if wp.Wdre != nil {
		Wd.Wdre = wp.Wdre[0]
	} else {
		Wd.Wdre = 0.0
	}

	if Wd.X > 0.0 && Wd.RH < 0.0 {
		Wd.RH = FNRhtx(Wd.T, Wd.X)
	} else if Wd.X < 0.0 && Wd.RH > 0.0 {
		Wd.X = FNXtr(Wd.T, Wd.RH)
	}

	if Wd.X > 0.0 {
		Wd.H = FNH(Wd.T, Wd.X)
	}

	if Wd.X > 0.0 && Wd.CC > 0.0 || Wd.RN < 0.0 {
		Br = 0.51 + 0.209*math.Sqrt(FNPwx(Wd.X))
		Wd.RN = (1.0 - 0.62*Wd.CC/10.0) * (1.0 - Br) * Sgm * math.Pow(Wd.T+273.15, 4.0)
		Wd.Rsky = ((1.0-0.62*Wd.CC/10.0)*Br + 0.62*Wd.CC/10.0) * Sgm * math.Pow(Wd.T+273.15, 4.0)
	} else {
		Wd.Rsky = Sgm*math.Pow(Wd.T+273.15, 4.0) - Wd.RN
	}
}

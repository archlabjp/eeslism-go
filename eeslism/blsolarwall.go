//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

// 空気集熱建材のシミュレーション

package eeslism

import (
	"fmt"
	"math"
)

// 集熱器の放射取得熱量
func FNScol(ta, I, EffPV, Ku, ao, Eo, Fs, RN float64) float64 {
	return (ta-EffPV)*I - Ku/ao*Eo*Fs*RN
}

// 建材一体型空気集熱器の相当外気温度を計算する
func CalcSolarWallTe(Rmvls *RMVLS, Wd *WDAT, Exsfs *EXSFS) {
	for i := range Rmvls.Rdpnl {
		rdpnl := &Rmvls.Rdpnl[i]
		Sd := rdpnl.sd[0]
		if Sd.mw != nil && Sd.mw.wall.WallType == WallType_C {
			Sd.Tcole = FNTcoleContrl(Sd, Wd, Exsfs)
		}
	}
}

// 集熱器相当外気温度の計算（制御用）
// 集熱器裏面温度は前時刻の値を使用する
func FNTcoleContrl(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) float64 {
	var Cidf float64
	var Wall *WALL
	var Exs *EXSF
	var Glsc float64
	var Ksu, alo, ku, kd float64

	if Sd.mw.wall.ColType != "" {
		Wall = Sd.mw.wall
		Exs = &Exsfs.Exs[Sd.exs]
		Glsc = 1.0
		Cidf = 1.0
		if Sd.mw.wall.ColType != "" &&
			Sd.mw.wall.ColType[:2] != "A2" &&
			Sd.mw.wall.ColType[:2] != "W2" {
			Glsc = Glscid(Exs.Cinc)
			Cidf = 0.91
		}

		//熱貫流率等の計算
		if Wall.chrRinput {
			FNKc(Wd, Exsfs, Sd)
			Ksu = Sd.dblKsu
			alo = Sd.dblao
			ku = Sd.ku
			kd = Sd.kd
		} else {
			//熱貫流率が入力値の場合
			Ksu = Wall.Ksu
			alo = *Exs.Alo
			ku = Wall.ku
			kd = Wall.kd
		}

		//Sd.Scol = FNScol(Wall.tra, Glsc*Exs.Idre+Cidf*Exs.Idf, Sd.PVwall.Eff, Ksu, alo, Wall.Eo, Exs.Fs, Wd.RN)
		//fmt.Printf("Ko=%g Scol=%g Ku=%g Ta=%g Kd=%g Tx=%g\n", Wall.Ko, Sd.Scol, Wall.Ku, Wd.T, Wall.Kd, Sd.oldTx)

		// Satoh DEBUG 2018/2/26  壁体一体型空気集熱器への影の影響を考慮するように修正
		Sd.Tcoled = Sd.oldTx
		Idre := Exs.Idre
		Idf := Exs.Idf
		RN := Exs.Rn

		Sd.dblSG = (Wall.tra - Sd.PVwall.Eff) * (Glsc*Idre + Cidf*Idf)
		if Sd.mw.wall.ColType[:2] == "A3" {
			Sd.Tcoled += Sd.dblSG / Sd.dblKsd
		}

		Sd.Tcoleu = Sd.dblSG/Ksu - Wall.Eo*RN/alo + Wd.T

		//return Wall.ku*(Sd.Scol/Wall.Ksu+Wd.T) + Wall.kd*Sd.oldTx

		//fmt.Printf("name=%s ku=%.2f kd=%.2f Tcoleu=%.2f Tcoled=%.2f\n", Sd.name, ku, kd, Sd.Tcoleu, Sd.Tcoled)
		return ku*Sd.Tcoleu + kd*Sd.Tcoled
	} else {
		return 0
	}
}

// 建材一体型空気集熱パネルの境界条件計算
func FNBoundarySolarWall(Sd *RMSRF, ECG, ECt, CFc *float64) {
	Wall := Sd.mw.wall
	Pnl := Sd.rpnl

	// 戻り値の初期化
	*ECG = 0.0

	// 各種熱貫流率の設定
	var Kc, Kcd, ku, kd float64
	if Wall.chrRinput {
		Kc = Sd.dblKc
		//Kcu = Sd.dblKcu
		Kcd = Sd.dblKcd
		ku = Sd.ku
		kd = Sd.kd
	} else {
		Kc = Wall.Kc
		//Kcu = Wall.Kcu
		Kcd = Wall.Kcd
		ku = Wall.ku
		kd = Wall.kd
	}

	// パネル動作時
	if Pnl.cG > 0.0 {
		*ECG = Pnl.cG * Pnl.Ec / (Kc * Sd.A)
	}

	*ECt = Kcd * ((1.0-*ECG)*ku - 1.0)
	*CFc = Kcd * (1.0 - *ECG) * kd
}

// 熱媒平均温度の計算
func FNTf(Tcin, Tcole, ECG float64) float64 {
	return (1.0-ECG)*Tcole + ECG*Tcin
}

// 外表面の総合熱伝達率の計算
func FNSolarWallao(Wd *WDAT, Sd *RMSRF, Exsfs *EXSFS) float64 {
	var dblac, dblar, dblao float64
	var Exs *EXSF
	var dblWdre float64
	var dblWa float64
	var dblu float64

	// 放射熱伝達率
	// 屋根の表面温度は外気温度で代用
	dblar = Sd.Eo * 4.0 * Sgm * math.Pow((Wd.T+Sd.Tg)/2.0+273.15, 3.0)

	// 対流熱伝達率
	Exs = &Exsfs.Exs[Sd.exs]

	// 外部風向の計算（南面0゜に換算）
	dblWdre = Wd.Wdre*22.5 - 180.0
	// 外表面と風向のなす角
	dblWa = float64(Exs.Wa) - dblWdre

	// 風上の場合
	if math.Cos(dblWa*math.Pi/180.0) > 0.0 {
		if Wd.Wv <= 2.0 {
			dblu = 2.0
		} else {
			dblu = 0.25 * Wd.Wv
		}
	} else {
		dblu = 0.3 + 0.05*Wd.Wv
	}

	// 対流熱伝達率
	dblac = 3.5 + 5.6*dblu

	// 総合熱伝達率
	dblao = dblar + dblac

	return dblao
}

// 通気層の放射熱伝達率の計算
func VentAirLayerar(dblEsu, dblEsd, dblTsu, dblTsd float64) float64 {
	var dblEs float64

	// 放射率の計算
	dblEs = 1.0 / (1.0/dblEsu + 1.0/dblEsd - 1.0)

	return 4.0 * dblEs * Sgm * math.Pow((dblTsu+dblTsd)/2.0+273.15, 3.0)
}

// 通気層の強制対流熱伝達率（ユルゲスの式）
func FNJurgesac(Sd *RMSRF, dblV, a, b float64) float64 {
	var Dh, Nu, Re, lam float64

	//if(fabs(dblV) <= 1.0e-3)
	//	dblTemp = 3.0 ;
	//else if(dblV <= 5.0)
	//	dblTemp = 7.1 * pow(dblV, 0.78) ;
	//else
	//	dblTemp = 5.8 + 3.9 * dblV ;

	// 長方形断面の相当直径の計算
	Dh = 1.232 * (a * b) / (a + b)
	// レイノルズ数の計算
	Re = dblV * Dh / FNanew(Sd.dblTf)
	// 空気の熱伝導率
	lam = FNalam(Sd.dblTf)
	// ヌセルト数の計算
	Nu = 0.0158 * math.Pow(Re, 0.8)
	return Nu / Dh * lam
}

// 屋根一体型空気集熱器の熱伝達率、熱貫流率の計算
func FNKc(Wd *WDAT, Exsfs *EXSFS, Sd *RMSRF) {
	var dblDet, dblWsuWsd, Ru, Cr, Cc float64
	//g := 9.81 // Assuming the value of gravity
	M_rad := math.Pi / 180.0

	Wall := Sd.mw.wall
	Exs := &Exsfs.Exs[Sd.exs]
	rad := M_rad

	// 外表面の総合熱伝達率の計算
	Sd.dblao = FNSolarWallao(Wd, Sd, Exsfs)

	// 通気層の対流熱伝達率の計算
	if Wall.air_layer_t < 0.0 {
		fmt.Printf("%s  Ventilation layer thickness is undefined\n", Sd.Name)
	}
	if Sd.rpnl.cmp.Elouts[0].G > 0.0 {
		Sd.dblacc = FNJurgesac(Sd, Sd.rpnl.cmp.Elouts[0].G/Roa/((Sd.dblWsd+Sd.dblWsu)/2.0*Wall.air_layer_t),
			(Sd.dblWsd+Sd.dblWsu)/2.0, Wall.air_layer_t)
	} else {
		Sd.dblacc = FNVentAirLayerac(Sd.dblTsu, Sd.dblTsd, Wall.air_layer_t, Exs.Wb*rad)
	}

	if math.Abs(Sd.dblacc) > 100.0 || Sd.dblacc < 0.0 {
		fmt.Printf("xxxxxx <FNKc> name=%s acc=%f\n", Sd.Name, Sd.dblacc)
	}

	// 通気層の放射熱伝達率の計算
	Sd.dblacr = VentAirLayerar(Wall.dblEsu, Wall.dblEsd, Sd.dblTsu, Sd.dblTsd)

	if math.Abs(Sd.dblacr) > 100.0 || Sd.dblacr < 0.0 {
		fmt.Printf("xxxxxx <FNKc> name=%s acr=%f\n", Sd.Name, Sd.dblacr)
	}

	// 通気層上面、下面から境界までの熱貫流率の計算
	if Wall.Ru >= 0.0 {
		Ru = Wall.Ru
	} else {
		// 空気層のコンダクタンスの計算
		// 放射成分
		Cr = VentAirLayerar(Wall.Eg, Wall.Eb, Sd.Tg, Sd.dblTsu)
		// 対流成分 component
		Cc = FNVentAirLayerac(Sd.Tg, Sd.dblTsu, Wall.ta, Exs.Wb*rad)
		Ru = 1.0 / (Cc + Cr)
		Sd.ras = Ru
	}

	Sd.dblKsu = 1.0 / (Ru + 1.0/Sd.dblao)
	Sd.dblKsd = 1.0 / Wall.Rd

	dblWsuWsd = Sd.dblWsu / Sd.dblWsd

	dblDet = (Sd.dblacr*dblWsuWsd+Sd.dblacc+Sd.dblKsd)*(Sd.dblacr+Sd.dblacc+Sd.dblKsu) - Sd.dblacr*Sd.dblacr*dblWsuWsd
	Sd.dblb11 = (Sd.dblacr + Sd.dblacc + Sd.dblKsu) / dblDet
	Sd.dblb12 = Sd.dblacr / dblDet
	Sd.dblb21 = Sd.dblacr * dblWsuWsd / dblDet
	Sd.dblb22 = (Sd.dblacr*dblWsuWsd + Sd.dblacc + Sd.dblKsd) / dblDet

	Sd.dblfcu = Sd.dblacc/dblWsuWsd*Sd.dblb12 + Sd.dblacc*dblWsuWsd*Sd.dblb22
	Sd.dblfcd = Sd.dblacc/dblWsuWsd*Sd.dblb11 + Sd.dblacc*dblWsuWsd*Sd.dblb21

	Sd.dblKcu = Sd.dblKsu * Sd.dblfcu
	Sd.dblKcd = Sd.dblKsd * Sd.dblfcd

	// 集熱器の総合熱貫流率の計算
	Sd.dblKc = Sd.dblKcu + Sd.dblKcd

	Sd.ku = Sd.dblKcu / Sd.dblKc
	Sd.kd = Sd.dblKcd / Sd.dblKc
}

// 空気集熱器の通気層上面、下面表面温度の計算
func FNTsuTsd(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) {
	//var dblTf float64 // 集熱空気の平均温度
	Rdpnl := Sd.rpnl
	Wall := Sd.mw.wall
	Exs := &Exsfs.Exs[Sd.exs]
	cG := Sd.rpnl.cG

	if Wall.chrRinput {
		Kc := Sd.dblKc
		Ksu := Sd.dblKsu
		Ksd := Sd.dblKsd
		ECG := cG * Rdpnl.Ec / (Kc * Sd.A)
		Sd.dblTf = (1.0-ECG)*Sd.Tcole + ECG*Rdpnl.Tpi
		Sd.dblTsd = Sd.dblb11*Ksd*(Sd.Tcoled-Sd.dblTf) + Sd.dblb12*Ksu*(Sd.Tcoleu-Sd.dblTf) + Sd.dblTf
		Sd.dblTsu = Sd.dblb21*Ksd*(Sd.Tcoled-Sd.dblTf) + Sd.dblb22*Ksu*(Sd.Tcoleu-Sd.dblTf) + Sd.dblTf

		if Sd.dblTsd < -100 {
			fmt.Println("Error")
		}

		// 空気層の熱抵抗が未定義の場合はガラスの温度を計算する
		if Wall.Ru < 0.0 {
			Sd.Tg = (Sd.dblao*Wd.T + 1.0/Sd.ras*Sd.dblTsu + Wall.ag*Exs.Iw - Wall.Eo*Exs.Fs*Wd.RN) / (Sd.dblao + 1.0/Sd.ras)
		}
	}
}

// 通気層の集熱停止時の熱コンダクタンス[W/m2K]
func FNVentAirLayerac(Tsu, Tsd, air_layer_t, Wb float64) float64 {
	var Tas, Ra, anew, a, RacosWb, lama float64
	g := 9.81 // Assuming the value of gravity

	var Tsud float64
	if math.Abs(Tsu-Tsd) < 1.0e-4 {
		Tsud = Tsu + 0.1
	} else {
		Tsud = Tsu
	}

	// 通気層の温度
	Tas = (Tsud + Tsd) / 2.0
	// 空気の動粘性係数
	anew = FNanew(Tas)
	// 空気の熱拡散率
	a = FNaa(Tas)
	// 空気の熱伝導率
	lama = FNalam(Tas)
	// レーリー数
	Ra = g * (1.0 / Tas) * math.Abs(Tsud-Tsd) * math.Pow(air_layer_t, 3.0) / (anew * a)

	RacosWb = Ra * math.Cos(Wb)

	dblTemp := (1.0 + 1.44*math.Max(0.0, 1.0-1708.0/RacosWb)*(1.0-math.Pow(math.Sin(1.8*Wb), 1.6)*1708.0/RacosWb) + math.Max(math.Pow(RacosWb/5830.0, 1.0/3.0)-1.0, 0.0)) * air_layer_t / lama

	return dblTemp
}

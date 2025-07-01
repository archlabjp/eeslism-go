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

/*
FNScol (Function for Collector Solar Gain)

この関数は、空気集熱建材（ソーラーウォールなど）における集熱器の放射取得熱量を計算します。
これは、太陽エネルギーを熱として利用するパッシブソーラーシステムの性能評価において重要な要素です。

建築環境工学的な観点:
- **太陽エネルギーの利用**: ソーラーウォールは、建物の外壁に設置された集熱層で太陽光を吸収し、
  その熱を室内に取り込むことで暖房負荷を軽減するシステムです。
  この関数は、その集熱器がどれだけの太陽エネルギーを熱として取得できるかを定量的に評価します。
- **放射取得熱量の構成要素**: 放射取得熱量は、主に以下の要素によって決定されます。
  - `(ta - EffPV) * I`: 透過率（`ta`）と太陽電池の発電効率（`EffPV`）を考慮した、
    集熱器が吸収する日射量（`I`）からの熱取得。
    `EffPV`は、太陽電池が組み込まれている場合に、日射エネルギーの一部が電力に変換されるため、
    熱として利用できない部分を差し引くことを示唆します。
  - `Ku / ao * Eo * Fs * RN`: 集熱器からの熱損失。
    `Ku`は集熱器の熱損失係数、`ao`は外表面の熱伝達率、`Eo`は放射率、
    `Fs`は天空に対する形態係数、`RN`は夜間放射量を示唆します。
    集熱器の性能は、太陽エネルギーの取得量を最大化しつつ、熱損失を最小限に抑えることで向上します。
- **パッシブソーラーシステムの評価**: この関数で計算される放射取得熱量は、
  ソーラーウォールが建物に供給する熱量の主要な部分を占めます。
  これにより、ソーラーウォールが建物の暖房負荷をどれだけ削減できるか、
  また、その熱的性能が地域の気象条件や建物の設計にどのように依存するかを評価できます。

この関数は、空気集熱建材の熱的性能を評価し、
パッシブソーラーシステムを導入した建物の省エネルギー性能を予測するための基礎となります。
*/
func FNScol(ta, I, EffPV, Ku, ao, Eo, Fs, RN float64) float64 {
	return (ta-EffPV)*I - Ku/ao*Eo*Fs*RN
}

/*
CalcSolarWallTe (Calculate Solar Wall Equivalent Outdoor Temperature)

この関数は、建材一体型空気集熱器（ソーラーウォール）の「相当外気温度」を計算します。
相当外気温度は、集熱器の熱的挙動を評価し、その熱取得性能をモデル化するために用いられる重要な概念です。

建築環境工学的な観点:
- **相当外気温度の概念**: 相当外気温度（Equivalent Outdoor Temperature, Te）は、
  日射や夜間放射などの影響を、あたかも外気温度が変化したかのように見なして、
  熱伝達計算を簡略化するために導入される仮想的な温度です。
  ソーラーウォールの場合、集熱器表面が受ける日射熱や、周囲への放射熱損失を、
  この相当外気温度に含めることで、集熱器からの熱取得量をより簡単に計算できるようになります。
- **集熱器の熱的境界条件**: ソーラーウォールは、外気と接する面、集熱層、そして室内と接する面から構成されます。
  この関数で計算される`Sd.Tcole`は、集熱器の熱的境界条件の一つとして機能し、
  集熱器内部の温度分布や、集熱された熱が室内へどれだけ供給されるかを計算する上で不可欠です。
- **制御との関連**: 関数名に`Contrl`とあるように、この相当外気温度は、
  ソーラーウォールの運転制御（例: ファンによる空気循環の開始・停止）にも利用される可能性があります。
  例えば、相当外気温度が一定値を超えた場合に集熱を開始するなど、
  システムの効率的な運用に貢献します。
- **WallType_Cの条件**: `Sd.mw.wall.WallType == WallType_C`という条件は、
  この計算が特定の壁タイプ（おそらく集熱機能を持つ壁）にのみ適用されることを示唆しています。
  これは、建物の様々な部位における熱的特性の違いを考慮した、詳細なモデル化の一環です。

この関数は、ソーラーウォールのようなパッシブソーラーシステムの熱的性能を正確に評価し、
建物のエネルギーシミュレーションにおいて、自然エネルギーの利用効果を定量的に把握するために不可欠です。
*/
func CalcSolarWallTe(Rmvls *RMVLS, Wd *WDAT, Exsfs *EXSFS) {
	for i := range Rmvls.Rdpnl {
		rdpnl := Rmvls.Rdpnl[i]
		Sd := rdpnl.sd[0]
		if Sd.mw != nil && Sd.mw.wall.WallType == WallType_C {
			Sd.Tcole = FNTcoleContrl(Sd, Wd, Exsfs)
		}
	}
}

/*
FNTcoleContrl (Function for Collector Equivalent Outdoor Temperature Control)

この関数は、建材一体型空気集熱器（ソーラーウォール）の「集熱器相当外気温度」を計算します。
この温度は、集熱器の熱的性能を評価し、特に集熱システムの運転制御に利用されることを想定しています。

建築環境工学的な観点:
- **集熱器の熱的性能評価**: ソーラーウォールは、日射を吸収して空気を加熱し、その熱を室内に供給します。
  この関数で計算される集熱器相当外気温度は、集熱器がどれだけの熱を効率的に取得できるかを示す指標となります。
  日射量、外気温度、集熱器の熱的特性（透過率、熱貫流率など）を総合的に考慮して算出されます。
- **日射取得と熱損失のバランス**: `Sd.dblSG`は日射による熱取得、`Wall.Eo*RN/alo`は夜間放射による熱損失、
  `Wd.T`は外気温度の影響を表しています。
  集熱器の性能は、日射取得を最大化しつつ、熱損失を最小限に抑えることで向上します。
- **熱貫流率等の計算 (FNKc)**:
  集熱器の熱的特性を正確に評価するためには、その熱貫流率（`Ksu`, `ku`, `kd`など）を適切に設定する必要があります。
  `FNKc`関数を呼び出すことで、これらの熱的パラメータが計算され、集熱器相当外気温度の算出に用いられます。
  `Wall.chrRinput`は、これらの値が入力値として与えられるか、計算によって求められるかを制御します。
- **影の影響の考慮**: `Sd.dblSG = (Wall.tra - Sd.PVwall.Eff) * (Glsc*Idre + Cidf*Idf)` の部分で、
  `Glsc`や`Cidf`といった係数が用いられていることから、
  日射の入射角や影の影響が集熱器の熱取得に与える影響を考慮していることが伺えます。
  これは、実際の建物における日射取得の複雑さをモデル化する上で重要です。
- **制御への応用**: この関数で得られる集熱器相当外気温度は、
  集熱システムの運転判断基準として利用されることが考えられます。
  例えば、この温度が室内設定温度を上回る場合にファンを稼働させて集熱空気を室内に供給するなど、
  システムの自動制御に貢献します。

この関数は、建材一体型空気集熱器の熱的挙動を詳細にモデル化し、
パッシブソーラーシステムの設計最適化や、エネルギーマネジメント戦略の策定に不可欠な情報を提供します。
*/
func FNTcoleContrl(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) float64 {
	var Cidf float64
	var Wall *WALL
	var Exs *EXSF
	var Glsc float64
	var Ksu, alo, ku, kd float64

	if Sd.mw.wall.ColType != "" {
		Wall = Sd.mw.wall
		Exs = Exsfs.Exs[Sd.exs]
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

/*
FNBoundarySolarWall (Function for Boundary Conditions of Solar Wall)

この関数は、建材一体型空気集熱パネルの熱的境界条件に関する係数を計算します。
特に、集熱パネル内部の空気流と熱伝達の相互作用をモデル化するために用いられます。

建築環境工学的な観点:
- **集熱パネルの熱伝達**: ソーラーウォールのような集熱パネルでは、
  パネル表面での日射吸収、パネル内部の空気流による熱輸送、
  そしてパネルを構成する材料の熱伝導が複雑に絡み合って熱が伝達されます。
  この関数は、これらの熱伝達メカニズムを簡略化された係数（`ECG`, `ECt`, `CFc`）として表現します。
- **熱媒（空気）の熱輸送効率 (ECG)**:
  `ECG`は、集熱パネル内部を流れる熱媒（空気）が、
  パネルが取得した熱をどれだけ効率的に輸送できるかを示す係数です。
  `Pnl.cG`は空気の流量や比熱、`Pnl.Ec`は集熱効率、`Kc`はパネルの熱貫流率、`Sd.A`はパネル面積を示唆します。
  この係数が高いほど、集熱された熱が効率的に利用されることを意味します。
- **熱貫流率の設定**: `Wall.chrRinput`の条件分岐は、
  パネルの熱貫流率（`Kc`, `Kcd`, `ku`, `kd`など）が、
  外部から入力されるか、あるいは内部で計算されるかを示しています。
  これにより、様々なパネル構造や材料特性に対応できるようになります。
- **熱収支方程式への組み込み**: 計算された`ECG`, `ECt`, `CFc`は、
  集熱パネルを含む建物の熱収支方程式に組み込まれ、
  パネルからの熱供給が室内環境に与える影響を評価するために用いられます。
  これらの係数は、パネルの熱的性能を簡潔に表現し、
  より大きなスケールの建物全体の熱シミュレーションに統合することを可能にします。

この関数は、建材一体型空気集熱パネルの熱的挙動をモデル化し、
パッシブソーラーシステムの設計や性能評価において、
パネル内部の熱輸送効率を考慮した正確な熱収支計算を行うために不可欠です。
*/
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

/*
FNTf (Function for Mean Fluid Temperature)

この関数は、空気集熱器内部を流れる熱媒（空気）の平均温度を計算します。
熱媒の平均温度は、集熱器の熱的性能を評価し、
集熱された熱がどれだけ効率的に利用されるかを判断する上で重要です。

建築環境工学的な観点:
- **熱媒の温度変化**: 空気集熱器では、空気が集熱層を通過する際に日射熱を吸収し、温度が上昇します。
  この関数は、集熱器の入口温度（`Tcin`）と、集熱器の熱的境界条件を代表する温度（`Tcole`）を考慮し、
  熱媒がパネル内でどれだけ加熱されるかを平均的な温度として表現します。
- **熱輸送効率 (ECG) の影響**: `ECG`は、集熱器が取得した熱を熱媒がどれだけ効率的に輸送できるかを示す係数です。
  `ECG`が高いほど、熱媒の温度上昇が大きくなり、より多くの熱が輸送されることを意味します。
  ` (1.0 - ECG) * Tcole + ECG * Tcin` という式は、
  集熱器の熱的性能と熱媒の流量のバランスを考慮した平均温度の計算を示唆します。
- **集熱システムの性能評価**: 熱媒の平均温度は、集熱器から供給される熱量や、
  その熱が室内の暖房にどれだけ貢献できるかを評価する上で不可欠です。
  例えば、熱媒の温度が低いと、室内への熱供給が不十分になる可能性があります。

この関数は、空気集熱器の熱的性能を評価し、
パッシブソーラーシステムにおける熱媒の挙動をモデル化するために用いられます。
*/
func FNTf(Tcin, Tcole, ECG float64) float64 {
	return (1.0-ECG)*Tcole + ECG*Tcin
}

/*
FNSolarWallao (Function for Solar Wall Outdoor Overall Heat Transfer Coefficient)

この関数は、ソーラーウォール（空気集熱器）の外表面における総合熱伝達率を計算します。
外表面の総合熱伝達率は、集熱器から外気への熱損失を評価する上で非常に重要なパラメータです。

建築環境工学的な観点:
- **外表面の熱伝達メカニズム**: 建物の外表面からの熱伝達は、
  主に「対流熱伝達」と「放射熱伝達」の二つのメカニズムによって行われます。
  - **放射熱伝達率 (dblar)**:
    表面からの放射熱伝達のしやすさを表します。
    ステファン・ボルツマンの法則に基づき、表面の放射率（`Sd.Eo`）と、
    表面温度および周囲の平均放射温度（ここでは外気温度`Wd.T`と地盤温度`Sd.Tg`の平均）の4乗に比例します。
    夜間放射や周囲の建物からの放射の影響を考慮する上で重要です。
  - **対流熱伝達率 (dblac)**:
    表面と周囲の空気との間の熱伝達のしやすさを表します。
    主に風速（`Wd.Wv`）に依存し、風速が速いほど対流熱伝達率は大きくなります。
    `dblu`は風速を考慮した係数であり、風向（`Wd.Wdre`）と壁面の向き（`Exs.Wa`）の関係も考慮されています。
    風上側では風速の影響が大きく、風下側では小さくなるという物理現象をモデル化しています。
- **総合熱伝達率 (dblao)**:
  対流熱伝達率と放射熱伝達率の合計として計算されます。
  この値が大きいほど、外表面からの熱損失が大きくなることを意味します。
- **集熱器の熱損失評価**: ソーラーウォールは日射を吸収して熱を取得しますが、
  同時に外気への熱損失も発生します。
  この関数で計算される`dblao`は、この熱損失を定量的に評価するために用いられ、
  集熱器の効率を決定する重要な要素となります。

この関数は、ソーラーウォールのようなパッシブソーラーシステムの熱的性能を正確に評価し、
熱損失を最小限に抑えるための設計最適化に不可欠な情報を提供します。
*/
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
	Exs = Exsfs.Exs[Sd.exs]

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

/*
VentAirLayerar (Ventilated Air Layer Radiant Heat Transfer Coefficient)

この関数は、通気層における放射熱伝達率を計算します。
通気層は、建物の断熱性能向上や、ソーラーウォールのような集熱システムにおいて重要な役割を果たします。

建築環境工学的な観点:
- **通気層の役割**: 通気層は、壁体内部の湿気排出、日射熱の遮蔽、そして断熱性能の向上に寄与します。
  特に、ソーラーウォールでは、通気層を介して空気を循環させることで、日射熱を室内に取り込みます。
- **放射熱伝達のメカニズム**: 通気層内の放射熱伝達は、
  通気層を挟む二つの表面（例えば、外壁の裏面と内壁の表面）間の放射熱交換によって生じます。
  この熱伝達は、各表面の放射率と表面温度に依存します。
- **実効放射率 (dblEs)**:
  `dblEs = 1.0 / (1.0/dblEsu + 1.0/dblEsd - 1.0)` の式は、
  通気層を挟む二つの表面の実効放射率を計算するものです。
  `dblEsu`と`dblEsd`は、それぞれ上流側と下流側の表面の放射率を示唆します。
  この実効放射率を用いることで、複雑な多重反射を考慮せずに、
  二つの表面間の放射熱伝達を簡潔にモデル化できます。
- **温度依存性**: 放射熱伝達率は、表面温度の4乗に比例するため、
  通気層内の温度が高いほど放射による熱移動が大きくなります。
  `4.0 * dblEs * Sgm * math.Pow((dblTsu+dblTsd)/2.0+273.15, 3.0)` の式は、
  この温度依存性を考慮した放射熱伝達率の計算を示しています。

この関数は、通気層を持つ建物の熱的性能を正確に評価し、
特にソーラーウォールのようなシステムにおける熱移動メカニズムを詳細にモデル化するために不可欠です。
*/
func VentAirLayerar(dblEsu, dblEsd, dblTsu, dblTsd float64) float64 {
	var dblEs float64

	// 放射率の計算
	dblEs = 1.0 / (1.0/dblEsu + 1.0/dblEsd - 1.0)

	return 4.0 * dblEs * Sgm * math.Pow((dblTsu+dblTsd)/2.0+273.15, 3.0)
}

/*
FNJurgesac (Function for Jurges' Convective Heat Transfer Coefficient in Ventilated Air Layer)

この関数は、通気層における強制対流熱伝達率を、ユルゲスの式に基づいて計算します。
通気層内の空気流による熱移動を評価する上で重要な要素です。

建築環境工学的な観点:
- **通気層内の対流熱伝達**: 通気層では、空気の自然対流や強制対流によって熱が移動します。
  特に、ソーラーウォールのようにファンを用いて空気を循環させる場合（強制対流）、
  その熱伝達率を正確に評価することが、集熱効率や熱回収量を予測する上で不可欠です。
- **ユルゲスの式**: ユルゲスの式は、管内流やダクト内の強制対流熱伝達率を計算するための経験式の一つです。
  この関数では、通気層をダクトと見なし、その断面形状（`a`, `b`）から相当直径（`Dh`）を計算し、
  空気の流速（`dblV`）に基づいてレイノルズ数（`Re`）を算出します。
  レイノルズ数は、流れの様式（層流か乱流か）を示す無次元数であり、
  熱伝達率に大きく影響します。
- **ヌセルト数 (Nu)**:
  ヌセルト数は、対流熱伝達の効率を示す無次元数であり、
  熱伝達率、代表長さ、熱伝導率の関係を表します。
  この関数では、レイノルズ数に基づいてヌセルト数を計算し、
  それから対流熱伝達率（`dblTemp`）を導出します。
- **空気の物性値**: 空気熱伝導率（`lam`）や動粘性係数（`anew`）など、
  空気の物性値は温度によって変化するため、
  `FNanew`や`FNalam`といった関数を用いて、通気層内の空気温度に応じた値を適用します。

この関数は、通気層を持つ建物の熱的性能を詳細に評価し、
特にソーラーウォールのような強制換気システムにおける熱移動メカニズムを正確にモデル化するために不可欠です。
*/
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

/*
FNKc (Function for Overall Heat Transfer Coefficient of Roof-Integrated Air Collector)

この関数は、屋根一体型空気集熱器（ソーラーウォール）の熱伝達率および総合熱貫流率を計算します。
集熱器の熱的性能を総合的に評価し、熱取得量や熱損失量を正確に算出するために不可欠です。

建築環境工学的な観点:
- **集熱器の熱的モデル化**: 屋根一体型空気集熱器は、
  外表面、通気層、そして室内側表面から構成される複雑な熱伝達システムです。
  この関数は、各層における熱伝達（対流、放射）を考慮し、
  集熱器全体の熱的性能を代表する熱伝達率や熱貫流率を算出します。
- **外表面の総合熱伝達率 (Sd.dblao)**:
  `FNSolarWallao`関数を用いて計算される外表面の総合熱伝達率は、
  集熱器から外気への熱損失を評価する上で重要です。
  風速や放射の影響が考慮されます。
- **通気層の熱伝達率 (Sd.dblacc, Sd.dblacr)**:
  通気層内の空気流による対流熱伝達率（`Sd.dblacc`）と、
  通気層を挟む表面間の放射熱伝達率（`Sd.dblacr`）を計算します。
  通気層の厚さや空気流速、表面の放射率などが影響します。
  特に、強制対流（ファンによる送風）がある場合は`FNJurgesac`を、
  自然対流の場合は`FNVentAirLayerac`を使用するなど、
  空気流の様式に応じた適切なモデルが適用されます。
- **熱貫流率の構成要素 (Sd.dblKsu, Sd.dblKsd, Sd.dblKc)**:
  - `Sd.dblKsu`: 集熱器の上流側（外側）から通気層への熱貫流率。
  - `Sd.dblKsd`: 集熱器の下流側（室内側）から通気層への熱貫流率。
  - `Sd.dblKc`: 集熱器全体の総合熱貫流率。
  これらの熱貫流率は、集熱器を介した熱の出入りを定量的に評価するために用いられます。
- **熱流の分配 (Sd.ku, Sd.kd)**:
  `Sd.ku`と`Sd.kd`は、集熱器が取得した熱が、
  それぞれ上流側と下流側のどちらにどれだけ分配されるかを示す係数です。
  これにより、集熱された熱が効率的に室内へ供給されるかを評価できます。

この関数は、屋根一体型空気集熱器の熱的性能を詳細にモデル化し、
パッシブソーラーシステムの設計最適化や、
建物のエネルギーシミュレーションにおける熱取得量の正確な予測に不可欠な役割を果たします。
*/
func FNKc(Wd *WDAT, Exsfs *EXSFS, Sd *RMSRF) {
	var dblDet, dblWsuWsd, Ru, Cr, Cc float64
	//g := 9.81 // Assuming the value of gravity
	M_rad := math.Pi / 180.0

	Wall := Sd.mw.wall
	Exs := Exsfs.Exs[Sd.exs]
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

/*
FNTsuTsd (Function for Ventilated Air Layer Upper and Lower Surface Temperatures)

この関数は、空気集熱器の通気層における上面（外側）と下面（内側）の表面温度を計算します。
これらの温度は、通気層内の熱移動を詳細にモデル化し、
集熱器の性能を評価する上で重要な中間変数となります。

建築環境工学的な観点:
- **通気層内の温度分布**: 通気層は、日射熱の取得や外気との熱交換によって、
  その内部に温度勾配が生じます。
  この関数は、通気層の熱的特性（熱貫流率など）と、
  集熱器の熱的境界条件（集熱器相当外気温度`Sd.Tcole`、熱媒入口温度`Rdpnl.Tpi`）を考慮して、
  通気層の上面と下面の温度を算出します。
- **熱媒平均温度 (Sd.dblTf)**:
  通気層を流れる空気の平均温度（`Sd.dblTf`）は、
  通気層内の熱伝達率の計算や、集熱された熱がどれだけ効率的に輸送されるかを評価する上で重要です。
  この平均温度は、集熱器相当外気温度と熱媒入口温度、そして熱輸送効率（`ECG`）に基づいて計算されます。
- **ガラス表面温度 (Sd.Tg)**:
  通気層の上面がガラスである場合、その表面温度（`Sd.Tg`）も計算されます。
  ガラス表面温度は、結露の発生予測や、ガラスを介した日射熱取得量の評価に影響します。
  特に、`Wall.Ru < 0.0`の条件は、空気層の熱抵抗が未定義の場合にガラスの温度を計算することを示唆しており、
  これはモデルの柔軟性を示しています。
- **熱的挙動の可視化**: 通気層の上面と下面の温度を計算することで、
  通気層内部の熱的な挙動をより詳細に把握することができます。
  例えば、通気層が効果的に日射熱を遮蔽しているか、
  あるいは熱が効率的に室内へ供給されているかなどを評価するのに役立ちます。

この関数は、空気集熱器の熱的性能を詳細にモデル化し、
パッシブソーラーシステムの設計最適化や、
建物のエネルギーシミュレーションにおける熱取得量の正確な予測に不可欠な役割を果たします。
*/
func FNTsuTsd(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) {
	//var dblTf float64 // 集熱空気の平均温度
	Rdpnl := Sd.rpnl
	Wall := Sd.mw.wall
	Exs := Exsfs.Exs[Sd.exs]
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

/*
FNVentAirLayerac (Function for Ventilated Air Layer Convective Heat Transfer Coefficient)

この関数は、通気層における自然対流熱伝達率を計算します。
通気層内の空気流がファンなどによる強制的なものではなく、
温度差によって生じる自然な空気の動き（自然対流）によって熱が伝達される場合に適用されます。

建築環境工学的な観点:
- **自然対流のメカニズム**: 通気層内の自然対流は、
  通気層を挟む二つの表面間の温度差によって空気の密度差が生じ、
  それによって空気が循環することで熱が伝達される現象です。
  この熱伝達は、通気層の厚さ、傾斜角、そして温度差に依存します。
- **レイリー数 (Ra)**:
  レイリー数は、自然対流の強さを示す無次元数であり、
  浮力と粘性力、熱拡散の相対的な影響を表します。
  この関数では、通気層の温度（`Tas`）、温度差（`Tsud-Tsd`）、
  通気層の厚さ（`air_layer_t`）、空気の動粘性係数（`anew`）、
  熱拡散率（`a`）などを用いてレイリー数を計算します。
  レイリー数が大きいほど、自然対流が活発になり、熱伝達率も大きくなります。
- **傾斜角の影響 (Wb)**:
  `Wb`は通気層の傾斜角を示唆しており、
  `math.Cos(Wb)`や`math.Sin(1.8*Wb)`といった項が含まれていることから、
  通気層の傾斜が自然対流熱伝達に与える影響を考慮していることが伺えます。
  垂直な通気層と水平な通気層では、自然対流の挙動が大きく異なります。
- **空気の物性値**: 空気熱伝導率（`lama`）や動粘性係数（`anew`）、熱拡散率（`a`）など、
  空気の物性値は温度によって変化するため、
  `FNanew`や`FNaa`、`FNalam`といった関数を用いて、通気層内の空気温度に応じた値を適用します。

この関数は、通気層を持つ建物の熱的性能を正確に評価し、
特に自然換気やパッシブソーラーシステムにおける熱移動メカニズムを詳細にモデル化するために不可欠です。
*/
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

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

/*   bl_panel.c  */

package eeslism

import "math"

/*
FNKPT (Function for Temperature Correction Factor of PV)

この関数は、太陽電池モジュールの温度補正係数を計算します。
太陽電池の発電効率は、その表面温度に大きく依存するため、
この補正係数は発電量予測において非常に重要です。

建築環境工学的な観点:
  - **太陽電池の温度特性**: 太陽電池は、温度が上昇すると発電効率が低下するという特性を持っています。
    これは、半導体のバンドギャップが温度によって変化するためです。
    一般的に、結晶シリコン系太陽電池の場合、モジュール温度が1℃上昇するごとに、
    出力が約0.3%〜0.5%低下すると言われています。
  - **温度補正係数 (KPT)**:
    この関数で計算される`KPT`は、基準温度（通常25℃）における発電効率を1.0とした場合の、
    実際のモジュール温度（`TPV`）における発電効率の相対的な変化率を表します。
    `apmax`は、温度による出力低下率を示すパラメータであり、
    太陽電池モジュールの種類やメーカーによって異なります。
    `KPT = 1.0 + apmax * (TPV - 25.0) / 100.0` の式は、
    この温度特性を線形近似でモデル化したものです。
  - **発電量予測の精度向上**: 太陽電池の発電量は、日射量だけでなく、モジュール温度にも大きく左右されます。
    特に、夏季の高温時や、屋根・壁一体型など放熱がしにくい設置条件では、モジュール温度が高くなりやすく、
    発電効率の低下が顕著になります。
    この温度補正係数を適切に考慮することで、より現実的な発電量を予測し、
    建物のエネルギー収支計算や、太陽光発電システムの経済性評価の精度を向上させることができます。

この関数は、太陽光発電システムの性能評価において、
日射量だけでなく温度の影響を考慮した正確な発電量予測を行うための基礎となります。
*/
func FNKPT(TPV, apmax float64) float64 {
	return 1.0 + apmax*(TPV-25.0)/100.0
}

/*
PVwallcatinit (PV Wall Category Initialization)

この関数は、太陽電池が設置された壁面（PVウォール）の初期パラメータを設定します。
これらのパラメータは、太陽電池の発電性能や熱的特性を定義するために用いられます。

建築環境工学的な観点:
  - **太陽電池の性能パラメータ**: 太陽電池の発電性能は、その種類や設置方法によって大きく異なります。
    この関数で設定されるパラメータは、以下のような太陽電池の特性を定義します。
  - `Type`: 太陽電池の種類（例: 結晶シリコン、薄膜など）。
  - `Apmax`: 温度による出力低下率（% / ℃）。`FNKPT`関数で用いられる重要なパラメータです。
  - `KHD`, `KPD`, `KPM`, `KPA`, `EffINO`: それぞれ、日射強度、部分影、モジュールミスマッチ、
    配線・変換効率など、様々な要因による発電量の損失係数を示唆します。
    これらの係数は、理想的な条件下での発電量から、実際の設置環境下での発電量を算出するために用いられます。
  - `Ap`: 太陽電池モジュールの熱的特性に関連するパラメータ（熱伝達係数など）を示唆します。
  - `Rcoloff`, `Kcoloff`: 太陽電池の放熱特性に関連するパラメータを示唆します。
  - **システム設計の基礎**: これらの初期パラメータは、太陽光発電システムのシミュレーションを行う上での基礎となります。
    適切なパラメータを設定することで、建物のエネルギー消費量と太陽光発電による創エネルギー量のバランスを評価し、
    ネット・ゼロ・エネルギー・ビル（ZEB）などの目標達成に向けた設計検討に役立ちます。
  - **デフォルト値の重要性**: この関数は、特定の太陽電池モジュールが指定されない場合に適用されるデフォルト値を設定します。
    これらのデフォルト値は、一般的な太陽電池の特性を反映している必要があります。

この関数は、太陽光発電システムのシミュレーションモデルを構築する上で、
太陽電池の物理的・電気的特性を定義するための初期設定を行う重要な役割を担います。
*/
func PVwallcatinit(PVwallcat *PVWALLCAT) {
	PVwallcat.Type = 'C'
	PVwallcat.Apmax = -0.41
	PVwallcat.KHD = 1.0
	PVwallcat.KPD = 0.95
	PVwallcat.KPM = 0.94
	PVwallcat.KPA = 0.97
	PVwallcat.EffINO = 0.9
	PVwallcat.Ap = 10.0
	PVwallcat.Rcoloff = FNAN
	PVwallcat.Kcoloff = FNAN
}

/*
PVwallPreCalc (PV Wall Pre-Calculation)

この関数は、太陽電池の発電量計算において、時間によらず一定とみなせる係数を事前に計算します。
これにより、シミュレーションの各時間ステップでの計算負荷を軽減し、効率化を図ります。

建築環境工学的な観点:
  - **総合効率係数 (PVwallcat.KConst)**:
    この関数で計算される`PVwallcat.KConst`は、太陽電池の発電量に影響を与える様々な損失要因（日射強度、部分影、
    モジュールミスマッチ、配線・変換効率など）を統合した総合的な効率係数です。
    これは、太陽電池モジュールの種類や設置条件によって決まる固定的な特性を表します。
    `KConst = KHD * KPD * KPM * KPA * EffINO` のように、各損失係数を乗算することで算出されます。
  - **シミュレーション効率の向上**: 太陽光発電システムのシミュレーションでは、
    日射量や温度など、刻々と変化する外部環境条件に基づいて発電量を計算します。
    この`KConst`のように、時間的に変化しないパラメータを事前に計算しておくことで、
    シミュレーションの計算量を削減し、より高速な解析を可能にします。
    これは、特に長期間のシミュレーションや、多数のケーススタディを行う場合に有効です。
  - **発電量予測の基礎**: `KConst`は、太陽電池の基準条件下での発電性能と、
    設置環境における様々な損失を考慮した実効的な発電性能を関連付ける重要な係数です。
    この係数に、日射量や温度補正係数（`FNKPT`で計算される`KPT`）を乗じることで、
    実際の発電量を予測することができます。

この関数は、太陽光発電システムのシミュレーションモデルにおいて、
計算効率と発電量予測の正確性を両立させるための重要な前処理ステップとなります。
*/
func PVwallPreCalc(PVwallcat *PVWALLCAT) {
	PVwallcat.KConst = PVwallcat.KHD * PVwallcat.KPD * PVwallcat.KPM * PVwallcat.KPA * PVwallcat.EffINO
}

/*
FNTPV (Function for PV Temperature Calculation)

この関数は、太陽電池モジュールの表面温度（`TPV`）を計算します。
太陽電池の発電効率は温度に大きく依存するため、正確なモジュール温度の推定は、
発電量予測の精度を向上させる上で非常に重要です。

建築環境工学的な観点:
  - **太陽電池の熱収支**: 太陽電池モジュールの温度は、日射吸収による熱取得と、
    周囲への熱放散（対流、放射）のバランスによって決定されます。
    この関数は、モジュールが受ける日射量（`Ipv`）、周囲の空気温度（`Wd.T`）、
    そしてモジュールと周囲間の熱伝達係数（`wall.PVwallcat.Ap`, `*Exs.Alo`）を考慮して、
    熱収支方程式を解くことでモジュール温度を推定します。
  - **設置条件の影響**: 太陽電池の設置方法（屋根置き、壁一体型、架台設置など）や、
    背面空間の有無、通風条件などは、モジュールからの放熱特性に大きく影響します。
    例えば、壁一体型（BIPV）のようにモジュール背面が密閉されている場合、
    放熱がしにくいためモジュール温度が高くなりやすく、発電効率が低下する傾向があります。
    `Sd.rpnl != nil && Sd.rpnl.cG > 0.0` の条件分岐は、
    このような設置条件の違いによる熱的特性の変化を考慮している可能性があります。
  - **発電量への影響**: 計算された`TPV`は、`FNKPT`関数に渡され、
    温度補正係数（`KPT`）の算出に用いられます。
    これにより、モジュール温度の上昇による発電効率の低下が発電量予測に反映され、
    より現実的な発電量評価が可能となります。

この関数は、太陽光発電システムの発電量予測において、
日射量だけでなく、設置環境における熱的挙動を考慮した正確なモジュール温度を推定し、
発電量予測の信頼性を高めるために不可欠な役割を果たします。
*/
func FNTPV(Sd *RMSRF, Wd *WDAT, Exsfs *EXSFS) float64 {
	wall := Sd.mw.wall
	Exs := Exsfs.Exs[Sd.exs]
	Ipv := (wall.tra - Sd.PVwall.Eff) * Sd.Iwall

	var TPV float64
	if Sd.rpnl != nil && Sd.rpnl.cG > 0.0 {
		TPV = (wall.PVwallcat.Ap*Sd.Tf + *Exs.Alo*Wd.T + Ipv) / (wall.PVwallcat.Ap + *Exs.Alo)
	} else {
		TPV = (wall.PVwallcat.Kcoloff*Sd.oldTx + *Exs.Alo*Wd.T + Ipv) / (wall.PVwallcat.Kcoloff + *Exs.Alo)
	}

	return TPV
}

/*
CalcPowerOutput (Calculate PV Power Output)

この関数は、各太陽電池が設置された壁面（PVウォール）からの発電量を計算します。
日射量、モジュール温度、そして様々な損失係数を考慮して、
実際の運用条件下での発電量を予測します。

建築環境工学的な観点:
  - **実発電量の予測**: 太陽光発電システムの設計や評価において、
    理論的な最大発電量だけでなく、実際の運用条件下での発電量を正確に予測することが重要です。
    この関数は、以下の要素を総合的に考慮することで、より現実的な発電量を算出します。
  - **日射量 (Sd[i].Iwall)**: 太陽電池が受ける日射量。これは発電量の最も基本的な要素です。
  - **モジュール温度 (pvwall.TPV)**: `FNTPV`関数で計算されたモジュール温度。
    温度上昇による発電効率の低下を考慮します。
  - **温度補正係数 (pvwall.KPT)**: `FNKPT`関数で計算された、モジュール温度による発電効率の補正係数。
  - **総合効率係数 (wall.PVwallcat.KConst)**: `PVwallPreCalc`関数で計算された、
    日射強度、部分影、モジュールミスマッチ、配線・変換効率など、
    様々な要因による損失を統合した係数。
  - **定格出力 (Sd[i].PVwall.PVcap)**: 太陽電池モジュールの最大定格出力。
  - **エネルギー収支への貢献**: 計算された発電量（`pvwall.Power`）は、
    建物のエネルギー消費量と対比され、建物のエネルギー収支（ZEB評価など）に組み込まれます。
    これにより、太陽光発電システムが建物のエネルギー自立性や環境負荷低減にどれだけ貢献するかを定量的に評価できます。
  - **発電効率の評価 (pvwall.Eff)**:
    発電効率は、太陽電池が受けた日射エネルギーに対して、どれだけの電力を生成できたかを示す指標です。
    この効率を計算することで、システムの性能を評価し、改善点を見つけることができます。
    特に、日射量が少ない場合や、モジュール温度が高い場合に効率がどのように変化するかを把握することは重要です。

この関数は、太陽光発電システムが建物全体のエネルギー性能に与える影響を評価し、
持続可能な建築設計やエネルギーマネジメント戦略を策定するための重要な情報を提供します。
*/
func CalcPowerOutput(Sd []*RMSRF, Wd *WDAT, Exsfs *EXSFS) {
	for i := range Sd {
		if Sd[i].mw != nil {
			wall := Sd[i].mw.wall

			/// 太陽電池が設置されているときのみ
			if Sd[i].PVwallFlg {
				pvwall := &Sd[i].PVwall

				pvwall.TPV = FNTPV(Sd[i], Wd, Exsfs)
				pvwall.KPT = FNKPT(pvwall.TPV, wall.PVwallcat.Apmax)
				pvwall.KTotal = wall.PVwallcat.KConst * pvwall.KPT

				pvwall.Power = math.Max(Sd[i].PVwall.PVcap*pvwall.KTotal*Sd[i].Iwall/1000.0, 0.0)

				// 発電効率の計算
				if Sd[i].Iwall > 0.0 {
					pvwall.Eff = pvwall.Power / (Sd[i].Iwall * Sd[i].A)
				} else {
					pvwall.Eff = 0.0
				}
			}
		}
	}
}

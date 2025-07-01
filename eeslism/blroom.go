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

/*  room.c       */
package eeslism

import "fmt"

/* ----------------------------------------------------------- */

/*
RMcf (Room Coefficient Calculation)

この関数は、室内の熱収支計算における各種係数を設定します。
特に、壁体（不透明部、窓）の熱応答に関する係数（FI, FO, FP）や、
室内表面間の放射熱交換を考慮した係数（XA, WSR, WSRN, WSPL）を算出します。

建築環境工学的な観点:
- **熱応答係数 (FI, FO, FP)**: 壁体や窓を介した熱の出入りを動的に評価するために不可欠です。
  FIは室内側からの熱流入、FOは室外側への熱流出、FPはパネルからの熱影響を表します。
  これらの係数は、壁体の熱容量や熱伝導率、表面熱伝達率、そして時間遅れ効果を考慮して計算されます。
  特に、PCM（相変化材料）が組み込まれた壁体の場合、その潜熱蓄熱効果が熱応答に大きく影響するため、
  PCMの有無や特性に応じて係数計算が分岐します。
- **室内表面間の放射熱交換 (XA, WSR, WSRN, WSPL)**:
  室内の各表面（壁、床、天井、窓など）は、互いに放射熱を交換します。
  XAは、各表面の温度が他の表面の温度にどのように影響するかを示す行列であり、
  形態係数や表面の放射率が考慮されます。
  WSRは、室内空気温度が各表面温度に与える影響、
  WSRNは、隣室からの熱伝達が各表面温度に与える影響、
  WSPLは、放射パネルからの熱が各表面温度に与える影響を表す係数です。
  これらの係数は、室内熱環境の均一性や快適性を評価する上で重要であり、
  特に放射冷暖房システムや日射侵入時の表面温度上昇の予測に用いられます。
- **家具の熱容量 (FunCoeffへの連携)**:
  室内の家具や備品も熱容量を持ち、室温変動を緩和する効果があります。
  この関数は、その計算を`FunCoeff`関数に委ねることで、室全体の熱収支モデルに組み込みます。
  特に、家具にPCMが内蔵されている場合、その潜熱蓄熱効果が室温安定化に寄与します。

この関数で計算される係数は、室内の熱収支方程式を解くための基礎となり、
室温、表面温度、熱負荷などの予測精度に直結します。
*/

/* ----------------------------------------------------------- */

func RMcf(Room *ROOM) {
	N := Room.N
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]

		if Sdn.mrk == '*' || Sdn.PCMflg {

			// 壁体（窓以外の場合）
			if Sdn.typ == RMSRFType_H || Sdn.typ == RMSRFType_E || Sdn.typ == RMSRFType_e {
				Mw := Sdn.mw
				M := Mw.M
				mp := Mw.mp

				if Sdn.mwside == RMSRFMwSideType_i {
					Sdn.FI = Mw.uo * Mw.UX[0]

					if Sdn.mw.wall.WallType != WallType_C {
						Sdn.FO = Mw.um * Mw.UX[M-1]
					} else {
						Sdn.FO = Sdn.ColCoeff * Mw.UX[M-1]
					}

					if Sdn.rpnl != nil {
						Sdn.FP = Mw.Pc * Mw.UX[mp] * Sdn.rpnl.Wp
					} else {
						Sdn.FP = 0.0
					}
				} else {
					MM := (M - 1) * M

					if Sdn.mw.wall.WallType != WallType_C {
						Sdn.FI = Mw.um * Mw.UX[MM+M-1]
					} else {
						Sdn.FI = Sdn.ColCoeff * Mw.UX[MM+M-1]
					}

					Sdn.FO = Mw.uo * Mw.UX[MM]
					if Sdn.rpnl != nil {
						Sdn.FP = Mw.Pc * Mw.UX[MM+mp] * Sdn.rpnl.Wp
					} else {
						Sdn.FP = 0.0
					}
				}
			} else {
				// 窓の場合
				/***      K = Sdn.K = 1.0/(Sdn.Rwall + rai + rao); ***/
				K := Sdn.K
				ali := Sdn.ali
				Sdn.FI = 1.0 - K/ali
				Sdn.FO = K / ali
				Sdn.FP = 0.0
			}
		}
	}

	alr := Room.alr
	XA := Room.XA
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]

		for j := 0; j < N; j++ {
			XA[n*N+j] = -Sdn.FI * alr[n*N+j] / Sdn.ali
		}

		XA[n*N+n] = 1.0
	}

	E := fmt.Sprintf("<RMcf> name=%s", Room.Name)
	Matinv(XA, N, N, E)

	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]

		Sdn.WSR = 0.0

		for j := 0; j < N; j++ {
			sdj := Room.rsrf[j]
			kc := sdj.alic / sdj.ali
			Sdn.WSR += XA[n*N+j] * sdj.FI * kc
		}

		for j := 0; j < Room.Ntr; j++ {
			wrn := &Sdn.WSRN[j]
			trn := Room.trnx[j]
			sdk := trn.sd

			// Find the index of sdk in Room.rsrf
			var kk int
			for kk = 0; kk < Room.N; kk++ {
				if sdk == Room.rsrf[kk] {
					break
				}
			}
			*wrn = XA[n*N+kk] * sdk.FO * sdk.nxsd.alic / sdk.nxsd.ali
		}

		for j := 0; j < Room.Nrp; j++ {
			sdk := Room.rmpnl[j].sd

			// Find the index of sdk in Room.rsrf
			var kk int
			for kk = 0; kk < Room.N; kk++ {
				if sdk == Room.rsrf[kk] {
					break
				}
			}

			// XA：室内表面温度計算のためのマトリックス
			// FP：パネルの係数
			Sdn.WSPL[j] = XA[n*N+kk] * sdk.FP
		}
	}

	Room.AR = 0.0
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]
		Room.AR += Sdn.A * Sdn.alic * (1.0 - Sdn.WSR)
	}

	// 室内空気の総合熱収支式の係数
	for j := 0; j < Room.Ntr; j++ {
		arn := 0.0
		for n := 0; n < N; n++ {
			sdk := Room.rsrf[n]
			arn += sdk.A * sdk.alic * sdk.WSRN[j]
		}
		Room.ARN[j] = arn
	}

	for j := 0; j < Room.Nrp; j++ { // 室のパネル総数
		rpnl := 0.0
		for n := 0; n < N; n++ {
			sdk := Room.rsrf[n]
			rpnl += sdk.A * sdk.alic * sdk.WSPL[j] // WSPL：パネルに関する係数
		}
		Room.RMP[j] = rpnl
	}

	// 室温の係数
	// 家具の熱容量の計算
	FunCoeff(Room)
}

/*
FunCoeff (Furniture Coefficient Calculation)

この関数は、室内の家具が持つ熱容量と、家具に内蔵されたPCM（相変化材料）の潜熱蓄熱効果を計算し、
室全体の熱収支モデルに組み込むための係数を設定します。

建築環境工学的な観点:
- **家具の熱容量 (Room.FunHcap)**:
  室内の家具や備品は、その質量と比熱に応じて熱を蓄える能力（熱容量）を持ちます。
  この熱容量は、室温の急激な変動を緩和する効果（熱的慣性）をもたらします。
  例えば、日中の日射取得による室温上昇を抑制したり、夜間の外気温度低下による室温降下を緩やかにしたりする役割があります。
  特に、軽量な建物において、家具の熱容量は室温安定化に寄与する重要な要素となります。
- **PCM（相変化材料）の潜熱蓄熱効果 (Room.PCMQl, Room.mPCM)**:
  PCMは、特定の温度範囲で相変化（固体から液体、またはその逆）する際に、大量の熱を潜熱として吸収・放出する材料です。
  家具にPCMを内蔵することで、室温がPCMの相変化温度に達した際に、
  その温度を一定に保ちながら熱を蓄えたり放出したりすることが可能になります。
  これにより、室温のピークカットや谷埋め効果が期待でき、空調負荷の低減や快適性の向上が図られます。
  `FNPCMStatefun`や`FNPCMstate_table`は、PCMの相変化特性（温度履歴、潜熱量）をモデル化し、
  現在の室温と前時刻の室温に基づいてPCMがどれだけの熱を吸収・放出しているかを計算します。
- **室温の係数 (Room.FMT, Room.RMt)**:
  これらの係数は、家具の熱容量やPCMの蓄熱効果を考慮した上で、
  室空気の熱収支方程式における室温項の動的な挙動を表現するために用いられます。
  `Room.FMT`は、家具の熱容量が室温変動に与える影響の度合いを示し、
  熱容量が大きいほど（`Room.FunHcap`が大きいほど）、室温変動が緩やかになることを示唆します。

この関数は、室内の熱的慣性や蓄熱性能を正確に評価するために不可欠であり、
特に省エネルギー設計や快適性評価において重要な役割を果たします。
*/
func FunCoeff(Room *ROOM) {
	// 室温の係数
	// 家具の熱容量の計算
	Room.FunHcap = 0.0
	if Room.CM != nil && *Room.CM > 0.0 {
		if Room.MCAP != nil && *Room.MCAP > 0.0 {
			Room.FunHcap += *Room.MCAP
		}
		if Room.PCM != nil {
			if Room.PCM.Spctype == 'm' {
				Room.PCMQl = FNPCMStatefun(Room.PCM.Ctype, Room.PCM.Cros, Room.PCM.Crol, Room.PCM.Ql,
					Room.PCM.Ts, Room.PCM.Tl, Room.PCM.Tp, Room.oldTM, Room.TM, Room.PCM.DivTemp, &Room.PCM.PCMp)
			} else {
				Room.PCMQl = FNPCMstate_table(&Room.PCM.Chartable[0], Room.oldTM, Room.TM, Room.PCM.DivTemp)
			}
			Room.FunHcap += Room.mPCM * Room.PCMQl
		}
	}
	if Room.FunHcap > 0.0 {
		Room.FMT = 1.0 / (Room.FunHcap/DTM/(*Room.CM) + 1.0)
	} else {
		Room.FMT = 1.0
	}

	Room.RMt = Room.MRM/DTM + Room.AR

	if Room.FunHcap > 0.0 {
		Room.RMt -= *Room.CM * (Room.FMT - 1.0)
	}
}

/*
RMrc (Room Constant Term Calculation)

この関数は、室内の熱収支計算における定数項（時間によらず一定とみなせる熱負荷や、
前時刻の室温に依存する項など）を計算します。
具体的には、壁体からの熱伝達、内部発熱、日射熱取得、隣室からの熱伝達、
そして家具の蓄熱効果などを考慮した係数を算出します。

建築環境工学的な観点:
- **壁体からの熱伝達 (Sdn.CF)**:
  壁体（不透明部、窓）を介して室内へ伝わる熱量を計算します。
  これは、壁体の熱貫流率、内外温度差、および前時刻の壁体内部温度履歴に依存します。
  特に、壁体の熱容量が大きい場合や、外気温度が大きく変動する場合に、
  壁体からの熱伝達が室温に与える影響は大きくなります。
- **内部発熱 (Room.HGc)**:
  室内の人体発熱、照明、機器発熱など、室内に直接加えられる熱負荷の合計です。
  これらは室温上昇の主要因の一つであり、空調設備の設計において重要な要素となります。
  `Room.Qeqp*Room.eqcv`は、機器発熱の係数と機器発熱量を表していると考えられます。
- **日射熱取得 (Sdn.FI*Sdn.RS/Sdn.ali)**:
  窓や透明な壁面を透過して室内に入り込む日射熱量を計算します。
  日射熱は、特に夏季の冷房負荷や冬季の暖房負荷軽減に大きく影響します。
  `Sdn.RS`は日射吸収量、`Sdn.ali`は室内側表面熱伝達率、`Sdn.FI`は日射熱取得係数を示唆します。
- **隣室からの熱伝達 (Sdn.FO*Sdn.Te)**:
  隣接する室や外部空間との温度差によって生じる熱伝達を計算します。
  `Sdn.Te`は隣室の温度や相当外気温度を表し、`Sdn.FO`は熱伝達係数を示唆します。
  これは、特に集合住宅やオフィスビルなど、複数の室が隣接する建物において重要です。
- **相互放射の計算 (Sdn.WSC, Room.CA)**:
  室内の各表面からの放射熱交換を考慮した係数を計算します。
  これは、室内の表面温度が室温に与える影響を評価するために用いられます。
  `XA`は放射形態係数を含む行列であり、各表面からの放射が他の表面にどのように影響するかを考慮します。
- **家具の影響項 (Room.FMC)**:
  `FunCoeff`関数で計算された家具の熱容量やPCMの蓄熱効果が、
  室空気の熱収支方程式の定数項に与える影響を考慮します。
  これにより、家具の蓄熱効果による室温変動の緩和が熱収支モデルに組み込まれます。

この関数で計算される定数項は、室内の熱収支方程式を解く際に、
室温の変動を引き起こす外部要因や内部要因を総合的に評価するために不可欠です。
*/
func RMrc(Room *ROOM) {
	N := Room.N
	XA := Room.XA
	CRX := make([]float64, N)

	for n := 0; n < N; n++ { // N：表面総数
		Sdn := Room.rsrf[n]
		Sdn.CF = 0.0
		if Sdn.typ == RMSRFType_H || Sdn.typ == RMSRFType_E || Sdn.typ == RMSRFType_e { // 壁の場合
			Mw := Sdn.mw
			M := Mw.M
			if Sdn.mwside != RMSRFMwSideType_M { // 室内側
				for j := 0; j < M; j++ {
					Sdn.CF += Mw.UX[j] * Mw.Told[j]
				}
			} else {
				MM := M * (M - 1)
				UX := Mw.UX[MM:]
				for j := 0; j < M; j++ {
					Sdn.CF += UX[j] * Mw.Told[j]
				}
			}
		}
	}

	Room.HGc = Room.Hc + Room.Lc + Room.Ac + Room.Qeqp*Room.eqcv

	// 表面熱収支に関係する係数の計算
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]
		CRX[n] = Sdn.CF + Sdn.FO*Sdn.Te + Sdn.FI*Sdn.RS/Sdn.ali
	}

	// 相互放射の計算
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]
		Sdn.WSC = 0.0
		for j := 0; j < N; j++ {
			Sdn.WSC += XA[n*N+j] * CRX[j]
		}
	}

	Room.CA = 0.0
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]
		Room.CA += Sdn.A * Sdn.alic * Sdn.WSC
	}

	// 室空気の熱収支の係数計算
	// 家具の影響項の追加
	if Room.FunHcap > 0.0 {
		dblTemp := DTM / Room.FunHcap
		Room.FMC = 1.0 / (dblTemp**Room.CM + 1.0) * (Room.oldTM + dblTemp*Room.Qsolm)
	} else {
		Room.FMC = 0.0
	}

	Room.RMC = Room.MRM/DTM*Room.Trold + Room.HGc + Room.CA
	if Room.FunHcap > 0.0 {
		Room.RMC += *Room.CM * Room.FMC
	}

}

/* ----------------------------------------------------- */
/*
RMsrt (Room Surface Temperature Calculation)

この関数は、室内の各表面（壁、床、天井、窓など）の温度を計算します。
表面温度は、室内空気温度、隣室温度、放射パネル温度、そして表面間の放射熱交換の影響を受けて決定されます。
また、計算された表面温度を用いて、各表面からの対流熱伝達量、放射熱伝達量、および総熱伝達量を算出します。

建築環境工学的な観点:
- **表面温度 (Sdn.Ts)**:
  室内の表面温度は、居住者の快適性（特に放射快適性）に直接影響を与える重要な要素です。
  例えば、冬期に窓表面温度が低いとコールドドラフトが発生し、不快感を引き起こします。
  夏期に日射が当たる壁面や窓の表面温度が高いと、放射熱によって不快感が増します。
  この関数では、以下の要素を考慮して表面温度を算出します。
  - `Sdn.WSR*Room.Tr`: 室内空気温度からの対流熱伝達の影響。
  - `Sdn.WSC`: 表面間の放射熱交換の影響。
  - `Sdn.WSRN[j]*trn.nextroom.Tr`: 隣室からの熱伝達の影響。
  - `Sdn.WSPL[j]*rmpnl.pnl.Tpi`: 放射パネルからの熱伝達の影響。
- **平均放射温度 (Sdn.Tmrt)**:
  平均放射温度（Mean Radiant Temperature, MRT）は、居住者が感じる放射熱環境を代表する温度です。
  室内の各表面温度とその表面に対する形態係数を考慮して計算されます。
  MRTは、作用温度やPMV（予測平均申告）などの快適性指標の算出に不可欠であり、
  特に放射冷暖房システムや日射侵入時の快適性評価において重要な指標となります。
  `alr`は形態係数を含む行列であり、各表面が他の表面から受ける放射の影響を考慮します。
- **表面からの熱伝達量 (Sdn.Qc, Sdn.Qr, Sdn.Qi)**:
  - `Sdn.Qc`: 表面から室内空気への対流熱伝達量。
  - `Sdn.Qr`: 表面から他の表面への放射熱伝達量。
  - `Sdn.Qi`: 表面からの総熱伝達量（対流＋放射－日射吸収）。
  これらの熱伝達量は、室内の熱負荷計算や、各表面の熱収支を詳細に分析するために用いられます。
  特に、結露の発生予測（表面温度が露点温度を下回るかどうかの判断）や、
  壁体内部の熱・湿気移動の評価において重要な情報となります。

この関数は、室内の熱環境を詳細に把握し、居住者の快適性評価や省エネルギー対策の検討に不可欠な情報を提供します。
*/
func RMsrt(Room *ROOM) {
	N := Room.N

	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]

		Sdn.Ts = Sdn.WSR*Room.Tr + Sdn.WSC

		for j := 0; j < Room.Ntr; j++ {
			trn := Room.trnx[j]
			Sdn.Ts += Sdn.WSRN[j] * trn.nextroom.Tr
		}

		for j := 0; j < Room.Nrp; j++ {
			rmpnl := Room.rmpnl[j]
			Sdn.Ts += Sdn.WSPL[j] * rmpnl.pnl.Tpi
		}
	}

	alr := Room.alr
	for n := 0; n < N; n++ {
		Sdn := Room.rsrf[n]
		Sdn.Tmrt = 0.0

		for j := 0; j < N; j++ {
			Sd := Room.rsrf[j]
			if j != n {
				Sdn.Tmrt += Sd.Ts * alr[n*N+j]
			}
		}
		Sdn.Tmrt /= alr[n*N+n]
	}

	for n := 0; n < N; n++ {
		Sd := Room.rsrf[n]
		Sd.Qc = Sd.alic * Sd.A * (Sd.Ts - Room.Tr)
		Sd.Qr = Sd.alir * Sd.A * (Sd.Ts - Sd.Tmrt)
		Sd.Qi = Sd.Qc + Sd.Qr - Sd.RS*Sd.A
	}
}

/*
RMwlc (Room Wall Coefficient Matrix Creation)

この関数は、室内の壁体（特に重量壁）の熱伝達に関する係数行列を作成します。
壁体の熱応答を正確にモデル化するために、室内側・室外側表面熱抵抗、壁体内部の熱容量、
そして壁体に組み込まれたパネルからの熱影響を考慮します。

建築環境工学的な観点:
- **重量壁の熱応答**: 重量壁は、その大きな熱容量により、室温変動を緩和する効果があります。
  この関数は、壁体内部の温度分布や熱流を計算するための基礎となる係数行列を構築します。
  これにより、壁体を介した熱の出入りが時間的にどのように遅れて室内に影響するかを評価できます。
- **表面熱抵抗 (rai, rao)**:
  室内側表面熱抵抗（rai）と室外側表面熱抵抗（rao）は、壁表面と空気間の熱伝達のしやすさを表します。
  これらの抵抗値は、対流と放射の複合的な効果を考慮したもので、
  壁体の熱貫流率や熱負荷計算において重要なパラメータとなります。
- **壁体内部の熱容量 (Mw.cap)**:
  壁体を構成する材料の熱容量は、壁体の蓄熱性能に直結します。
  この関数で作成される係数行列は、壁体内部の各層の熱容量を考慮し、
  壁体内部の温度が時間とともにどのように変化するかをモデル化します。
- **パネルからの熱影響 (Wp)**:
  壁体に放射パネルなどが組み込まれている場合、そのパネルからの熱が壁体内部の温度分布に影響を与えます。
  `Wp`はその影響度合いを示す係数であり、壁体全体の熱収支に組み込まれます。
- **Wallfdc関数への連携**: 最終的に、これらのパラメータを用いて`Wallfdc`関数を呼び出し、
  壁体の熱応答を記述する係数行列（`Mw.UX`, `Mw.uo`, `Mw.um`, `Mw.Pc`など）を生成します。
  これらの係数は、壁体内部の温度計算や、壁体を介した熱流の計算に用いられます。

この関数は、建物の熱的性能を詳細に評価し、特に蓄熱効果を考慮した省エネルギー設計や、
快適な室内環境の実現に向けた壁体設計の検討に不可欠な役割を果たします。
*/
func RMwlc(Mw []*MWALL, Exsfs *EXSFS, Wd *WDAT) {
	for i := range Mw {
		var Mw *MWALL = Mw[i]
		var Wall *WALL = Mw.wall

		var Sd *RMSRF = Mw.sd
		rai := 1.0 / Sd.ali // 室内側表面熱抵抗
		rao := 1.0 / Sd.alo // 室外側表面熱抵抗

		Mw.res[0] = rai
		if Sd.typ == 'H' {
			Mw.res[Mw.M] = rao
		}

		// 壁体にパネルがある場合
		var Wp float64
		if Sd.rpnl != nil {
			Wp = Sd.rpnl.Wp
		} else {
			Wp = 0.0
		}

		// 行列作成
		Wallfdc(Mw.M, Mw.mp, Mw.res, Mw.cap, Wp, Mw.UX,
			&Mw.uo, &Mw.um, &Mw.Pc, Wall.WallType, Sd, Wd, Exsfs, Wall,
			Mw.Told, Mw.Toldd, Mw.sd.pcmstate)
	}
}

/*
RMwlt (Room Wall Temperature Calculation)

この関数は、室内の壁体（不透明部）の内部温度を計算します。
壁体内部の温度分布は、室内空気温度、隣室温度（または相当外気温度）、
壁体表面からの日射吸収、そして壁体に組み込まれたパネルからの熱影響を考慮して決定されます。

建築環境工学的な観点:
- **壁体内部温度の動的挙動**: 壁体内部の温度は、外部環境や室内環境の変化に対して時間遅れを伴って応答します。
  この動的な挙動を正確にモデル化することは、建物の熱的慣性や蓄熱効果を評価する上で非常に重要です。
  特に、日射の侵入や外気温度の変動が室温に与える影響を予測するために不可欠です。
- **境界条件の設定 (Tie, Tee)**:
  - `Tie` (室内側相当温度): 室内空気温度、室内表面からの放射熱、および室内側表面での日射吸収を考慮した、
    壁体室内側の熱的境界条件を表します。
  - `Tee` (室外側相当温度): 共用壁の場合は隣室の温度、専用壁の場合は相当外気温度（外気温度、日射、夜間放射などを考慮した温度）を考慮した、
    壁体室外側の熱的境界条件を表します。
  これらの相当温度は、壁体を介した熱流の駆動力となります。
- **パネルからの熱影響 (WTp)**:
  壁体に放射パネルなどが組み込まれている場合、そのパネルからの熱が壁体内部の温度分布に影響を与えます。
  `WTp`はパネルからの熱供給量を示し、壁体内部の温度計算に組み込まれます。
- **Twall関数への連携**: 最終的に、これらの境界条件と壁体の熱的特性（熱応答係数など）を用いて`Twall`関数を呼び出し、
  壁体内部の各層の温度（`Mw.Tw`）を計算します。
  この計算は、壁体内部の熱伝導方程式を解くことで行われます。
- **温度履歴の更新**: 計算された壁体内部温度は、次の時間ステップの計算のために`Mw.Told`に更新されます。
  これにより、壁体の熱的履歴が考慮され、より正確な動的熱応答のシミュレーションが可能となります。

この関数は、建物の熱的性能評価において、壁体内部の複雑な熱挙動をモデル化するために不可欠であり、
特に高断熱・高気密住宅や、蓄熱性能を重視した建築物の設計・評価に重要な役割を果たします。
*/
func RMwlt(Mw []*MWALL) {
	for i := range Mw {
		Mw := Mw[i]
		Sd := Mw.sd

		// 壁体の反対側の表面温度 ?
		var Tee float64
		if Sd.mwtype == RMSRFMwType_C {
			// 共用壁の場合
			nxsd := Sd.nxsd
			Tee = (nxsd.alic*nxsd.room.Tr + nxsd.alir*nxsd.Tmrt + nxsd.RS) / nxsd.ali
		} else if Sd.mwtype == RMSRFMwType_I {
			// 専用壁の場合 => 外表面の相当外気温度
			Tee = Sd.Te
		} else {
			panic(Sd.mwtype)
		}

		Room := Sd.room
		Tie := (Sd.alic*Room.Tr + Sd.alir*Sd.Tmrt + Sd.RS) / Sd.ali

		if DEBUG {
			fmt.Printf("----- RMwlt i=%d room=%s ble=%c %s  Tie=%f Tee=%f\n", i, Sd.room.Name, Sd.ble, get_string_or_null(Sd.Name), Tie, Tee)
		}

		var WTp float64
		if Sd.rpnl != nil {
			WTp = Sd.rpnl.Wp * Sd.rpnl.Tpi
		} else {
			WTp = 0.0
		}

		// 壁体表面、壁体内部温度の計算
		Twall(Mw.M, Mw.mp, Mw.UX, Mw.uo, Mw.um, Mw.Pc, Tie, Tee, WTp, Mw.Told, Mw.Tw, Sd, Mw.wall.PCMLyr)

		// 壁体表面温度、壁体内部温度の更新
		for m := 0; m < Mw.M; m++ {
			// 前時刻の壁体内部温度を更新
			Mw.Told[m] = Mw.Tw[m]
			// 収束過程初期値の壁体内部温度を更新
			Mw.Twd[m] = Mw.Tw[m]
			Mw.Told[m] = Mw.Tw[m]
		}
	}
}

/*
RMwltd (Room Wall Temperature Temporary Calculation)

この関数は、壁体内部温度の仮計算を行います。
主に、PCM（相変化材料）が組み込まれた壁体において、
PCMの相変化を考慮した壁体内部温度の収束計算の初期段階で用いられます。

建築環境工学的な観点:
- **PCMの相変化と収束計算**: PCMは、相変化する際に大量の潜熱を吸収・放出するため、
  壁体内部の温度分布に非線形な影響を与えます。
  そのため、壁体内部温度の計算には反復的な収束計算が必要となる場合があります。
  この`RMwltd`関数は、その収束計算の初期値や中間段階での仮の温度分布を計算するために使用されます。
- **境界条件の再評価**: `RMwlt`関数と同様に、室内側相当温度（Tie）と室外側相当温度（Tee）を計算し、
  壁体内部の熱流の駆動力とします。
  PCMの相変化に伴い、これらの境界条件も動的に変化する可能性があるため、
  仮計算の段階でこれらを再評価することは重要です。
- **Twalld関数への連携**: `Twalld`関数は、PCMの特性を考慮した壁体内部の熱伝導方程式を解き、
  仮の壁体内部温度（`Mw.Twd`）を算出します。
  この仮温度は、次の反復計算の入力として使用され、最終的な収束解に近づけていきます。

この関数は、PCMを組み込んだ壁体の熱的挙動を正確にモデル化するために不可欠であり、
特にPCMの潜熱蓄熱効果を最大限に活用した省エネルギー設計や、
室温の安定化を図るためのシミュレーションにおいて重要な役割を果たします。
*/
func RMwltd(Mw []*MWALL) {
	for i := range Mw {
		var Mw *MWALL = Mw[i]
		var Sd *RMSRF = Mw.sd
		var nxsd *RMSRF = Sd.nxsd
		var Room *ROOM = Sd.room

		if Sd.PCMflg {
			// Tee
			var Tee float64
			if Sd.mwtype == RMSRFMwType_C {
				Tee = (nxsd.alic*nxsd.room.Tr + nxsd.alir*nxsd.Tmrt + nxsd.RS) / nxsd.ali
			} else if Sd.mwtype == RMSRFMwType_I {
				Tee = Sd.Te
			} else {
				panic(Sd.mwtype)
			}

			// Tie
			Tie := (Sd.alic*Room.Tr + Sd.alir*Sd.Tmrt + Sd.RS) / Sd.ali

			if DEBUG {
				fmt.Printf("----- RMwlt i=%d room=%s ble=%c %s  Tie=%f Tee=%f\n",
					i, Sd.room.Name, Sd.ble, get_string_or_null(Sd.Name), Tie, Tee)
			}

			// WTp
			var WTp float64
			if Sd.rpnl != nil {
				WTp = Sd.rpnl.Wp * Sd.rpnl.Tpi
			} else {
				WTp = 0.0
			}

			// 壁体内部温度の仮計算
			Twalld(Mw.M, Mw.mp, Mw.UX, Mw.uo, Mw.um, Mw.Pc,
				Tie, Tee, WTp, Mw.Told, Mw.Twd, Sd)
		}
	}
}

/*
RTsav (Room Temperature of Surface - Average)

この関数は、室内の全表面（壁、床、天井、窓など）の面積加重平均表面温度を計算します。

建築環境工学的な観点:
- **平均表面温度の重要性**: 平均表面温度は、室内の放射環境を評価する上で重要な指標の一つです。
  居住者が感じる快適性は、空気温度だけでなく、周囲の表面温度にも大きく影響されます。
  特に、放射冷暖房システムが導入されている場合や、窓からの日射侵入が大きい場合など、
  表面温度が空気温度と大きく異なる状況では、平均表面温度が快適性評価の鍵となります。
- **平均放射温度 (MRT) との関係**: 平均表面温度は、平均放射温度（MRT）の計算における基礎となります。
  MRTは、各表面の温度とその表面に対する形態係数を考慮して計算されるため、
  平均表面温度はMRTの簡易的な指標として、あるいはMRT計算の途中段階として利用されることがあります。
- **熱負荷計算への応用**: 室内の熱負荷計算において、壁体からの熱伝達量を評価する際に、
  室内空気温度と平均表面温度の差を用いることで、より実態に近い熱伝達量を算出できる場合があります。

この関数は、室内の熱環境の全体的な評価や、居住者の快適性評価、
そして熱負荷計算の補助的な指標として利用されます。
*/
func RTsav(N int, Sd []*RMSRF) float64 {
	var Tav, Aroom float64
	for n := 0; n < N; n++ {
		Tav += Sd[n].Ts * Sd[n].A
		Aroom += Sd[n].A
	}
	return Tav / Aroom
}

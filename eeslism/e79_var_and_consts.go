/*
e79_var_and_consts.go (Global Variables and Constants for Building Energy Simulation)

このファイルは、建物のエネルギーシミュレーション全体で共通して使用されるグローバル変数と定数を定義します。
これらの定数や変数は、熱計算の基礎となる物理定数、シミュレーションの制御パラメータ、
およびデバッグフラグなどを統一的に管理するために用いられます。

建築環境工学的な観点:
- **物理定数**: 熱計算や流体計算において、
  - `Sgm`: ステファン・ボルツマン定数 [W/(m2・K4)]。放射熱伝達の計算に用いられます。
  - `Ca`: 空気の定圧比熱 [J/(kg・K)]。顕熱計算に用いられます。
  - `Cv`: 水蒸気の定圧比熱 [J/(kg・K)]。潜熱計算に用いられます。
  - `Roa`: 空気の密度 [kg/m3]。流量から質量流量への変換などに用いられます。
  - `Cw`: 水の比熱 [J/(kg・K)]。水系の熱計算に用いられます。
  - `Row`: 水の密度 [kg/m3]。
  - `Ro`: 水の蒸発潜熱 [J/kg]。潜熱計算に用いられます。
  - `G`: 重力加速度 [m/s2]。
  - `ALO`: 外表面熱伝達率のデフォルト値 [W/(m2・K)]。
  - `UNIT`: 単位系（SI単位系など）。
  - `PI`: 円周率。
    これらの定数は、シミュレーションの正確性と物理的な妥当性を確保するために重要です。

- **シミュレーション制御パラメータ**:
  - `DTM`: 時間ステップの長さ [秒]。シミュレーションの細かさを決定します。
  - `Cff_kWh`: エネルギー単位変換係数（ジュールからkWhへ）。
  - `VAVCountMAX`: VAVシステムの収束計算の最大反復回数。

- **デバッグと出力制御**:
  - `DEBUG`: デバッグモードのON/OFFを切り替えるフラグ。
  - `dayprn`: 日ごとの詳細出力のON/OFFを切り替えるフラグ。
  - `Ferr`: エラーメッセージの出力先。
  - `SETprint`: SET（作用温度）の出力のON/OFFを切り替えるフラグ。
    これらのフラグは、シミュレーションの実行中に詳細な情報を取得し、
    モデルの検証や問題の特定を効率的に行うために用いられます。
  - **共通データ**: `Fbmlist`（壁材料リストファイル名）など、
    複数のモジュールで共有されるファイル名やデータが定義されます。

このファイルは、建物のエネルギーシミュレーションの基盤を形成し、
計算の正確性、安定性、および保守性を向上させるための重要な役割を果たします。
*/
package eeslism

import "os"

const (
	ALO = 23.0

	// Uncomment these lines if you want to use these constants
	// MAXBDP  = 100
	// MAXOBS  = 100
	// MAXTREE = 10 // Maximum number of trees
	// MAXPOLY = 50

	UNIT = "SI"
	PI   = 3.141592654
	FNAN = -999.0
	INAN = -999
)

var (
	Sgm = 5.67e-8
	Ca  = 1005.0
	Cv  = 1846.0
	Roa = 1.29
	Cw  = 4186.0

	Row = 1000.0
	Ro  = 2501000.0

	G           = 9.8
	DTM         = 0.0 // Assign the value of dTM here
	Cff_kWh     = 0.0 // Assign the value of cff_kWh here
	VAVCountMAX = 0   // Assign the value of VAV_Count_MAX here

	Fbmlist = "" // 壁材料リストファイル名 ref: wbmlist.md

	DEBUG   = false
	dayprn  = false
	DAYweek = [8]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", "Hol"}

	Ferr = os.Stderr // Assuming you want to write errors to standard error
	//DISPLAY_DELAY = 0 // Assign the value of DISPLAY_DELAY here
	SETprint = false //  SET(体感温度)を出力する場合は true
)

// 月の末日かどうかをチェックする
func IsEndDay(Mon, Day, Dayend, SimDayend int) bool {
	Nde := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	nday := SimDayend
	if nday > 365 {
		nday -= 365
	}
	if Day == Nde[Mon-1] || Dayend == nday {
		return true
	}

	return false
}

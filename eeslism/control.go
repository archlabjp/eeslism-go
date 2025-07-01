/*
control.go (Control Logic Data Structures for Building Energy Simulation)

このファイルは、建物のエネルギーシミュレーションにおける制御ロジックを定義するためのデータ構造を提供します。
これらの構造体は、空調システム、熱源設備、換気システムなどの運転制御をモデル化するために用いられます。

建築環境工学的な観点:
- **制御ロジックのモデル化**: 建物のエネルギー消費量や室内環境は、
  空調システムなどの運転制御によって大きく左右されます。
  このファイルで定義される`CONTL`、`CTLIF`、`CTLTYP`、`CTLST`、`VPTR`などの構造体は、
  - **条件分岐**: 特定の条件（例: 室温が設定値を超えた場合）に基づいて運転モードを切り替える。
  - **比較演算**: 温度、湿度、流量などの変数を比較する。
  - **変数へのポインター**: 制御対象となる変数（例: 室温、設定温度）にアクセスする。
  といった制御ロジックを柔軟に記述することを可能にします。
- **省エネルギー運転の実現**: 適切な制御ロジックをモデル化することで、
  - **デマンド制御**: 実際の熱負荷に応じて機器の運転を調整し、無駄なエネルギー消費を削減します。
  - **最適制御**: 快適性を維持しつつ、エネルギー消費を最小化する運転戦略を検討します。
  - **スケジュール運転**: 時間帯や曜日、季節に応じた運転モードの切り替えを可能にします。
- **快適性の維持**: 室内温度や湿度などの環境変数を目標値に維持するための制御をモデル化することで、
  居住者の快適性を確保できます。
- **システム統合**: これらの制御構造は、
  建物の様々な設備システム（空調、換気、熱源など）を統合的に制御し、
  建物全体のエネルギーマネジメント戦略を評価するための基盤となります。

このファイルは、建物のエネルギーシミュレーションにおいて、
複雑な制御ロジックを正確にモデル化し、
省エネルギー、快適性、および運用効率の向上を図るための重要な役割を果たします。
*/
package eeslism

type VPtrType rune

const (
	VAL_CTYPE    VPtrType = 'v' // 設定値 (float64)
	SW_CTYPE     VPtrType = 's' // 切替状態(ControlSWType)
	MAIN_CPTYPE  VPtrType = 'M' // For *MPATH ?
	LOCAL_CPTYPE VPtrType = 'L' // For *PLIST ?
)

type CONTL struct {
	Type      byte // 条件付き 'c', 無条件 ' '
	Lgv       int  // True:1、False:0
	Cif       *CTLIF
	AndCif    *CTLIF
	AndAndCif *CTLIF
	OrCif     *CTLIF
	Cst       *CTLST
}

// union V or S alternative
type CTLTYP struct {
	V *float64
	S *ControlSWType
}

// ControlIf
type CTLIF struct {
	// 演算対象の変数の型
	Type VPtrType // 'v' or 's'

	// 演算の種類
	Op byte // 比較演算の種類: 'g' for  ">",  'G' for ">=", 'l' for "<",  'L' for "<=", '=' for "==" and 'N' for "!="

	// 演算対象の変数
	Nlft int    // 演算対象の左変数の数 1 or 2. 1の場合はLft1のみ使用し、2の場合はLft1とLft2の使用する。
	Lft1 CTLTYP // 左辺その1
	Lft2 CTLTYP // 左辺その2
	Rgt  CTLTYP // 右辺
}

type CTLST struct {
	Type     VPtrType
	PathType VPtrType // 'M', 'L'
	Path     interface{}
	Lft      CTLTYP
	Rgt      CTLTYP
}

type VPTR struct {
	Type VPtrType
	Ptr  interface{}
}
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

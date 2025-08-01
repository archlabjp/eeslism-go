# SYSCMP（システム機器定義）

## 概要
SYSCMP（System Component）は、建物の設備システムを構成する個々の機器や要素の特性を定義するデータセットです。ボイラー、ポンプ、熱交換器、空調機、分岐・合流要素など、システム内のすべての構成要素の性能特性、設置条件、接続仕様を詳細に設定することで、実際の設備システムの動作を正確にシミュレーションできます。

## データ形式
```
SYSCMP
    <機器名|室名>
        -c <機器カタログ名>
        [ -L <配管長> ]
        [ -env <周辺温度> | -room <設置室名> ]
        [ -Tinit <初期温度> | -Tinit (<項目1> <温度1> ... <項目N> <温度N>) ]
        [ -exs <方位名> ]
        [ -roomheff <室名> <効率> ]
        [ Ac=<面積> ]
        [ type=<機器タイプ> ]
        [ -Nin <入口数> ]
        [ -Nout <出口数> ]
    ;
    ... 繰り返し
*
```

## パラメータ説明

| パラメータ | 単位 | 説明 | デフォルト値 | 参照 |
|:---|:---|:---|:---|:---|
| 機器名/室名 | - | 機器の識別名または室名 | - | 必須 |
| 機器カタログ名 | - | [EQPCAT](EQPCAT.md)で定義された機器仕様名 | - | 必須 |
| 配管長 | m | 機器に関連する配管の長さ | - | 任意 |
| 周辺温度 | ℃ | 機器周辺の環境温度（定数またはスケジュール名） | - | 任意 |
| 設置室名 | - | 機器が設置される室の名前 | - | 任意 |
| 初期温度 | ℃ | 蓄熱槽等の初期水温 | - | 任意 |
| 方位名 | - | [EXSRF](EXSRF.md)で定義された方位名 | - | 任意 |
| 効率 | - | 室に対する機器の効率 | - | 任意 |
| 面積 | m² | 機器の伝熱面積等 | - | 任意 |
| 機器タイプ | - | 特殊機器のタイプ指定（下記参照） | - | 任意 |
| 入口数 | - | 機器の入口数（分岐・合流要素用） | - | 任意 |
| 出口数 | - | 機器の出口数（分岐・合流要素用） | - | 任意 |

## 機器タイプ一覧

### 分岐・合流要素
| タイプ | 説明 | 用途 |
|:---|:---|:---|
| `B` | 水系統の分岐要素 | 配管の分岐点 |
| `BA` | 空気系統の分岐要素 | ダクトの分岐点 |
| `C` | 水系統の合流要素 | 配管の合流点 |
| `CA` | 空気系統の合流要素 | ダクトの合流点 |

### 制御機器
| タイプ | 説明 | 用途 |
|:---|:---|:---|
| `V` | 弁・ダンパー | 流量制御 |
| `VT` | 温調弁（水系統のみ） | 温度制御 |

### 空調機器
| タイプ | 説明 | 用途 |
|:---|:---|:---|
| `HCLD` | 仮想空調機コイル（直膨） | 直膨式冷暖房 |
| `HCLDW` | 仮想空調機コイル（冷温水） | 冷温水式空調 |
| `RMAC` | ルームエアコン | 個別空調 |

### 計測・境界条件
| タイプ | 説明 | 用途 |
|:---|:---|:---|
| `QMEAS` | カロリーメータ | 熱量計測 |
| `FLI` | 流入境界条件 | システム境界 |

## 使用例

### 基本的な機器定義
```
SYSCMP
    Boiler1
        -c GasBoiler_100kW
        -room MechanicalRoom
        -env 20.0 ;
    
    Pump1
        -c CirculationPump_5kW
        -room MechanicalRoom ;
    
    HeatExchanger1
        -c PlateHX_50kW
        -env 15.0 ;
*
```

### 蓄熱槽の定義
```
SYSCMP
    ThermalStorage
        -c WaterTank_10000L
        -room TankRoom
        -Tinit (C 60.0 D 40.0)
        -env 25.0 ;
*
```

### 空調システム機器
```
SYSCMP
    AirHandler1
        -c AHU_10000CMH
        -room MechanicalRoom
        -env 25.0 ;
    
    VAV_Zone1
        -c VAV_2000CMH
        -room Zone1 ;
    
    FCU_Room1
        -c FanCoil_1000CMH
        -room Room1
        -roomheff Room1 0.85 ;
*
```

### 分岐・合流要素
```
SYSCMP
    WaterBranch1
        type=B
        -Nin 1
        -Nout 3 ;
    
    WaterMerge1
        type=C
        -Nin 3
        -Nout 1 ;
    
    AirBranch1
        type=BA
        -Nin 1
        -Nout 2 ;
*
```

### 制御機器
```
SYSCMP
    TempControlValve1
        type=VT
        -c ThreeWayValve_DN50
        -room MechanicalRoom ;
    
    Damper1
        type=V
        -c MotorizedDamper_600x400
        -room DuctSpace ;
*
```

### 室への流入経路
```
SYSCMP
    LivingRoom
        -Nin 3 ;
    
    Bedroom
        -Nin 2 ;
*
```

## 初期温度設定の詳細

### 単一温度設定
```
-Tinit 50.0
```

### 複数部位の温度設定
```
-Tinit (C 60.0 D 40.0 M 50.0)
```
- `C`: 蓄熱側初期温度
- `D`: 放熱側初期温度  
- `M`: 中間部初期温度

## 設置条件の考慮

### 周辺温度の設定
- **定数**: `-env 25.0`
- **スケジュール**: `-env AmbientTemp`
- **設置室**: `-room MechanicalRoom`

### 方位の影響
外気に面する機器の場合：
```
-exs South
```

## 機器性能の定義

機器の詳細性能は[EQPCAT](EQPCAT.md)で定義し、SYSCMPでは以下を指定：
- 機器カタログの参照
- 設置条件
- 初期状態
- 接続仕様

## 設計上の注意事項

1. **機器選定**: 実際の設備仕様に基づいた適切な機器カタログの選択
2. **設置環境**: 機器の設置環境（温度、湿度）の正確な設定
3. **接続仕様**: 入出口数や配管仕様の整合性確保
4. **制御連携**: [CONTL](CONTL.md)データとの制御ロジック整合性

## 出力データ

各機器に関する以下の項目が出力されます：
- 機器の運転状態
- 入出口温度・湿度・流量
- エネルギー消費量
- 効率・COP
- 制御信号

## 関連データセット

- [SYSPTH](SYSPTH.md): システム経路の定義
- [EQPCAT](EQPCAT.md): 機器カタログの定義
- [CONTL](CONTL.md): 制御ロジックの定義
- [EXSRF](EXSRF.md): 外表面・方位の定義
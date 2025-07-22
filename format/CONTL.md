# CONTL（制御ロジック定義）

## 概要
CONTL（Control Logic）は、建物の設備システムにおける制御ロジックを定義するデータセットです。温度、湿度、流量などの変数に基づく条件分岐、比較演算、制御信号の生成を詳細に設定することで、空調システム、熱源設備、換気システムなどの複雑な運転制御をモデル化できます。省エネルギー運転と快適性維持の両立を図る制御戦略の評価に不可欠です。

## データ形式
```
CONTL
    <制御名>
        [ -c ]
        <条件式>
        [ AND <条件式> ]
        [ ANDAND <条件式> ]
        [ OR <条件式> ]
        <設定式>
    ;
    ... 繰り返し
*
```

## パラメータ説明

| パラメータ | 説明 | 設定値 | 参照 |
|:---|:---|:---|:---|
| 制御名 | 制御ロジックの識別名 | 任意の文字列 | 必須 |
| -c | 条件付き制御フラグ | フラグ | 任意 |
| 条件式 | 制御の判定条件 | 比較式 | 必須 |
| AND条件式 | 追加のAND条件 | 比較式 | 任意 |
| ANDAND条件式 | さらなるAND条件 | 比較式 | 任意 |
| OR条件式 | 代替のOR条件 | 比較式 | 任意 |
| 設定式 | 制御実行内容 | 代入式 | 必須 |

## 条件式の書式

### 比較演算子
| 演算子 | 意味 | 例 |
|:---|:---|:---|
| `>` | より大きい | `Tr > 25.0` |
| `>=` | 以上 | `Tr >= 24.0` |
| `<` | より小さい | `Tr < 20.0` |
| `<=` | 以下 | `Tr <= 26.0` |
| `==` | 等しい | `Mode == 1` |
| `!=` | 等しくない | `Status != OFF` |

### 変数の種類
- **温度変数**: `Tr`（室温）、`To`（外気温）、`Tw`（水温）等
- **湿度変数**: `RH`（相対湿度）、`X`（絶対湿度）等
- **流量変数**: `G`（流量）、`V`（風量）等
- **制御変数**: `Mode`（運転モード）、`Status`（状態）等

## 使用例

### 基本的な温度制御
```
CONTL
    HeatingControl
        -c
        Tr < 20.0
        Heater = ON ;
    
    CoolingControl
        -c
        Tr > 26.0
        Cooler = ON ;
    
    SystemOff
        -c
        Tr >= 20.0 AND Tr <= 26.0
        Heater = OFF
        Cooler = OFF ;
*
```

### 複合条件制御
```
CONTL
    EconomicControl
        -c
        To < 15.0 AND Tr < 22.0
        HeatingMode = 1
        CoolingMode = 0 ;
    
    ComfortControl
        -c
        Occupancy == 1 AND (Tr < 21.0 OR Tr > 25.0)
        ComfortMode = 1 ;
*
```

### 時刻・スケジュール制御
```
CONTL
    DayTimeOperation
        -c
        Time >= 8.0 AND Time <= 18.0
        SystemMode = AUTO ;
    
    NightTimeOperation
        -c
        Time < 8.0 OR Time > 18.0
        SystemMode = SETBACK ;
*
```

### 外気冷房制御
```
CONTL
    EconomizerControl
        -c
        To < Tr AND To > 10.0 AND Tr > 24.0
        OutdoorAirDamper = 100
        ReturnAirDamper = 0 ;
    
    NormalOperation
        -c
        To >= Tr OR To <= 10.0 OR Tr <= 24.0
        OutdoorAirDamper = 20
        ReturnAirDamper = 80 ;
*
```

### VAV制御
```
CONTL
    VAVControl_Zone1
        -c
        Tr_Zone1 > Tset_Zone1 + 1.0
        VAV_Zone1_Damper = 100 ;
    
    VAVControl_Zone1_Normal
        -c
        Tr_Zone1 <= Tset_Zone1 + 1.0 AND Tr_Zone1 >= Tset_Zone1 - 1.0
        VAV_Zone1_Damper = 50 ;
    
    VAVControl_Zone1_Min
        -c
        Tr_Zone1 < Tset_Zone1 - 1.0
        VAV_Zone1_Damper = 20 ;
*
```

### 熱源機器制御
```
CONTL
    BoilerStaging
        -c
        HeatLoad > 80.0
        Boiler1 = ON
        Boiler2 = ON ;
    
    BoilerStaging_Single
        -c
        HeatLoad > 40.0 AND HeatLoad <= 80.0
        Boiler1 = ON
        Boiler2 = OFF ;
    
    BoilerOff
        -c
        HeatLoad <= 40.0
        Boiler1 = OFF
        Boiler2 = OFF ;
*
```

## 制御ロジックの設計指針

### 温度制御
- **デッドバンド**: 頻繁な切り替えを防ぐため適切な不感帯を設定
- **オーバーシュート防止**: 制御応答の調整
- **設定温度**: 快適性と省エネのバランス

### 湿度制御
- **除湿制御**: 夏季の過湿防止
- **加湿制御**: 冬季の乾燥防止
- **結露防止**: 表面温度との関係考慮

### 流量制御
- **最小流量**: 機器保護のための下限設定
- **最大流量**: 機器容量の上限設定
- **変動制御**: 負荷に応じた流量調整

## 省エネルギー制御戦略

### 外気冷房
```
CONTL
    Freecooling
        -c
        To < Tr - 2.0 AND To > 5.0
        EconomizerMode = ON ;
*
```

### ナイトパージ
```
CONTL
    NightPurge
        -c
        Time >= 22.0 AND To < 20.0 AND Tr > 25.0
        NightVentilation = ON ;
*
```

### デマンド制御
```
CONTL
    DemandLimit
        -c
        PowerDemand > 80.0
        LoadShedding = ON ;
*
```

## 制御の優先順位

制御ロジックは記述順に評価されるため、優先度の高い制御を先に記述します：

1. **安全制御**: 機器保護、異常時対応
2. **快適性制御**: 室内環境維持
3. **省エネ制御**: エネルギー最適化
4. **デフォルト制御**: 通常運転

## 変数の参照方法

### システム変数
- 室温、外気温等の環境変数
- 機器状態、運転モード等の制御変数
- 時刻、曜日等の時間変数

### ユーザー定義変数
- スケジュールで定義された変数
- 計算により導出された変数

## デバッグとテスト

### 制御ロジックの検証
1. **条件の網羅性**: すべての運転状況をカバー
2. **論理の整合性**: 矛盾する制御の排除
3. **境界条件**: 設定値近傍での動作確認

### シミュレーション結果の確認
- 制御信号の時系列変化
- 制御による省エネ効果
- 快適性への影響

## 出力データ

制御に関する以下の項目が出力されます：
- 制御信号の状態
- 条件式の評価結果
- 制御による機器運転状況
- 省エネ効果

## 関連データセット

- [SYSPTH](SYSPTH.md): システム経路の定義
- [SYSCMP](SYSCMP.md): システム機器の定義
- [SCHNM](SCHNM.md): スケジュール名の定義
- [ROOM](ROOM.md): 室の定義
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

### 例1: オフィスビルの統合制御システム
**建物概要**: 10階建てオフィスビル、営業時間7:00-19:00
**制御目標**: 快適性確保、省エネルギー運転、設備保護

```
CONTL
    !  基本運転モード制御：営業時間に応じた運転切替
    BusinessHours_Control
        -c
        Time >= 7.0 AND Time <= 19.0 AND Weekday == 1
        SystemMode = NORMAL          !  通常運転モード
        SetTemp_Cooling = 26.0       !  冷房設定温度26℃
        SetTemp_Heating = 22.0       !  暖房設定温度22℃
        MinOutdoorAir = 30 ;         !  最小外気量30%
    
    !  夜間・休日セットバック制御：省エネ運転
    Setback_Control
        -c
        (Time < 7.0 OR Time > 19.0) OR Weekday == 0
        SystemMode = SETBACK         !  セットバック運転
        SetTemp_Cooling = 28.0       !  冷房設定温度緩和
        SetTemp_Heating = 20.0       !  暖房設定温度緩和
        MinOutdoorAir = 10 ;         !  最小外気量削減
    
    !  予冷・予熱制御：始業前の快適環境準備
    PreCooling_Control
        -c
        Time >= 6.0 AND Time < 7.0 AND Weekday == 1 AND To > 25.0
        SystemMode = PRECOOL         !  予冷運転
        SetTemp_Cooling = 24.0       !  予冷設定温度
        ChillerStart = ON ;          !  冷凍機早期起動
    
    PreHeating_Control
        -c
        Time >= 6.0 AND Time < 7.0 AND Weekday == 1 AND To < 15.0
        SystemMode = PREHEAT         !  予熱運転
        SetTemp_Heating = 24.0       !  予熱設定温度
        BoilerStart = ON ;           !  ボイラー早期起動
*
```

**制御ロジックの詳細説明**:
- **営業時間判定**: `Time >= 7.0 AND Time <= 19.0 AND Weekday == 1`で平日営業時間を判定
- **温度設定**: 営業時間は快適性重視（22-26℃）、夜間は省エネ重視（20-28℃）
- **外気量制御**: 営業時間は換気基準遵守、夜間は最小限に削減
- **予冷予熱**: 外気温度条件と時刻で判定し、始業時の快適性を確保

### 例2: 病院の精密空調制御
**建物概要**: 総合病院、手術室・ICU・一般病棟
**制御目標**: 高精度温湿度制御、24時間安定運転、感染対策

```
CONTL
    !  手術室精密制御：±0.5℃、±2%RH制御
    OR_PrecisionControl
        -c
        RoomType == OR AND Tr_OR < SetTemp_OR - 0.5
        HeatingValve_OR = 100        !  加熱弁全開
        CoolingValve_OR = 0          !  冷却弁閉
        Humidifier_OR = AUTO ;       !  加湿器自動制御
    
    OR_CoolingControl
        -c
        RoomType == OR AND Tr_OR > SetTemp_OR + 0.5
        HeatingValve_OR = 0          !  加熱弁閉
        CoolingValve_OR = 100        !  冷却弁全開
        Dehumidifier_OR = AUTO ;     !  除湿器自動制御
    
    !  ICU陽圧制御：感染防止のための圧力管理
    ICU_PressureControl
        -c
        RoomType == ICU AND Pressure_ICU < 5.0
        SupplyFan_ICU = 110          !  給気ファン風量増加
        ExhaustFan_ICU = 90 ;        !  排気ファン風量減少
    
    !  一般病棟省エネ制御：患者快適性と省エネの両立
    Ward_ComfortControl
        -c
        RoomType == WARD AND Occupancy_WARD > 0
        SetTemp_WARD = 24.0          !  在室時快適温度
        AirChange_WARD = 6 ;         !  換気回数6回/h
    
    Ward_EnergyControl
        -c
        RoomType == WARD AND Occupancy_WARD == 0
        SetTemp_WARD = 26.0          !  不在時省エネ温度
        AirChange_WARD = 2 ;         !  換気回数削減
*
```

**制御ロジックの詳細説明**:
- **部屋タイプ判定**: `RoomType == OR`で手術室、ICU、一般病棟を区別
- **精密制御**: 手術室は±0.5℃の高精度制御で医療機器の安定動作を確保
- **陽圧制御**: ICUは5Pa以上の陽圧維持で感染リスクを低減
- **在室連動**: 一般病棟は在室センサーと連動した省エネ制御

### 例3: 工場の生産連動制御
**建物概要**: 自動車部品工場、3交代24時間稼働
**制御目標**: 生産効率最大化、品質管理、エネルギー最適化

```
CONTL
    !  生産ライン連動制御：製造工程に応じた環境制御
    Production_LineA_Control
        -c
        ProductionLine_A == RUNNING AND ProcessTemp_A > 80.0
        ExhaustFan_A = 100           !  排気ファン全開
        CoolingCoil_A = 100          !  冷却コイル全開
        MakeupAir_A = 80 ;           !  補給空気80%
    
    Production_LineA_Standby
        -c
        ProductionLine_A == STANDBY
        ExhaustFan_A = 20            !  排気ファン最小
        CoolingCoil_A = 0            !  冷却停止
        MakeupAir_A = 20 ;           !  補給空気最小
    
    !  品質管理制御：製品品質に影響する環境要因の制御
    QualityControl_PaintBooth
        -c
        ProcessType == PAINT AND RH > 60.0
        Dehumidifier_Paint = ON      !  除湿機運転
        ExhaustRate_Paint = 120 ;    !  排気量増加
    
    QualityControl_Assembly
        -c
        ProcessType == ASSEMBLY AND Dust > 0.1
        AirFilter_Assembly = HIGH    !  高性能フィルター
        CleanRoom_Assembly = ON ;    !  クリーンルーム運転
    
    !  エネルギー最適化制御：電力デマンド制御
    DemandControl_Peak
        -c
        PowerDemand > 1800.0 AND Time >= 13.0 AND Time <= 16.0
        NonEssential_Equipment = OFF !  非必須設備停止
        Lighting_Dimming = 80        !  照明調光80%
        Compressor_Staging = 2 ;     !  コンプレッサー台数制限
    
    !  夜間メンテナンス制御：保守作業時の安全確保
    Maintenance_Control
        -c
        MaintenanceMode == ON AND Time >= 1.0 AND Time <= 5.0
        EmergencyLighting = ON       !  非常照明点灯
        VentilationRate = 150        !  換気量増加
        SafetySystem = ACTIVE ;      !  安全システム作動
*
```

**制御ロジックの詳細説明**:
- **生産連動**: `ProductionLine_A == RUNNING`で生産状態を判定し、環境制御を最適化
- **品質管理**: 塗装工程では湿度60%以下、組立工程では清浄度管理
- **デマンド制御**: 電力需要1800kW超過時に非必須設備を制限
- **安全制御**: 夜間メンテナンス時は安全性を最優先した制御

### 例4: データセンターの高効率制御
**建物概要**: クラウドサービス用データセンター、年間稼働率99.9%
**制御目標**: IT機器冷却、高効率運転、障害時継続運転

```
CONTL
    !  外気冷房制御：外気条件に応じた高効率運転
    Economizer_FullMode
        -c
        To < 18.0 AND RH_outside < 70.0 AND IT_Load > 50.0
        OutdoorAir_Damper = 100      !  外気ダンパー全開
        MechanicalCooling = OFF      !  機械冷房停止
        FreeCoiling_Mode = FULL ;    !  完全外気冷房
    
    Economizer_PartialMode
        -c
        To >= 18.0 AND To <= 25.0 AND RH_outside < 80.0
        OutdoorAir_Damper = 60       !  外気ダンパー60%
        MechanicalCooling = PARTIAL  !  機械冷房部分運転
        FreeCoiling_Mode = PARTIAL ; !  部分外気冷房
    
    !  IT負荷連動制御：サーバー負荷に応じた冷却制御
    IT_HighLoad_Control
        -c
        CPU_Utilization > 80.0 AND IT_Power > 500.0
        CoolingCapacity = 120        !  冷却能力120%
        AirFlow_Rate = 110           !  風量110%
        Chiller_Staging = 3 ;        !  冷凍機3台運転
    
    IT_LowLoad_Control
        -c
        CPU_Utilization < 30.0 AND IT_Power < 200.0
        CoolingCapacity = 80         !  冷却能力80%
        AirFlow_Rate = 90            !  風量90%
        Chiller_Staging = 1 ;        !  冷凍機1台運転
    
    !  障害時制御：単一故障時の継続運転
    Chiller_Failure_Control
        -c
        Chiller1_Status == FAULT AND IT_Load > 70.0
        Chiller2_Capacity = 120      !  予備機過負荷運転
        Emergency_Cooling = ON       !  非常用冷却起動
        IT_Load_Shedding = 10 ;      !  IT負荷10%削減
    
    !  電源障害時制御：UPS・発電機連動
    Power_Failure_Control
        -c
        MainPower_Status == FAULT
        UPS_Mode = ACTIVE            !  UPS運転開始
        Generator_Start = ON         !  発電機起動
        NonCritical_Load = OFF ;     !  非重要負荷遮断
*
```

**制御ロジックの詳細説明**:
- **外気冷房**: 外気温18℃以下で完全外気冷房、25℃以下で部分外気冷房
- **負荷連動**: CPU使用率80%超過時は冷却能力120%で確実な冷却を実現
- **障害対応**: 冷凍機故障時は予備機過負荷運転とIT負荷削減で継続運転
- **電源障害**: UPS・発電機の自動起動で無停電運転を維持

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
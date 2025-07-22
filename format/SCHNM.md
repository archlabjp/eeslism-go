# SCHNM（スケジュール名定義）

## 概要
SCHNM（Schedule Name）は、日スケジュールの季節・曜日による変更を定義するデータセットです。基本となる日スケジュール（SCHTB）を季節や曜日に応じて切り替えることで、年間を通じた詳細なスケジュール制御を実現できます。空調運転、照明制御、機器稼働など、建物の運用パターンを現実的にモデル化するために不可欠な機能です。

## データ形式
```
SCHNM
    -v <スケジュール名> <設定値スケジュール組合せ> <設定値スケジュール組合せ> ... ;
    -s <スケジュール名> <切替スケジュール組合せ> <切替スケジュール組合せ> ... ;
    ... 繰り返し
*
```

## パラメータ説明

| パラメータ | 説明 | 設定値 | 参照 |
|:---|:---|:---|:---|
| スケジュール名 | 他のデータセットから参照される識別名 | 任意の文字列 | 必須 |
| 設定値スケジュール組合せ | 季節・曜日・日スケジュールの組合せ | `<季節>:<曜日>:<日スケジュール>` | -v時必須 |
| 切替スケジュール組合せ | 季節・曜日・切替スケジュールの組合せ | `<季節>:<曜日>:<切替スケジュール>` | -s時必須 |

## スケジュール組合せの書式

### 基本書式
```
<季節名>:<曜日名>:<日スケジュール名>
```

### 省略形式
- `*:<曜日名>:<日スケジュール名>` - 全季節共通
- `<季節名>:*:<日スケジュール名>` - 全曜日共通  
- `*:*:<日スケジュール名>` - 全期間共通

## 使用例

### 例1: 大規模オフィスビルの年間運用スケジュール
**建物概要**: 20階建てオフィスビル、従業員1,500名、営業時間8:00-18:00
**運用方針**: 快適性確保、省エネルギー、働き方改革対応

```
SCHNM
    !  室温設定スケジュール：季節・曜日別の最適温度管理
    -v OfficeTemperature_Main
        Winter:Weekday:WinterOffice_Weekday     !  冬季平日：20-24℃設定
        Winter:Weekend:WinterOffice_Weekend     !  冬季休日：18-22℃省エネ設定
        Summer:Weekday:SummerOffice_Weekday     !  夏季平日：24-28℃設定
        Summer:Weekend:SummerOffice_Weekend     !  夏季休日：26-30℃省エネ設定
        Inter:Weekday:InterOffice_Weekday       !  中間期平日：22-26℃自然換気併用
        Inter:Weekend:InterOffice_Weekend ;     !  中間期休日：20-28℃大幅省エネ
    
    !  照明制御スケジュール：昼光利用と省エネの両立
    -v OfficeLighting_Control
        Winter:Weekday:WinterLight_Weekday      !  冬季平日：7:00-19:00点灯
        Winter:Weekend:WinterLight_Weekend      !  冬季休日：必要時のみ点灯
        Summer:Weekday:SummerLight_Weekday      !  夏季平日：8:00-18:00点灯（昼光利用）
        Summer:Weekend:SummerLight_Weekend      !  夏季休日：セキュリティ照明のみ
        Inter:*:InterLight_Standard ;           !  中間期：季節共通スケジュール
    
    !  空調運転モード：効率的な空調制御
    -s HVAC_OperationMode
        Winter:Weekday:WinterHVAC_Normal        !  冬季平日：暖房メイン運転
        Winter:Weekend:WinterHVAC_Setback       !  冬季休日：セットバック運転
        Summer:Weekday:SummerHVAC_Normal        !  夏季平日：冷房メイン運転
        Summer:Weekend:SummerHVAC_Setback       !  夏季休日：セットバック運転
        Inter:Weekday:InterHVAC_Economizer      !  中間期平日：外気冷房活用
        Inter:Weekend:InterHVAC_Natural ;       !  中間期休日：自然換気メイン
*
```

**スケジュール設計の詳細説明**:
- **温度設定**: 冬季は暖房効率重視で20-24℃、夏季は冷房効率重視で24-28℃
- **休日省エネ**: 週末は設定温度を2-4℃緩和し、大幅な省エネを実現
- **中間期活用**: 外気冷房・自然換気で機械空調の負荷を最小化
- **照明制御**: 昼光利用センサーと連動し、夏季は点灯時間を1時間短縮

### 例2: 病院の24時間運用スケジュール
**建物概要**: 総合病院500床、手術室・ICU・一般病棟・外来
**運用方針**: 患者安全最優先、医療機器安定動作、感染対策

```
SCHNM
    !  病棟温度管理：患者快適性と医療安全の両立
    -v Hospital_PatientArea_Temp
        Winter:*:WinterHospital_Patient         !  冬季：24℃一定（患者快適性）
        Summer:*:SummerHospital_Patient         !  夏季：26℃一定（感染リスク考慮）
        Inter:*:InterHospital_Patient ;         !  中間期：25℃一定（安定性重視）
    
    !  手術室精密制御：医療機器の安定動作確保
    -v OR_PrecisionControl
        Winter:*:WinterOR_Precision             !  冬季：22±0.5℃、50±2%RH
        Summer:*:SummerOR_Precision             !  夏季：24±0.5℃、45±2%RH
        Inter:*:InterOR_Precision ;             !  中間期：23±0.5℃、48±2%RH
    
    !  外来エリア：診療時間連動制御
    -v Outpatient_AreaControl
        Winter:Weekday:WinterOutpatient_Open    !  冬季平日：8:00-17:00診療時間
        Winter:Weekend:WinterOutpatient_Closed  !  冬季休日：最小限運転
        Summer:Weekday:SummerOutpatient_Open    !  夏季平日：診療時間中快適制御
        Summer:Weekend:SummerOutpatient_Closed  !  夏季休日：省エネ運転
        Inter:Weekday:InterOutpatient_Open      !  中間期平日：自然換気併用
        Inter:Weekend:InterOutpatient_Closed ;  !  中間期休日：停止可能
    
    !  感染対策換気：病室種別別制御
    -s InfectionControl_Ventilation
        *:*:IC_Standard_Ventilation ;           !  全期間：感染対策基準換気
*
```

**スケジュール設計の詳細説明**:
- **患者エリア**: 年間を通じて安定した温度環境で患者の回復を支援
- **手術室**: ±0.5℃の高精度制御で医療機器の誤動作を防止
- **外来エリア**: 診療時間のみ快適制御、休日は大幅省エネ
- **感染対策**: 全期間で基準以上の換気量を確保

### 例3: 製造工場の生産連動スケジュール
**建物概要**: 自動車部品工場、3交代24時間稼働、品質管理重要
**運用方針**: 生産効率最大化、品質安定、エネルギー最適化

```
SCHNM
    !  生産エリア環境制御：製品品質に直結する環境管理
    -v Production_EnvironmentControl
        Winter:Weekday:WinterProd_3Shift        !  冬季平日：3交代対応制御
        Winter:Weekend:WinterProd_2Shift        !  冬季休日：2交代運転
        Summer:Weekday:SummerProd_3Shift        !  夏季平日：冷却強化運転
        Summer:Weekend:SummerProd_2Shift        !  夏季休日：標準冷却運転
        Inter:Weekday:InterProd_3Shift          !  中間期平日：外気利用運転
        Inter:Weekend:InterProd_2Shift ;        !  中間期休日：省エネ運転
    
    !  品質管理エリア：精密環境制御
    -v QualityControl_Environment
        Winter:*:WinterQC_Precision             !  冬季：20±1℃、45±3%RH
        Summer:*:SummerQC_Precision             !  夏季：23±1℃、50±3%RH
        Inter:*:InterQC_Precision ;             !  中間期：22±1℃、48±3%RH
    
    !  塗装ブース：溶剤管理と品質確保
    -v PaintBooth_Control
        Winter:*:WinterPaint_Optimal            !  冬季：25℃、40%RH（乾燥促進）
        Summer:*:SummerPaint_Optimal            !  夏季：22℃、45%RH（品質安定）
        Inter:*:InterPaint_Optimal ;            !  中間期：24℃、42%RH（標準）
    
    !  電力デマンド制御：ピークカット対応
    -s PowerDemand_Management
        Winter:Weekday:WinterDemand_Control     !  冬季平日：暖房ピーク制御
        Summer:Weekday:SummerDemand_Control     !  夏季平日：冷房ピーク制御
        *:Weekend:WeekendDemand_Relaxed ;       !  休日：デマンド制御緩和
*
```

**スケジュール設計の詳細説明**:
- **生産連動**: 3交代・2交代に応じた環境制御で生産効率を最大化
- **品質管理**: 製品検査に必要な±1℃、±3%RHの精密制御
- **塗装ブース**: 溶剤の乾燥特性を考慮した季節別最適環境
- **デマンド制御**: 電力ピーク時間帯の負荷制限で電力コスト削減

### 例4: 学校の教育環境最適化スケジュール
**建物概要**: 小中一貫校、児童生徒800名、地域開放施設併設
**運用方針**: 学習環境最適化、健康配慮、地域利用対応

```
SCHNM
    !  教室環境制御：学習効率向上のための環境管理
    -v Classroom_LearningEnvironment
        Winter:Weekday:WinterClass_Learning     !  冬季授業日：20-22℃（集中力向上）
        Winter:Weekend:WinterClass_Community    !  冬季休日：18-24℃（地域開放）
        Summer:Weekday:SummerClass_Learning     !  夏季授業日：26-28℃（熱中症予防）
        Summer:Weekend:SummerClass_Community    !  夏季休日：24-30℃（地域開放）
        Inter:Weekday:InterClass_Natural        !  中間期授業日：自然換気活用
        Inter:Weekend:InterClass_Community ;    !  中間期休日：地域開放対応
    
    !  体育館：運動強度に応じた環境制御
    -v Gymnasium_ActivityControl
        Winter:Weekday:WinterGym_PE             !  冬季授業：15-18℃（運動適温）
        Winter:Weekend:WinterGym_Event          !  冬季イベント：18-20℃（観客考慮）
        Summer:Weekday:SummerGym_PE             !  夏季授業：28-30℃（熱中症対策）
        Summer:Weekend:SummerGym_Event          !  夏季イベント：26-28℃（快適性）
        Inter:*:InterGym_Natural ;              !  中間期：自然換気メイン
    
    !  給食室：食品衛生管理
    -v Kitchen_HygieneControl
        Winter:*:WinterKitchen_Sanitary         !  冬季：18℃以下（食品安全）
        Summer:*:SummerKitchen_Sanitary         !  夏季：25℃以下（食中毒防止）
        Inter:*:InterKitchen_Sanitary ;         !  中間期：22℃以下（衛生管理）
    
    !  夏季休暇対応：長期休暇中の省エネ運転
    -s SummerVacation_Operation
        Summer:*:SummerVacation_Minimal ;       !  夏季休暇：最小限運転
*
```

**スケジュール設計の詳細説明**:
- **学習環境**: 冬季20-22℃、夏季26-28℃で児童生徒の集中力と健康を両立
- **体育館**: 運動強度を考慮し、授業時は低め、イベント時は観客快適性重視
- **給食室**: 食品衛生法に基づく温度管理で食中毒リスクを最小化
- **長期休暇**: 夏季休暇中は最小限運転で大幅な省エネを実現

## スケジュール設計の考慮事項

### 季節区分の設定
- **冬季**: 暖房期間（11月～4月）
- **夏季**: 冷房期間（6月～9月）  
- **中間期**: 空調軽負荷期間（4月～5月、10月～11月）

### 曜日区分の設定
- **平日**: 月曜～金曜（通常運転）
- **週末**: 土曜・日曜・祝日（軽負荷運転）
- **特別日**: 年末年始、夏季休暇等

### 建物用途別の特徴

#### オフィスビル
- 平日・週末の運転パターンが大きく異なる
- 季節による設定温度の調整
- 照明・OA機器の使用パターン

#### 住宅
- 在宅時間帯の考慮
- 季節による生活パターンの変化
- 省エネ運転の重視

#### 商業施設
- 営業時間に応じた運転
- 季節商品による負荷変動
- 客数変動への対応

## 優先順位と適用ルール

### 適用優先順位
1. **具体的指定**: `Winter:Weekday:Schedule1`
2. **曜日省略**: `Winter:*:Schedule2`
3. **季節省略**: `*:Weekday:Schedule3`
4. **全省略**: `*:*:Schedule4`

### 未定義時の動作
指定された組合せが見つからない場合は、より一般的な設定が適用されます。

## エネルギー管理への活用

### ピークカット制御
```
SCHNM
    -s PeakCutOperation
        Summer:Weekday:SummerPeakCut
        *:*:NormalOperation ;
*
```

### 夜間電力活用
```
SCHNM
    -v NightTimeOperation
        Winter:*:WinterNightHeating
        Summer:*:SummerNightCooling ;
*
```

### デマンド制御
```
SCHNM
    -s DemandControl
        Summer:Weekday:SummerDemandLimit
        Winter:Weekday:WinterDemandLimit
        *:Weekend:WeekendDemand ;
*
```

## 設計上の注意事項

1. **整合性確保**: 季節・曜日定義（SCHTB）との整合性
2. **網羅性**: すべての期間・曜日の組合せをカバー
3. **現実性**: 実際の建物運用に即したパターン設定
4. **省エネ性**: エネルギー効率を考慮したスケジュール

## デバッグとテスト

### スケジュール確認方法
1. **年間カレンダー**: 各日のスケジュール適用状況
2. **季節移行**: 季節変更時の動作確認
3. **祝日処理**: 祝日の適切な処理

### 出力データでの確認
- スケジュール値の時系列変化
- 季節・曜日による切替状況
- エネルギー消費パターン

## 関連データセット

- [SCHTB](SCHTB.md): 日スケジュールテーブルの定義
- [WEEK](WEEK.md): 週間カレンダーの定義
- [RESI](RESI.md): 居住者スケジュールでの参照
- [AAPL](AAPL.md): 照明・機器スケジュールでの参照
- [VENT](VENT.md): 換気スケジュールでの参照


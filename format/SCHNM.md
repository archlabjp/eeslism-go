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

### 基本的な設定値スケジュール
```
SCHNM
    -v RoomTempSchedule
        Winter:Weekday:WinterWeekdayTemp
        Winter:Weekend:WinterWeekendTemp
        Summer:Weekday:SummerWeekdayTemp
        Summer:Weekend:SummerWeekendTemp
        Inter:*:InterSeasonTemp ;
    
    -v LightingSchedule
        *:Weekday:WeekdayLighting
        *:Weekend:WeekendLighting ;
*
```

### 切替スケジュール
```
SCHNM
    -s HVACOperation
        Winter:Weekday:WinterWeekdayHVAC
        Winter:Weekend:WinterWeekendHVAC
        Summer:Weekday:SummerWeekdayHVAC
        Summer:Weekend:SummerWeekendHVAC
        Inter:*:InterSeasonHVAC ;
    
    -s VentilationMode
        *:Weekday:WeekdayVent
        *:Weekend:WeekendVent ;
*
```

### オフィスビルの包括的スケジュール
```
SCHNM
    -v OfficeTemperature
        Winter:Weekday:WinterOfficeTemp_WD
        Winter:Weekend:WinterOfficeTemp_WE
        Summer:Weekday:SummerOfficeTemp_WD
        Summer:Weekend:SummerOfficeTemp_WE
        Inter:Weekday:InterOfficeTemp_WD
        Inter:Weekend:InterOfficeTemp_WE ;
    
    -v OfficeLighting
        *:Weekday:OfficeLighting_WD
        *:Weekend:OfficeLighting_WE ;
    
    -s OfficeHVAC
        Winter:Weekday:WinterHVAC_WD
        Winter:Weekend:WinterHVAC_WE
        Summer:Weekday:SummerHVAC_WD
        Summer:Weekend:SummerHVAC_WE
        Inter:*:InterHVAC ;
*
```

### 住宅の詳細スケジュール
```
SCHNM
    -v ResidentialHeating
        Winter:Weekday:WinterHeating_WD
        Winter:Weekend:WinterHeating_WE
        Inter:*:InterHeating ;
    
    -v ResidentialLighting
        Winter:*:WinterLighting
        Summer:*:SummerLighting
        Inter:*:InterLighting ;
    
    -s VentilationControl
        *:Weekday:WeekdayVent
        *:Weekend:WeekendVent ;
*
```

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


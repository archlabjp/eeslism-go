# RESI（居住者スケジュール）

## 概要
RESI（Resident Schedule）は、建物内の居住者の在室状況、活動レベル、熱的快適性に関する設定を定義するデータセットです。人体発熱量、代謝率、着衣量、室内風速などのパラメータを時間スケジュールと組み合わせて設定することで、居住者の影響を考慮した詳細な熱環境シミュレーションが可能になります。

## データ形式
```
RESI
    <室名>
        [ H=(<基準人数>,<在室スケジュール名>,<作業強度スケジュール名>) ]
        [ comfrt=(<代謝率スケジュール名>,<着衣量スケジュール名>,<風速スケジュール名>) ]
    ;
    ... 繰り返し
*
```

## パラメータ説明

| パラメータ | 単位 | 説明 | 参照 |
|:---|:---|:---|:---|
| 室名 | - | 対象となる室の名称。[ROOMデータセット](ROOM.md)で定義された室名を指定 | 必須 |
| 基準人数 | 人 | 人体発熱計算の基準となる人数 | H設定時必須 |
| 在室スケジュール名 | - | 基準人数に対する在室率のスケジュール名（0.0-1.0の比率） | H設定時必須 |
| 作業強度スケジュール名 | - | 人体発熱の作業強度を表すスケジュール名 | H設定時必須 |
| 代謝率スケジュール名 | met | 熱的快適性評価用の代謝率スケジュール名 | comfrt設定時必須 |
| 着衣量スケジュール名 | clo | 熱的快適性評価用の着衣量スケジュール名 | comfrt設定時必須 |
| 風速スケジュール名 | m/s | 熱的快適性評価用の室内風速スケジュール名 | comfrt設定時必須 |

※ すべてのスケジュール名は[SCHNMデータセット](SCHNM.md)で事前に定義する必要があります。

## 使用例

### 基本的な居住者設定
```
RESI
    LivingRoom
        H=(4,OccupancySchedule,ActivitySchedule) ;
    
    Bedroom
        H=(2,BedroomOccupancy,SleepActivity) ;
*
```

### 快適性評価を含む設定
```
RESI
    Office
        H=(10,OfficeOccupancy,OfficeActivity)
        comfrt=(MetabolicRate,ClothingLevel,AirVelocity) ;
    
    MeetingRoom
        H=(8,MeetingSchedule,MeetingActivity)
        comfrt=(MeetingMetabolic,BusinessClothing,StandardAirflow) ;
*
```

### 住宅の詳細設定例
```
RESI
    LivingDiningKitchen
        H=(4,FamilySchedule,DailyActivity)
        comfrt=(HomeMetabolic,SeasonalClothing,NaturalVentilation) ;
    
    MasterBedroom
        H=(2,BedroomSchedule,RestActivity)
        comfrt=(SleepMetabolic,NightClothing,QuietAirflow) ;
    
    ChildRoom
        H=(1,ChildSchedule,StudyActivity)
        comfrt=(ChildMetabolic,LightClothing,StandardAirflow) ;
*
```

## 人体発熱の計算

### 発熱量の算出
人体からの発熱量は以下の式で計算されます：
```
Q_human = 基準人数 × 在室率 × 作業強度係数 × 基準発熱量
```

### 作業強度の目安
- **安静時**: 0.8-1.0（睡眠、読書）
- **軽作業**: 1.2-1.6（事務作業、軽い家事）
- **中作業**: 1.6-2.0（掃除、料理）
- **重作業**: 2.0-3.0（運動、重労働）

## 熱的快適性評価

### PMV/PPD計算
comfrtパラメータが設定されている場合、以下の快適性指標が計算されます：
- **PMV**（Predicted Mean Vote）: 温熱感覚の予測平均申告
- **PPD**（Predicted Percentage of Dissatisfied）: 不満足者率

### 代謝率（Metabolic Rate）
- **安静**: 0.8-1.0 met
- **軽作業**: 1.2-1.6 met
- **中作業**: 1.6-2.4 met
- **重作業**: 2.4-4.0 met

### 着衣量（Clothing Insulation）
- **夏季軽装**: 0.3-0.5 clo
- **中間期**: 0.5-0.7 clo
- **冬季**: 0.8-1.2 clo
- **厚着**: 1.2-2.0 clo

## スケジュール設定の注意事項

1. **在室率**: 0.0（不在）から1.0（定員）の範囲で設定
2. **時間精度**: 計算時間間隔に応じた適切な時間分解能を設定
3. **季節変動**: 着衣量は季節に応じて調整が必要
4. **用途別設定**: 建物用途に応じた適切な活動レベルを設定

## 出力データ

居住者に関する以下の項目が出力されます：
- 在室人数
- 人体発熱量（顕熱・潜熱）
- PMV値
- PPD値
- 代謝率
- 着衣量
- 室内風速

## 関連データセット

- [ROOM](ROOM.md): 室の定義
- [SCHNM](SCHNM.md): スケジュール名の定義
- [SCHTB](SCHTB.md): スケジュールテーブルの定義
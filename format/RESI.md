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

### 例1: 高層オフィスビルの在室者管理
**建物概要**: 30階建てオフィスタワー、従業員2,500名、フレックスタイム制
**管理目標**: 在室率の正確な把握、快適性評価、省エネルギー運転

```
RESI
    !  一般オフィスフロア：標準的な事務作業エリア
    GeneralOffice_Floor10to25
        H=(80,FlexTime_Occupancy,OfficeWork_Activity)    !  基準人数80名/フロア
        comfrt=(OfficeWork_Metabolic,Business_Clothing,AC_Airflow) ;
        !  在室率：フレックス制で8:00-10:00に出社ピーク、17:00-19:00に退社ピーク
        !  活動強度：PC作業メインで1.2met相当
        !  快適性：ビジネス服装0.7clo、空調風速0.15m/s
    
    !  役員フロア：高品質な執務環境
    ExecutiveFloor_29to30
        H=(20,Executive_Occupancy,ExecutiveWork_Activity)
        comfrt=(Executive_Metabolic,Executive_Clothing,Premium_Airflow) ;
        !  在室率：役員の不規則な在室パターン
        !  活動強度：会議・判断業務で1.3met
        !  快適性：高級スーツ0.8clo、個別空調0.1m/s
    
    !  会議室：利用予約システム連動
    ConferenceRoom_Large
        H=(20,Meeting_Reservation,Meeting_Activity)
        comfrt=(Meeting_Metabolic,Meeting_Clothing,Conference_Airflow) ;
        !  在室率：予約システムと連動した確実な在室管理
        !  活動強度：プレゼン・議論で1.4met
        !  快適性：多様な服装0.6-0.8clo、会議用空調0.2m/s
    
    !  カフェテリア：食事時間帯の集中利用
    Cafeteria_Restaurant
        H=(200,Lunch_Rush,Dining_Activity)
        comfrt=(Dining_Metabolic,Casual_Clothing,Dining_Airflow) ;
        !  在室率：12:00-13:00に利用集中、その他は低利用
        !  活動強度：食事で1.0met
        !  快適性：軽装0.5clo、食事環境用空調0.25m/s
*
```

**設計のポイント**:
- **在室率管理**: フレックス制に対応した時間帯別在室パターン
- **活動レベル**: 業務内容に応じた代謝率設定（1.0-1.4met）
- **服装季節変動**: 夏季0.5clo、冬季0.8cloの季節別着衣量
- **快適性評価**: PMV±0.5以内を目標とした環境制御

### 例2: 総合病院の患者・職員管理
**建物概要**: 500床総合病院、医師・看護師800名、患者・家族1,000名/日
**管理目標**: 患者快適性確保、医療従事者の作業環境、感染対策

```
RESI
    !  一般病棟：患者の療養環境
    GeneralWard_4bed
        H=(4,Patient_Occupancy,Patient_Activity)
        comfrt=(Patient_Metabolic,Patient_Clothing,Hospital_Airflow) ;
        !  在室率：入院患者は24時間在室、面会時間で変動
        !  活動強度：安静時0.8met（病床での療養）
        !  快適性：病衣0.3clo、療養環境0.1m/s
    
    !  手術室：医療チームの高度医療環境
    OperatingRoom_Main
        H=(8,Surgery_Schedule,Surgery_Activity)
        comfrt=(Surgery_Metabolic,Surgical_Clothing,OR_Airflow) ;
        !  在室率：手術スケジュールに完全連動
        !  活動強度：手術執刀で2.0met（高い集中力と体力）
        !  快適性：手術着0.4clo、清浄空調0.05m/s
    
    !  ICU：重篤患者の集中治療
    ICU_Bed
        H=(1,ICU_Patient,ICU_Activity)
        comfrt=(ICU_Metabolic,ICU_Clothing,ICU_Airflow) ;
        !  在室率：患者24時間、医療スタッフ交代制
        !  活動強度：重篤患者0.7met、医療スタッフ1.5met
        !  快適性：最小限着衣0.2clo、精密空調0.08m/s
    
    !  外来待合室：患者・家族の一時滞在
    Outpatient_Waiting
        H=(50,Outpatient_Flow,Waiting_Activity)
        comfrt=(Waiting_Metabolic,Outpatient_Clothing,Waiting_Airflow) ;
        !  在室率：診療時間に応じた患者流動
        !  活動強度：待機・移動で1.1met
        !  快適性：一般服装0.6clo、快適空調0.15m/s
*
```

**設計のポイント**:
- **患者中心**: 療養に最適な環境条件（PMV-0.2～+0.2）
- **医療従事者**: 高い活動強度に対応した環境制御
- **感染対策**: 適切な気流制御で院内感染リスク低減
- **24時間対応**: 夜勤体制を考慮した連続的な環境管理

### 例3: 製造工場の作業者安全管理
**建物概要**: 自動車部品工場、作業者500名、3交代24時間稼働
**管理目標**: 作業安全確保、生産性向上、熱中症予防

```
RESI
    !  組立ライン：精密作業エリア
    AssemblyLine_A
        H=(25,ThreeShift_Assembly,Assembly_Activity)
        comfrt=(Assembly_Metabolic,Work_Clothing,Assembly_Airflow) ;
        !  在室率：3交代で24時間稼働、休憩時間で変動
        !  活動強度：組立作業で1.8met（中程度の肉体労働）
        !  快適性：作業服0.6clo、作業環境0.3m/s
    
    !  溶接エリア：高温作業環境
    WeldingArea_B
        H=(15,Welding_Schedule,Welding_Activity)
        comfrt=(Welding_Metabolic,Protective_Clothing,Welding_Airflow) ;
        !  在室率：溶接作業スケジュールに連動
        !  活動強度：溶接作業で2.5met（重労働）
        !  快適性：防護服1.2clo、強制換気0.5m/s
    
    !  品質管理室：検査作業環境
    QualityControl_Lab
        H=(8,QC_Schedule,Inspection_Activity)
        comfrt=(Inspection_Metabolic,Lab_Clothing,QC_Airflow) ;
        !  在室率：品質管理スケジュールに応じた在室
        !  活動強度：検査作業で1.3met（軽作業）
        !  快適性：作業着0.5clo、精密環境0.1m/s
    
    !  休憩室：作業者の休息空間
    RestArea_Common
        H=(30,Break_Schedule,Rest_Activity)
        comfrt=(Rest_Metabolic,Casual_Clothing,Rest_Airflow) ;
        !  在室率：交代時間・休憩時間に集中利用
        !  活動強度：休息で0.9met
        !  快適性：私服0.7clo、快適空調0.2m/s
*
```

**設計のポイント**:
- **安全第一**: 高温・有害環境での作業者保護（WBGT管理）
- **生産性**: 適切な環境で作業効率を最大化
- **疲労軽減**: 休憩エリアでの十分な回復環境
- **3交代対応**: 24時間を通じた一定品質の作業環境

### 例4: 学校の教育環境最適化
**建物概要**: 小中一貫校、児童生徒800名、教職員80名
**管理目標**: 学習効率向上、健康管理、省エネルギー

```
RESI
    !  普通教室：基本的な学習環境
    Classroom_Elementary
        H=(30,School_Schedule,Learning_Activity)
        comfrt=(Student_Metabolic,School_Clothing,Class_Airflow) ;
        !  在室率：授業時間割に完全連動、夏休み等長期休暇考慮
        !  活動強度：授業受講で1.2met（軽い知的活動）
        !  快適性：制服0.6clo、学習環境0.2m/s
    
    !  体育館：運動時の環境管理
    Gymnasium_Sports
        H=(40,PE_Schedule,Sports_Activity)
        comfrt=(Sports_Metabolic,PE_Clothing,Sports_Airflow) ;
        !  在室率：体育授業・部活動・行事に応じた利用
        !  活動強度：体育授業で3.0met（激しい運動）
        !  快適性：体操服0.3clo、運動環境0.4m/s
    
    !  図書室：静寂な学習環境
    Library_Study
        H=(50,Library_Hours,Study_Activity)
        comfrt=(Study_Metabolic,Study_Clothing,Library_Airflow) ;
        !  在室率：開館時間内の自由利用、試験期間で増加
        !  活動強度：読書・自習で1.0met（静的活動）
        !  快適性：制服0.6clo、静寂環境0.1m/s
    
    !  給食室：調理作業環境
    Kitchen_Cooking
        H=(12,Cooking_Schedule,Cooking_Activity)
        comfrt=(Cooking_Metabolic,Kitchen_Clothing,Kitchen_Airflow) ;
        !  在室率：給食調理時間に集中、衛生管理重要
        !  活動強度：調理作業で2.2met（中重労働）
        !  快適性：調理服0.4clo、厨房環境0.6m/s
*
```

**設計のポイント**:
- **学習効率**: 集中力向上に最適な温熱環境（PMV±0.3以内）
- **健康管理**: 熱中症予防と風邪予防の両立
- **年齢考慮**: 児童生徒の体温調節機能に配慮した設定
- **教育活動**: 授業内容に応じた環境の最適化

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
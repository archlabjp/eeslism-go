# PCM（潜熱蓄熱材）

## 概要
PCM（Phase Change Material：潜熱蓄熱材）は、相変化時の潜熱を利用して熱エネルギーを蓄積・放出する材料です。EESLISM Go版では、壁体内や家具内に組み込まれたPCMの熱的挙動をシミュレーションできます。

## データ形式
```
PCM
    <PCM名称>
        [ Condl=<value> ]
        [ Conds=<value> ]
        [ Crol=<value> ]
        [ Cros=<value> ]
        [ Ql=<value> ]
        [ Tl=<value> ]
        [ Ts=<value> ]
        [ Tp=<value> ]
        [ Ctype=<value> ]
        [ DivTemp=<value> ]
        [ spcheattable=<filename> ]
        [ conducttable=<filename> ]
        [ table=<e|h> ]
        [ minTempChange=<value> ]
        [ nWieght=<value> ]
        [ IterateJudge=<value> ]
        [ a=<value> ]
        [ b=<value> ]
        [ c=<value> ]
        [ d=<value> ]
        [ e=<value> ]
        [ f=<value> ]
        [ B=<value> ]
        [ T=<value> ]
        [ bs=<value> ]
        [ bl=<value> ]
        [ skew=<value> ]
        [ omega=<value> ]
        [ -iterate ]
        [ -pcmnode ]
    ;
    ... 繰り返し
*
```

### パラメータ説明

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| PCM名称 | - | PCM材料の識別名。必須。 | |
| Condl | W/(m·K) | 液相状態での熱伝導率 | `FNAN` |
| Conds | W/(m·K) | 固相状態での熱伝導率 | `FNAN` |
| Crol | J/m³K | 液相状態での容積比熱 | `FNAN` |
| Cros | J/m³K | 固相状態での容積比熱 | `FNAN` |
| Ql | J/m³ | 潜熱量 | `FNAN` |
| Tl | ℃ | 液体から凝固が始まる温度 | `FNAN` |
| Ts | ℃ | 固体から融解が始まる温度 | `FNAN` |
| Tp | ℃ | 見かけの比熱のピーク温度 | `FNAN` |
| Ctype | - | 見かけの比熱の特性曲線番号 | 2 (二等辺三角形) |
| DivTemp | - | 比熱数値積分時の温度分割数 | 1 |
| spcheattable | - | 比熱テーブルファイル名 | - |
| conducttable | - | 熱伝導率テーブルファイル名 | - |
| table | - | テーブルタイプ (e:エンタルピー, h:見かけの比熱) | e |
| minTempChange | ℃ | 最小温度変化幅 | 0.5 |
| nWieght | - | 収束計算時の現在ステップ温度の重み係数 | 0.5 |
| IterateJudge | - | 収束判定条件（前ステップ見かけの比熱の割合） | 0.05 |
| a | - | 数学的モデルパラメータ（振幅等） | `FNAN` |
| b | - | 数学的モデルパラメータ（標準偏差等） | `FNAN` |
| c | - | 数学的モデルパラメータ（多項式係数） | `FNAN` |
| d | - | 数学的モデルパラメータ（多項式係数） | `FNAN` |
| e | - | 数学的モデルパラメータ（多項式係数） | `FNAN` |
| f | - | 数学的モデルパラメータ（多項式係数） | `FNAN` |
| B | - | 数学的モデルパラメータ（形状・線形項係数） | `FNAN` |
| T | - | 数学的モデルパラメータ（温度パラメータ） | `FNAN` |
| bs | - | 非対称ガウス関数の左側標準偏差パラメータ | `FNAN` |
| bl | - | 非対称ガウス関数の右側標準偏差パラメータ | `FNAN` |
| skew | - | 誤差関数の歪度パラメータ | `FNAN` |
| omega | - | 誤差関数の分散パラメータ | `FNAN` |
| -iterate | - | 収束計算を行う場合に指定 | `false` |
| -pcmnode | - | PCM温度を節点温度で計算する場合に指定 | `false` (平均温度) |


## 使用例

### 基本的なPCM定義
```
PCM
    ParaffinWax28
        Condl=0.15 Conds=0.20 Crol=2000000 Cros=1800000
        Ql=180000 Ts=26 Tl=30 Tp=28 ;
*
```

### 収束計算ありのPCM
```
PCM
    PCM_HighPerf
        Condl=0.18 Conds=0.22 Ql=200000
        Ts=24 Tl=26 Tp=25
        -iterate ;
*
```

## 壁体での使用方法

PCMを壁体に組み込む場合は、WALL定義でPCM名称を指定します：

```
WALL
    wall_with_pcm
    3
    concrete 0.15
    PCM_layer ParaffinWax28 0.02
    insulation 0.10
```

## 計算方法

### 相変化の判定
- **融解過程**: 温度が融解温度を超えた場合
- **凝固過程**: 温度が凝固温度を下回った場合
- **相変化中**: 融解温度と凝固温度の間で潜熱を考慮

### 熱容量の計算
相変化中の見かけの熱容量は、潜熱と温度幅から計算されます：
```
C_apparent = C_sensible + L / ΔT
```

### 収束計算
収束計算フラグが'y'の場合、PCMの相変化による非線形性を考慮した反復計算が実行されます。

## 出力データ

PCMの状態は以下の項目で出力されます：
- PCM温度
- 液相率
- 蓄熱量
- 放熱量
- 相変化状態

## 注意事項

1. **温度範囲**: 融解温度 ≤ ピーク温度 ≤ 凝固温度の関係を満たす必要があります
2. **収束性**: 収束計算を有効にすると計算時間が増加しますが、精度が向上します
3. **材料特性**: 実際のPCM材料の物性値を正確に入力することが重要です
4. **厚さ**: PCM層の厚さは壁体定義で指定します

## 物性値読み込み機能

### 概要
PCMの物性値をファイルから読み込む機能。温度依存の物性値を詳細に定義可能。

### パラメータ
- `spcheattable`: 比熱テーブルファイル名
- `conducttable`: 熱伝導率テーブルファイル名
- `table`: テーブルタイプ
  - `e`: エンタルピーテーブル（デフォルト）
  - `h`: 見かけの比熱テーブル

### テーブルファイル形式
```
温度[℃] 特性値
20.0 50000.0
22.0 52000.0
24.0 55000.0
26.0 180000.0
28.0 185000.0
30.0 58000.0
32.0 60000.0
```

### 使用例
```
PCM
    PCM_WithTable
        spcheattable=enthalpy_data.txt table=e
        conducttable=conductivity_data.txt
        Condl=0.15 Conds=0.20 Ql=180000 Ts=26 Tl=30 Tp=28 ;
    
    PCM_SpecificHeat
        spcheattable=specific_heat_data.txt table=h
        Condl=0.15 Conds=0.20 ;
*
```

## 見かけの比熱の特性曲線の数学的モデル (Ctype)

PCMの見かけの比熱特性を表現する数学的モデル：

### Ctype=0: 熱伝導率計算
感熱のみの線形補間

### Ctype=1: 一定潜熱
矩形分布の潜熱

### Ctype=2: 二等辺三角形
三角形分布の潜熱
- T < Tp: 左側の線形増加
- T > Tp: 右側の線形減少

### Ctype=3: 双曲線関数
cosh関数による潜熱分布

### Ctype=4: 対称ガウス関数
正規分布による潜熱分布

### Ctype=5: 非対称ガウス関数
左右で異なる標準偏差を持つガウス分布
- T ≤ Tp: bs パラメータ使用
- T > Tp: bl パラメータ使用

### Ctype=6: 誤差関数（歪度付き）
skew パラメータによる非対称分布

### Ctype=7: 有理関数
多項式の比による複雑な分布

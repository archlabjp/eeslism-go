# 7. チュートリアル

## 7.1 はじめに

このチュートリアルでは、EESLISM Go版を使用して建物のエネルギーシミュレーションを行う方法を、段階的に学習します。簡単な1室モデルから始めて、徐々に複雑なシステムまで扱えるようになることを目標とします。

## 7.2 基本的な使い方

### 7.2.1 プログラムの実行

EESLISM Go版の基本的な実行方法：

```bash
# 基本的な実行
./eeslism input_file.txt

# EFLファイルのディレクトリを指定
./eeslism input_file.txt --efl Base

# ヘルプの表示
./eeslism --help
```

### 7.2.2 ファイル構成

EESLISM Go版では以下のファイルが必要です：

```
project/
├── input_file.txt          # メイン入力ファイル
├── Base/                   # 基礎データディレクトリ
│   ├── tokyo_3column_SI.has # 気象データ
│   ├── wbmlist.efl         # 壁材料リスト
│   ├── reflist.efl         # 圧縮機特性リスト
│   ├── dayweek.efl         # 曜日・祝日設定
│   └── supw.efl            # 給水温・地中温データ
└── output/                 # 出力ファイル（自動生成）
```

## 7.3 チュートリアル1: 簡単な1室モデル

最初に、最もシンプルな1室モデルを作成してみましょう。

### 7.3.1 入力ファイルの作成

以下の内容で `simple_room.txt` を作成します：

```
! 簡単な1室モデル

WEEK
    1/1=Sun ;

TITLE
    Simple Room Model ;

GDAT
    FILE    w=tokyo_3column_SI.has ;
    RUN     1/1-1/7 dTime=3600 ;
    PRINT   1/1-1/7 *wd ;
*

SCHTB
    %s -s    AllDay   001-(1)-2400 ;
*

EXSRF
    r=0.1 ;
    south   a=0. alo=25 ;
    west    a=90. alo=25 ;
    north   a=180. alo=25 ;
    east    a=270. alo=25 ;
    Hor     alo=23.3 ;
*

WALL
    exterior_wall
    3
    concrete    0.15
    insulation  0.10
    gypsum      0.012
    ;
    
    interior_wall
    1
    gypsum      0.012
    ;
*

WINDOW
    single_glass
    glass=6 frame=20 ;
*

ROOM
    SimpleRoom
    Rg=1.0 Cg=1000000 Ag=50.0
    
    south   exterior_wall   10.0    single_glass    5.0     south
    west    exterior_wall   8.0     -               -       west
    north   exterior_wall   10.0    -               -       north
    east    interior_wall   8.0     -               -       -
    roof    exterior_wall   50.0    -               -       Hor
    floor   interior_wall   50.0    -               -       -
    ;
*

SCHNM
    Room_Schedule   SimpleRoom  ;
*

VCFILE
    Room_Schedule
    %v -t   RoomTemp    AllDay  20.0    ;
*
```

### 7.3.2 入力ファイルの解説

#### WEEK（週設定）
```
WEEK
    1/1=Sun ;
```
- 1月1日を日曜日として設定

#### TITLE（タイトル）
```
TITLE
    Simple Room Model ;
```
- シミュレーションのタイトルを設定

#### GDAT（気象データと実行条件）
```
GDAT
    FILE    w=tokyo_3column_SI.has ;
    RUN     1/1-1/7 dTime=3600 ;
    PRINT   1/1-1/7 *wd ;
*
```
- `FILE`: 東京の気象データを使用
- `RUN`: 1月1日から1週間、1時間間隔で計算
- `PRINT`: 結果を出力

#### SCHTB（スケジュールテーブル）
```
SCHTB
    %s -s    AllDay   001-(1)-2400 ;
*
```
- 24時間一定値のスケジュールを定義

#### EXSRF（外皮面）
```
EXSRF
    r=0.1 ;
    south   a=0. alo=25 ;
    west    a=90. alo=25 ;
    north   a=180. alo=25 ;
    east    a=270. alo=25 ;
    Hor     alo=23.3 ;
*
```
- 各方位の日射面を定義
- `r=0.1`: 地面反射率
- `a`: 方位角、`alo`: 緯度

#### WALL（壁体構成）
```
WALL
    exterior_wall
    3
    concrete    0.15
    insulation  0.10
    gypsum      0.012
    ;
```
- 外壁：コンクリート150mm + 断熱材100mm + 石膏ボード12mm

#### WINDOW（窓）
```
WINDOW
    single_glass
    glass=6 frame=20 ;
```
- 単板ガラス6mm、フレーム幅20mm

#### ROOM（室）
```
ROOM
    SimpleRoom
    Rg=1.0 Cg=1000000 Ag=50.0
    
    south   exterior_wall   10.0    single_glass    5.0     south
    ...
    ;
```
- 室名：SimpleRoom
- 室内発熱：1W、熱容量：1000000J/K、床面積：50m²
- 各面の構成を定義

### 7.3.3 実行と結果確認

```bash
# シミュレーション実行
./eeslism simple_room.txt

# 出力ファイルの確認
ls simple_room_*.es
```

出力ファイルの例：
- `simple_room_rm.es`: 室温・湿度
- `simple_room_sf.es`: 表面温度
- `simple_room_wd.es`: 気象データ

## 7.4 チュートリアル2: 空調システム付きモデル

次に、空調システムを追加したモデルを作成します。

### 7.4.1 入力ファイルの拡張

先ほどの `simple_room.txt` に以下を追加：

```
EQPCAT
    HCLOAD
    SimpleAC D 5000 3.0 ;
*

SYSCMP
    HCLOAD  AC_SimpleRoom   SimpleAC ;
*

SYSPTH
    AC_SimpleRoom   SimpleRoom ;
*

CONTL
    AC_SimpleRoom
    %v -t   SetTemp_Cool    Summer  26.0    ;
    %v -t   SetTemp_Heat    Winter  20.0    ;
*
```

### 7.4.2 追加部分の解説

#### EQPCAT（機器仕様）
```
EQPCAT
    HCLOAD
    SimpleAC D 5000 3.0 ;
*
```
- 直膨コイル型空調機、容量5000W、COP3.0

#### SYSCMP（システム構成要素）
```
SYSCMP
    HCLOAD  AC_SimpleRoom   SimpleAC ;
*
```
- 空調機器をシステムに組み込み

#### SYSPTH（システム経路）
```
SYSPTH
    AC_SimpleRoom   SimpleRoom ;
*
```
- 空調機器と室を接続

#### CONTL（制御）
```
CONTL
    AC_SimpleRoom
    %v -t   SetTemp_Cool    Summer  26.0    ;
    %v -t   SetTemp_Heat    Winter  20.0    ;
*
```
- 夏季26℃、冬季20℃で温度制御

## 7.5 チュートリアル3: 太陽光発電システム

太陽光発電システムを追加してみましょう。

### 7.5.1 PVシステムの追加

```
EQPCAT
    PV
    CrystalSi_4kW 4000 20.0 0.97 0.95 0.94 0.96 0.95 -0.45 20.0 C 0.0175 0.0 A ;
*

SYSCMP
    PV  PV_South    CrystalSi_4kW   south ;
*
```

### 7.5.2 解説

- 4kWの結晶系太陽電池
- 南面に設置
- 架台設置型（設置方式A）

## 7.6 チュートリアル4: PCM（潜熱蓄熱材）

PCMを壁体に組み込んでみましょう。

### 7.6.1 PCMの定義

```
PCM
    ParaffinWax28
    Ql=200000000
    Condl=0.15
    Conds=0.20
    Crol=1000000
    Cros=1000000
    Ts=26.0
    Tl=30.0
    Tp=28.0
    -iterate
    ;
*
```

### 7.6.2 PCM入り壁体

```
WALL
    pcm_wall
    4
    concrete    0.15
    insulation  0.05
    ParaffinWax28   0.02
    gypsum      0.012
    ;
*
```

## 7.7 標準プランの解析

付属の標準プランサンプルを使用してみましょう。

### 7.7.1 標準プランの実行

```bash
# 標準プランの実行
./eeslism samples/standard-plan-no-hcap-PCM-CM-fsolm.txt

# 実行時間の測定
time ./eeslism samples/standard-plan-no-hcap-PCM-CM-fsolm.txt
```

### 7.7.2 標準プランの特徴

- **建物**: 19室の戸建住宅
- **計算期間**: 1ヶ月助走期間 + 12ヶ月本計算
- **時間間隔**: 30分間隔
- **機能**: PCM、空調システム、太陽熱利用

### 7.7.3 出力ファイルの解析

主要な出力ファイル：
- `*_rm.es`: 室温・湿度・負荷
- `*_sf.es`: 表面温度
- `*_mr.es`: 月集計（室）
- `*_mt.es`: 月集計（全体）

## 7.8 結果の可視化

### 7.8.1 データの抽出

出力ファイルはテキスト形式なので、表計算ソフトやプログラムで処理できます：

```bash
# 室温データの抽出例
grep "SimpleRoom" simple_room_rm.es | head -24
```

### 7.8.2 グラフ化の例

Python を使用した例：

```python
import pandas as pd
import matplotlib.pyplot as plt

# データ読み込み
data = pd.read_csv('simple_room_rm.es', sep='\s+', header=None)

# 室温のグラフ化
plt.figure(figsize=(12, 6))
plt.plot(data[2])  # 3列目が室温
plt.title('Room Temperature')
plt.xlabel('Time (hours)')
plt.ylabel('Temperature (°C)')
plt.grid(True)
plt.show()
```

## 7.9 よくある問題と解決方法

### 7.9.1 実行エラー

#### ファイルが見つからない
```
Error: Cannot open file: tokyo_3column_SI.has
```
**解決方法**: EFLディレクトリを正しく指定
```bash
./eeslism input.txt --efl Base
```

#### 構文エラー
```
Error: Syntax error in line 25
```
**解決方法**: 入力ファイルの構文を確認、セクション終了の `*` を確認

### 7.9.2 計算が収束しない

#### 収束しない場合
```
Warning: Calculation did not converge
```
**解決方法**: 
- 時間間隔を小さくする（dTime=1800 → dTime=900）
- 最大反復回数を増やす（MaxIterate=200）

### 7.9.3 メモリ不足

#### 大規模モデルでのメモリ不足
**解決方法**:
- 計算期間を短くする
- 出力項目を減らす
- 時間間隔を大きくする

## 7.10 次のステップ

### 7.10.1 応用例

1. **複数室モデル**: 住宅全体のモデル化
2. **設備システム**: 太陽熱集熱、蓄熱槽の組み合わせ
3. **制御システム**: 複雑な制御ロジックの実装
4. **パラメータスタディ**: 断熱性能の影響分析

### 7.10.2 参考資料

- [入力データ形式詳細](../format/README.md)
- [機器仕様一覧](../format/EQPCAT.md)
- [Go言語版特有機能](./6_go_features.md)

### 7.10.3 コミュニティ

- GitHub リポジトリ: https://github.com/archlabjp/eeslism-go
- 問題報告: GitHub Issues
- 機能要求: GitHub Discussions

## 7.11 練習問題

### 問題1: 基本モデルの作成
2室（リビング・寝室）の住宅モデルを作成し、各室の温度変化を比較してください。

### 問題2: 断熱性能の比較
断熱材厚さを50mm、100mm、150mmで変更し、暖房負荷の違いを分析してください。

### 問題3: 太陽光発電の効果
太陽光発電システムありとなしで、年間の電力収支を比較してください。

### 問題4: PCMの効果
PCM入り壁体と通常壁体で、夏季の室温変動を比較してください。

これらの練習問題を通じて、EESLISM Go版の様々な機能を習得できます。
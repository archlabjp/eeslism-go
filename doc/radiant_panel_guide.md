# 放射パネル（輻射パネル）入力ガイド

## 概要

EESLISMでは、床暖房、天井放射冷暖房、壁面放射などの放射パネルシステムをモデル化できます。放射パネルは建築部位（壁体）として定義し、設備システムの一部として熱源機器と接続します。

## 入力の流れ

放射パネルの入力は以下の順序で行います：

1. **WALL**: パネル付き壁体の定義
2. **ROOM**: 室でのパネル部位指定
3. **SYSCMP**: 熱源機器の定義
4. **SYSPTH**: システム経路（配管）の定義
5. **CONTL**: 制御設定

## 1. WALL - 壁体定義

### 基本書式
```
WALL
    -<ble>:<panelname>  <layer1> <P[:効率]> <layer2> ... ;
*
```

### パラメータ
- `<ble>`: 部位記号（f:床、c:天井、i:内壁）
- `<panelname>`: パネル壁体名（英字で始まる）
- `<P>`: 放射パネル発熱面位置の指定
- `<P:効率>`: パネル熱通過有効度（省略時は0.7）

### 例

#### 床暖房パネル
```
WALL
    -f:FloorPanel  WD-12 <P> WDB-9 FPS-100 ;
*
```

#### 天井放射パネル
```
WALL
    -c:CeilPanel   GPB-12 <P:0.8> WDB-9 FPS-100 ;
*
```

#### 壁面放射パネル
```
WALL
    -i:WallPanel   GPB-12 <P> WDB-9 FPS-100 ;
*
```

## 2. ROOM - 室構成

### 基本書式
```
ROOM
    <roomname>
        <direction>: -<ble> <panelname> <area> i=<elementname> ;
    *
*
```

### パラメータ
- `<panelname>`: WALL で定義したパネル壁体名
- `<area>`: パネル面積 [m²]
- `i=<elementname>`: SYSPTH で使用する要素名

### 例

#### 床暖房の場合
```
ROOM
    LivingRoom  Vol=50.0
        Hor: -f FloorPanel 20.0 i=LR_FloorPanel ;
    *
*
```

#### 天井放射の場合
```
ROOM
    Office  Vol=40.0
        Hor: -c CeilPanel 16.0 i=Office_CeilPanel ;
    *
*
```

#### 壁面放射の場合
```
ROOM
    Bedroom  Vol=30.0
        (LivingRoom): -i WallPanel 15.0 i=BR_WallPanel ;
    *
*
```

## 3. SYSCMP - システム構成要素

### 熱源機器の定義
```
SYSCMP
    <heatername> -c <catalogname> ;
*
```

### 例
```
SYSCMP
    Boiler -c BoilerCat ;
*
```

## 4. SYSPTH - システム経路

### 基本書式
```
SYSPTH
    <pathname>  -sys A  -f W
        > (<flowrate>) <heater> <panelelement> > ;
*
```

### パラメータ
- `<pathname>`: システム経路名
- `-sys A`: 空調・暖房システム
- `-f W`: 水系統
- `<flowrate>`: 流量 [kg/s]
- `<panelelement>`: ROOM で指定した要素名

### 例

#### 床暖房システム
```
SYSPTH
    FloorHeatingPath  -sys A  -f W
        > (0.1) Boiler LR_FloorPanel > ;
*
```

#### 天井放射システム
```
SYSPTH
    CeilRadiantPath  -sys A  -f W
        > (0.08) Chiller Office_CeilPanel > ;
*
```

## 5. CONTL - 制御設定

### 制御方式

#### 1. 表面温度制御
```
CONTL
    if (<roomname>_<panelname>_Ts < <setpoint>)
        <pathname>=<schedule> ;
*
```

#### 2. 室温制御
```
CONTL
    if (<roomname>_Tr < <setpoint>)
        <pathname>=<schedule> ;
*
```

#### 3. 負荷制御
```
CONTL
    LOAD:H -e <roomname> <roomname>_Tr=<schedule> ;
    if (<condition>)
        <pathname>=<schedule> ;
*
```

### 制御例

#### 床暖房の表面温度制御
```
CONTL
    if (LivingRoom_LR_FloorPanel_Ts < 28)
        FloorHeatingPath=HeatingSchedule ;
    
    LOAD -e Boiler Boiler_Tout=45 ;
    
    %s -s HeatingSchedule 601-(-)-900 1601-(-)-2300 ;
*
```

#### 天井放射の室温制御
```
CONTL
    if (Office_Tr > 26)
        CeilRadiantPath=CoolingSchedule ;
    
    LOAD -e Chiller Chiller_Tout=15 ;
    
    %s -s CoolingSchedule 800-(-)-1800 ;
*
```

## 完全な入力例

### 床暖房システムの例
```
WALL
    -f:FloorPanel  WD-12 <P> WDB-9 FPS-100 ;
*

ROOM
    LivingRoom  Vol=50.0
        south: -E 20.0 ;
        north: -E 20.0 ;
        east:  -E 15.0 ;
        west:  -E 15.0 ;
        Hor:   -R 25.0 ;
        Hor:   -f FloorPanel 25.0 i=LR_FloorPanel ;
    *
*

EQPCAT
    BoilerCat boi Qo=5000 ;
*

SYSCMP
    Boiler -c BoilerCat ;
*

SYSPTH
    FloorHeatingPath  -sys A  -f W
        > (0.12) Boiler LR_FloorPanel > ;
*

CONTL
    if (LivingRoom_LR_FloorPanel_Ts < 28)
        FloorHeatingPath=HeatingSchedule ;
    
    LOAD -e Boiler Boiler_Tout=45 ;
    
    %s -s HeatingSchedule 601-(-)-900 1601-(-)-2300 ;
*
```

## パネルタイプ別設定

### 床暖房（Floor Heating）
- **部位記号**: `-f`
- **設置位置**: 床面
- **供給温度**: 35-45℃
- **表面温度**: 25-30℃
- **制御**: 表面温度制御が一般的

### 天井放射（Ceiling Radiant）
- **部位記号**: `-c`
- **設置位置**: 天井面
- **供給温度**: 冷房15-18℃、暖房35-40℃
- **表面温度**: 冷房18-22℃、暖房28-32℃
- **制御**: 室温制御が一般的

### 壁面放射（Wall Radiant）
- **部位記号**: `-i`
- **設置位置**: 内壁面
- **供給温度**: 冷房15-18℃、暖房35-40℃
- **表面温度**: 冷房18-22℃、暖房28-32℃
- **制御**: 室温制御が一般的

## 制御パラメータ

### 出力変数
- `<roomname>_<panelname>_Ts`: パネル表面温度 [℃]
- `<roomname>_<panelname>_Q`: パネル供給熱量 [W]
- `<roomname>_Tr`: 室温 [℃]
- `<heatername>_Tout`: 熱源出口温度 [℃]

### スケジュール設定
```
%s -s <schedulename> <start_time>-(<value>)-<end_time> ;
%s -v <schedulename> <start_time>-(<value>)-<end_time> ;
```

## トラブルシューティング

### よくあるエラー

1. **パネル名の不一致**
   - WALL、ROOM、SYSPTHでのパネル名が一致していることを確認

2. **要素名の未定義**
   - ROOMで`i=<elementname>`が正しく指定されていることを確認

3. **流量設定エラー**
   - SYSPTHでの流量が適切に設定されていることを確認

4. **制御条件の設定ミス**
   - CONTRLでの条件式と変数名が正しいことを確認

### デバッグのヒント

1. **出力ファイルの確認**
   - `*s`オプションで表面温度出力を有効にする
   - パネル表面温度の変化を確認

2. **段階的な確認**
   - まず単純な制御から始める
   - 複雑な制御は段階的に追加

3. **サンプルファイルの参照**
   - `eeslism/radiant_ceiling_cooling.txt`を参考にする

## 関連ファイル

- サンプル: `eeslism_sample/radiant_*.txt`
- 材料定義: `Base/wbmlist.efl`
- 基礎データ: `Base/supw.efl`

## 参考

- [2.4 建築データ](2_4.md) - WALL、ROOM定義の詳細
- [2.5 設備機器データ](2_5.md) - SYSCMP定義の詳細  
- [2.6 システム経路](2_6.md) - SYSPTH定義の詳細
- [2.7 制御データ](2_7.md) - CONTL定義の詳細
# ROOM（室）

## 概要
ROOM（室）は、建物内の室空間の熱的特性と構成要素を定義するためのデータ形式です。EESLISM Go版では、室容積、熱容量、日射分配、各種出力設定、および室を構成する壁体・窓・設備などの部位データを詳細に設定できます。

## データ形式
```
ROOM
    <rmname>
        [ Vol=<vol> ]
        [ flrsr=<flrsr> ]
        [ fsolm=<fsolm> ]
        [ Hcap=<hcap> ]
        [ Mxcap=<mxcap> ]
        [ MCAP=<mcap> ]
        [ CM=<cm> ]
        [ alc=<alc> ]
        [ PCMFurn=<PCMFurn> ]
        [ OTc=<otc> ]
        [ rsrnx ]
        [ *s ]
        [ *q ]
        [ *sfe ]
        <component>
        <component>
        ... 繰り返し
    *
    <rmname>
        Vol=<vol>
        ...
    *
    以下繰り返し
*
```

### 室パラメータ説明

#### 基本パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| rmname | - | 室名称。必須。 | 必須 |
| Vol | m³ | 室容積。直方体の場合は「間口*奥行き*高さ」で自動計算可能 | 必須 |

#### 熱容量パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| Hcap | J/K | 室内空気に付加する熱容量 | - |
| Mxcap | kg/(kg/kg) | 室内空気に付加する湿気容量 | - |
| MCAP | J/K | 室内に置かれた物体の熱容量（スケジュール名） | - |
| CM | W/K | 室内物体と室内空気間の熱コンダクタンス（スケジュール名） | - |
| PCMFurn | - | PCM家具の指定 | - |

#### 日射・表面特性パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| flrsr | - | 床の日射吸収比率（スケジュール名） | - |
| fsolm | - | 家具への日射吸収割合（スケジュール名） | - |
| alc | W/m²K | 室内表面熱伝達率（スケジュール名） | - |
| OTc | - | 作用温度設定時の対流成分重み係数 | - |

#### 計算オプション

| パラメータ | 説明 | デフォルト値 |
|:---|:---|:---|
| rsrnx | 共用壁でない内壁で隣室側表面放射を考慮 | false |

#### 出力設定

| パラメータ | 説明 | 出力ファイル |
|:---|:---|:---|
| *s | 室内表面温度・熱流・熱伝達率の出力 | outfile_sf.es, outfile_sfq.es, outfile_sfa.es |
| *q | 日射熱取得・室内発熱・隙間風熱取得の出力 | outfile_rq.es, outfile_dqr.es |
| *sfe | 要素別壁体表面温度の出力 | - |

## 部位データ形式
```
[<ename>:[<rmp>[-<wname>]]]
    -<ble>
    [ <name> ]
    [ alc=<schnamealc> ]
    { <area> | <width>*<height> | A=<area> }

    [ fsol=<fsol> ]
    [ alc=<alc> ]
    [ alr=<alr> ]
    [ PVap=<pvap> ]
    [ PVcap=<pvcap> ]
    [ Wsu=<wsu> ]
    [ Wsd=<wsd> ]
    [ Ndiv=<ndiv> ]
    [ i=<panelname> ]
    [ sb=<sbnamedd> ]
    [ rmp=<winname> ]
    [ tnxt=<tnxt> ]
    [ c=<c> ]
    [ r=<r> ]
    [ sw=<sw> ]
    [ rmp=<sname>]
    [ e=<e>]
    [ *p ]
    [ *sfe ]
    [ *shd ]
;
```
### 部位パラメータ説明

#### 基本パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| ename | - | 部位要素名 | - |
| ble | - | 部位コード（E/R/F/i/c/f/W） | 必須 |
| n | - | 部位番号 | - |
| area | m² | 面積 | 必須 |
| width | m | 幅 | - |
| height | m | 高さ | - |

#### 部位コード（ble）

| コード | 説明 |
|:---|:---|
| E | 外壁 |
| R | 屋根 |
| F | 床（外部） |
| i | 内壁 |
| c | 天井（内部） |
| f | 床（内部） |
| W | 窓 |

#### 表面・日射特性パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| fsol | - | 部位室内表面の日射吸収比率 | - |
| alc | W/m²K | 室内側熱伝達率（スケジュール名） | - |
| alr | W/m²K | 室外側熱伝達率 | - |

#### 隣接・境界条件パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| r | - | 隣室名 | - |
| e | - | 外表面名（[EXSRF](EXSRF.md)参照） | - |
| c | - | 隣室温度係数 | - |
| tnxt | - | 隣接空間への日射分配 | - |

#### 太陽電池・集熱器パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| PVap | - | 太陽電池関連パラメータ | - |
| PVcap | W | 太陽電池アレイの定格発電量 | - |
| Wsu | m | 屋根一体型空気集熱器の通気層上面幅 | - |
| Wsd | m | 屋根一体型空気集熱器の通気層下面幅 | - |
| Ndiv | - | 空気式集熱器の流れ方向分割数 | - |

#### 設備・付属要素パラメータ

| パラメータ | 単位 | 説明 | デフォルト値 |
|:---|:---|:---|:---|
| i | - | 放射暖冷房パネル名 | - |
| sb | - | 日よけ名 | - |
| rmp | - | RMP名または窓名 | - |
| sw | - | 窓変更設定番号 | - |

#### 出力設定

| パラメータ | 説明 |
|:---|:---|
| *p | 壁体内部温度出力指定 |
| *sfe | 要素別壁体表面温度出力指定 |
| *shd | 日よけの影面積出力 |

## 使用例

### 基本的な室定義
```
ROOM
    LivingRoom
        Vol=3.5*4.2*2.7
        Hcap=50000
        fsolm=furniture_schedule
        *s *q
        
        wall_south:
            -E 15.0
            fsol=0.7 e=south_wall ;
        
        wall_north:
            -E 10.0
            r=Kitchen ;
        
        window_south:
            -W 8.0
            rmp=double_glass ;
    *
*
```

### 集熱器付き屋根を持つ室
```
ROOM
    TopFloor
        Vol=50.0
        *s
        
        roof_collector:
            -R 25.0
            PVcap=3000 Wsu=1.5 Wsd=1.2 Ndiv=5
            i=roof_panel ;
    *
*
```

### PCM家具付きの室
```
ROOM
    Bedroom
        Vol=4.0*3.0*2.5
        PCMFurn=PCM_furniture
        MCAP=furniture_capacity
        CM=furniture_conductance
        *sfe
        
        wall_ext:
            -E 12.0
            fsol=0.6 *sfe ;
    *
*
```

## 計算方法

### 室温計算
室温は以下の熱収支式から計算されます：
```
ρ_air * V * C_air * dT/dt = Q_wall + Q_window + Q_solar + Q_internal + Q_ventilation
```

### 日射分配
室内に入射した日射は以下の順序で分配されます：
1. 床面への直接日射
2. 家具への日射吸収
3. 各壁面への拡散日射

### 湿度計算
室内湿度は水蒸気収支から計算され、湿気容量も考慮されます。

## 出力データ

室の計算結果として以下の項目が出力されます：
- 室温・湿度
- 各部位の表面温度・熱流
- 日射熱取得量
- 内部発熱量
- 換気による熱損失

## 注意事項

1. **室容積**: 正確な室容積の設定が重要です
2. **部位面積**: 各部位の面積は実際の寸法に基づいて設定してください
3. **隣室設定**: 隣室がある場合は適切な隣室名を指定してください
4. **スケジュール**: 時間変動する値はスケジュール名で指定してください
5. **出力設定**: 必要な出力項目のみを指定して計算負荷を軽減してください
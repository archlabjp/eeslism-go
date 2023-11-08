# EQPCAT 機器仕様データ一覧

- [BOI](#boi) ボイラー
- [REFA](#refa) 冷温水方式の圧縮式電動ヒートポンプ,仮想熱源
- [COL](#col) 架台設置型太陽熱集熱器
- [STANK](#stank) 蓄熱槽(熱交換型内蔵型含む)
- [HEX](#hex) 熱交換器
- [HCC](#hcc) 冷温水コイル
- [PIPE](#pipe) 配管
- [DUCT](#duct) ダクト
- [PUMP](#pipe) ポンプ
- [FAN](#fan) ファン
- [VAV](#vav) VAVユニット
- [STHEAT](#stheat) 電気蓄熱式暖房器
- [THEX](#thex) 全熱交換器
- [PV](#pv) 架台設置型太陽電池
- [OMVAV](#omvav)
- [DESI](#desi) デシカント槽  (Ver7.2時点では未定義)
- [EVAC](#evac) 気化冷却器  (Ver7.2時点では未定義)


NOTE: 屋根一体型太陽熱集熱器または屋根一体型PVは建築部位として [WALL](WALL.md)で入力する

書式:
```
EQPCAT
   <type> <catname> <key>=<value>  <key>=<value> ... <option> ;
   <type> <catname> <key>=<value>  <key>=<value> ... <option> ;
   ...
*
```
- type: 機器種別 ex) `BOI`, `REFA`
- catname: カタログ名
- key: 機器種別ごとに定められる入力項目名
- value: `key` に対応した値
- option: 機器種別ごとに定められた追加の指定


例)
```
EQPCAT
   STHEAT ETS4	Q=4000  KA=20  eff=0.3  PCM=PCMcat  ;
*
```

## BOI
ボイラー機器仕様
- p
- en: 燃料種類 G：ガス、O：灯油、E：電気
- Qo: 定格能力 [W]
- Qmin: 最小出力 [W]
- blwQmin: 最小出力以下の時のON/OFF指定
- eff: 定格時効率 [-]
- Ph: 補機動力 [W]

# REFA
チラー、ヒートポンプチラー（空気熱源）
- c
- a
- Ph
- m
- Qo
- Go
- Two
- eo
- Qex
- Gex
- eex
- Tex
- W

## COL
太陽熱集熱器
- b0
- b1

## STANK
蓄熱槽
- Vol
- KAtop
- KAbtm
- KAside
- gxr
- ri:in_eff
- ri:in_KA
- ri:in_dL

## HEX
熱交換器
- eff
- KA

## THEX
全熱交換器
- et
- eh

## HCC
冷温水コイル
- et: コイル温度効率 [-]
- KA: コイルの熱通過率と伝熱面積の積 [W/K]
- eh: コイルエンタルピー効率
NOTE:
- etもしくはKAのどちらかを必ず設定する
- エンタルピー効率を指定しないときの冷水コイルは乾きコイルとなる。
- 温水コイルの場合にはet、ehは指定せずにKAのみ指定する。


## PIPE
配管
- Ko: 配管線熱通過率 [W/mK]

## DUCT
ダクト
- Ko: ダクト熱通過率 [W/mK]

## PUMP
ポンプ
- type
- Go
- Wo
- qref
- a0
- a1
- a2
- Ic
- G
- qef

## FAN
ファン(送風機)
- type: 
  - C:定風量ファン
  - Vd：変風量ファン（ダンパ制御）
  - Vs: 変風量ファン（サクションベーン制御）
  - Vp: 変風量ファン（可変ピッチ制御）
  - Vr：変風量ファン（可変速制御）
- Go: 風量[kg/s]
- Wo: モーター入力電力[W]
- qef: ファン発熱比率

## VAV
VAVユニット
- Gmin: 最小風量[kg/s]
- Gmax: 最大風量[kg/s]

## STHEAT
電気蓄熱式暖房機
- Q: 電気蓄熱式暖房器カタログ名
- KA: 電気蓄熱式暖房器ヒーター容量[W]
- eff: 電気蓄熱式暖房器表面の熱損失係数 [W/K]
- Hcap: 電気蓄熱式暖房器の蓄熱体熱容量 [J/K]
- PCM: 

## PV
太陽電池
- KHD: 日射量気候変動補正係数[-] (1.0)
- KPD: 経時変化補正係数[-](0.95)
- KPM: アレイ負荷整合補正係数[-] (0.94)
- KPA: アレイ回路補正係数[-] (0.97)
- EffInv: インバータ実効効率[-] (0.97)
- Apmax: 最大出力温度係数
- PVcap: 太陽電池アレイ設置容量
- Area: アレイ面積
- InstallType:PV設置方法（A:架台設置形、B:屋根置き形、C:屋根材型（裏面通風構造があるタイプ）

## OMVAV
屋根一体型空気集熱器の出口温度設定変風量制御ユニット
- Gmin: 最小風量[kg/s]
- Gmax: 最大風量[kg/s]

## OAVAV
[OMVAV](#omvav)の別名。

## DESI
デシカント槽

## EVAC
気化冷却器
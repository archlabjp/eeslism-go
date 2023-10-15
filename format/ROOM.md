# ROOM 室

書式:
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
- vol: 室容積 [m3]入力室が直方体の場合には間口、奥行き、高さを'*'でつなげると、EESLISM内部で室容積を計算する。
- flrsr: 床の日射吸収比率の指定。設定値のスケジュール名
- fsolm: 家具への日射吸収割合。設定値のスケジュール名
- hcap: 室内空気に付加する熱容量 [J/K]
- mxcap: 室内空気に付加する湿気容量 [kg/(kg/Kg)]
- mcap: 室内に置かれた物体の熱容量 [J/K]。設定値のスケジュール名
- cm: 室内に置かれた物体と室内空気との間の熱コンダクタンス [W/K]。設定値のスケジュール名
- alc: 室内表面熱伝達率[W/m2K]。設定値のスケジュール名（設定値を指定するときのみ必要。この部屋の全ての室内表面に適用される。）
- otc: 作用温度設定時の対流成分重み係数
- rsrnx: 共用壁でない内壁で隣室側表面に吸収される放射を考慮するときに指定。
- *s: outfile_sf.es への室内表面温度、outfile_sfq.esへの部位別表面熱流、outfile_sfa.esへの部位別表面熱伝達率の出力指定
- *q: outfile_rq.es、outfile_dqr.es への日射熱取得、室内発熱、隙間風熱取得要素の出力指定
* *sfe: 各部位に `*sfe` を指定するのと同じです。 要素別壁体表面温度を出力します。
- compoment: 部位データ

部位データの書式:
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
- ble: E、R、F、I、c、fなど 窓の場合はW

- name: 壁体名または窓名
- sname: RMP名
- fsol: 部位室内表面の日射吸収比率 [-]
- r: 隣室名
- e: 外表面名 (See: [EXSRF](EXSRF.md))
- c: 隣室温度係数
- sbnamedd: 日よけ名
- area: 面積 [m2]
- width: 幅 [m]
- height: 高さ [m]
- alc:
- alr:
- sw: 窓変更設定番号
- pvcap: 空気式集熱器で太陽電池(PV)付のときのアレイの定格発電量[W],これを記入すると太陽電池付部位となる。（例：屋根一体型太陽光発電アレイ）
- wsu: 屋根一体型空気集熱器の通気層上面の幅
- wsd: 屋根一体型空気集熱器の通気層下面の幅
- ndiv: 空気式集熱器のときの流れ方向（入口から出口）の分割数
- tnxt: 当該部位への入射日射の隣接空間への日射分配（連続空間の隣室への日射分
- panelname: 壁体名。放射暖冷房パネル、部位一体型集熱器のときにSYSPTHでの要素名で使用する。
- xxx: 面積 [m2]
- sbnamedd: 日除け名
- winname: 外表面定義COORDNTで定義したwinnameの指定。窓面の影を考慮する場合に使用する。
- *p:  壁体内部温度出力指定
- *sfe: 要素別壁体表面温度出力指定
- *shd: 日よけの影面積出力
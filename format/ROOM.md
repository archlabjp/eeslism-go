# ROOM 室

書式:
```
ROOM
    <rmname>
        Vol=<vol>
        Hcap=<hcap>
        Mxcap=<mxcap>
        [ MCAP=<mcap> ]
        [ CM=<cm> ]
        [ *s ]
        [ *q ]
        [ alc=<schnamealc> ]
        [ rsrnx ]
        <component> ;
        <component> ;
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
- hcap: 室内空気に付加する熱容量 [J/K]
- mxcap: 室内空気に付加する湿気容量 [kg/(kg/Kg)]
- mcap: 室内に置かれた物体の熱容量 [J/K]
- cm: 室内に置かれた物体と室内空気との間の熱コンダクタンス [W/K]
- schnamealc: alc 室内表面熱伝達率[W/m2K]。schnamealcはalcの設定値のスケジュール名（設定値を指定するときのみ必要。この部屋の全ての室内表面に適用される。）
- rsrnx: 共用壁でない内壁で隣室側表面に吸収される放射を考慮するときに指定。
- *s: outfile_sf.es への室内表面温度、outfile_sfq.esへの部位別表面熱流、outfile_sfa.esへの部位別表面熱伝達率の出力指定
- *q: outfile_rq.es、outfile_dqr.es への日射熱取得、室内発熱、隙間風熱取得要素の出力指定
- compoment: 部位データ

部位データの書式:
```
[<ename>[:<rmp>[-<wname>]]] -<ble> [<name>] [alc=<schnamealc>] <xxx> [*p] [fsol=<fsol>] [PVap=<pvap>] [PVcap=<pvcap>] [Ndiv=<ndiv>] [i=<panelname>]  [ sb=sbnamedd ] [ rmp=winname ]
```
- ble: E、R、F、I、c、fなど 窓の場合はW
- name: 壁体名または窓名
- xxx: 面積 [m2]
- sbnamedd: 日除け名
- winname: 外表面定義COORDNTで定義したwinnameの指定。窓面の影を考慮する場合に使用する。
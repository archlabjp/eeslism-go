# WALL 壁体構成

書式:
```
WALL
    [wbmlist=<filename>]
    <ewname>
    [\<C\>]
    [ Ei=<ei> ] 
    [ Eo=<eo> ] 
    [ as=<as> ] 
    [ type=<type> ] 
    [ tra=<tra> ]
    [ Ksu=<ksu> ]
    [ Ksd=<ksd> ]
    [ Ru=<ru> ]
    [ Rd=<rd> ]
    [ fcu=<fcu> ]
    [ fcd=<fcd> ]
    [ Kc=<kc> ]
    [ Esu=<esu> ]
    [ Esd=<esd> ]
    [ Eg=<eg> ]
    [ Eb=<eb> ]
    [ ag=<ag> ]
    [ ta=<ta> ]
    [ tnxt=<tnxt>]
    [ t=<t> ]
    [ KHD=<khd> ]
    [ KPD=<kpd> ]
    [ KPM=<kpm> ]
    [ KPA=<kpa> ]
    [ EffInv=<effinv> ]
    [ apmax=<apmax> ]
    [ ap=<ap>]
    [ Rcoloff=<rcoloff> ]
    <layer> <layer> ･･･････････
    （外表面、隣室側、室内側） 屋根、天井のとき
    （室内側、外表面、隣室側） 部位無指定、外壁、内壁、床のとき
    ;
    <ewname> ･････ ; 以下繰り返し
    .
    .
*
```

- filename: 壁体の材料定義リストを指定します。省略すると、[wbmlist.efl](wbmlist.md)が使用されます。
- ewname: 部位・壁体名。`-<ble>[:<wname>]`
  - ble: 部位。`E`, `R`, `F`, `i`, `c`, `f` or `R`
    - E: 外壁
    - R: 屋根
    - F: 床(外部)
    - i: 内壁
    - c: 天井(内部)
    - f: 床(内部)
    - R: 屋根一体型集熱器および太陽電池
  - wname: 壁体名。ただし、最初の1文字は英字とする。省略された場合は、外壁の既定値と見なす。
- ei: 室内表面放射率 0.9
- eo: 外表面放射率 0.9
- as: 外表面日射吸収率 0.7
- type: 集熱器のタイプ。 A1 or A2 or A2P or A3 or W2
- tra: τα
- ksu: 通気層内上側から屋外までの熱貫流率 [W/m2K]
- ksd: 通気層内下側から裏面までの熱貫流率 [W/m2K]
- ru: 通気層から上面までの熱抵抗 [m2K/W]
- rd: 通気層から裏面までの熱抵抗 [m2K/W]
- fcu: Kcu / Ksu
- fcd: Kcd / Ksd
- kc: Kcu + Kcd
- esu: 通気層内上側の放射率
- esd: 通気層内下側の放射率
- eg: 透過体の中空層側表面の放射率
- eb: 集熱版の中空層側表面の放射率
- ag: 透過体の日射吸収率
- ta: 中空層の厚さ [mm]
- tnxt:
- t: 通気層の厚さ [mm]
- KHD:　日射量年変動補正係数（安全率）(1.0)
- KPD:　経時変化補正係数(0.95)
- KPM:　アレイ負荷整合補正係数(0.94)
- KPA:　アレイ回路補正係数(0.97)
- EffInv:　インバータ実効効率(0.9)
- apmax:　最大出力温度係数 [%/℃] (結晶系のとき-0.41)
- ap: 太陽電池裏面の対流熱伝達率
- rcoloff: 通気停止時の集熱器裏面からアレイまでの熱抵抗 [m2K/W] または　集熱停止時の太陽電池から集熱器裏面までの熱抵抗[m2K/W]
- layer: `<code>[-<L>[/<ND>]]` or `P[:<yyy>]` の何れか。
  - code: 材料コード。[材料定義リスト](wbmlist.md)の「材料名」に相当する。ただし、内表面熱伝達率="ali"、外表面伝達率="alo"とする。
  - L: 厚み[mm]
  - ND: 分割総数
  - P: 放射暖冷房パネルの発熱面位置の指定
  - yyy: パネルの熱通過有効度(未指定の場合:0.7)

例1:
```
WALL
    -E:Ewall RC-150 RWL-50 as1 ALM-2 ;
*
```

例2:
```
WALL
    ! 空気集熱器
    -R:Panel <C> type=A1 FPS-120 WDB-22 tra=0.8 Eg=0.9 Eb=0.9 ta=40 ag=0.05 Rd=1.0 t=50 Esu=0.9 Esd=0.9 ;
    ! 予備集熱（鋼板のみなのでRu=0とした）
    -R:Panel2 <C> type=A2 FPS-120 WDB-22 tra=0.9 Ru=0.0 Rd=1.0 t=50 Esu=0.9 Esd=0.9 ;
*
```
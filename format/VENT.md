# VENT(RAICH) 外気導入量および室間相互換気量の設定スケジュ－ル入力


```
VENT
    <Rmname> [ Vent=(vent,vschname) ] [ Inf=(inf,ischname)] ;
    .... 繰り返し
*
```
- Rmname: 室名 ([ROOMデータセット](ROOM.md)で定義)
- vent: 換気量基準値 [kg/s]
- vschname: 換気量設定名 ([SCHNMデータセット](SCHNM.md)で定義)
- inf: 隙間風邪基準値[kg/s]
- ischname: 隙間風設定名 ([SCHNMデータセット](SCHNM.md)で定義)


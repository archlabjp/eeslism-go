# WINDOW 窓リスト

```
WINDOW
        <Fname> 窓名
        [ t=<Ttn> ] 日射透過率
        [ B=<Bn> ] Bn吸収日射取得率
        R=<R> 窓部材熱抵抗 [m2K/W]
        [ Ei=<Ei> ] 室内表面放射率(0.9)
        [ Eo=<Eo> ] 外表面放射率(0.9)
    ;

*
```
- Fname: 窓名
- Ttn: Ttn 日射透過率
- Bn: Bn吸収日射取得率
- R: 窓部材熱抵抗 [m2K/W]
- Ei: 室内表面放射率(0.9)
- Eo: 外表面放射率(0.9)

例:
```
WINDOW
    C6mm t=0.79 B=0.04 R=0.0 ;
*
```
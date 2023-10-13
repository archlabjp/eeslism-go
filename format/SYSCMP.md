# SYSCMP

```
SYSCMP
    {elmname|rmname}
    -c <catname>
    [ -L <L> ]
    [ -env <envtemp> | -room <roomname> ]
    [ -Tinit <itemp> | -Tinit (<item> <itemp> ... <itemp>)]
    [ -exs <exsname> ]
    [ -roomheff <roomname> <roomheff> ]
    [ Ac=<Ac> ]
    [ type=<type> ]
    -Nin <Nin>
    -Nout <Nout>
    ;
    以下繰り返し
*
```
- elmname: 要素名
- rmname: 室名
- catname: 機器カタログ名([EQPCAT](EQPCAT.md)データに定義した名前を引用)
- envtemp: 周辺温度。定数または...
- roomname: 機器設置空間の名前
- itemp: 蓄熱槽の初期水温
- exsname: 方位名([EXSRF](EXSRF.md)で定義した名前を引用)
- roomheff
- Ac
- type:
  - `B`: 通過流体が水の分岐要素
  - `BA`: 通過流体が空気の分岐要素
  - `C`: 通過流体が水の合流要素
  - `CA`: 通過流体が空気の合流要素
  - `V`: 弁およびダンパー
  - `VT`: 温調弁（水系統のみ）
  - `HCLD`: 仮想空調機コイル(直膨コイル)
  - `HCLDW`: 仮想空調機コイル(冷・温水コイル)
  - `RMAC`: ルームエアコン
  - `QMEAS`: カロリーメータ
  - `FLI`: システム経路への流入条件
- Nin: 合流数
- Nout: 分岐数


例1: 室への流入経路
```
SYSCMP
    LD		-Nin	3 ;
*
```
# AAPL 照明・機器仕様スケジュール

```
AAPL
    <rmname>
        [ L=(l,type,lschname) ]
        [ As=(xxx,xxx,aschname)]
        [ Al=(xxx,alschname)]
        [ AE=(xxx,aeschname)]
        [ AG=(xxx,agschname)]
    ;
    <rmname> .... ;
    繰り返し
*
```
- rmname: 室名  ([ROOMデータセット](ROOM.md)で定義)
- l: 照明の基準値 [W]
- type: 照明器具タイプ (予約: `x` を入れること)
- lschname: 照明入力設定値名
- xxx: 機器顕熱の対流成分基準値 [W]
- xxx: 機器顕熱の放射成分基準値 [W]
- aschname: 機器顕熱の設定値名
- xxx: 機器潜熱の基準値 [W]
- alschname: 機器潜熱の設定値名
- xxx: 電力に関する集計の基準値 [W]
- aeschname: 電力に関する集計の電力設定値名
- xxx: ガスに関する集計の基準値 [W]
- agschname: ガスに関する集計のガス設定値名

設定値は([SCHNMデータセット](SCHNM.md)で定義

AE、AGは電力およびガスの使用量であるが、L, As, Alで入力した照明、機器発熱との関連はない。照明、機器で使用したエネルギー量は、出力指定することにより日積算値、月積算値などが得られるので、AE, AGは特に必要ない場合が多い。L, As, AlおよびAE, AGの日積算値、月積算値の集計結果はoutfile__rq.es , outfile_dqr.eaファイルに出力される。
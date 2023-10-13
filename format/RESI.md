# RESI 居住者スケジュール

```
RESI <Rmname> [ H=(h,schnamea,hschnameb) ] [ comfrt=(mschenamem,cschnamec,wschnamev)] ;
    .... 繰り返し
*
```
- Rmname: 室名  ([ROOMデータセット](ROOM.md)で定義)
- h: 人体発熱の基準人数 [人]
- schnamea: 人体発熱の在室設定値名(基準人数に対する比率)  
- hschnameb: 人体発熱の作業強度設定値名
- mschenamem: 熱環境条件の代謝率設定値名 [met]
- cshnamec: 熱環境条件の着衣量設定名 [clo] 
- wschnamev: 熱環境条件の室内風速設定値名 [m/s]

設定値は([SCHNMデータセット](SCHNM.md)で定義
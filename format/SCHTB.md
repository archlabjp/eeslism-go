# SCHTB 一日の設定値、切換スケジュールおよび季節、曜日の指定


## 定義の種類
- -v or VL: 設定値スケジュール定義(引用)
- -s or SW: 切替設定スケジュール定義(引用)
- -ssn or SSN: 季節設定
- -wkd or WKD: 曜日設定

NOTE: 設定値および切り替え設定は他から引用されます。

### 設定値スケジュール定義

```
-v vdname ttmm-(xxx)-ttmm ttmm-(xxx)-ttmm ... ;
```
- vdname: 設定値名
- ttmm: 開始または終了の時分
- xxx: 設定値

例: `-v Temp 800-(20.0)-2300 ; ` 設定値`Temp`で8時から23時まで20℃に設定する

### 切替設定スケジュール定義

```
-s wdname ttmm-(mode)-ttmm ttmm-(mode)-ttmm ... ;
```
- wdname: 切替設定名
- ttmm: 開始または終了の時分
- mode: 切替名

例: `-s WSCH 001-(N)-800 801-(D)-1700 1701-(N)-2400 ;` 0時1分から8時まではモードN、17時まではモードD、それ以降はモードN

### 季節設定

`mm`月`dd日`から`mm`月`dd日`を季節 `sname` と定義する。
```
-ssn sname mm/dd-mm/dd mm/dd-mm/dd ... ;
```
- sname: 季節名
- mm: 月
- dd: 日

例: `-ssn	Winter		11/4-4/21 ;`

### 曜日設定

`wday`(複数)に `wname` という曜日設定名を定義する。
```
-wkd wname wday wday wday ... wday ; 
```
- wname: 曜日名
- wday: 曜日を表す記号(Hol Sun Mon Tue Wed Thu Fri Sat)

例: `-wkd Weekend Sat Sun Hol ;` 

NOTE: 曜日および祝日の設定は[dayweek.efl](dayweek.md)または[WEEK](WEEK/md)に基づく。

## 他のデータセット中でのスケジュールデータ定義

先頭に "%s" を付け加えると、入力データ中の任意の場所で定義をすることができる。この場合は、データセット名 `SCHTB` のくくりは不要である。

## 例1 SCHTBデータセットの中で
```
SCHTB
	-s    WSCH     001-(N)-800 801-(D)-1700 1701-(N)-2400 ;
	-wkd  Weekday  Mon Tue Wed Thu Fri ;
	-wkd  Weekend  Sat Sun Hol ;
	-ssn  Winter   11/4-4/21 ;
	-ssn  Summer   5/30-9/23 ;
	-ssn  Inter	   4/22-5/29  9/24-11/3 ;
*
```


## 例2 任意の場所で
```
	%s -s    WSCH     001-(N)-800 801-(D)-1700 1701-(N)-2400 ;
	%s -wkd  Weekday  Mon Tue Wed Thu Fri ;
	%s -wkd  Weekend  Sat Sun Hol ;
	%s -ssn	Winter		11/4-4/21 ;
	%s -ssn	Summer		5/30-9/23 ;
	%s -ssn	Inter		4/22-5/29  9/24-11/3 ;
```
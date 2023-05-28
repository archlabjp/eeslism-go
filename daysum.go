package main

type SVDAY struct {
	M      float64 // 平均
	Mn     float64 // 最高
	Mx     float64 // 最低
	Hrs    int64   // 平均値の母数
	Mntime int64   // 最高値発生時刻
	Mxtime int64   // 最低値発生時刻
}

type QDAY struct {
	H       float64 // 加熱積算値
	C       float64 // 冷却積算値
	Hmx     float64 // 加熱最大値
	Cmx     float64 // 冷却最大値
	Hhr     int64   // 加熱時間回数
	Chr     int64   // 冷却時間回数
	Hmxtime int64
	Cmxtime int64
}

type EDAY struct {
	D      float64 // 積算値
	Mx     float64 // 最大値
	Hrs    int64   // 運転時間回数
	Mxtime int64
}

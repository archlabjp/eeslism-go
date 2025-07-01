package eeslism

const MAXINT_DAY = -999.0
const MININT_DAY = 999.0

/*
svdyint (Daily Summary Initialization for Scalar Values)

この関数は、日単位で集計されるスカラー値（温度、湿度など）のデータ構造（`SVDAY`）をリセットします。
これは、新しい日の集計を開始する前に、
前日のデータをクリアするために用いられます。

建築環境工学的な観点:
- **日単位の集計の準備**: 建物のエネルギー消費量や室内環境を日単位で評価するためには、
  各日の開始時に集計値をゼロにリセットする必要があります。
  この関数は、平均値（`M`）、最高値（`Mx`）、最低値（`Mn`）、
  およびそれらの発生時刻（`Mxtime`, `Mntime`）を初期化します。
- **正確なデータ分析の確保**: 日積算値が適切にリセットされることで、
  日ごとのエネルギー消費量や室内環境を正確に比較分析することが可能になります。
  これにより、特定の日のエネルギー消費が多かった原因を特定したり、
  省エネルギー対策の効果を日単位で評価したりする際の信頼性が向上します。

この関数は、建物のエネルギー消費量や室内環境を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func svdyint(vd *SVDAY) {
	vd.M = 0.0
	vd.Mn = MININT_DAY
	vd.Mx = MAXINT_DAY
	vd.Hrs = 0
	vd.Mntime = -1
	vd.Mxtime = -1
}

/*
svdaysum (Daily Summary Accumulation for Scalar Values)

この関数は、日単位で集計されるスカラー値（温度、湿度など）のデータ構造（`SVDAY`）に、
現在の時刻の値を加算し、平均値、最高値、最低値を更新します。

建築環境工学的な観点:
- **日単位のデータ集計**: シミュレーションの各時間ステップで計算された温度や湿度などの値を、
  日単位で集計します。
  `vd.M += v` で合計値を、`vd.Hrs++` でデータ数をカウントし、
  最終的に日平均値を算出します。
- **最高値・最低値の記録**: `minmark`と`maxmark`関数を呼び出すことで、
  日中の最高値と最低値、およびそれらの発生時刻を記録します。
  これにより、日中の温湿度変動の幅や、
  快適性が最も悪化する時間帯などを把握できます。
- **運転状態の考慮**: `if control != '0'` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。

この関数は、建物のエネルギー消費量や室内環境を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func svdaysum(time int64, control ControlSWType, v float64, vd *SVDAY) {
	if control != '0' {
		vd.M += v
		vd.Hrs++
		minmark(&vd.Mn, &vd.Mntime, v, time)
		maxmark(&vd.Mx, &vd.Mxtime, v, time)
	}
	if time == 2400 && vd.Hrs > 0 {
		vd.M /= float64(vd.Hrs)
	}
}

/*
svmonsum (Monthly Summary Accumulation for Scalar Values)

この関数は、月単位で集計されるスカラー値（温度、湿度など）のデータ構造（`SVDAY`）に、
現在の時刻の値を加算し、平均値、最高値、最低値を更新します。

建築環境工学的な観点:
- **月単位のデータ集計**: シミュレーションの各時間ステップで計算された温度や湿度などの値を、
  月単位で集計します。
  `vd.M += v` で合計値を、`vd.Hrs++` でデータ数をカウントし、
  最終的に月平均値を算出します。
- **最高値・最低値の記録**: `minmark`と`maxmark`関数を呼び出すことで、
  月中の最高値と最低値、およびそれらの発生時刻を記録します。
  これにより、月中の温湿度変動の幅や、
  快適性が最も悪化する期間などを把握できます。
- **運転状態の考慮**: `if control != '0'` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。
- **月末判定**: `IsEndDay`関数を用いて月末を判定し、
  月末であれば月平均値を計算します。

この関数は、建物のエネルギー消費量や室内環境を月単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func svmonsum(Mon int, Day int, time int, control ControlSWType, v float64, vd *SVDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if control != '0' {
		vd.M += v
		vd.Hrs++
		minmark(&vd.Mn, &vd.Mntime, v, MoNdTt)
		maxmark(&vd.Mx, &vd.Mxtime, v, MoNdTt)
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && vd.Hrs > 0 && time == 2400 {
		vd.M /= float64(vd.Hrs)
	}
}

/*
qdyint (Daily Summary Initialization for Heat/Cooling Quantities)

この関数は、日単位で集計される熱量（加熱、冷却）のデータ構造（`QDAY`）をリセットします。
これは、新しい日の集計を開始する前に、
前日のデータをクリアするために用いられます。

建築環境工学的な観点:
- **日単位の熱量集計の準備**: 建物の熱負荷を日単位で評価するためには、
  各日の開始時に集計値をゼロにリセットする必要があります。
  この関数は、加熱積算値（`H`）、冷却積算値（`C`）、
  加熱最大値（`Hmx`）、冷却最大値（`Cmx`）、
  およびそれらの発生時刻（`Hmxtime`, `Cmxtime`）を初期化します。
- **正確なデータ分析の確保**: 日積算値が適切にリセットされることで、
  日ごとの熱負荷を正確に比較分析することが可能になります。
  これにより、特定の日の熱負荷が大きかった原因を特定したり、
  空調システムの運転状況を評価したりする際の信頼性が向上します。

この関数は、建物の熱負荷を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func qdyint(Qd *QDAY) {
	Qd.H = 0.0
	Qd.C = 0.0
	Qd.Hmx = 0.0
	Qd.Cmx = 0.0
	Qd.Hhr = 0
	Qd.Chr = 0
	Qd.Hmxtime = -1
	Qd.Cmxtime = -1
}

/*
qdaysum (Daily Summary Accumulation for Heat/Cooling Quantities)

この関数は、日単位で集計される熱量（加熱、冷却）のデータ構造（`QDAY`）に、
現在の時刻の値を加算し、最大値を更新します。

建築環境工学的な観点:
- **日単位の熱量集計**: シミュレーションの各時間ステップで計算された熱量を、
  日単位で集計します。
  `Qd.H += Q` で加熱熱量を、`Qd.C += Q` で冷却熱量を合計します。
- **最大値の記録**: `maxmark`と`minmark`関数を呼び出すことで、
  日中の加熱最大値と冷却最大値、およびそれらの発生時刻を記録します。
  これにより、日中のピーク負荷の大きさや、
  空調設備の容量設計に必要な情報を把握できます。
- **運転状態の考慮**: `if control != '0'` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。
- **単位変換**: `time == 2400`（一日の終わり）の場合に、
  熱量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物の熱負荷を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func qdaysum(time int64, control ControlSWType, Q float64, Qd *QDAY) {
	if control != '0' {
		if Q > 0.0 {
			Qd.H += Q
			maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, time)
			Qd.Hhr++
		} else if Q < 0.0 {
			Qd.C += Q
			minmark(&Qd.Cmx, &Qd.Cmxtime, Q, time)
			Qd.Chr++
		}
	}

	if time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

/*
qmonsum (Monthly Summary Accumulation for Heat/Cooling Quantities)

この関数は、月単位で集計される熱量（加熱、冷却）のデータ構造（`QDAY`）に、
現在の時刻の値を加算し、最大値を更新します。

建築環境工学的な観点:
- **月単位の熱量集計**: シミュレーションの各時間ステップで計算された熱量を、
  月単位で集計します。
  `Qd.H += Q` で加熱熱量を、`Qd.C += Q` で冷却熱量を合計します。
- **最大値の記録**: `maxmark`と`minmark`関数を呼び出すことで、
  月中の加熱最大値と冷却最大値、およびそれらの発生時刻を記録します。
  これにより、月中のピーク負荷の大きさや、
  空調設備の容量設計に必要な情報を把握できます。
- **運転状態の考慮**: `if control != '0'` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。
- **単位変換**: `IsEndDay`関数を用いて月末を判定し、
  月末であれば熱量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物の熱負荷を月単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func qmonsum(Mon int, Day int, time int, control ControlSWType, Q float64, Qd *QDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if control != '0' {
		if Q > 0.0 {
			Qd.H += Q
			maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, MoNdTt)
			Qd.Hhr++
		} else if Q < 0.0 {
			Qd.C += Q
			minmark(&Qd.Cmx, &Qd.Cmxtime, Q, MoNdTt)
			Qd.Chr++
		}
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

/*
qdaysumNotOpe (Daily Summary Accumulation for Heat/Cooling Quantities, Including Non-Operating Time)

この関数は、日単位で集計される熱量（加熱、冷却）のデータ構造（`QDAY`）に、
現在の時刻の値を加算し、最大値を更新します。
この関数は、機器が運転していない時間帯のデータも集計対象とします。

建築環境工学的な観点:
- **日単位の熱量集計（非運転時含む）**: シミュレーションの各時間ステップで計算された熱量を、
  日単位で集計します。
  `Qd.H += Q` で加熱熱量を、`Qd.C += Q` で冷却熱量を合計します。
- **最大値の記録**: `maxmark`と`minmark`関数を呼び出すことで、
  日中の加熱最大値と冷却最大値、およびそれらの発生時刻を記録します。
  これにより、日中のピーク負荷の大きさや、
  空調設備の容量設計に必要な情報を把握できます。
- **運転状態によらない集計**: `control`パラメータのチェックがないため、
  機器の運転状態に関わらず、常に熱量を集計します。
  これは、例えば、太陽光発電の発電量のように、
  機器の運転状態とは独立して発生する熱量（エネルギー）を集計する際に有用です。
- **単位変換**: `time == 2400`（一日の終わり）の場合に、
  熱量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物の熱負荷を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func qdaysumNotOpe(time int64, Q float64, Qd *QDAY) {
	if Q > 0.0 {
		Qd.H += Q
		maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, time)
		Qd.Hhr++
	} else if Q < 0.0 {
		Qd.C += Q
		minmark(&Qd.Cmx, &Qd.Cmxtime, Q, time)
		Qd.Chr++
	}

	if time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

/*
qmonsumNotOpe (Monthly Summary Accumulation for Heat/Cooling Quantities, Including Non-Operating Time)

この関数は、月単位で集計される熱量（加熱、冷却）のデータ構造（`QDAY`）に、
現在の時刻の値を加算し、最大値を更新します。
この関数は、機器が運転していない時間帯のデータも集計対象とします。

建築環境工学的な観点:
- **月単位の熱量集計（非運転時含む）**: シミュレーションの各時間ステップで計算された熱量を、
  月単位で集計します。
  `Qd.H += Q` で加熱熱量を、`Qd.C += Q` で冷却熱量を合計します。
- **最大値の記録**: `maxmark`と`minmark`関数を呼び出すことで、
  月中の加熱最大値と冷却最大値、およびそれらの発生時刻を記録します。
  これにより、月中のピーク負荷の大きさや、
  空調設備の容量設計に必要な情報を把握できます。
- **運転状態によらない集計**: `control`パラメータのチェックがないため、
  機器の運転状態に関わらず、常に熱量を集計します。
  これは、例えば、太陽光発電の発電量のように、
  機器の運転状態とは独立して発生する熱量（エネルギー）を集計する際に有用です。
- **単位変換**: `IsEndDay`関数を用いて月末を判定し、
  月末であれば熱量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物の熱負荷を月単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func qmonsumNotOpe(Mon int, Day int, time int, Q float64, Qd *QDAY, Dayend int, SimDayend int) {
	MoNdTt := int64(1000000*Mon + 10000*Day + time)

	if Q > 0.0 {
		Qd.H += Q
		maxmark(&Qd.Hmx, &Qd.Hmxtime, Q, MoNdTt)
		Qd.Hhr++
	} else if Q < 0.0 {
		Qd.C += Q
		minmark(&Qd.Cmx, &Qd.Cmxtime, Q, MoNdTt)
		Qd.Chr++
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Qd.H *= Cff_kWh
		Qd.C *= Cff_kWh
	}
}

/*
edyint (Daily Summary Initialization for Energy Quantities)

この関数は、日単位で集計されるエネルギー量（電力消費量など）のデータ構造（`EDAY`）をリセットします。
これは、新しい日の集計を開始する前に、
前日のデータをクリアするために用いられます。

建築環境工学的な観点:
- **日単位のエネルギー集計の準備**: 建物のエネルギー消費量を日単位で評価するためには、
  各日の開始時に集計値をゼロにリセットする必要があります。
  この関数は、積算値（`D`）、最大値（`Mx`）、
  および運転時間回数（`Hrs`）を初期化します。
- **正確なデータ分析の確保**: 日積算値が適切にリセットされることで、
  日ごとのエネルギー消費量を正確に比較分析することが可能になります。
  これにより、特定の日のエネルギー消費が多かった原因を特定したり、
  省エネルギー対策の効果を日単位で評価したりする際の信頼性が向上します。

この関数は、建物のエネルギー消費量や室内環境を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための基礎的な役割を果たします。
*/
func edyint(Ed *EDAY) {
	Ed.D = 0.0
	Ed.Mx = 0.0
	Ed.Hrs = 0
	Ed.Mxtime = -1
}

/*
edaysum (Daily Summary Accumulation for Energy Quantities)

この関数は、日単位で集計されるエネルギー量（電力消費量など）のデータ構造（`EDAY`）に、
現在の時刻の値を加算し、最大値を更新します。

建築環境工学的な観点:
- **日単位のエネルギー集計**: シミュレーションの各時間ステップで計算されたエネルギー量を、
  日単位で集計します。
  `Ed.D += E` で合計値を、`Ed.Hrs++` で運転時間回数をカウントします。
- **最大値の記録**: `maxmark`関数を呼び出すことで、
  日中の最大エネルギー消費量と、その発生時刻を記録します。
  これにより、日中のピーク電力需要の大きさや、
  電力系統への影響を把握できます。
- **運転状態の考慮**: `if control != '0'` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。
- **単位変換**: `time == 2400`（一日の終わり）の場合に、
  エネルギー量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物のエネルギー消費量を日単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func edaysum(time int, control ControlSWType, E float64, Ed *EDAY) {
	if control != '0' {
		Ed.D += E
		maxmark(&Ed.Mx, &Ed.Mxtime, E, int64(time))
		Ed.Hrs++
	}

	if time == 2400 {
		Ed.D *= Cff_kWh
	}
}

/*
emonsum (Monthly Summary Accumulation for Energy Quantities)

この関数は、月単位で集計されるエネルギー量（電力消費量など）のデータ構造（`EDAY`）に、
現在の時刻の値を加算し、最大値を更新します。

建築環境工学的な観点:
- **月単位のエネルギー集計**: シミュレーションの各時間ステップで計算されたエネルギー量を、
  月単位で集計します。
  `Ed.D += E` で合計値を、`Ed.Hrs++` で運転時間回数をカウントします。
- **最大値の記録**: `maxmark`関数を呼び出すことで、
  月中の最大エネルギー消費量と、その発生時刻を記録します。
  これにより、月中のピーク電力需要の大きさや、
  電力系統への影響を把握できます。
- **運転状態の考慮**: `if control != OFF_SW` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。
- **単位変換**: `IsEndDay`関数を用いて月末を判定し、
  月末であればエネルギー量を`Cff_kWh`（kWhへの変換係数）で乗じることで、
  最終的な出力単位をkWhに変換します。

この関数は、建物のエネルギー消費量を月単位で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func emonsum(Mon, Day, time int, control ControlSWType, E float64, Ed *EDAY, Dayend, SimDayend int) {
	var MoNdTt int64 = int64(1000000*Mon + 10000*Day + time)

	if control != OFF_SW {
		Ed.D += E
		maxmark(&Ed.Mx, &Ed.Mxtime, E, MoNdTt)
		Ed.Hrs++
	}

	if IsEndDay(Mon, Day, Dayend, SimDayend) && time == 2400 {
		Ed.D *= Cff_kWh
	}
}

/*
emtsum (Monthly-Time-of-Day Summary Accumulation for Energy Quantities)

この関数は、月・時刻別で集計されるエネルギー量（電力消費量など）のデータ構造（`EDAY`）に、
現在の時刻の値を加算します。

建築環境工学的な観点:
- **月・時刻別のエネルギー集計**: シミュレーションの各時間ステップで計算されたエネルギー量を、
  月と時刻の組み合わせで集計します。
  これにより、特定の月における時間帯ごとのエネルギー消費量の傾向を把握できます。
- **デマンドサイドマネジメント**: 月・時刻別のエネルギー消費量データは、
  デマンドサイドマネジメント（DSM）戦略を策定する上で非常に有用です。
  例えば、ピーク時間帯の電力消費量を削減するための運転戦略を検討したり、
  蓄熱システムや再生可能エネルギーの導入効果を評価したりする際に役立ちます。
- **運転状態の考慮**: `if control != OFF_SW` の条件は、
  機器が運転中のデータのみを集計対象とすることを示唆します。
  これにより、機器の実際の運転状況を反映した集計が可能になります。

この関数は、建物のエネルギー消費量を月・時刻別で詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func emtsum(Mon, Day, time int, control ControlSWType, E float64, Ed *EDAY) {
	if control != OFF_SW {
		Ed.D += E
	}
}

/*
minmark (Minimum Value Marking)

この関数は、与えられた値`v`が現在の最小値`minval`よりも小さい場合に、
最小値を更新し、その時刻を記録します。

建築環境工学的な観点:
- **最低値の記録**: シミュレーションでは、
  室温、湿度、熱負荷などの最低値を記録することが重要です。
  これにより、
  - **快適性評価**: 室内環境が最も不快になる時間帯や、
    設定温度を下回る時間帯を特定できます。
  - **熱負荷の評価**: 暖房負荷のピーク時など、
    熱需要が最も高まる時間帯を把握できます。
- **時刻の記録**: `timemin`に最低値が発生した時刻を記録することで、
  その現象が発生した具体的な時間帯を特定し、
  原因分析や対策検討に役立てることができます。

この関数は、建物のエネルギー消費量や室内環境を詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func minmark(minval *float64, timemin *int64, v float64, time int64) {
	if v <= *minval {
		*timemin = time
		*minval = v
	}
}

/*
maxmark (Maximum Value Marking)

この関数は、与えられた値`v`が現在の最大値`maxval`よりも大きい場合に、
最大値を更新し、その時刻を記録します。

建築環境工学的な観点:
- **最高値の記録**: シミュレーションでは、
  室温、湿度、熱負荷などの最高値を記録することが重要です。
  これにより、
  - **快適性評価**: 室内環境が最も不快になる時間帯や、
    設定温度を上回る時間帯を特定できます。
  - **熱負荷の評価**: 冷房負荷のピーク時など、
    熱需要が最も高まる時間帯を把握できます。
- **時刻の記録**: `timemax`に最高値が発生した時刻を記録することで、
  その現象が発生した具体的な時間帯を特定し、
  原因分析や対策検討に役立てることができます。

この関数は、建物のエネルギー消費量や室内環境を詳細に分析し、
運用改善や省エネルギー対策の効果評価を行うための重要な役割を果たします。
*/
func maxmark(maxval *float64, timemax *int64, v float64, time int64) {
	if v >= *maxval {
		*timemax = time
		*maxval = v
	}
}

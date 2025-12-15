// reset_state.go - グローバル状態変数のリセット
// テスト間の独立性を確保するために使用されます。

package eeslism

// resetPrintStates は印刷・日集計に関するグローバル状態変数をリセットします。
func resetPrintStates() {
	// spline.go
	__Intgtsup_ic = 0

	// blroomday.go
	__Roomday_oldday = INAN
	__Roomday_oldMon = INAN
	__Rmdyprint_id = 0
	__Rmmonprint_id = 0

	// wdprint.go
	__Wdtsum_oldday = 0
	__Wdtsum_oldMon = 0
	__Wdtsum_hrs = 0
	__Wdtsum_hrsm = 0
	__Wdtsum_cffWh = 0.0

	// blrmprint.go
	__Rmsfprint_ic = 0
	__Rmsfqprint_ic = 0
	__Rmsfaprint_ic = 0
	__Dysfprint_ic = 0
	__Shdprint_ic = 0
	__Wallprint_ic = 0
	__PCMprint_ic = 0
	__Qrmprint_ic = 0
	__Dyqrmprint_ic = 0
	__Qrmsum_oldday = 0

	// eeprint_s.go
	__Eeprintd_ic = 0
	__Wdtdprint_ic = 0
	__Wdtprint_ic = 0
	__Wdtmprint_ic = 0

	// blrmaceqcf.go
	__Rmhtrcf_count = 0

	// eecmpday_s.go
	__Compoday_OldDay = 0
	__Compoday_OldMon = 0
	__Compodyprt_id = 0
	__Compomonprt_id = 0
	__Compomtprt_id = 0

	// blrmqprt.go
	__Rmpnlprint_id = 0

	// eevcdat.go
	__Vcfinput_Mon = 0
	__Vcfinput_Day = 0
	__Vcfinput_Time = 0

	// eepthprt.go
	__Pathprint_id = 0

	// mcrefas.go
	__Refadata_hpch = nil

	// eecmpprt_s.go
	__Hcmpprint_id = 0
	__Hstkprint_id = 0

	// u_sun.go
	__Solpos_Sdecl = 0.0
	__Solpos_Sld = 0.0
	__Solpos_Cld = 0.0
	__Solpos_Ttprev = 25.0

	// bltcomfrt.go
	__Rmotset_Pint = 0
	__Fotinit_init = 'i'

	// blhelm.go
	__Helmprint_id = 0
	__Helmsurfprint_id = 0
	__Helmdy_oldday = -1
	__Helmdyprint_id = 0

	// wdread.go
	__Weatherdt_ptt = 25
	__Weatherdt_nc = 0
	__Weatherdt_decl = 0.0
	__Weatherdt_E = 0.0
	__Weatherdt_tas = 0.0
	__Weatherdt_timedg = 0.0
	// __Weatherdt_dt, __Weatherdt_dtL は配列なのでゼロクリアが必要
	for i := range __Weatherdt_dt {
		for j := range __Weatherdt_dt[i] {
			__Weatherdt_dt[i][j] = 0.0
			__Weatherdt_dtL[i][j] = 0.0
		}
	}
	__gtsupw_ic = 0
	__hspwdread_ic = 0
	__hspwdread_recl = 0

	// blsrprint.go
	__Pmvprint_count = 0
	__Rmevprint_count = 0
}

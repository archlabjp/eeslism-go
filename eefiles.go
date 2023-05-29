package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

/* ----------------------------------------------------- */

// ファイルのオープン
func (Simc *SIMCONTL) eeflopen(Nflout int, Flout []*FLOUT) {
	// 気象データファイルを開く
	if Simc.Wdtype == 'H' {
		var err error
		Simc.Fwdata, err = os.Open(Simc.Wfname)
		if err != nil {
			Eprint("<eeflopen>", Simc.Wfname)
			os.Exit(EXIT_WFILE)
		}
		Simc.Fwdata2, err = os.Open(Simc.Wfname)
		if err != nil {
			Eprint("<eeflopen>", Simc.Wfname)
			os.Exit(EXIT_WFILE)
		}
		Simc.Ftsupw, err = ioutil.ReadFile("supw.efl")
		if err != nil {
			Eprint("<eeflopen>", "supw.efl")
			os.Exit(EXIT_SUPW)
		}
	}

	// Fname = Simc.ofname + ".log"
	// ferr, err := os.Create(Fname)
	// if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	// }

	// 出力ファイルを開く
	for i := 0; i < Nflout; i++ {
		fl := Flout[i]
		Fname := Simc.Ofname + string(fl.Idn) + ".es"

		var err error
		fl.F, err = os.Create(Fname)
		if err != nil {
			Eprint("<eeflopen>", Fname)
			os.Exit(EXIT_WFILE)
		}
	}
}

/* ----------------------------------------------------- */

func Eeflclose(Nflout int, Flout []*FLOUT) {
	var fl *FLOUT

	if Ferr != nil {
		Ferr.Close()
	}

	for i := 0; i < Nflout; i++ {
		fl = Flout[i]
		fmt.Fprintln(fl.F, "-999")
		// fl.F.Sync()
		// fl.F.Close()
	}
}

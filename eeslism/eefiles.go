package eeslism

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/* ----------------------------------------------------- */

// ファイルのオープン
func (Simc *SIMCONTL) eeflopen(Flout []*FLOUT) {
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
	for _, fl := range Flout {
		fl.Fname = Simc.Ofname + string(fl.Idn) + ".es"
		fl.F = new(strings.Builder)
	}
}

/* ----------------------------------------------------- */

func Eeflclose(Flout []*FLOUT) {
	var fl *FLOUT

	if Ferr != nil {
		Ferr.Close()
	}

	for _, fl = range Flout {
		fo, err := os.Create(fl.Fname)
		if err != nil {
			Eprint("<eeflopen>", fl.Fname)
			os.Exit(EXIT_WFILE)
		}

		fmt.Fprint(fo, fl.F)
		fmt.Fprintln(fo, "-999")
		fo.Close()
	}
}

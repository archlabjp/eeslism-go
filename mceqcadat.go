//This file is part of EESLISM.
//
//Foobar is free software : you can redistribute itand /or modify
//it under the terms of the GNU General Public License as published by
//the Free Software Foundation, either version 3 of the License, or
//(at your option) any later version.
//
//Foobar is distributed in the hope that it will be useful,
//but WITHOUT ANY WARRANTY; without even the implied warranty of
//MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.See the
//GNU General Public License for more details.
//
//You should have received a copy of the GNU General Public License
//along with Foobar.If not, see < https://www.gnu.org/licenses/>.

/*  eqcadat.c  */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

/* ----------------------------------------- */

/*  機器仕様入力       */

func Eqcadata(f *EeTokens, dsn string, Eqcat *EQCAT) {
	var (
		s       string
		ss      string
		ce      int
		E       string
		N       int
		NBOI    int
		NREFA   int
		NCOL    int
		NSTANK  int
		NHCC    int
		NHEX    int
		NPIPE   int
		NPUMP   int
		NVAV    int
		NSTHEAT int
		NTHEX   int
		NPV     int
		NOMVAV  int
		NDESI   int
		NEVAC   int
		frf     *os.File
	)

	NBOI = 0
	NREFA = 0
	NCOL = 0
	NSTANK = 0
	NHCC = 0
	NHEX = 0
	NPIPE = 0
	NPUMP = 0
	NVAV = 0
	NSTHEAT = 0
	NTHEX = 0
	NPV = 0
	NOMVAV = 0
	NDESI = 0
	NEVAC = 0

	Eqpcount(f, &NBOI, &NREFA, &NCOL, &NSTANK, &NHCC, &NHEX,
		&NPIPE, &NPUMP, &NVAV, &NSTHEAT, &NTHEX, &NPV, &NOMVAV, &NDESI, &NEVAC)

	N = NHCC
	Eqcat.Hccca = nil
	if N > 0 {
		Eqcat.Hccca = make([]HCCCA, 0, N)
	}

	N = NBOI
	Eqcat.Boica = nil
	if N > 0 {
		Eqcat.Boica = make([]BOICA, 0, N)
	}

	N = NREFA
	Eqcat.Refaca = nil
	if N > 0 {
		Eqcat.Refaca = make([]REFACA, 0, N)
	}

	N = NCOL
	Eqcat.Collca = nil
	if N > 0 {
		Eqcat.Collca = make([]COLLCA, 0, N)
	}

	N = NPV
	Eqcat.PVca = nil
	if N > 0 {
		Eqcat.PVca = make([]PVCA, 0, N)
	}

	N = NPIPE
	Eqcat.Pipeca = nil
	if N > 0 {
		Eqcat.Pipeca = make([]PIPECA, 0, N)
	}

	N = NSTANK
	Eqcat.Stankca = nil
	if N > 0 {
		Eqcat.Stankca = make([]STANKCA, 0, N)
	}

	N = NHEX
	Eqcat.Hexca = nil
	if N > 0 {
		Eqcat.Hexca = make([]HEXCA, 0, N)
	}

	N = NPUMP
	Eqcat.Pumpca = nil
	if N > 0 {
		Eqcat.Pumpca = make([]PUMPCA, 0, N)
	}

	N = NVAV
	Eqcat.Vavca = nil
	if N > 0 {
		Eqcat.Vavca = make([]VAVCA, 0, N)
	}

	N = NSTHEAT
	Eqcat.Stheatca = nil
	if N > 0 {
		Eqcat.Stheatca = make([]STHEATCA, 0, N)
	}

	N = NTHEX
	Eqcat.Thexca = nil
	if N > 0 {
		Eqcat.Thexca = make([]THEXCA, 0, N)
	}

	N = NOMVAV
	Eqcat.OMvavca = nil
	if N > 0 {
		Eqcat.OMvavca = make([]OMVAVCA, 0, N+1)
	}

	N = NDESI
	Eqcat.Desica = nil
	if N > 0 {
		Eqcat.Desica = make([]DESICA, 0, N+1)
	}

	N = NEVAC
	Eqcat.Evacca = nil
	if N > 0 {
		Eqcat.Evacca = make([]EVACCA, 0, N+1)
	}

	frf, err := os.Open("reflist.efl")
	if err != nil {
		Eprint(" file ", "reflist.efl")
	}

	const RFCMPLSTMX = 5
	N = RFCMPLSTMX
	Eqcat.Rfcmp = nil
	if N > 0 {
		Eqcat.Rfcmp = make([]RFCMP, N)
	} else {
		Rf := Eqcat.Rfcmp
		for i := 0; i < N; i++ {
			Rf[i].cname = ""
			for j := 0; j < 4; j++ {
				Rf[i].d[j] = 0.0
				Rf[i].e[j] = 0.0
				Rf[i].w[j] = 0.0
			}
			Rf[i].Meff = 0.0
			Rf[i].name = ""
			for j := 0; j < 2; j++ {
				Rf[i].Tco[j] = 0.0
				Rf[i].Teo[j] = 0.0
			}
		}
	}

	Refcmpdat(frf, &Eqcat.Nrfcmp, Eqcat.Rfcmp)
	frf.Close()

	frf, err = os.Open("pumpfanlst.efl")
	if err != nil {
		Eprint(" file ", "pumpfanlst.efl")
	}
	N = pflistcount(frf)
	if N > 0 {
		Eqcat.Pfcmp = make([]PFCMP, N)
	}
	PFcmpInit(N, Eqcat.Pfcmp)
	PFcmpdata(frf, &Eqcat.Npfcmp, Eqcat.Pfcmp)
	frf.Close()

	E = fmt.Sprintf(ERRFMT, dsn)

	for f.IsEnd() == false {
		s = f.GetToken()
		if s[0] == '*' {
			break
		}

		eqpType := EqpType(s)

		if eqpType == HCCOIL_TYPE {
			Eqcat.Hccca = append(Eqcat.Hccca, HCCCA{})
			Hccca := &Eqcat.Hccca[len(Eqcat.Hccca)-1]

			Hccca.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Hccdata(s, Hccca) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == BOILER_TYPE {
			Eqcat.Boica = append(Eqcat.Boica, BOICA{})
			Boica := &Eqcat.Boica[len(Eqcat.Boica)-1]

			Boica.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce = strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Boidata(s, Boica) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == COLLECTOR_TYPE || eqpType == ACOLLECTOR_TYPE {
			Eqcat.Collca = append(Eqcat.Collca, COLLCA{})
			Collca := &Eqcat.Collca[len(Eqcat.Collca)-1]

			Collca.name = ""
			Collca.Fd = 0.9
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce = strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if Colldata(eqpType, ss, Collca) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == PV_TYPE {
			Eqcat.PVca = append(Eqcat.PVca, PVCA{})
			PVca := &Eqcat.PVca[len(Eqcat.PVca)-1]

			PVca.Name = ""
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce := strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if PVcadata(ss, PVca) != 0 {
					fmt.Printf("%s %s\n", E, ss)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == REFACOMP_TYPE {
			Eqcat.Refaca = append(Eqcat.Refaca, REFACA{})
			Refaca := &Eqcat.Refaca[len(Eqcat.Refaca)-1]

			Refaca.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexByte(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Refadata(s, Refaca, Eqcat.Nrfcmp, Eqcat.Rfcmp) != 0 {
					fmt.Println(E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == PIPEDUCT_TYPE || s == DUCT_TYPE {
			Eqcat.Pipeca = append(Eqcat.Pipeca, PIPECA{})
			Pipeca := &Eqcat.Pipeca[len(Eqcat.Pipeca)-1]

			Pipeca.name = ""
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce := strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if Pipedata(s, ss, Pipeca) != 0 {
					fmt.Printf("%s %s\n", E, ss)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == STANK_TYPE {
			Eqcat.Stankca = append(Eqcat.Stankca, STANKCA{})
			Stankca := &Eqcat.Stankca[len(Eqcat.Stankca)-1]

			Stankca.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Stankdata(f, s, Stankca) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == HEXCHANGR_TYPE {
			Eqcat.Hexca = append(Eqcat.Hexca, HEXCA{})
			Hexca := &Eqcat.Hexca[len(Eqcat.Hexca)-1]

			Hexca.Name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Hexdata(s, Hexca) != 0 {
					fmt.Printf("%s %s\n", E, s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == PUMP_TYPE || eqpType == FAN_TYPE {
			Eqcat.Pumpca = append(Eqcat.Pumpca, PUMPCA{})
			Pumpca := &Eqcat.Pumpca[len(Eqcat.Pumpca)-1]

			Pumpca.name = ""
			Pumpca.Type = ""
			Pumpca.val = nil
			Pumpca.pfcmp = nil
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce := strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if Pumpdata(s, ss, Pumpca, Eqcat.Npfcmp, Eqcat.Pfcmp) != 0 {
					fmt.Printf("%s %s\n", E, ss)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == VAV_TYPE || eqpType == VWV_TYPE {
			Eqcat.Vavca = append(Eqcat.Vavca, VAVCA{})
			vavca := &Eqcat.Vavca[len(Eqcat.Vavca)-1]

			vavca.dTset = -999.0
			vavca.Name = ""
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce := strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if VAVdata(eqpType, ss, vavca) != 0 {
					Eprint("<Eqcadata> VAV", ss)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == OMVAV_TYPE || eqpType == OAVAV_TYPE {
			Eqcat.OMvavca = append(Eqcat.OMvavca, OMVAVCA{})
			OMvavca := &Eqcat.OMvavca[len(Eqcat.OMvavca)-1]

			OMvavca.Name = ""
			OMvavca.Gmax = -999.0
			OMvavca.Gmin = -999.0
			for f.IsEnd() == false {
				ss = f.GetToken()
				if ss[0] == ';' {
					break
				}
				if ce := strings.IndexRune(ss, ';'); ce != -1 {
					s = s[:ce]
				}
				if OMVAVdata(ss, OMvavca) != 0 {
					Eprint("<Eqcadata> OMVAV", ss)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == STHEAT_TYPE {
			Eqcat.Stheatca = append(Eqcat.Stheatca, STHEATCA{})
			stheatca := &Eqcat.Stheatca[len(Eqcat.Stheatca)-1]

			stheatca.Name = ""
			stheatca.PCMName = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Stheatdata(s, stheatca) != 0 {
					Eprint("<Eqcadata> STHEAT", s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == THEX_TYPE {
			Eqcat.Thexca = append(Eqcat.Thexca, THEXCA{})
			Thexca := &Eqcat.Thexca[len(Eqcat.Thexca)-1]

			Thexca.Name = ""
			Thexca.et = -999.0
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Thexdata(s, Thexca) != 0 {
					Eprint("<Eqcadata> THEX", s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == DESI_TYPE {
			Eqcat.Desica = append(Eqcat.Desica, DESICA{})
			Desica := &Eqcat.Desica[len(Eqcat.Desica)-1]

			Desica.name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Desiccantdata(s, Desica) != 0 {
					Eprint("<Eqcadata> DESICCANT", s)
				}
				if ce != -1 {
					break
				}
			}
		} else if eqpType == EVAC_TYPE {
			Eqcat.Evacca = append(Eqcat.Evacca, EVACCA{})
			Evacca := &Eqcat.Evacca[len(Eqcat.Evacca)-1]

			Evacca.Name = ""
			for f.IsEnd() == false {
				s = f.GetToken()
				if s[0] == ';' {
					break
				}
				if ce := strings.IndexRune(s, ';'); ce != -1 {
					s = s[:ce]
				}
				if Evacdata(s, Evacca) != 0 {
					Eprint("<Eqcadata> EVAC", s)
				}
				if ce != -1 {
					break
				}
			}
		} else {
			fmt.Printf("%s %s\n", E, s)
		}
	}
	Eqcat.Nhccca = len(Eqcat.Hccca)
	Eqcat.Nboica = len(Eqcat.Boica)
	Eqcat.Ncollca = len(Eqcat.Collca)
	Eqcat.Nrefaca = len(Eqcat.Refaca)
	Eqcat.Npipeca = len(Eqcat.Pipeca)
	Eqcat.Nstankca = len(Eqcat.Stankca)
	Eqcat.Nhexca = len(Eqcat.Hexca)
	Eqcat.Npumpca = len(Eqcat.Pumpca)
	Eqcat.Nvavca = len(Eqcat.Vavca)
	Eqcat.Nstheatca = len(Eqcat.Stheatca)
	Eqcat.Nthexca = len(Eqcat.Thexca)
	Eqcat.Npvca = len(Eqcat.PVca)
	Eqcat.Nomvavca = len(Eqcat.OMvavca)
	Eqcat.Ndesica = len(Eqcat.Desica) // Satoh追加　デシカント空調機　2013/10/20
	Eqcat.Nevacca = len(Eqcat.Evacca) // Satoh追加　気化冷却器　2013/10/26
}

/****************************************************************************/
func Eqpcount(fi *EeTokens, NBOI, NREFA, NCOL, NSTANK, NHCC, NHEX, NPIPE, NPUMP, NVAV, NSTHEAT, NTHEX, NPV, NOMVAV, NDESI, NEVAC *int) {
	ad := fi.GetPos()

	for fi.IsEnd() == false {
		s := fi.GetToken()

		if s == "*" {
			break
		} else if s == string(HCCOIL_TYPE) {
			*NHCC++
		} else if s == string(BOILER_TYPE) {
			*NBOI++
		} else if s == string(COLLECTOR_TYPE) || s == string(ACOLLECTOR_TYPE) {
			*NCOL++
		} else if s == string(REFACOMP_TYPE) {
			*NREFA++
		} else if s == string(PIPEDUCT_TYPE) || s == string(DUCT_TYPE) {
			*NPIPE++
		} else if s == string(STANK_TYPE) {
			*NSTANK++
		} else if s == string(HEXCHANGR_TYPE) {
			*NHEX++
		} else if s == string(PUMP_TYPE) || s == string(FAN_TYPE) {
			*NPUMP++
		} else if s == string(VAV_TYPE) || s == string(VWV_TYPE) {
			*NVAV++
		} else if s == string(STHEAT_TYPE) {
			*NSTHEAT++
		} else if s == string(THEX_TYPE) {
			*NTHEX++
		} else if s == string(PV_TYPE) {
			*NPV++
		} else if s == string(OMVAV_TYPE) || s == string(OAVAV_TYPE) {
			*NOMVAV++
		} else if s == string(DESI_TYPE) {
			*NDESI++
		} else if s == string(EVAC_TYPE) {
			*NEVAC++
		}
	}

	fi.RestorePos(ad)
}

func pflistcount(fl io.ReadSeeker) int {
	N := 0
	reader := bufio.NewReader(fl)

	for {
		s, err := reader.ReadString(' ')
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		s = strings.TrimSpace(s)

		if s == "*" {
			break
		} else if s == "!" {
			_, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
		} else if s == ";" {
			N++
		}

		if err == io.EOF {
			break
		}
	}

	_, err := fl.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	return N
}

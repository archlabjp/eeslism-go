package eeslism

import (
	"fmt"
	"io"
	"strings"
)

/* 標題、注記の出力（時刻別計算値ファイル） */

func __replace_dir_sep(path *string) {
	// Replace directory separator to unify them
	*path = strings.Replace(*path, "/", "\\", -1)
}

func ttlprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	// Replace directory separator to unify them
	__replace_dir_sep(&simc.File)

	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)
	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid h\n")

	fmt.Fprint(fo, "-tmid ")
	for i := 0; i < len(simc.Timeid); i++ {
		fmt.Fprint(fo, string(simc.Timeid[i]))
	}
	fmt.Fprint(fo, "\n")

	fmt.Fprintf(fo, "-u %s ;\n", simc.Unit)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprintf(fo, "-Ntime %d\n", simc.Ntimehrprt)
}

/* ---------------------------------------------------- */

/* 標題、注記の出力（日集計値ファイル） */

func ttldyprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)
	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid d\n")

	fmt.Fprint(fo, "-tmid ")
	for i := 0; i < len(simc.Timeid)-1; i++ {
		fmt.Fprint(fo, string(simc.Timeid[i]))
	}
	fmt.Fprint(fo, "\n")

	fmt.Fprintf(fo, "-u %s %s ;\n", simc.Unit, simc.Unitdy)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprintf(fo, "-Ntime %d\n", simc.Ntimedyprt)
}

/* ---------------------------------------------------- */

/* 標題、注記の出力（日集計値ファイル） */

func ttlmtprint(fo io.Writer, fileid string, simc *SIMCONTL) {
	fmt.Fprintf(fo, "%s#\n", fileid)
	fmt.Fprintf(fo, "-ver %s\n", EEVERSION)
	fmt.Fprintf(fo, "-t %s ;\n", simc.Title)
	fmt.Fprintf(fo, "-dtf %s\n", simc.File)

	fmt.Fprintf(fo, "-w %s\n", simc.Wfname)
	fmt.Fprint(fo, "-tid h\n")

	fmt.Fprint(fo, "-tmid MT\n")

	fmt.Fprintf(fo, "-u %s %s ;\n", simc.Unit, simc.Unitdy)
	fmt.Fprintf(fo, "-dtm %d\n", simc.DTm)
	fmt.Fprint(fo, "-Ntime 288\n") // 24 * 12
}

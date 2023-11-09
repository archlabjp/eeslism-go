package eeslism

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

/*        注釈文の除去             */

// Eesprera removes comments from the input file.
func Eesprera(file string) string {
	// 設定ファイルを開く
	fi, err := os.Open(file)
	if err != nil {
		fmt.Printf("File not found '%s'\n", file)
		os.Exit(1)
	}
	defer fi.Close()

	// 注釈文の除去語の設定ファイルを作成
	RET := strings.TrimSuffix(file, filepath.Ext(file))
	fb := new(strings.Builder)

	scanner := bufio.NewScanner(fi)

	// 各行を処理
	for scanner.Scan() {
		processedLine := processLine(scanner.Text())
		if processedLine != "" {
			_, err := fb.WriteString(processedLine + "\n")
			if err != nil {
				panic(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// var s string
	// var c byte
	// for scanner.Scan() {
	// 	s = scanner.Text()

	// 	if s == "!" {
	// 		//改行するまで読み進める
	// 		for _, err := fmt.Fscanf(fi, "%c", &c); err == nil && c != '\n'; _, err = fmt.Fscanf(fi, "%c", &c) {
	// 		}
	// 	} else {
	// 		if s == "　" {
	// 			//全角スペースは半角に置き換える
	// 			fmt.Fprint(fb, "  \n")
	// 		} else if s != "" {
	// 			fmt.Fprintf(fb, " %s ", strings.TrimSpace(s))
	// 		}

	// 		if s == ";" || strings.HasSuffix(s, ";") {
	// 			fmt.Fprintln(fb)
	// 		} else if s == "#" || strings.HasSuffix(s, "#") {
	// 			fmt.Fprintln(fb)
	// 		} else if s == "*" || strings.HasSuffix(s, "*") {
	// 			fmt.Fprintln(fb)
	// 		} else if s == "*debug" {
	// 			DEBUG = true
	// 		} else {
	// 			fmt.Fprintln(fb)
	// 		}
	// 	}
	// }

	// fmt.Fprintln(fb, " ")

	//互換性のために出力
	fbo, err := os.Create(strings.Join([]string{RET, "bdata0.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fbo, fb)
	}
	defer fbo.Close()

	return fb.String()
}

func processLine(line string) string {
	// "!"以降を削除
	if index := strings.Index(line, "!"); index != -1 {
		line = line[:index]
	}
	return line
}

/* ---------------------------------------------------------- */

/*              スケジュ－ルデ－タファイルの作成    */

// Eespre creates a schedule data file.
// 入力:
//   bdata0: コメント除去済みの入力テキスト
// 出力:
//   (1) %s または %sn から始まる論理行を除いた入力テキスト -> bdata.ewk
//   (2) %sから始まる論理行のみを収録したテキスト -> schtba.ewk
//   (3) %snから始まる論理行のみを収録したテキスト -> schenma.ewk
//   (4) WEEKデータセット -> week.ewk
func Eespre(bdata0 string, Ipath string, key *int) (string, string, string, string) {
	fi := strings.NewReader(bdata0) //bdata0.ewk 相当

	syspth := 0
	syscmp := 0

	Syscheck(fi, &syspth, &syscmp)

	fb := new(strings.Builder)  //bdata.ewk 相当
	fs := new(strings.Builder)  //schtba.ewk 相当 => %s を拾う
	fsn := new(strings.Builder) //schenma.ewk 相当 => %sn を拾う
	fw := new(strings.Builder)  //week.ewk 相当

	var st int

	var section_marker = []string{"TITLE", "GDAT", "RUN", "PRINT", "SCHTB", "EXSRF", "PCM",
		"WALL", "WINDOW", "SUNBRK", "ROOM", "RESI", "APPL", "VENT", "SYSCMP", "SYSPTH", "CONTL"}

	// %s, %sn を先に分離しておく。
	// 理由: 微妙な位置にある場合にうまく分離できない。
	tokens := NewEeTokens(bdata0)
	sb := new(strings.Builder)
	for !tokens.IsEnd() {

		line := tokens.GetLogicalLine()
		if line[0] == "*" {
			sb.WriteString("*\n")
			line = line[1:]
		}

		if slices.Contains(section_marker, line[0]) {
			sb.WriteString(line[0])
			sb.WriteString("\n")
			line = line[1:]
		}

		if line[0] == "%s" {
			for _, item := range line[1:] {
				fmt.Fprintf(fs, "%s ", item)
			}
			fs.WriteString(";\n")
		} else if line[0] == "%sn" {
			for _, item := range line[1:] {
				fmt.Fprintf(fsn, "%s ", item)
			}
			fsn.WriteString(";\n")
		} else {
			for _, item := range line {
				fmt.Fprintf(sb, "%s ", item)
			}
			sb.WriteString(";\n")
		}
	}

	// トークン分割
	tokens = NewEeTokens(sb.String())
	for !tokens.IsEnd() {

		s := tokens.GetToken()
		if s == "\n" {
			continue
		}

		// 壁体の材料定義リストを指定
		if strings.HasPrefix(s, "wbmlist=") {
			if st = strings.IndexRune(s, ';'); st != -1 {
				s = s[:st+1]
			} else {
				fmt.Fscanf(fi, "%*s")
			}

			Fbmlist = s[8:]
		} else if s == "WEEK" {
			*key = 1
			line := tokens.GetLogicalLine()
			for _, item := range line {
				fmt.Fprintf(fw, "%s ", item)
			}
			fw.WriteByte(';')
		} else if s == "*" {
			fb.WriteString("*\n")
		} else {
			fb.WriteString(s)
			line := tokens.GetLogicalLine()

			for _, item := range line {
				fmt.Fprintf(fb, " %s", item)
			}

			fb.WriteString(" ; \n")
		}
	}

	fmt.Fprintln(fb, "")
	//fb.Seek(0, 1)

	fmt.Fprintln(fb, " *")
	fmt.Fprintln(fw, "")
	fmt.Fprintln(fs, "*")
	fmt.Fprintln(fsn, "*")

	// ファイルに保存する(互換性のため)
	fbo, err := os.Create(strings.Join([]string{Ipath, "bdata.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fbo, fb)
		defer fbo.Close()
	}

	// SCHTBデータセット
	fso, err := os.Create(strings.Join([]string{Ipath, "schtba.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fso, fs)
		defer fso.Close()
	}

	// SCHNMAデータセット
	fsno, err := os.Create(strings.Join([]string{Ipath, "schnma.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fsno, fsn)
		defer fsno.Close()
	}

	// WEEKデータセット
	fwo, err := os.Create(strings.Join([]string{Ipath, "week.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fwo, fw)
		defer fwo.Close()
	}

	return fb.String(), fs.String(), fsn.String(), fw.String()
}

/******************************************************************/

func Syscheck(fi io.ReadSeeker, syspth *int, syscmp *int) {
	var s string
	for _, err := fmt.Fscanf(fi, "%s", &s); err == nil; _, err = fmt.Fscanf(fi, "%s", &s) {
		if s == "SYSPTH" {
			*syspth = 1
		} else if s == "SYSCMP" {
			*syscmp = 1
		}
	}
	fi.Seek(0, 0)
}

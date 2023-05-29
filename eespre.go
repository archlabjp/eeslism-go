package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

/*        注釈文の除去             */

// Eesprera removes comments from the input file.
func Eesprera(file string) string {
	// 設定ファイルを開く
	fi, err := os.Open(file)
	if err != nil {
		fmt.Printf("<eesprera> %s\n", file)
		os.Exit(1)
	}
	defer fi.Close()

	// 注釈文の除去語の設定ファイルを作成
	RET := strings.TrimSuffix(file, filepath.Ext(file))
	fb := new(strings.Builder)

	scanner := bufio.NewScanner(fi)
	var s string
	var c byte
	for scanner.Scan() {
		s = scanner.Text()

		if s == "!" {
			//改行するまで読み進める
			for _, err := fmt.Fscanf(fi, "%c", &c); err == nil && c != '\n'; _, err = fmt.Fscanf(fi, "%c", &c) {
			}
		} else {
			if s == "　" {
				//全角スペースは半角に置き換える
				fmt.Fprint(fb, "  \n")
			} else if s != "" {
				fmt.Fprintf(fb, " %s ", strings.TrimSpace(s))
			}

			if s == ";" || strings.HasSuffix(s, ";") {
				fmt.Fprintln(fb)
			} else if s == "#" || strings.HasSuffix(s, "#") {
				fmt.Fprintln(fb)
			} else if s == "*" || strings.HasSuffix(s, "*") {
				fmt.Fprintln(fb)
			} else if s == "*debug" {
				DEBUG = true
			} else {
				fmt.Fprintln(fb)
			}
		}
	}

	fmt.Fprintln(fb, " ")

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

/* ---------------------------------------------------------- */

/*              スケジュ－ルデ－タファイルの作成    */

// Eespre creates a schedule data file.
func Eespre(bdata0 string, Ipath string, key *int) (string, string, string, string) {
	fi := strings.NewReader(bdata0) //bdata0.ewk 相当

	syspth := 0
	syscmp := 0

	Syscheck(fi, &syspth, &syscmp)

	fb := new(strings.Builder)  //bdata.ewk 相当
	fs := new(strings.Builder)  //schtba.ewk 相当
	fsn := new(strings.Builder) //schenma.ewk 相当
	fw := new(strings.Builder)  //week.ewk 相当

	var t string
	var st int
	var stt int
	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {

		for _, s := range strings.Fields(scanner.Text()) {
			if s == "TITLE" {
				scanner.Scan()
				s = strings.TrimSuffix(scanner.Text(), ";")
				fmt.Fprintf(fb, "TITLE  %s ;\n", s)
			} else if strings.HasPrefix(s, "wbmlist=") {
				if st = strings.IndexRune(s, ';'); st != -1 {
					s = s[:st+1]
				} else {
					fmt.Fscanf(fi, "%*s")
				}

				Fbmlist = s[8:]
			} else if s == "WEEK" {
				*key = 1
				fmt.Fscanf(fi, " %[^;];", &s)
				fmt.Fprintf(fw, "%s ;\n", s)
			} else if s == "%s" {
				fmt.Fscanf(fi, " %[^;];", &s)
				fmt.Fprintf(fs, "%s ;\n", s)
			} else if s == "%sn" {
				fmt.Fscanf(fi, " %[^;];", &s)
				fmt.Fprintf(fsn, "%s ;\n", s)
			} else if strings.Contains(s, `"`) {
				fmt.Fprintf(fb, " %s", s)
				st = strings.Index(s, "\"")
				for st != -1 {
					stt = strings.Index(s[st+1:], "\"")
					if stt == -1 {
						break
					}
					stt = st + stt + 1
					t = s[st+1 : stt]
					if unicode.IsLetter(rune(t[0])) || t == "-" || t == "+" {
						fmt.Fprintf(fs, "-s %s\"  000-(%c)-2400 ;\n", t, t[1])
					} else {
						fmt.Fprintf(fs, "-v %s\"  000-(%s)-2400 ;\n", t, t[1:])
					}
					st = strings.Index(s[stt+1:], "\"")
				}
			} else {
				if strings.HasSuffix(s, ";") {
					t = s[:len(s)-1]
					fmt.Fprintf(fb, " %s ;", t)
				} else {
					fmt.Fprintf(fb, " %s", s)
				}

				if s == ";" || s[len(s)-1] == ';' {
					//fmt.Fprintln(fb)
				} else if s == "#" || s[len(s)-1] == '#' {
					//fmt.Fprintln(fb)
				} else if s == "*" || s[len(s)-1] == '*' {
					//fmt.Fprintln(fb)
				}
			}
		}
		fmt.Fprintln(fb)
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

	fso, err := os.Create(strings.Join([]string{Ipath, "schtba.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fso, fs)
		defer fso.Close()
	}

	fsno, err := os.Create(strings.Join([]string{Ipath, "schnma.ewk"}, ""))
	if err != nil {
		fmt.Println("Error creating file: ", err)
	} else {
		fmt.Fprint(fsno, fsn)
		defer fsno.Close()
	}

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

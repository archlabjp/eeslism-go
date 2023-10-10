// EESLISM Go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	eeslism "github.com/archlabjp/eeslism-go/eeslism"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	// コマンドライン引数の処理
	parser := argparse.NewParser("EESLISIM Go", "a general-purpose simulation program for building thermal-environmental control systems consisting of both buildings and facilities")

	filename := parser.StringPositional(&argparse.Options{
		Required: true,
		Help:     "Input data file name"})

	efl_path := parser.String("", "efl", &argparse.Options{
		Default: "Base",
		Help:    "EFLファイルのディレクトリ"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	if len(*efl_path) > 0 {
		os.Chdir(*efl_path)
	}

	eeslism.Entry(*filename)
}

module github.com/archlabjp/eeslism-go

go 1.23.0

toolchain go1.23.10

replace github.com/archlabjp/eeslism-go => ./eeslism

require (
	github.com/akamensky/argparse v1.4.0
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b
	gotest.tools v2.2.0+incompatible
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

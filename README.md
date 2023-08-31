# EESLISM Go

This program is a port of "EESLISM", a general-purpose simulation program for building thermal-environmental 
control systems consisting of both buildings and facilities, to the Go language.

The original EESLISM is written in C and is available at https://github.com/satoh-er/open_eeslism.

## Why porting to Go ?

The original EESLISM is littered with improper string memory handling. It was considered difficult to fix these completely. It seemed appropriate to port it to a C-compatible language with garbage collection, so we considered Carbon, Zig, and Go, and based on the popularity of the languages, we decided to try porting it to Go.

## Quick Start

For Ubuntu/Debian user
```
sudo apt install golang # if you didnot install golang
git clone https://github.com/archlabjp/eeslism-go
cd eeslism-go
go build
./eeslism
```

## How to build

We assume that the Go compiler is available.If not available, please refer to https://go.dev/doc/install to install Go.

Run next command.
```
go build
```

If you build as Windows Executable (64bit), run next command. You will get `eeslism.exe`.
```
GOOS=windows GOARCH=amd64 go build -o eeslism.exe
```

If you build as WebAssembly, run next command. You will get `eeslism.wasm`.
```
GOOS=js GOARCH=wasm go build -o eeslism.wasm
```

For other compilation targets, please refer to [here](https://go.dev/doc/install/source#environment
).

## Porting Policy

In porting from C to Go language, we keep the changes to a minimum. We also try to keep the source code names as one-to-one as possible. For example, if the original source code name is name.c, the ported source code name is name.go. This is to facilitate verification in case of mistakes by maintaining the correspondence.

All code is stored in the main module, which is a private function in the Go language if the function name starts with a lowercase letter. In order to maintain identity with the original function name, it was necessary to store all code in the main module.

### Accuracy of porting

We have confirmed that the same calculation results are obtained for a minimum sample. However, EESLISM is a very versatile program with a long history and requires much validation.

## Internal Structure

Please refer to [this picture](eeslism_data_structure.png) for data structure.

## Author

Wataru Uda

## License

Distributed under the GPL-2.0 License. See [LICENSE](LICENSE) for more information.

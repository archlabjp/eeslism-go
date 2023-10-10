# How to DEBUG

## Get Coverage 

```
$ cd eeslism
$ go test -cover -coverprofile=cover.out
$ go tool cover -html=cover.out -o cover.html
$ open cover.html
```
# How to DEBUG

## Get Coverage 

```
$ cd eeslism
$ go test -cover -coverprofile=cover.out
$ go tool cover -html=cover.out -o cover.html
$ open cover.html
```

## Run Equipment Tests

```
# Run all equipment tests
$ cd eeslism
$ go test -run "Test.*" -v

# Run specific equipment tests
$ go test -run TestPV -v
$ go test -run TestVAV -v
$ go test -run TestCOL -v
$ go test -run TestSTANK -v

# Run equipment tests with benchmarks
$ go test -run TestPV -bench BenchmarkPVCalculation -v
```

## Test Status

| Data Type | Status |
| --------- | ------ |
| WEEK  | OK |
| TITLE | OK |
| GDAT  | Work |
| SCHTB | Work |
| VCFILE | No |
| SCHNM | No |
| EXSRF | Work |
| WALL | Work |
| WINDOW | Work |
| SUNBRK | Work |
| ROOM | Work |
| VENT | Work |
| RESI | Work |
| APPL | Work |
| PCM | No |
| EQPCAT | Not |
| SYSCMP | Work |
| SYSPTH | Work |
| CONTL | Work |
| COORDNT | No |
| OBS | No |
| POLYGON | No |
| TREE | No |
| SHDSCHTB | No |
| DIVID | No |

### Equipment Test Status

| Equipment Type | Status |
| -------------- | ------ |
| BOI | OK |
| REFA | OK |
| COL | OK |
| STANK | OK |
| HEX | OK |
| HCC | OK |
| PIPE | OK |
| DUCT | OK |
| PUMP | OK |
| FAN | OK |
| VAV | OK |
| STHEAT | OK |
| THEX | OK |
| PV | OK |
| OMVAV | OK |
| DESI | OK |
| EVAC | OK |

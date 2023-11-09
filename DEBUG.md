# How to DEBUG

## Get Coverage 

```
$ cd eeslism
$ go test -cover -coverprofile=cover.out
$ go tool cover -html=cover.out -o cover.html
$ open cover.html
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
| BOI | No |
| REFA | No |
| COL | No |
| STANK | No |
| HEX | No |
| HCC | No |
| PIPE | No |
| DUCT | No |
| PUMP | No |
| FAN | No |
| VAV | No |
| STHEAT | No |
| THEX | No |
| PV | No |
| OMVAV | No |
| DESI | No |
| EVAC | No |

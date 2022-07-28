## Trade Executor

### Requirements
- golang 1.18
- sqlite3

### Build
```
go build -o sell cmd/sell/main.go
```

### Help
```
NAME:
   sell - create market sell order

USAGE:
   sell [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h      show help (default: false)
   --price value   price of trade (default: 0)
   --size value    size of trade (default: 0)
   --symbol value  ticker symbol

```

### Example
```
./sell --symbol BNBUSDT --price 266.1 --size 100
```

### Test
```
go test -v ./pkg/trader/**.go
```

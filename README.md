# pipe/line

## getting the library

```bash
go get github.com/MasteryConnect/pipe
```

## examples

### using pipe/line with unix pipes

Here is a basic script to count lines of input.

```golang
package main

import (
  l "github.com/MasteryConnect/pipe/line"
  "github.com/MasteryConnect/pipe/x"
)

func main() {
  l.New().Add(
    x.Head{N: 1000}.T,
    x.Count{}.T,
  ).Run()
}
```

run with
```bash
echo -e "foo\nbar" | go run foo.go
```

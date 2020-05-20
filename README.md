# pipe/line

pipe/line is an implementation of the pipeline design pattern. It follows the "Producer" -> "Transformer" -> "Consumer" nomenclature
to describe the various parts of the pipeline. It uses golang channels to connect those pieces together. Each part of the pipeline
runs in it's own goroutine so you can take advantage of the out-of-the-box concurrency without having to create your own channels
and start your own goroutines.


# getting the library

```bash
go get github.com/MasteryConnect/pipe
```

# examples

## Hello World

The most basic example here uses the primitives of the line. There are some convenience functions that can make this less verbose.
We will take a look at that a bit later.

```golang
// create the pipeline
pipe := line.New()

// add a producer (optional) (STDIN is used instead of not specified)
pipe.SetP(func(out chan<- interface{}, errs chan<- error) {
  out <- "Hello World"
})

// add a transformer to upper case messages
pipe.Add(func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
  for msg := range in {
    out <- strings.ToUpper(msg.(string))
  }
})

// add a consumer that prints out each message
pipe.SetC(func(in <-chan interface{}, errs chan<- error) {
  for msg := range in {
    fmt.Println(msg)
  }
})

// run it
pipe.Run()

// output: HELLO WORLD

```

Since the SetP, Add, and SetP functions all return the same line, you can chain them together.

```golang
line.New().SetP(func(out chan<- interface{}, errs chan<- error) {
  out <- "Hello World"
}).Add(func(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
  for msg := range in {
    out <- strings.ToUpper(msg.(string))
  }
}).SetC(func(in <-chan interface{}, errs chan<- error) {
  for msg := range in {
    fmt.Println(msg)
  }
}).Run()

// output: HELLO WORLD
```

There is a common pattern that emerges with the transformers where you want to just process each message
as the input to a function right inline instead of having to range over the in chan and put on the out chan.

```golang
line.New().SetP(func(out chan<- interface{}, errs chan<- error) {
  out <- "Hello World"
}).Add(line.Inline(func(msg interface{}) (interface{}, error) {
  return strings.ToUpper(msg.(string)), nil
}).Add(line.Inline(func(msg interface{}) (interface{}, error) {
  fmt.Println(msg)
  return nil, nil // returning nil for the value means nothing is sent to the next step in the pipeline
}).Run()

// output: HELLO WORLD

```

## using pipe/line with unix pipes

Here is a basic script to count lines of input. Since the producer is not set, STDIN is used.
A message is produced per line of input.

```golang
package main

import (
  "github.com/MasteryConnect/pipe/line"
  "github.com/MasteryConnect/pipe/x"
)

func main() {
  line.New().Add(
    x.Head{N: 1000}.T,
    x.Count{}.T,
  ).Run()
}
```

run with
```bash
echo -e "foo\nbar" | go run foo.go
```

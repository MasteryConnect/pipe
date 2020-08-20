# pipe/line

pipe/line is an implementation of the pipeline design pattern. It follows the "Producer" -> "Transformer" -> "Consumer" nomenclature
to describe the various parts of the pipeline. It uses golang channels to connect those pieces together. Each part of the pipeline
runs in it's own goroutine so you can take advantage of the out-of-the-box concurrency without having to create your own channels
and start your own goroutines.


# getting the library

```bash
go get github.com/MasteryConnect/pipe/line
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

To take this a step further, you can use some syntactic sugar to make this much more readable.

```golang
line.New().SetP(func(out chan<- interface{}, errs chan<- error) {
  out <- bytes.NewBufferString("Hello World")
}).Map(func(msg *bytes.Buffer) string {
  return strings.ToUpper(msg.String())
}).ForEach(func(msg string) {
  fmt.Println(msg)
}).Run()

// output: HELLO WORLD

```

It is recommended what you only use Map, ForEach, and Filter funcs with arg types other than interface{}
if you control the incoming message types explicitly.  Any mismatch between the message type coming in
and the type defined in the callback func will result in a panic. if you don't know or have strong confidence
in the source of the messages coming in to these callback funcs, it would be safest to use interface{}
as the arg type and do the type assertion explicitly as this will allow you to handle any type mismatches.

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

## syntactic sugar (Map,Filter,ForEach)

There are three sugar functions that can help with readability.

* Map
* Filter
* ForEach

These can be used directly on the pipeline. They are wrappers to the same func so they all behave the same way.
They can take a few different function signature shapes. The type 'interface{}' being used in the examples is a
placeholder for any type you want to use. The Map func does a type assertion to make the message
match the type in the func signature. Just like a type assertion like `foo := bar.(string)`, if the type
assertion fails, it panics so make sure the right type is being passed.

```golang
// only gets the messasge and the message will be automaticaly passed on after this function is done.
func(m interface{}) {}

// same as above but with a ctx when the pipeline is run with RunContext()
// The ctx can be added to any of these signatures.
func(ctx context.Context, m T) {}

// will send the returned error down the errs chan if not nil
func(m interface{}) error {}

// will send the returned value down stream. A nil value doesn't send anything downstream
func(m interface{}) interface{} {}

// a combination of the above two signatures
func(m interface{}) (interface{}, err) {}
```

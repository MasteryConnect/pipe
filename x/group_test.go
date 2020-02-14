package x_test

import (
	"fmt"

	l "github.com/MasteryConnect/pipe/line"
	"github.com/MasteryConnect/pipe/x"
)

func ExampleGroup() {
	l.New().SetP(func(out chan<- interface{}, errs chan<- error) {
		for i := 0; i < 10; i++ {
			out <- i
		}
	}).Add(
		x.NewGroup(3, func(msg interface{}) []string {
			if msg.(int)%2 == 0 {
				return []string{"even"} // put in the "even" group
			}
			return []string{"odd"} // put in the "odd" group
		}).T,
		l.I(func(msg interface{}) (interface{}, error) {
			group := msg.(*x.GroupMsg) // type cast to the group message
			return fmt.Sprintf("-- %s --\n%v\n", group.Name, group.Batch), nil
		}),
		l.Stdout,
	).Run()

	// Output:
	// -- even --
	// 0
	// 2
	// 4
	//
	// -- odd --
	// 1
	// 3
	// 5
	//
	// -- even --
	// 6
	// 8
	//
	// -- odd --
	// 7
	// 9
}

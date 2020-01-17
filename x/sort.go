package x

import "fmt"

// Sort will kind of sort the messages withing a size window.
// There is no guarantee of exact ordering, but better than nothing.
type Sort struct {
	N       int
	Compare func(x, y interface{}) bool // true if x <= y
}

// T is the Tfunc for Sort.
func (s Sort) T(in <-chan interface{}, out chan<- interface{}, errs chan<- error) {
	tree := &sortBTree{compare: s.Compare}
	cnt := 0

	for msg := range in {
		cnt++

		if cnt < s.N {
			tree.Insert(msg)
		} else {
			tree.Insert(msg)
			val := tree.RemoveLeftMost()
			out <- val
		}
	}

	// flush the rest out
	msg := tree.RemoveLeftMost()
	for msg != nil {
		out <- msg
		msg = tree.RemoveLeftMost()
	}
}

type sortCompareFunc func(x, y interface{}) bool // true if x <= y
type sortBTree struct {
	root      *sortNode
	leftMost  *sortNode
	rightMost *sortNode
	compare   sortCompareFunc
}
type sortNode struct {
	val   interface{}
	t     *sortBTree
	left  *sortNode
	right *sortNode
	up    *sortNode
}

func (t *sortBTree) Insert(val interface{}) error {
	if val == nil {
		return fmt.Errorf("attempting to insert nil value")
	}
	newNode := &sortNode{val: val, t: t}
	if t.root == nil {
		t.root = newNode
		t.leftMost = t.root
		t.rightMost = t.root
		return nil
	}

	return t.root.Insert(newNode)
}

// RemoveLowest will
func (t *sortBTree) RemoveLeftMost() interface{} {
	if t.root == nil {
		return nil
	}
	node := t.leftMost
	if t.leftMost != t.root {
		t.leftMost.up.left = t.leftMost.right
		if t.leftMost.right != nil {
			t.leftMost.right.up = t.leftMost.up
			t.leftMost = t.leftMost.right.LeftMost()
		} else {
			t.leftMost = t.leftMost.up
		}
	} else {
		// slide the "right" over to the root
		t.root = t.leftMost.right
		t.leftMost = t.root.LeftMost()
	}
	return node.val
}

func (n *sortNode) Insert(newNode *sortNode) error {
	if n.t.compare(newNode.val, n.val) {
		// new val is less
		if n.left == nil {
			n.left = newNode
			newNode.up = n
			if n.t.leftMost == n {
				n.t.leftMost = n.left
			}
			return nil
		}
		return n.left.Insert(newNode)
	}

	// new val is bigger
	if n.right == nil {
		n.right = newNode
		newNode.up = n
		if n.t.rightMost == n {
			n.t.rightMost = n.right
		}
		return nil
	}

	return n.right.Insert(newNode)
}

// LeftMost will recursively travers left until it can't and return that.
func (n *sortNode) LeftMost() *sortNode {
	if n == nil {
		return n
	}
	if n.left == nil {
		return n
	}
	return n.left.LeftMost()
}

package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	var walk func(t *tree.Tree)
	walk = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walk(t.Left)
		ch <- t.Value
		walk(t.Right)
	}

	walk(t)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for {
		select {
		case v1, ok1 := <-ch1:
			v2, ok2 := <-ch2
			if !ok1 || !ok2 {
				return ok1 == ok2
			}
			if v1 != v2 {
				return false
			}
		}
	}
}

func main() {
	t1, t2 := tree.New(1), tree.New(2)
	fmt.Println(Same(t1, t1))
	fmt.Println(Same(t1, t2))
	fmt.Println(Same(t2, t2))
}

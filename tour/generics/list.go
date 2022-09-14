package main

import "fmt"

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

func main() {
	l4 := List[string]{nil, "test"}
	l3 := List[string]{&l4, "a"}
	l2 := List[string]{&l3, "is"}
	l1 := List[string]{&l2, "this"}

	l := l1
	for {
		fmt.Println(l.val)
		if (l.next) == nil {
			break
		}
		l = *(l.next)
	}
}

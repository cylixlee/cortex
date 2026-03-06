# Iterator

Provide a way to access elements of a collection sequentially without exposing its underlying representation.

## Intent

- Access elements of a collection sequentially
- Decouple iteration from the collection
- Support multiple traversals simultaneously

## Implementation

```go
package main

import "fmt"

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type ArrayIterator struct {
	array   []interface{}
	position int
}

func NewArrayIterator(array []interface{}) *ArrayIterator {
	return &ArrayIterator{array: array}
}

func (i *ArrayIterator) HasNext() bool {
	return i.position < len(i.array)
}

func (i *ArrayIterator) Next() interface{} {
	item := i.array[i.position]
	i.position++
	return item
}

type ListIterator struct {
	list     []interface{}
	position int
}

func NewListIterator(list []interface{}) *ListIterator {
	return &ListIterator{list: list}
}

func (i *ListIterator) HasNext() bool {
	return i.position < len(i.list)
}

func (i *ListIterator) Next() interface{} {
	item := i.list[i.position]
	i.position++
	return item
}

func main() {
	array := []interface{}{"A", "B", "C"}
	arrayIterator := NewArrayIterator(array)
	
	for arrayIterator.HasNext() {
		fmt.Println(arrayIterator.Next())
	}
	
	list := []interface{}{"X", "Y", "Z"}
	listIterator := NewListIterator(list)
	
	for listIterator.HasNext() {
		fmt.Println(listIterator.Next())
	}
}
```

## When to Use

- When you want to access a collection's elements without exposing its internal structure
- When you want to provide multiple traversal methods for the same collection
- When you want to decouple collection logic from iteration logic

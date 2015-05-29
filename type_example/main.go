package main

import "fmt"

type Count int

func main() {
	var c Count = 1
	fmt.Println(c.Increment())
}

func(count *Count) Increment() Count {
	o := *count
	*count++
	return o
}
func(count *Count) Decrement() Count {
	o := *count
	*count--
	return o
}

func (count Count) String() string {
	return fmt.Sprintf("%d", int(count))
}


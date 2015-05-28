package main

import "fmt"

func main() {
	done := make(chan string)
	go routine(done)
	go routine(done)
	go routine(done)
	<-done
	<-done
	<-done
	close(done)
}

func routine(done chan string) {
	fmt.Println("Hello world!")
	done <- "done"
}

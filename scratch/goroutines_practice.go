package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("HI, I'm Aryan")
	go func() {
		fmt.Println("Random 1")
	}()
	go fmt.Println("Whats up")

	time.Sleep(time.Second * 2)

	func(text string) {
		fmt.Println(text)
	}("Now lets go!")

	messages := make(chan string, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		messages <- "Random 2"
		wg.Done()
	}()

	go func() {
		messages <- "Random 3"
		wg.Done()
	}()

	fmt.Println(<-messages)
	fmt.Println(<-messages)

	wg.Wait()

	fmt.Println("THANK YOU!")
}

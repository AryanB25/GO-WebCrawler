package main

import (
	"fmt"
	"time"
)

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for i := 1; i <= numJobs; i++ {
		go worker(i, jobs, results)
	}

	for a := 1; a <= numJobs; a++ {
		jobs <- a
	}
	close(jobs)

	for n := 1; n <= numJobs; n++ {
		fmt.Print(<-results)
	}

}

func worker(jobID int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker ", jobID, " has started his job!")
		time.Sleep(time.Second)
		fmt.Println("Worker ", jobID, " has finished his job!")
		results <- j * 2
	}
}

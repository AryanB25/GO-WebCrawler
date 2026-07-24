package tests

import (
	"GO-WebCrawler/internal/datastructures"
	"testing"
)

// TestPriorityQueue tests pushing a single element into the priority queue
// and verifies it is correctly stored and retrievable via Peek
func TestPriorityQueue(t *testing.T) {
	testVariables := []struct {
		operation string
		data      string
		score     int
		expected  string
	}{
		{"Add the URL and it's score to the priority queue", "https://www.ubc.ca", 20, "https://www.ubc.ca"},
		{"Add an empty string to the priority queue", "", 23, ""},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			priorityQueue := datastructures.NewPriorityQueue()
			priorityQueue.Push(tt.data, tt.score) // push the test element
			if priorityQueue.IsEmpty() {          // verify the element was added
				t.Errorf("Unable to push the element into the queue")
			}
			if priorityQueue.Peek() != tt.expected { // verify the correct element is at the top
				t.Errorf("Expected %q, got %q", tt.expected, priorityQueue.Peek())
			}
		})
	}
}

// TestPriorityQueueMultiple tests pushing multiple elements and verifies
// the max heap property — highest score always stays at the top
func TestPriorityQueueMultiple(t *testing.T) {
	priorityQueue := datastructures.NewPriorityQueue()
	testVariables := []struct {
		operation string
		data      string
		score     int
		expected  string
	}{
		{"Add the URL and it's score to the priority queue", "https://www.ubc.ca", 20, "https://www.ubc.ca"},
		{"Add an empty string to the priority queue", "", 20, "https://www.ubc.ca"},
		{"Add the URL and it's score to the priority queue", "https://news.ycombinator.com/item?id=19413348", 60, "https://news.ycombinator.com/item?id=19413348"},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			priorityQueue.Push(tt.data, tt.score) // push each element into the shared queue
			if priorityQueue.IsEmpty() {          // verify queue is not empty after push
				t.Errorf("Unable to push the element into the queue")
			}
			if priorityQueue.Peek() != tt.expected { // verify highest score element is at top
				t.Errorf("Expected %q, got %q", tt.expected, priorityQueue.Peek())
			}
		})
	}
}

// TestPriorityQueuePopMultiple tests that Pop returns elements in descending
// score order, confirming the max heap ordering is maintained across multiple pops
func TestPriorityQueuePopMultiple(t *testing.T) {
	priorityQueue := datastructures.NewPriorityQueue()
	priorityQueue.Push("https://www.ubc.ca", 20)                            // score 20
	priorityQueue.Push("", 10)                                              // score 10
	priorityQueue.Push("https://news.ycombinator.com/item?id=19413348", 60) // score 60 — should pop first

	testVariables := []struct {
		operation string
		expected  string
		status    bool
	}{
		{"Pop the first URL", "https://news.ycombinator.com/item?id=19413348", true}, // highest score pops first
		{"Pop the empty string", "https://www.ubc.ca", true},                         // second highest score
		{"Pop the final URL", "", false},                                             // queue empty after this
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			result, statusCode := priorityQueue.Pop()
			if result != tt.expected || statusCode != tt.status { // verify both value and status match
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

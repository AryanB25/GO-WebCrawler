package tests

import (
	"GO-WebCrawler/internal/datastructures"
	"testing"
)

// TestEnqueue verifies that individual elements can be successfully added
// to the queue and that the inserted value can be retrieved correctly.
func TestEnqueue(t *testing.T) {

	// Table-driven tests allow multiple enqueue scenarios to be tested
	// independently using different inputs and expected outputs.
	testVariables := []struct {
		operation string
		input     string
		expected  string
	}{
		{"One variable added to the queue", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/"},
		{"Empty string added to the queue", "", ""},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {

			// Each test case initializes a fresh queue to ensure tests
			// remain isolated and do not depend on previous executions.
			queue := datastructures.NewQueue()

			// Add the test input to the queue and validate the queue state.
			queue.Enqueue(tt.input)

			if queue.Size() != 1 {
				t.Errorf("Expected Answer (1), got the answer: %d", queue.Size())
			}

			// Confirm that the value stored in the queue matches the value inserted.
			if queue.Peek() != tt.expected {
				t.Errorf("Expected Answer %q, got the answer: %q", tt.expected, queue.Peek())
			}
		})
	}
}

// TestEnqueueMultiple validates that the queue correctly tracks its size
// when multiple elements are inserted.
func TestEnqueueMultiple(t *testing.T) {

	// Create a queue containing multiple URLs to simulate
	// storing crawler tasks waiting to be processed.
	queue := QueueDefinition()

	if queue.Size() != 3 {
		t.Errorf("Expected Answer (3), got the answer: %d", queue.Size())
	}

	// Ensures that adding elements prevents the queue from being empty.
	if queue.IsEmpty() {
		t.Errorf("Expected to have 3 elements")
	}
}

// TestDequeue verifies FIFO (First-In-First-Out) behavior by ensuring
// elements are removed in the same order they were inserted.
func TestDequeue(t *testing.T) {

	// Populate the queue with sample URLs representing pages
	// waiting to be crawled.
	queue := QueueDefinition()

	if queue.Size() != 3 {
		t.Errorf("Expected Answer (3), got the answer: %d", queue.Size())
	}

	testVariables := []struct {
		operation string
		input     string
		expected  string
	}{
		{"Dequeue the first link", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/"},
		{"Dequeue the second link", "https://www.ubc.ca", "https://www.ubc.ca"},
		{"Dequeue the third link", "https://news.ycombinator.com/item?id=19413348", "https://news.ycombinator.com/item?id=19413348"},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			expectedValue, isNotEmpty := queue.Dequeue()

			if !isNotEmpty {
				t.Errorf("Expected the queue to not be empty")
			}

			if expectedValue != tt.expected {
				t.Errorf("Expected Answer %q, got the answer: %q", tt.expected, expectedValue)
			}
		})
	}
}

// TestPeek verifies that Peek returns the next element without
// removing it from the queue.
func TestPeek(t *testing.T) {

	// Initialize queue with multiple elements to verify
	// that Peek returns the front element correctly.
	queue := QueueDefinition()

	if queue.Size() != 3 {
		t.Errorf("Expected Answer (3), got the answer: %d", queue.Size())
	}

	if queue.IsEmpty() {
		t.Errorf("Expected to have 3 elements")
	}

	// Remove the first element to expose the next item in the queue.
	expectedValue, isNotEmpty := queue.Dequeue()

	if !isNotEmpty {
		t.Errorf("Expected the queue to not be empty")
	}

	if expectedValue != "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/" {
		t.Errorf("Expected Answer %q, got the answer: %q", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/", expectedValue)
	}

	// Verify that Peek only reads the next queue element
	// without modifying the queue structure.
	peekValue := queue.Peek()

	if peekValue != "https://www.ubc.ca" {
		t.Errorf("Expected %q, but got %q", "https://www.ubc.ca", peekValue)
	}
}

// QueueDefinition creates a reusable test queue containing sample URLs.
// The URLs simulate pages that would be stored and processed by the web crawler.
func QueueDefinition() *datastructures.Queue {
	newQueue := datastructures.NewQueue()

	// Insert sample crawler targets to test queue ordering and retrieval.
	newQueue.Enqueue("https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/")
	newQueue.Enqueue("https://www.ubc.ca")
	newQueue.Enqueue("https://news.ycombinator.com/item?id=19413348")

	return newQueue
}

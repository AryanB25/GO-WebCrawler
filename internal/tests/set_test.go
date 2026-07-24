package tests

import (
	"GO-WebCrawler/internal/datastructures"
	"testing"
)

// TestSet tests that adding a single element to the set correctly updates its size.
// Also verifies that adding an empty string does not increase the size.
func TestSet(t *testing.T) {
	testVariables := []struct {
		operation string
		input     string
		expected  int
	}{
		{"Adding a visited URL to the set", "https://www.ubc.ca", 1},
		{"Adding an empty string to the set", "", 0},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			set := datastructures.NewSet()
			set.AddIfNotExists(tt.input)   // add the test input to a fresh set
			if set.Size() != tt.expected { // verify size matches expected after add
				t.Errorf("Expected the size to be %d, but the size is %d", tt.expected, set.Size())
			}
		})
	}
}

// TestAddSetMultiple tests adding multiple unique URLs to the same set
// and verifies the size increments correctly after each addition.
func TestAddSetMultiple(t *testing.T) {
	set := datastructures.NewSet() // shared set across all test cases
	testVariables := []struct {
		operation string
		input     string
		expected  int
	}{
		{"Adding the first URL to the set", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/", 1},
		{"Adding the second URL to the set", "https://www.ubc.ca", 2},
		{"Adding the third URL to the set", "https://news.ycombinator.com/item?id=19413348", 3},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			set.AddIfNotExists(tt.input)              // add each URL to the shared set
			if set.Size() != tt.expected { // verify size grows by 1 each time
				t.Errorf("Expected the size to be %d, but the size is %d", tt.expected, set.Size())
			}
		})
	}
}

// TestContains tests that Contains correctly returns true for an added URL
// and false for an empty string that was never added.
func TestContains(t *testing.T) {
	testVariables := []struct {
		operation string
		input     string
		expected  bool
	}{
		{"Check if the visited link is present in the set", "https://www.ubc.ca", true},
		{"Check if the empty link is present in the set", "", false},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			set := datastructures.NewSet()
			set.AddIfNotExists(tt.input)                          // add the input so we can check it exists
			if set.Contains(tt.input) != tt.expected { // verify Contains returns the expected result
				t.Errorf("Expected %t, got %t", tt.expected, set.AddIfNotExists(tt.input))
			}
		})
	}
}

// TestContainsMultiple tests Contains on a set with multiple URLs already added.
// Verifies that added URLs are found and URLs never added return false.
func TestContainsMultiple(t *testing.T) {
	set := datastructures.NewSet()
	set.AddIfNotExists("https://www.ubc.ca")                            // pre-populate the set
	set.AddIfNotExists("https://news.ycombinator.com/item?id=19413348") // pre-populate the set

	testVariables := []struct {
		operation string
		input     string
		expected  bool
	}{
		{"Check if the visited link is present in the set", "https://www.ubc.ca", true},
		{"Check if the empty link is present in the set", "", false},
		{"Check if the visited link is present in the set", "https://news.ycombinator.com/item?id=19413348", true},
		{"Check if the visited link is present in the set", "https://www.ubc.ca", true},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			if set.Contains(tt.input) != tt.expected { // verify each URL's presence matches expectation
				t.Errorf("Expected %t, got %t", tt.expected, set.AddIfNotExists(tt.input))
			}
		})
	}
}

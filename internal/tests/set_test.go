package tests

import (
	"GO-WebCrawler/internal/datastructures"
	"testing"
)

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
			set.Add(tt.input)
			if set.Size() != tt.expected {
				t.Errorf("Expected the size to be %d, but the size is %d", tt.expected, set.Size())
			}
		})
	}
}

func TestAddSetMultiple(t *testing.T) {
	set := datastructures.NewSet()
	testVariables := []struct {
		operation string
		input     string
		expected  int
	}{{"Adding the first URL to the set", "https://www.freecodecamp.org/news/unit-testing-in-go-a-beginners-guide/", 1},
		{"Adding the second URL to the set", "https://www.ubc.ca", 2},
		{"Adding the third URL to the set", "https://news.ycombinator.com/item?id=19413348", 3},
	}

	for _, tt := range testVariables {
		t.Run(tt.operation, func(t *testing.T) {
			set.Add(tt.input)
			if set.Size() != tt.expected {
				t.Errorf("Expected the size to be %d, but the size is %d", tt.expected, set.Size())
			}
		})
	}
}

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
			set.Add(tt.input)
			if set.Contains(tt.input) != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, set.Contains(tt.input))
			}
		})
	}
}

func TestContainsMultiple(t *testing.T) {
	set := datastructures.NewSet()
	set.Add("https://www.ubc.ca")
	set.Add("https://news.ycombinator.com/item?id=19413348")

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
			if set.Contains(tt.input) != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, set.Contains(tt.input))
			}
		})
	}
}

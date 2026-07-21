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
		{"Adding an empty string to the set", "", 1},
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
		t.Run(tt.operation, func (t *testing.T)  {
			set.Add(tt.input)
			if set.Size() != tt.expected {
				t.Errorf("Expected the size to be %d, but the size is %d", tt.expected, set.Size())
			}
		})
	}
}

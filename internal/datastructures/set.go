package datastructures

type Set struct { // structure of the set
	elements map[string]bool
}

func NewSet() *Set { // create a new set
	return &Set{elements: map[string]bool{}} // returns the address to a new set
}

func (set *Set) Add(url string) {
	set.elements[url] = true // adds the url to the set
}

func (set *Set) Contains(url string) bool {
	return set.elements[url] // checks if the url is in the set
}

func (set *Set) Size() int {
	return len(set.elements) // returns the size of the set
}

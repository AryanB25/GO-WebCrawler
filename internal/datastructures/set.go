package datastructures

type Set struct {
	elements map[string]bool
}

func NewSet() *Set {
	return &Set{elements: map[string]bool{}}
}

func (set *Set) Add(url string) {
	set.elements[url] = true
}

func (set *Set) Contains(url string) bool {
	return set.elements[url]
}

func (set *Set) Size() int {
	return len(set.elements)
}

package datastructures

import "sync"

type Set struct { // structure of the set
	mutex    sync.Mutex
	elements map[string]bool
}

func NewSet() *Set { // create a new set
	return &Set{elements: map[string]bool{}} // returns the address to a new set
}

func (set *Set) Add(url string) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	set.elements[url] = true // adds the url to the set
}

func (set *Set) Contains(url string) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	return set.elements[url] // checks if the url is in the set
}

func (set *Set) Size() int {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	return len(set.elements) // returns the size of the set
}

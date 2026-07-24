package datastructures

import "sync"

type Set struct { // structure of the set
	mutex    sync.Mutex
	elements map[string]bool
}

func NewSet() *Set { // create a new set
	return &Set{elements: map[string]bool{}} // returns the address to a new set
}

func (s *Set) AddIfNotExists(value string) bool {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, exists := s.elements[value]; exists {
        return false
    }

    s.elements[value] = true
    return true
}

func (set *Set) Size() int {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	return len(set.elements) // returns the size of the set
}

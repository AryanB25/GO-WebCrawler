package datastructures

import (
	"sync"
)

type Item struct {
	data  string
	score int
}

type priorityQueue struct {
	elements []Item
	mutex    sync.Mutex
}

func NewPriorityQueue() (pq *priorityQueue) {
	return &priorityQueue{} // creates a new priority queue
}

func (pq *priorityQueue) Push(data string, score int) {
	// ensures we prevent the goroutines from accessing all at once; prevents data race
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	priorityItem := Item{data, score}               // creates a new "item" for our priority queue
	pq.elements = append(pq.elements, priorityItem) // appends the "item" to the slice
	childIndex := pq.size() - 1                     // index of the child item
	parentIndex := (pq.size() - 1) / 2              // index of the parent item to the child

	for childIndex > 0 { // until we reach the root of the queue
		if score > pq.elements[parentIndex].score { // if the score is greater in value than it's parents value
			// swapping process of the child and parent
			previousItem := pq.elements[parentIndex]
			pq.elements[parentIndex] = priorityItem
			pq.elements[childIndex] = previousItem

			// updated child and parent index after swapping
			childIndex = parentIndex
			parentIndex = (childIndex - 1) / 2
		} else {
			break // break if the score is lower since there is nothing to do
		}
	}
}

func (pq *priorityQueue) Pop() (string, bool) {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	if pq.size() == 0 { // if the priority queue is empty
		return "", false
	}

	front := pq.elements[0].data // access the first element in the heap

	pq.elements[0] = pq.elements[pq.size()-1] // assign the last element to the first index
	pq.elements = pq.elements[:pq.size()-1]   // reduces the size of the array to accomodate for the removal
	parentIndex := 0                          // starts at the root of the heap
	childLeftndex := 0                        // index of the left child
	childRightIndex := 0                      // index of the right child
	replacingIndex := 0                       // stores the index with which it must swap

	for {
		childLeftndex = (2*parentIndex + 1)   // find the index of the left child relative to the parent index
		childRightIndex = (2*parentIndex + 2) // find the index of the right child relative to the parent index
		if childLeftndex >= pq.size() {
			break
		} else if childRightIndex >= pq.size() {
			replacingIndex = childLeftndex
		} else {
			if max(pq.elements[childLeftndex].score, pq.elements[childRightIndex].score) == pq.elements[childLeftndex].score {
				replacingIndex = childLeftndex
			} else {
				replacingIndex = childRightIndex
			} // find the minimum score index between the right and left child
		}
		if pq.elements[parentIndex].score < pq.elements[replacingIndex].score {
			// swapping process of child and parent
			tempItem := pq.elements[replacingIndex]
			pq.elements[replacingIndex] = pq.elements[parentIndex]
			pq.elements[parentIndex] = tempItem
			parentIndex = replacingIndex // the new index of the parent after swapping
		} else {
			break // if the element has the highest priority compared to it's children
		}
	}

	return front, true // returns the first element in the priority queue
}

func (pq *priorityQueue) Peek() string {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	if len(pq.elements) == 0 {
		return ""
	}
	return pq.elements[0].data // returns the first element in the priority queue
}

func (pq *priorityQueue) size() int {
	return len(pq.elements) // returns the size of the priority queue
}

func (pq *priorityQueue) IsEmpty() bool {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	return len(pq.elements) == 0 // checks if the priority queue is empty
}

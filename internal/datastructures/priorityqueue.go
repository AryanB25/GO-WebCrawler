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
	childIndex := len(pq.elements) - 1              // index of the child item
	parentIndex := (len(pq.elements) - 1) / 2       // index of the parent item to the child

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
	if len(pq.elements) == 0 {
		return "", false
	}
}

func (pq *priorityQueue) Peek() (data string) {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()
	return pq.elements[0].data // returns the first element in the priority queue
}

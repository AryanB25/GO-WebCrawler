package datastructures

type Queue struct {
	elements []string // stores the URL's
}

func NewQueue() *Queue {
	return &Queue{[]string{}} // creates a new queue and returns the address of the new queue
}

func (queue *Queue) Enqueue(url string) {
	queue.elements = append(queue.elements, url) // adding the URL to the list of URL's
}

func (queue *Queue) Dequeue() (string, bool) {
	if len(queue.elements) == 0 { // if the queue is empty or has no URL's
		return "", false
	}
	front := queue.elements[0]          // access the first URL of the queue
	queue.elements = queue.elements[1:] // remove the first URL of the queue
	return front, true
}

func (queue *Queue) Peek() string {
	return queue.elements[0] // peek at the first url present in the queue
}

func (queue *Queue) IsEmpty() bool {
	return len(queue.elements) == 0 // checks if the queue is empty
}

func (queue *Queue) Size() int {
	return len(queue.elements) // returns the number of URL's in the queue
}

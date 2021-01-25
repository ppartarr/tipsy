package main

import (
	"container/heap"
)

// An Item is something we manage in a weight queue.
type Item struct {
	value  string  // The value of the item; arbitrary.
	weight float64 // The weight of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the smallest, not highest, weight so we use smaller than here.
	return pq[i].weight < pq[j].weight
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq PriorityQueue) Find(value string) *Item {
	for _, item := range pq {
		if item.value == value {
			return item
		}
	}
	return nil
}

// update modifies the weight and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, weight float64) {
	item.value = value
	item.weight = weight
	heap.Fix(pq, item.index)
}

// // Some items and their priorities.
// items := map[string]int{
// 	"password": 59462, "iloveyou": 49952, "princess": 33291,
// }

// // Create a weight queue, put the items in it, and
// // establish the weight queue (heap) invariants.
// pq := make(PriorityQueue, len(items))
// i := 0
// for value, weight := range items {
// 	pq[i] = &Item{
// 		value:    value,
// 		weight: weight,
// 		index:    i,
// 	}
// 	i++
// }
// heap.Init(&pq)

// // Insert a new item and then modify its weight.
// item := &Item{
// 	value:    "rockyou",
// 	weight: 20901,
// }
// heap.Push(&pq, item)
// pq.update(item, item.value, 5)

// // Take the items out; they arrive in decreasing weight order.
// for pq.Len() > 0 {
// 	item := heap.Pop(&pq).(*Item)
// 	fmt.Printf("%.2d:%s ", item.weight, item.value)
// }

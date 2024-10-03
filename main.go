// Implement and Use a Stack in Go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Stack struct {
	items []interface{}
	// add mutex definition
	mu sync.Mutex
}

// IsEmpty returns true if the stack is empty otherwise returns false and the size of the stack
func (s *Stack) IsEmpty() (bool, int) {
	stackSize := len(s.items)
	if stackSize == 0 {
		return true, stackSize
	}

	return false, stackSize
}

// Push adds an item to the stack
func (s *Stack) Push(item interface{}) {
	//locks the mutex to ensure thread safety
	s.mu.Lock()
	//ensure the mutex is unlockd when the function exist
	defer s.mu.Unlock()
	s.items = append(s.items, item)
}

// Pop removes the top item from the stack and returns it
func (s *Stack) Pop() (interface{}, bool) {
	//locks the mutex to ensure thread safety
	s.mu.Lock()
	//ensure the mutex is unlockd when the function exist
	defer s.mu.Unlock()

	isEmpty, stackLen := s.IsEmpty()
	if isEmpty {
		return 0, false
	}
	stackLen = stackLen - 1
	topItem := s.items[stackLen] // it has the top item value to return
	// to remove that top item from the array
	s.items = s.items[:stackLen]

	return topItem, true
}

// Peek returns the top item without removing it
func (s *Stack) Peek() (interface{}, bool) {
	//locks the mutex to ensure thread safety
	s.mu.Lock()
	//ensure the mutex is unlockd when the function exist
	defer s.mu.Unlock()

	isEmpty, stackLen := s.IsEmpty()
	if isEmpty {
		return 0, false
	}
	topItem := s.items[stackLen-1]

	return topItem, true
}

func main() {
	//define the wait group in order to wait for the goroutines to finish excecuting
	var wg sync.WaitGroup
	//creates a chanel to comunicate between routines
	ch := make(chan int, 20)
	myStack := Stack{}

	// PUSH ROUTINE
	wg.Add(1) //Add the waitggroup
	//start go routine
	go func() {
		defer wg.Done()
		//push 20 numbers to mystack and set the values to be available in the chanel
		for i := 0; i <= 20; i++ {
			go myStack.Push(i)
			time.Sleep(100 * time.Millisecond) //add timer for the push of each value to let it be visible in the prints
			ch <- i

		}

		close(ch) //close the chanel

	}()

	//PEEK ROUTINE
	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range ch {
			fmt.Println("****Peeked values:", item)
			peekValue, k := myStack.Peek()
			if k {
				fmt.Println("**Peek value form the stack:", peekValue)
			}
			time.Sleep(200 * time.Millisecond) //add timer to wait the push of all the values
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for item := range ch {
			fmt.Printf("----Popped value: %v\n", item)
			poppedValue, k := myStack.Pop()
			if k {
				fmt.Printf("--Popped value from stack: %v\n", poppedValue)
			}
			time.Sleep(250 * time.Millisecond) //add timer to be able to see peek in order
		}
	}()
	//wait for the go routines to finish before proceding
	wg.Wait()

}

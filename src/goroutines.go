package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func counter(count int) {
	for i := 0; i < count; i++ {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}

func worker(workerID int, data chan int) {
	for x := range data {
		fmt.Printf("Worker %d: received %d\n", workerID, x)
		// Generates a random time between 1 and 5 seconds
		randSleep := time.Duration(rand.IntN(5)+1) * time.Second
		time.Sleep(randSleep) // Simulates a hard work
	}
}

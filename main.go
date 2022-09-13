package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/**
Calculate the sum of the 1's in the binary representation of al the numbers from
1 to 10_000_000
*/

const (
	MaxGoroutines = 10_000_000
)

func main() {
	fmt.Println("Information about RUNTIME")
	fmt.Printf("num. cpus: %d\n", runtime.NumCPU())
	fmt.Printf("num. max procs: %d\n", runtime.GOMAXPROCS(-1))
	measureFuncTime(calculateSumConcurrent, "concurrent function")   // 2.5 and 3 seconds
	measureFuncTime(calculateSumSynchronous, "Synchronous function") // 590 and 800 milliseconds
	measureFuncTime(calculateSumPaginated, "paginated function")     // between 170 and 190 milliseconds
	/**
	Overall, the paginated way of calculating this problem is the best, because it optimizes concurrency
	by knowing that this task is CPU-Bound, and we don't need to swap between go routines, we should let them
	run freely. That's why the most optimal way, is to only use the number of go routines that are able to run
	in parallel (the number returned by runtime.GOMAXPROCS(-1))
	*/
}

func measureFuncTime(f func(), desc string) {
	start := time.Now()
	f()
	fmt.Printf("Excecution lasted: %v for %s\n", time.Since(start), desc)
}

func calculateSumSynchronous() {
	sum := 0
	for i := 1; i <= MaxGoroutines; i++ {
		b := strconv.FormatInt(int64(i), 2)
		c := 0
		for _, v := range b {
			if v == '1' {
				c++
			}
		}
		sum += c
	}
	fmt.Println(sum)
}

func calculateSumConcurrent() {
	wg := sync.WaitGroup{}
	wg.Add(MaxGoroutines)
	var sum uint64 = 0

	for i := 1; i <= MaxGoroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			b := strconv.FormatInt(int64(i), 2)
			c := 0
			for _, v := range b {
				if v == '1' {
					c++
				}
			}
			atomic.AddUint64(&sum, uint64(c))
		}()
	}

	wg.Wait()
	fmt.Println(sum)
}

func calculateSumPaginated() {
	wg := sync.WaitGroup{}
	maxParallelProcesses := runtime.GOMAXPROCS(-1)

	wg.Add(maxParallelProcesses)
	sum := uint64(0)
	delta := uint64(MaxGoroutines / maxParallelProcesses)
	prevStart := uint64(0)
	for i := 0; i < maxParallelProcesses; i++ {
		prev := prevStart
		go func() {
			defer wg.Done()
			for j := uint64(math.Min(MaxGoroutines, float64(delta+prev))); j > prev; j-- {
				b := strconv.FormatInt(int64(j), 2)
				c := 0
				for _, v := range b {
					if v == '1' {
						c++
					}
				}
				atomic.AddUint64(&sum, uint64(c))
			}
		}()
		prevStart += delta
	}

	wg.Wait()
	fmt.Println(sum)
}

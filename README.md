# go-routines
Simple example of how to use and optimize go routines


## Problem description

Calculate the sum of the 1's in the binary representation of al the numbers from 1 to 10.000.000

This problem will be faced three different ways:
  - Synchronous way, iterating over all the numbers and summing the amount of 1's
  - Asynchronous way, but spawning a go routine for each number (10.000.000 go routines)
  - Asynchronous way, but only using the same amount of go routines as the value of runtime.GOMAXPROCS(-1), to only run parallel go routines
  
Measuring execution times, the most performant way is the third one by being (in average) 3 times faster than the second best.

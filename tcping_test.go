package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestGoRoutine(t *testing.T) {
	var wg sync.WaitGroup
	for i:= 0;  i < 5;i++ {
		var index = i*i
		wg.Add(1)
		go func() {
			fmt.Println(i, index)
			wg.Done()
		}()
	}
	wg.Wait()
}
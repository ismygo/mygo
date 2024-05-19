package mygo

import (
	"fmt"
	"sync"
	"testing"
)

func TestStack(t *testing.T) {
	fmt.Println(string(stack()))
}

func TestRecover(t *testing.T) {
	var err error

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer Panic.Recover(&err)
		panic("test panic")
	}()
	wg.Wait()

	t.Log(err)
}

package skiplist

import (
	"testing"

	"sync"
	"time"
)

func TestList(t *testing.T) {
	sl := New()

	if sl.Contains(2) {
		t.Fatal("list contains something we never added")
	}

	if sl.Add(2) == false {
		t.Fatal("failed to add new item to list, someone deleting ??????")
	}

	if sl.Contains(3) {
		t.Fatal("list contains something we never added")
	}

	if !sl.Contains(2) {
		t.Fatal("list doesnt contain what we just added")
	}

	if sl.Contains(3) {
		t.Fatal("list contains something we never added")
	}

	if sl.Contains(1) {
		t.Fatal("list contains something we never added")
	}

	if sl.Len() != 1 {
		t.Fatal("expected list to be of length 1")
	}

	if sl.Add(2) == true {
		t.Fatal("Add with already present value should have returned false")
	}

	if sl.Remove(2) == false {
		t.Fatal("failed to remove item from list, someone deleting it ??????")
	}

	if sl.Contains(2) {
		t.Fatal("list contains something we removed")
	}

	if sl.Contains(1) || sl.Contains(2) || sl.Contains(3) {
		t.Fatal("list contains something we never added")
	}

	if sl.Len() != 0 {
		t.Fatal("expected list to be of length 0")
	}
}

func TestParallel(t *testing.T) {
	// return
	c := make(chan bool)
	times := 10000
	values := 5
	sl := New()
	wg := sync.WaitGroup{}
	wg.Add(2)
	insert := func() {
		for j := 0; j < values; j++ {
			sl.Add(j)
		}
	}
	delete := func() {
		for j := 0; j < values; j++ {
			sl.Remove(j)
		}
	}
	go func() {
		defer wg.Done()
		<-c
		for i := 0; i <= times; i++ {
			insert()
		}
	}()

	go func() {
		defer wg.Done()
		<-c
		for i := 0; i < times; i++ {
			delete()
		}
	}()
	time.Sleep(time.Nanosecond * 10)
	close(c)
	wg.Wait()
	println("poulet: ", sl.Len())
	println("rotis !")
}

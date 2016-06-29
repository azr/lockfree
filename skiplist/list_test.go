package skiplist

import (
	"testing"

	"sync"
	"time"
	"unsafe"
)

func TestList(t *testing.T) {
	sl := New()
	two := 2
	three := "three"

	if sl.Contains(2) {
		t.Fatal("list contains something we never added")
	}

	if sl.Set(2, unsafe.Pointer(&two)) == false {
		t.Fatal("failed to add new item to list, someone deleting ??????")
	}

	if sl.Contains(3) {
		t.Fatal("list contains something we never added")
	}

	if !sl.Contains(2) {
		t.Fatal("list doesnt contain what we just added")
	}
	v, found := sl.Get(2)
	s := *(*int)(v)
	if found == false || s != two {
		t.Fatalf("could not get wat we stored, found '%d' instead of '%d'", s, two)
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

	if sl.Set(2, unsafe.Pointer(&three)) == true {
		t.Fatal("Add with already present value should have returned false")
	}

	v, found = sl.Get(2)
	_three := *(*string)(v)
	if found == false || _three != three {
		t.Fatalf("could not get wat we stored, found '%d' instead of '%d'", s, two)
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

func TestAlot(t *testing.T) {
	sl := New()
	in := 10000
	insert(t, sl, in, true)
	if sl.Len() != in {
		t.Fatal("inserted ", in, " items and size is: ", sl.Len())
	}
	remove(t, sl, in, true)
	if sl.Len() != 0 {
		t.Fatal("removed ", in, " items and size is: ", sl.Len())
	}
	if sl.Contains(0) {
		t.Fatal("list contains something we removed")
	}
}

func TestParallel(t *testing.T) {
	c := make(chan bool)
	times := 100
	values := 5
	sl := New()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-c
		for i := 0; i <= times; i++ {
			insert(t, sl, values, false)
		}
	}()

	go func() {
		defer wg.Done()
		<-c
		for i := 0; i < times; i++ {
			remove(t, sl, values, false)
		}
	}()
	time.Sleep(time.Nanosecond * 10)
	close(c)
	wg.Wait()
}

func insert(t *testing.T, sl *Header, values int, check bool) {
	for j := 0; j < values; j++ {
		sl.Set(j, unsafe.Pointer(nil))
		if check {
			checkList(t, sl)
		}
	}
}

func remove(t *testing.T, sl *Header, values int, check bool) {
	for j := 0; j < values; j++ {
		sl.Remove(j)
		if check {
			checkList(t, sl)
		}
	}
}

func checkList(t *testing.T, sl *Header) {
	//check that everything is in a valid state
	for i := range sl.leftSentinel.nexts {
		n := sl.leftSentinel.nexts.get(i)
		if n == nil {
			t.Fatalf("leftSentinel.next[%d] is nil ?", i)
		}
	}
	for curr := sl.leftSentinel; curr != nil; curr = curr.nexts.get(0) {
		curr.lock.Lock()
		curr.lock.Unlock()
	}
}

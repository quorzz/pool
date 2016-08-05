package pool

import (
	"errors"
	"sync"
	"testing"
)

type testStrcut struct {
	name string
}

func (ts testStrcut) GetName() string {
	return ts.name
}

var createFunc = func() (interface{}, error) {
	return &testStrcut{"quorzz"}, nil
}
var checkError = errors.New("errors on check")

var checkOnGetFunc = func(item interface{}) error {
	return checkError
}

func TestPool(t *testing.T) {

	p := NewPool(createFunc, 3)

	item, err := p.Get()

	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if p.Len() != 0 {
		t.Log("pool.Len() should be 0")
		t.Fail()
	}

	p.Put(item)
	p.Put(item)
	p.Put(item)
	if p.Len() != 3 {
		t.Log("pool.Len() should be 3")
		t.Fail()
	}

	p.Put(item)
	if p.Len() != 3 {
		t.Log("pool.Len() should be 3, maxidle error")
		t.Fail()
	}

	if ts, e := item.(*testStrcut); !e {
		t.Log("type error")
		t.Fail()
	} else if ts.GetName() != "quorzz" {
		t.Log("error:method")
		t.Fail()
	}

}

func TestCheckOnGet(t *testing.T) {
	p := NewPool(createFunc, 3)

	p.CheckOnGet = checkOnGetFunc

	_, err := p.Get()

	if err != checkError {
		t.Error(err)
		t.Fail()
	}
}

func TestConcurrent(t *testing.T) {
	testMaxIdle(t, 500, 1000)
	testMaxIdle(t, 1000, 1000)
	testMaxIdle(t, 1000, 2000)
	testMaxIdle(t, 1, 10000)
	testMaxIdle(t, 10000, 1)
}

func testMaxIdle(t *testing.T, maxIdle, addNum int) {
	p := NewPool(func() (interface{}, error) {
		return &testStrcut{"slina"}, nil
	}, maxIdle)

	item, _ := p.Get()

	var wg sync.WaitGroup
	wg.Add(addNum)
	for i := 0; i < addNum; i++ {
		go func() {
			p.Put(item)
			wg.Done()
		}()
	}

	wg.Wait()
	if p.Len() != addNum && p.Len() != p.MaxIdle {
		t.Log("mutex error", p.Len())
		t.Fail()
	}

	p.Clear()
	if p.Len() != 0 {
		t.Error("clear errors")
		t.Fail()
	}
}

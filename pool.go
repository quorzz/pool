package pool

import (
	"container/list"
	"sync"
)

type CreateFunc func() (interface{}, error)

type Pool struct {
	MaxIdle int
	Checker func(item interface{}) error

	idleList   *list.List
	mutex      sync.Mutex
	createFunc CreateFunc
}

func NewPool(createFunc CreateFunc, maxIdel int) *Pool {
	return &Pool{
		MaxIdle:    maxIdel,
		createFunc: createFunc,
		idleList:   list.New(),
		Checker: func(item interface{}) error {
			return nil
		},
	}
}

func (p *Pool) Get() (interface{}, error) {
	if p.MaxIdle <= 0 {
		goto CREATE_NEW
	}

	p.mutex.Lock()

	if p.idleList.Len() > 0 {
		item := p.idleList.Back()
		p.idleList.Remove(item)
		p.mutex.Unlock()

		return p.checkItem(item.Value)
	}
	p.mutex.Unlock()

CREATE_NEW:
	if newItem, err := p.createFunc(); err != nil {
		return nil, err
	} else {
		return newItem, nil
	}
}

func (p *Pool) checkItem(item interface{}) (interface{}, error) {
	if err := p.Checker(item); err != nil {
		return nil, err
	} else {
		return item, nil
	}
}

func (p *Pool) Put(item interface{}) {

	if p.MaxIdle <= 0 {
		return
	}
	p.mutex.Lock()
	if nil == item {
		p.mutex.Unlock()
		return
	}

	if p.idleList.Len() >= p.MaxIdle {
		p.idleList.Remove(p.idleList.Front())
	}

	p.idleList.PushBack(item)
	p.mutex.Unlock()
}

func (p *Pool) Clear() {
	p.idleList.Init()
}

func (p *Pool) Len() int {
	return p.idleList.Len()
}

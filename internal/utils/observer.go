package utils

import "sync"

type Observable struct {
	observers []chan any
	mu        sync.Mutex
}

func CreateObservable() *Observable {
	return &Observable{
		observers: make([]chan any, 0),
	}
}

func (o *Observable) RegisterObserver(observer chan any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.observers = append(o.observers, observer)
}

func (o *Observable) UnregisterObserver(observer chan any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for i, obs := range o.observers {
		if obs == observer {
			close(obs)
			o.observers = append(o.observers[:i], o.observers[i+1:]...)
			break
		}
	}
}

func (o *Observable) NotifyObservers(message any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, observer := range o.observers {
		go func(ch chan any) {
			ch <- message
		}(observer)
	}
}

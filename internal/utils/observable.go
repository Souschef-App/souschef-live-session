package utils

import (
	"fmt"
	"sync"
)

type Observable struct {
	observers []Observer
	mu        sync.Mutex
}

func CreateObservable() *Observable {
	return &Observable{
		observers: []Observer{},
	}
}

func (o *Observable) RegisterObserver(observer Observer) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.observers = append(o.observers, observer)
}

func (o *Observable) UnregisterObserver(observer Observer) {
	o.mu.Lock()
	defer o.mu.Unlock()
	for i, obs := range o.observers {
		if obs == observer {
			o.observers = append(o.observers[:i], o.observers[i+1:]...)
			break
		}
	}
}

func (o *Observable) NotifyObservers(message any) {
	o.mu.Lock()
	defer o.mu.Unlock()
	fmt.Println("Notifying:", len(o.observers))
	for _, observer := range o.observers {
		observer.Update(message)
	}
}

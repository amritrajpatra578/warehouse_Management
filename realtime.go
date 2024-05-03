package main

type Publisher interface {
	subscribe(Subscriber) error
	unsubscribe(Subscriber) error
	notify()
}

type Subscriber interface {
	Update([]Product)
	Id() string
}

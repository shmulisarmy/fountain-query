package singlevaluesubs

type Broadcaster[T any] struct {
	subscribers []func(T)
}

func (receiver *Broadcaster[T]) Subscribe(subscriber func(T)) *Broadcaster[T] {
	receiver.subscribers = append(receiver.subscribers, subscriber)
	return receiver
}

func (receiver *Broadcaster[T]) Broadcast(value T) {
	for _, subscriber := range receiver.subscribers {
		subscriber(value)
	}
}

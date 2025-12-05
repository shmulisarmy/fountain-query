package singlevaluesubs

// Count keeps track of the number of items and broadcasts the count
// T is the type of items being counted, but the count itself is always int64
type Count[T any] struct {
	Broadcaster[int64]
	count int64
}

func NewCount[T any]() *Count[T] {
	return &Count[T]{
		Broadcaster: Broadcaster[int64]{},
		count:       0,
	}
}

func (receiver *Count[T]) On_add(item T) {
	receiver.count++
	receiver.Broadcast(receiver.count)
}

func (receiver *Count[T]) On_remove(item T) {
	receiver.count--
	receiver.Broadcast(receiver.count)
}

func (receiver *Count[T]) On_update(oldItem, newItem T) {
	receiver.Broadcast(receiver.count)
}

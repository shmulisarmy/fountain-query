package utils

type CappedKeyValueArray[T any] struct {
	keys         []string
	values       []T
	constant_cap int
}

func NewKeyValueArray[T any](constant_cap int) *CappedKeyValueArray[T] {
	return &CappedKeyValueArray[T]{
		keys:         make([]string, 0, constant_cap),
		values:       make([]T, 0, constant_cap),
		constant_cap: constant_cap,
	}
}

func (this *CappedKeyValueArray[T]) Add(key string, value T) {
	if len(this.keys) >= this.constant_cap {
		panic("key_value_array is full, next time create it with a larger capacity")
	}
	this.keys = append(this.keys, key)
	this.values = append(this.values, value)
}

func (this *CappedKeyValueArray[T]) Get(key string) *T { //getting back a pointer is safe because the array is never resized
	for i, item := range this.keys {
		if item == key {
			return &this.values[i] //getting back a pointer is safe because the array is never resized
		}
	}
	panic("key " + key + " not found")
}

package singlevaluesubs

import (
	"reflect"
)

// Sum calculates the sum of a numeric field across all items
// T is the type of items containing the field to sum
// The sum is always int64
// FieldToSum must be set to the name of the field to sum
// The field must be of type int, int64, float32, or float64
// Panics if the field doesn't exist or is not a numeric type
type Sum[T any] struct {
	Broadcaster[int64]
	sum        int64
	FieldToSum string
}

func NewSum[T any](fieldToSum string) *Sum[T] {
	return &Sum[T]{
		Broadcaster: Broadcaster[int64]{},
		sum:         0,
		FieldToSum:  fieldToSum,
	}
}

func (receiver *Sum[T]) getFieldValue(item T) int64 {
	val := reflect.ValueOf(item)
	field := val.FieldByName(receiver.FieldToSum)
	if !field.IsValid() {
		panic("field not found: " + receiver.FieldToSum)
	}

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()
	case reflect.Float32, reflect.Float64:
		return int64(field.Float())
	default:
		panic("field is not a numeric type: " + field.Kind().String())
	}
}

func (receiver *Sum[T]) On_add(item T) {
	receiver.sum += receiver.getFieldValue(item)
	receiver.Broadcast(receiver.sum)
}

func (receiver *Sum[T]) On_remove(item T) {
	receiver.sum -= receiver.getFieldValue(item)
	receiver.Broadcast(receiver.sum)
}

func (receiver *Sum[T]) On_update(oldItem, newItem T) {
	oldValue := receiver.getFieldValue(oldItem)
	newValue := receiver.getFieldValue(newItem)
	receiver.sum = receiver.sum - oldValue + newValue
	receiver.Broadcast(receiver.sum)
}

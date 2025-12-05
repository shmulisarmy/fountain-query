package singlevaluesubs

import (
	"fmt"
	"reflect"
)

// Avg calculates the average of a numeric field across all items
// T is the type of items containing the field to average
// The average is always float64
// FieldToSum must be set to the name of the field to average
// The field must be of type int, int64, float32, or float64
// Panics if the field doesn't exist or is not a numeric type
type Avg[T any] struct {
	Broadcaster[float64]
	sum        int64
	totalCount int64
	FieldToSum string
}

func NewAvg[T any](fieldToSum string) *Avg[T] {
	return &Avg[T]{
		Broadcaster: Broadcaster[float64]{},
		sum:         0,
		totalCount:  0,
		FieldToSum:  fieldToSum,
	}
}

func (receiver *Avg[T]) getFieldValue(item T) int64 {
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

func (receiver *Avg[T]) calculateAverage() float64 {
	if receiver.totalCount == 0 {
		return 0 // Return 0 when there are no items
	}
	return float64(receiver.sum) / float64(receiver.totalCount)
}

func (receiver *Avg[T]) On_add(item T) {
	receiver.sum += receiver.getFieldValue(item)
	receiver.totalCount++
	receiver.Broadcast(receiver.calculateAverage())
}

func (receiver *Avg[T]) On_remove(item T) {
	if receiver.totalCount <= 0 {
		panic(fmt.Sprintf("inconsistent state: totalCount is %d but trying to remove an item", receiver.totalCount))
	}
	receiver.sum -= receiver.getFieldValue(item)
	receiver.totalCount--
	receiver.Broadcast(receiver.calculateAverage())
}

func (receiver *Avg[T]) On_update(oldItem, newItem T) {
	if receiver.totalCount < 1 {
		panic("cannot update item: no items in the collection")
	}
	oldValue := receiver.getFieldValue(oldItem)
	newValue := receiver.getFieldValue(newItem)
	receiver.sum = receiver.sum - oldValue + newValue
	receiver.Broadcast(receiver.calculateAverage())
}

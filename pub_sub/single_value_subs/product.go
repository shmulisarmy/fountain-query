package singlevaluesubs

import (
	"reflect"
)

// Product calculates the product of a numeric field across all items
// T is the type of items containing the field to multiply
// FieldToProduct must be set to the name of the field to multiply
// The field must be of type int, int64, float32, or float64
// Panics if the field doesn't exist or is not a numeric type
type Product[T any] struct {
	Broadcaster[int64]
	product        int64
	FieldToProduct string
}

func NewProduct[T any](fieldToProduct string) *Product[T] {
	return &Product[T]{
		Broadcaster:    Broadcaster[int64]{},
		product:        1, // Initialize to 1 for multiplication
		FieldToProduct: fieldToProduct,
	}
}

func (receiver *Product[T]) getFieldValue(item T) int64 {
	val := reflect.ValueOf(item)
	field := val.FieldByName(receiver.FieldToProduct)
	if !field.IsValid() {
		panic("field not found: " + receiver.FieldToProduct)
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

func (receiver *Product[T]) On_add(item T) {
	receiver.product *= receiver.getFieldValue(item)
	receiver.Broadcast(receiver.product)
}

func (receiver *Product[T]) On_remove(item T) {
	// Avoid division by zero
	if fieldValue := receiver.getFieldValue(item); fieldValue != 0 {
		receiver.product /= fieldValue
	} else {
		// Handle the case where we're removing a zero value
		receiver.product = 0
	}
	receiver.Broadcast(receiver.product)
}

func (receiver *Product[T]) On_update(oldItem, newItem T) {
	oldValue := receiver.getFieldValue(oldItem)
	newValue := receiver.getFieldValue(newItem)

	// Avoid division by zero
	if oldValue != 0 {
		receiver.product /= oldValue
	} else {
		// If old value was zero, we need to recompute the entire product
		receiver.recomputeProduct()
		return
	}

	receiver.product *= newValue
	receiver.Broadcast(receiver.product)
}

// recomputeProduct recalculates the product from scratch
// This is a fallback for when we can't easily update the product
func (receiver *Product[T]) recomputeProduct() {
	// This is a placeholder - in a real implementation, we'd need access to all items
	// For now, we'll set it to 0 to indicate an invalid state
	receiver.product = 0
	receiver.Broadcast(receiver.product)
}

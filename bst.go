package main

// Node represents a key-value pair in the BST
type Node[ValueType any] struct {
	Key   int
	Value ValueType
	Left  *Node[ValueType]
	Right *Node[ValueType]
}

// BSTMap represents the ordered map
type BSTMap[ValueType any] struct {
	Root *Node[ValueType]
}

// Insert inserts or updates a key-value pair
func (m *BSTMap[ValueType]) Insert(key int, value ValueType) {
	m.Root = m.Root.insertNode(key, value)
}

func (node *Node[ValueType]) insertNode(key int, value ValueType) *Node[ValueType] {
	if node == nil {
		return &Node[ValueType]{Key: key, Value: value}
	}
	if key < node.Key {
		node.Left = node.Left.insertNode(key, value)
	} else if key > node.Key {
		node.Right = node.Right.insertNode(key, value)
	} else {
		// update existing
		node.Value = value
	}
	return node
}

// Get retrieves a value by key
func (m *BSTMap[ValueType]) Get(key int) (ValueType, bool) {
	node := m.Root
	for node != nil {
		if key == node.Key {
			return node.Value, true
		} else if key < node.Key {
			node = node.Left
		} else {
			node = node.Right
		}
	}
	var zero ValueType
	return zero, false
}

// InOrder iterates over all key-value pairs in order
func (m *BSTMap[ValueType]) InOrder(f func(int, ValueType)) {
	m.Root.inOrderNode(f)
}

func (node *Node[ValueType]) inOrderNode(f func(int, ValueType)) {
	if node == nil {
		return
	}
	node.Left.inOrderNode(f)
	f(node.Key, node.Value)
	node.Right.inOrderNode(f)
}

// Range returns all key-value pairs within [low, high]
func (m *BSTMap[ValueType]) Range(low, high int) []struct {
	Key   int
	Value ValueType
} {
	var result []struct {
		Key   int
		Value ValueType
	}
	m.Root.rangeNode(low, high, &result)
	return result
}

func (node *Node[ValueType]) rangeNode(low, high int, result *[]struct {
	Key   int
	Value ValueType
}) {
	if node == nil {
		return
	}
	if node.Key > low {
		node.Left.rangeNode(low, high, result)
	}
	if node.Key >= low && node.Key <= high {
		*result = append(*result, struct {
			Key   int
			Value ValueType
		}{node.Key, node.Value})
	}
	if node.Key < high {
		node.Right.rangeNode(low, high, result)
	}
}

func (this *BSTMap[ValueType]) apply_to_range(low, high int, f func(*Node[ValueType])) {
	if this.Root == nil {
		return
	}
	if this.Root != nil && (this.Root.Key < low || this.Root.Key > high) {
		f(this.Root)
	}
	apply_to_child_nodes_in_range(this.Root, low, high, f)
}

func apply_to_child_nodes_in_range[ValueType any](node *Node[ValueType], low, high int, f func(*Node[ValueType])) {
	if node == nil {
		return
	}
	if node.Key > low {
		f(node.Left)
		apply_to_child_nodes_in_range(node.Left, low, high, f)
	}
	if node.Key < high {
		f(node.Right)
		apply_to_child_nodes_in_range(node.Right, low, high, f)
	}
}

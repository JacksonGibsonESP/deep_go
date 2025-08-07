package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key   int
	value int
	left  *Node
	right *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &Node{key: key, value: value}
		m.size++
		return
	}

	current := m.root
	for {
		if key == current.key {
			current.value = value
			return
		} else if key < current.key {
			if current.left == nil {
				current.left = &Node{key: key, value: value}
				m.size++
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &Node{key: key, value: value}
				m.size++
				return
			}
			current = current.right
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	var parent *Node
	current := m.root
	isLeftChild := false

	// find node and its parent
	for current != nil && current.key != key {
		parent = current
		if key < current.key {
			current = current.left
			isLeftChild = true
		} else {
			current = current.right
			isLeftChild = false
		}
	}

	if current == nil {
		return // not found
	}

	m.size--

	// case with one or no children
	if current.left == nil || current.right == nil {
		var child *Node
		if current.left == nil {
			child = current.right
		} else {
			child = current.left
		}

		if parent == nil { // it's root
			m.root = child
			return
		}

		if isLeftChild {
			parent.left = child
		} else {
			parent.right = child
		}
	} else {
		// case with two child nodes
		// find min element from right child tree
		// and place instead of current
		minParent := current
		minNode := current.right
		for minNode.left != nil {
			minParent = minNode
			minNode = minNode.left
		}

		current.key = minNode.key
		current.value = minNode.value

		// free min node object
		if minParent == current {
			minParent.right = minNode.right
		} else {
			minParent.left = minNode.right
		}
	}
}

func (m *OrderedMap) Contains(key int) bool {
	current := m.root
	for current != nil {
		if key == current.key {
			return true
		} else if key < current.key {
			current = current.left
		} else {
			current = current.right
		}
	}
	return false
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	stack := make([]*Node, 0)
	current := m.root

	for current != nil || len(stack) > 0 {
		// go to the left list and fill stack
		for current != nil {
			stack = append(stack, current)
			current = current.left
		}

		// processing
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		action(current.key, current.value)

		// go to the right tree
		current = current.right
	}
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

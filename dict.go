package rest

import (
	"sort"
)

// Filter tests if a value fulfills the given logic.
type Filter func(Model) bool

// Dict is a container with guaranteed key ordering.
type Dict struct {
	Keys   []string
	Values []Model
}

// NewDictWithCap creates a new Dict object with given capacity.
func NewDictWithCap(n int) *Dict {
	return &Dict{Keys: make([]string, 0, n), Values: make([]Model, 0, n)}
}

// NewDict creates a new Dict object.
func NewDict() *Dict {
	return NewDictWithCap(8)
}

// Index finds the index closest to the given key.
func (d *Dict) Index(key string) int {
	return sort.Search(len(d.Keys), func(i int) bool {
		return d.Keys[i] >= key
	})
}

// Get a value for a given key. Returns nil if the key does not exist.
func (d *Dict) Get(key string) Model {
	if i := d.Index(key); i < d.Len() && d.Keys[i] == key {
		return d.Values[i]
	}
	return nil
}

// Set a value for the given key. Returns false if the key does not exist.
func (d *Dict) Set(key string, value Model) bool {
	if i := d.Index(key); i < d.Len() && d.Keys[i] == key {
		d.Values[i] = value
		return true
	}
	return false
}

// Insert a value for the given key. Returns false if the key exists.
func (d *Dict) Insert(key string, value Model) bool {
	if d.Len() == 0 {
		d.Keys = append(d.Keys, key)
		d.Values = append(d.Values, value)
		return true
	}

	i := d.Index(key)

	if i == d.Len() {
		d.Keys = append(d.Keys, key)
		d.Values = append(d.Values, value)

		return true
	}

	if d.Keys[i] != key {
		d.Keys = append(d.Keys, "")
		copy(d.Keys[i+1:], d.Keys[i:])
		d.Keys[i] = key

		d.Values = append(d.Values, nil)
		copy(d.Values[i+1:], d.Values[i:])
		d.Values[i] = value

		return true
	}

	return false
}

// Remove a value for the given key. Returns nil if the key does not exist.
func (d *Dict) Remove(key string) Model {
	if i := d.Index(key); i < d.Len() && d.Keys[i] == key {
		copy(d.Keys[i:], d.Keys[i+1:])
		d.Keys[len(d.Keys)-1] = ""
		d.Keys = d.Keys[:len(d.Keys)-1]

		model := d.Values[i]

		copy(d.Values[i:], d.Values[i+1:])
		d.Values[len(d.Values)-1] = nil
		d.Values = d.Values[:len(d.Values)-1]

		return model
	}

	return nil
}

// ClearWithCap initializes the content with the given capacity.
func (d *Dict) ClearWithCap(n int) {
	d.Keys = make([]string, 0, n)
	d.Values = make([]Model, 0, n)
}

// Clear the entire content.
func (d *Dict) Clear() {
	d.ClearWithCap(8)
}

// Search indices for values matching the filter provided.
func (d *Dict) Search(f Filter) []int {
	indices := make([]int, 0, 8)

	for i := range d.Keys {
		if f(d.Values[i]) {
			indices = append(indices, i)
		}
	}

	return indices
}

// Len returns the length of the Dict (keys).
func (d *Dict) Len() int {
	return len(d.Keys)
}

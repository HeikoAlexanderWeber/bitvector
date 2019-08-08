package vector

import (
	"errors"
)

// component alias representing the most atomic block of information.
// Primarily used as definition of the scaling level (bit length) of the vector.
type component = byte

const (
	maxcomponent = 8                           // bit length of component type
	maxSize      = int(int32(^uint32(0) >> 1)) // due to array indices etc.
)

// The Bitvector struct stores booleans optimized for storage space efficiently
// inside it's fields. Instead of using large boolean types of 1 component per boolean,
// this type stores booleans as bits inside larger numbers.
// Default scaling size is 8bit (uint8, component).
type Bitvector struct {
	occupied int
	data     []component
}

// New func creates a new Bitvector and returns a reference to it.
func New() *Bitvector {
	return &Bitvector{
		occupied: 0,
		data:     []component{},
	}
}

func (v *Bitvector) resize(count int) {
	// optimization for clearing
	if count <= 0 {
		v.data = []component{}
		return
	}
	// calc dimensions
	neededSize := (count-1)/maxcomponent + 1
	actualSize := len(v.data)
	remaining := neededSize - actualSize
	if remaining == 0 { // both are the same
		return
	}
	if remaining > 0 { // need to grow
		for i := 0; i < remaining; i++ {
			v.data = append(v.data, component(0))
		}
	} else { // need to shrink
		v.data = v.data[0 : actualSize+remaining] // remaining is negative
	}
}

// calculate location
func (v *Bitvector) location(count int) (int, byte) {
	x := count / maxcomponent
	y := count % maxcomponent
	return x, byte(y)
}

// Push func appends a boolean to the imaginary stack.
// This method returns an error if the imaginary stack would exceeded
// its maximum size.
func (v *Bitvector) Push(val ...bool) error {
	if v.occupied+len(val) > maxSize {
		return errors.New("can not push, maximum size would be exceeded")
	}
	v.resize(v.occupied + len(val))
	for _, b := range val {
		x, y := v.location(v.occupied)
		if b {
			v.data[x] |= (1 << y)
		}
		v.occupied++
	}
	return nil
}

// Pop func removes the last n items from the imaginary stack.
// This method returns an error if the amount of requested items exceeds the imaginary stack size.
func (v *Bitvector) Pop(n int) ([]bool, error) {
	if n > v.occupied {
		return nil, errors.New("can not pop, too many items requested")
	}
	vals := []bool{}
	for i := 0; i < n; i++ {
		x, y := v.location(v.occupied - 1)
		vals = append(vals, (v.data[x]>>y) > 0)
		v.data[x] &= ^(1 << y)

		v.resize(v.occupied - 1)
		v.occupied--
	}
	return vals, nil
}

// Get func gets a specified boolean at the given index.
// This method returns an error if the index is invalid.
func (v *Bitvector) Get(index ...int) ([]bool, error) {
	vals := []bool{}
	for _, i := range index {
		if i < 0 || i >= v.occupied {
			return nil, errors.New("invalid index")
		}
		x, y := v.location(i)
		vals = append(vals, (v.data[x]&(1<<y)) > 0)
	}
	return vals, nil
}

// Set func sets a specific index in the imaginary stack to the given value.
// This method returns an error if the index is invalid.
func (v *Bitvector) Set(index int, val bool) error {
	if index < 0 || index >= v.occupied {
		return errors.New("invalid index")
	}
	x, y := v.location(index)
	if val {
		v.data[x] |= (1 << y)
	} else {
		v.data[x] &= ^(1 << y)
	}
	return nil
}

// Clear func clears all data and resets the instance to an empty state.
func (v *Bitvector) Clear() {
	v.occupied = 0
	v.resize(0)
}

// AsArray func returns an array representation of the current state.
func (v *Bitvector) AsArray() []bool {
	result := []bool{}
	for i := 0; i < v.occupied; i++ {
		x, y := v.location(i)
		result = append(result, (v.data[x]&(1<<y)) > 0)
	}
	return result
}

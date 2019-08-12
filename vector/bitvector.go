package vector

import (
	"errors"
)

// component alias representing the most atomic block of information.
// Primarily used as definition of the scaling level (bit length) of the vector.
type component = byte

const (
	maxcomponent = 8 // bit length of component type
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
	if n < 0 || n > v.occupied {
		return nil, errors.New("can not pop, too many items requested")
	}

	indices := []int{}
	i := 0
	start := v.occupied - n
	for i < n {
		indices = append(indices, start+i)
		i++
	}
	data, err := v.Get(indices...)

	v.resize(v.occupied - n)
	v.occupied -= n

	return data, err
}

// PopOne is a convenience wrapper-function of Pop() to only pop one item.
func (v *Bitvector) PopOne() (bool, error) {
	d, e := v.Pop(1)
	if e != nil {
		return false, e
	}
	return d[0], e
}

// Insert func inserts the given data into the storage at a given index.
// It returns an error if the index is invalid.
func (v *Bitvector) Insert(index int, vals ...bool) error {
	if index < 0 || index > v.occupied {
		return errors.New("can not insert, invalid index was provided")
	}

	v.resize(v.occupied + len(vals))
	v.occupied += len(vals)

	offset := 0
	for offset < len(vals) {
		oldV, _ := v.Get(index + offset)
		v.Set(index+offset, vals[offset])
		v.Set(index+offset+len(vals), oldV[0])
		offset++
	}

	return nil
}

// Delete func deletes all elements with the given indices. Returns error
// if an index is invalid.
func (v *Bitvector) Delete(indices ...int) error {
	for _, i := range indices {
		if i < 0 || i > v.occupied {
			return errors.New("can not insert, invalid index was provided")
		}
	}
	deletedItems := 0
	for _, i := range indices {
		realIndex := i - deletedItems
		for realIndex < v.occupied-1 {
			next, _ := v.Get(realIndex + 1)
			v.Set(realIndex, next[0])
			realIndex++
		}
		deletedItems++
	}
	v.resize(v.occupied - len(indices))
	v.occupied -= len(indices)
	return nil
}

// DeleteRange func for convenience. This method deletes chunks of data.
// The first element to be deleted is at the given index. Every item in the range of
// [index, count) is deleted.
func (v *Bitvector) DeleteRange(index int, count int) error {
	if index < 0 || index+count > v.occupied {
		return errors.New("can not insert, invalid index was provided")
	}
	indices := []int{}
	for i := index; i < index+count; i++ {
		indices = append(indices, i)
	}
	return v.Delete(indices...)
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

// GetOne func is a convenience wrapper-function for Get() to only retrieve one item.
func (v *Bitvector) GetOne(index int) (bool, error) {
	d, e := v.Get(index)
	if e != nil {
		return false, e
	}
	return d[0], e
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

// Length func returns the amount of occupied bits / the amount of booleans stored.
func (v *Bitvector) Length() int {
	return v.occupied
}

// Size func returns the byte-size of the storage. This method only takes the storage
// itself into account. The constant overhead of the fields of this type are ignored.
func (v *Bitvector) Size() int {
	return len(v.data) * maxcomponent / 8 // from bit to byte
}

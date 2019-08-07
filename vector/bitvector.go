package vector

import "errors"

const (
	maxComponent = 8 // bit length of byte
	maxSize      = ^int(0)
)

// Bitvector struct
type Bitvector struct {
	occupied int
	data     []byte
}

// New func
func New() *Bitvector {
	return &Bitvector{
		occupied: 0,
		data:     []byte{},
	}
}

func (v *Bitvector) ensureSize(count int) {
	// optimization for clearing
	if count == 0 {
		v.data = []byte{}
		return
	}
	neededSize := (count-1)/maxComponent + 1
	actualSize := len(v.data)
	remaining := neededSize - actualSize
	if remaining == 0 { // both are the same
		return
	}
	if remaining > 0 { // need to grow
		for i := 0; i < remaining; i++ {
			v.data = append(v.data, byte(0))
		}
	} else { // need to shrink
		v.data = v.data[0 : len(v.data)+remaining]
	}
}

func (v *Bitvector) location(count int) (int, byte) {
	x := count / maxComponent
	y := count % maxComponent
	return x, byte(y)
}

// Push func
func (v *Bitvector) Push(val bool) error {
	if v.occupied == maxSize {
		return errors.New("maximum size exceeded")
	}
	v.ensureSize(v.occupied + 1)
	x, y := v.location(v.occupied)
	if val {
		v.data[x] |= (1 << y)
	}
	v.occupied++
	return nil
}

// Pop func
func (v *Bitvector) Pop() (bool, error) {
	if v.occupied <= 0 {
		return false, errors.New("vector is empty")
	}
	x, y := v.location(v.occupied - 1)
	value := (v.data[x] >> y) > 0
	v.data[x] &= ^(1 << y)

	v.ensureSize(v.occupied - 1)
	v.occupied--
	return value, nil
}

// Set func
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

// Clear func
func (v *Bitvector) Clear() {
	v.occupied = 0
	v.ensureSize(0)
}

// AsArray func
func (v *Bitvector) AsArray() []bool {
	result := []bool{}
	cnt := 0
	for _, x := range v.data {
		for y := byte(0); y < maxComponent; y++ {
			cnt++
			if cnt > v.occupied {
				return result
			}
			result = append(result, (x&(1<<y)) > 0)
		}
	}
	return result
}

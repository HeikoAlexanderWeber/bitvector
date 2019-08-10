package vector

import (
	"errors"
	"testing"
)

var casesTestPush = []struct {
	nr   int
	size int
}{
	{0, 0},
	{1, 1},
	{2, 1},
	{7, 1},
	{8, 1},
	{9, 2},
	{15, 2},
	{16, 2},
	{17, 3},
	{64, 8},
}

var casesTestPop = []struct {
	nrIn         int
	nrOut        int
	failsAt      int
	sizeAfterPop int
}{

	{1, 1, -1, 0},
	{2, 1, -1, 1},
	{9, 1, -1, 1},
	{8, 8, -1, 0},
	{0, 1, 0, 0},
	{8, 9, 8, 0},
}

func TestPush(t *testing.T) {
	for _, c := range casesTestPush {
		v := New()
		for i := 0; i < c.nr; i++ {
			v.Push(true)
		}
		if len(v.data) != c.size || v.occupied != c.nr {
			t.Errorf("wrong internal size (%v instead of %v) for %v (case %v)",
				len(v.data), c.size, c.nr, c)
		}
	}
}

func TestPop(t *testing.T) {
	for _, c := range casesTestPop {
		v := New()
		for i := 0; i < c.nrIn; i++ {
			v.Push(true)
		}
		for i := 0; i < c.nrOut; i++ {
			b, err := v.Pop(1)
			if c.failsAt != i && err != nil {
				t.Error(err)
			} else if c.failsAt == i && err == nil {
				t.Errorf("did not throw exception as expected by case %v", c)
			} else if c.failsAt != i && !b[0] {
				t.Error(errors.New("retrieved wrong value from pop"))
			}
		}
		expectedOccupied := c.nrIn - c.nrOut
		if expectedOccupied < 0 {
			expectedOccupied = 0
		}
		if len(v.data) != c.sizeAfterPop || v.occupied != expectedOccupied {
			t.Errorf("wrong size for case %v", c)
		}
	}
}

func TestPopMany(t *testing.T) {
	v := New()
	expected := []bool{true, false, true, false, true, false, true, false, true, false}
	v.Push(expected...)
	arr, err := v.Pop(10)
	if err != nil {
		t.Error(err)
	}
	if len(arr) != len(expected) {
		t.Error(errors.New("array does not match expected array"))
	}
	for i, e := range expected {
		if arr[i] != e {
			t.Error(errors.New("array does not match expected array"))
		}
	}
}

func TestGet(t *testing.T) {
	v := New()
	v.Push(true, false, true)

	raise := func(e error) {
		if e != nil {
			t.Error(e)
		}
	}

	b, err := v.Get(0, 1, 2)
	raise(err)
	if !b[0] || b[1] || !b[2] {
		t.Error(errors.New("received invalid data"))
	}

	b, err = v.Get(-1, 0)
	if err == nil {
		t.Error(errors.New("did not return expected error for invalid index"))
	}
}

func TestSet(t *testing.T) {
	v := New()
	v.Push(false, false, false, false)
	v.Set(2, true)

	if v.occupied != 4 || len(v.data) != 1 {
		t.Error(errors.New("unexpected size"))
	}
	if v.data[0] != 1<<2 {
		t.Error(errors.New("unexpeted state"))
	}

	v.Set(0, true)
	v.Set(1, true)
	v.Set(2, false)
	v.Set(3, true)

	if v.occupied != 4 || len(v.data) != 1 {
		t.Error(errors.New("unexpected size"))
	}
	if v.data[0] != 11 { // 11010000
		t.Error(errors.New("unexpeted state"))
	}

	if v.Set(-1, false) == nil {
		t.Error(errors.New("did not return expected invalid index error"))
	}
}

func TestAsArray(t *testing.T) {
	v := New()
	if len(v.AsArray()) != 0 {
		t.Error(errors.New("empty array expected"))
	}
	expected := []bool{true, false, true, false, true, false, true, false, true}
	v.Push(expected...)
	arr := v.AsArray()
	if len(arr) != len(expected) {
		t.Error(errors.New("array does not match expected array"))
	}
	for i, e := range expected {
		if arr[i] != e {
			t.Error(errors.New("array does not match expected array"))
		}
	}
}

func TestClear(t *testing.T) {
	v := New()
	for i := 0; i < 9; i++ {
		v.Push(true)
	}
	v.Clear()
	if len(v.data) > 0 || v.occupied > 0 {
		t.Error(errors.New("unexpected size"))
	}
}

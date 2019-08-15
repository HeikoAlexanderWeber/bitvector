package vector

import (
	"errors"
	"testing"
)

// =====
// Tests
// =====

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
	initial := []bool{true, false, true}
	expected := []bool{false, true}
	v.Push(initial...)
	arr, err := v.Pop(len(expected))
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

func TestPopOne(t *testing.T) {
	v := New()
	d, e := v.PopOne()
	if d || e == nil {
		t.Error(errors.New("did not throw expected error"))
	}
	v.Push(false, true)
	d, e = v.PopOne()
	if e != nil {
		t.Error(e)
	}
	if !d {
		t.Error(errors.New("received wrong data"))
	}
	d, e = v.PopOne()
	if e != nil {
		t.Error(e)
	}
	if d {
		t.Error(errors.New("received wrong data"))
	}
}

func TestInsert(t *testing.T) {
	v := New()
	if v.Insert(0, true) == nil {
		t.Error(errors.New("did not return expected error"))
	}
	v.Push(false, false)
	if v.Insert(-1, true) == nil {
		t.Error(errors.New("did not return expected error"))
	}
	if v.Insert(3, true) == nil {
		t.Error(errors.New("did not return expected error"))
	}

	if e := v.Insert(1, true, true); e != nil {
		t.Error(e)
	}
	expected := []bool{false, true, true, false}
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

func TestDelete(t *testing.T) {
	v := New()
	if v.Delete(-1) == nil {
		t.Error(errors.New("did not return expected error"))
	}
	v.Push(false, true, true, false, true)
	if v.Delete(10) == nil {
		t.Error(errors.New("did not return expected error"))
	}
	if e := v.Delete(1, 2); e != nil {
		t.Error(e)
	}
	if e := v.Delete(2); e != nil {
		t.Error(e)
	}
	// make sure that one invalid index results in the vector not being modified at all
	if v.Delete(0, 5) == nil { // 5 is invalid
		t.Error(errors.New("did not return expected error"))
	}

	expected := []bool{false, false}
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

func TestDeleteRange(t *testing.T) {
	v := New()
	if v.DeleteRange(0, 1) == nil {
		t.Error(errors.New("did not return expected error"))
	}
	v.Push(false, true, true, false, true)
	if v.Delete(0, 6) == nil { // 6 is invalid
		t.Error(errors.New("did not return expected error"))
	}
	if e := v.DeleteRange(1, 2); e != nil {
		t.Error(e)
	}
	if e := v.DeleteRange(2, 1); e != nil {
		t.Error(e)
	}
	expected := []bool{false, false}
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

func TestGetOne(t *testing.T) {
	v := New()
	if d, e := v.GetOne(0); d || e == nil {
		t.Error(errors.New("did not return expected error"))
	}
	v.Push(false, true)
	if d, e := v.GetOne(2); d || e == nil {
		t.Error(errors.New("did not return expected error"))
	}
	d, e := v.GetOne(1)
	if e != nil {
		t.Error(e)
	}
	if !d {
		t.Error(errors.New("received wrong data"))
	}
	d, e = v.GetOne(0)
	if e != nil {
		t.Error(e)
	}
	if d {
		t.Error(errors.New("received wrong data"))
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

func TestLength(t *testing.T) {
	v := New()
	if v.Length() != 0 {
		t.Error(errors.New("returned wrong length"))
	}
	v.Push(true, false, true)
	if v.Length() != 3 {
		t.Error(errors.New("returned wrong length"))
	}
}

func TestSize(t *testing.T) {
	v := New()
	if v.Size() != 0 {
		t.Error(errors.New("returned wrong size"))
	}
	v.Push(true, false, true, false, true, false, true, false, true)
	if v.Size() != 2 {
		t.Error(errors.New("returned wrong size"))
	}
}

// ==========
// Benchmarks
// ==========

func BenchmarkNew(b *testing.B) {
	x := New()
	for n := 0; n < b.N; n++ {
		x = New()
	}
	x.Length()
	b.ReportAllocs()
}

func BenchmarkPush(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ReportAllocs()
}

func BenchmarkPop(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Pop(1)
	}
	b.ReportAllocs()
}

func BenchmarkPopOne(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.PopOne()
	}
	b.ReportAllocs()
}

func BenchmarkInsertFront(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Insert(0, true)
	}
	b.ReportAllocs()
}

func BenchmarkInsertTail(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Insert(n, true)
	}
	b.ReportAllocs()
}

func BenchmarkDeleteFront(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Delete(0)
	}
	b.ReportAllocs()
}

func BenchmarkDeleteTail(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Delete(b.N - n - 1)
	}
	b.ReportAllocs()
}

func BenchmarkDeleteRangeFront(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.DeleteRange(0, 1)
	}
	b.ReportAllocs()
}

func BenchmarkDeleteRangeTail(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.DeleteRange(n, 1)
	}
	b.ReportAllocs()
}

func BenchmarkGet(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Get(n)
	}
	b.ReportAllocs()
}

func BenchmarkGetOne(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.GetOne(n)
	}
	b.ReportAllocs()
}

func BenchmarkSet(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Set(n, false)
	}
	b.ReportAllocs()
}

func BenchmarkClear(b *testing.B) {
	vecs := []*Bitvector{}
	for n := 0; n < b.N; n++ {
		v := New()
		v.Push(true, true, true, true)
		vecs = append(vecs, v)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		vecs[n].Clear()
	}
	b.ReportAllocs()
}

func BenchmarkAsArray(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true, true, true, true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.AsArray()
	}
	b.ReportAllocs()
}

func BenchmarkLength(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Length()
	}
	b.ReportAllocs()
}

func BenchmarkSize(b *testing.B) {
	v := New()
	for n := 0; n < b.N; n++ {
		v.Push(true)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Size()
	}
	b.ReportAllocs()
}

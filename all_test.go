package sortutil

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	day = 24 * time.Hour
)

type Item struct {
	Id    int64
	Name  string
	Date  time.Time
	Valid bool
}

type SortableItems []Item

func (s SortableItems) Len() int {
	return len(s)
}

func (s SortableItems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortableItems) Less(i, j int) bool {
	return s[i].Id > s[j].Id
}

type SortablePointers []*Item

func (s SortablePointers) Len() int {
	return len(s)
}

func (s SortablePointers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortablePointers) Less(i, j int) bool {
	return s[i].Id > s[j].Id
}

func names() []string {
	return []string{"A", "C", "a", "b", "d", "g", "h", "y", "z"}
}

func namesInsensitive() []string {
	return []string{"A", "a", "b", "C", "d", "g", "h", "y", "z"}
}

var now = time.Now()

func dates() []time.Time {
	return []time.Time{
		now.Add(-3 * day),
		now.Add(-2 * day),
		now.Add(-1 * day),
		now,
		now.Add(1 * day),
		now.Add(2 * day),
		now.Add(4 * day),
		now.Add(5 * day),
		now.Add(7 * day),
	}
}

func items() []Item {
	n := names()
	d := dates()
	is := []Item{
		{6, n[4], d[0], true},
		{1, n[3], d[5], true},
		{9, n[1], d[6], true},
		{3, n[8], d[2], false},
		{7, n[7], d[8], true},
		{2, n[2], d[4], false},
		{8, n[0], d[1], false},
		{5, n[5], d[7], false},
		{4, n[6], d[3], true},
	}
	return is
}

func pointers() []*Item {
	n := names()
	d := dates()
	is := []*Item{
		{6, n[4], d[0], true},
		{1, n[3], d[5], true},
		{9, n[1], d[6], true},
		{3, n[8], d[2], false},
		{7, n[7], d[8], true},
		{2, n[2], d[4], false},
		{8, n[0], d[1], false},
		{5, n[5], d[7], false},
		{4, n[6], d[3], true},
	}
	return is
}

func nestedIntSlice() [][]int {
	return [][]int{
		{4, 5, 1},
		{2, 1, 7},
		{9, 3, 3},
		{1, 6, 2},
	}
}

func TestSortReverseInterface(t *testing.T) {
	is := items()
	SortReverseInterface(SortableItems(is))
	for i, v := range is {
		if v.Id != int64(i+1) {
			t.Errorf("is[%d].Id is not %d, but %d", i, i+1, v.Id)
		}
	}
}

func TestSortByStringFieldAscending(t *testing.T) {
	is := items()
	AscByField(is, "Name")
	c := names()
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name is not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldDescending(t *testing.T) {
	is := items()
	DescByField(is, "Name")
	c := names()
	ReverseInterface(sort.StringSlice(c))
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name is not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldCaseInsensitiveAscending(t *testing.T) {
	is := items()
	CiAscByField(is, "Name")
	c := namesInsensitive()
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name is not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldCaseInsensitiveDescending(t *testing.T) {
	is := items()
	CiDescByField(is, "Name")
	c := namesInsensitive()
	ReverseInterface(sort.StringSlice(c))
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name is not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByInt64FieldAscending(t *testing.T) {
	is := items()
	AscByField(is, "Id")
	for i, v := range is {
		if v.Id != int64(i+1) {
			t.Errorf("is[%d].Id is not %d, but %d", i, i+1, v.Id)
		}
	}
}

func TestSortByInt64FieldDescending(t *testing.T) {
	is := items()
	DescByField(is, "Id")
	l := len(is)
	for i, v := range is {
		if v.Id != int64(l-i) {
			t.Errorf("is[%d].Id is not %d, but %d", i, l-i, v.Id)
		}
	}
}

func TestSortByIntIndexAscending(t *testing.T) {
	is := nestedIntSlice()
	AscByIndex(is, 2)
	if !sort.IntsAreSorted([]int{is[0][2], is[1][2], is[2][2], is[3][2]}) {
		t.Errorf("Nested int slice is not sorted by index 2 in child slices: %v", is)
	}
}

func TestSortIntArray(t *testing.T) {
	return // TEMP: Disabled
	ints := [...]int{4, 3, 1, 5, 2}
	Asc(ints)
	if !sort.IntsAreSorted(ints[:]) {
		t.Errorf("Array ints weren't sorted: %v", ints)
	}
}

func TestSortByTimeFieldAscending(t *testing.T) {
	is := items()
	AscByField(is, "Date")
	c := dates()
	for i, v := range is {
		if !v.Date.Equal(c[i]) {
			t.Errorf("is[%d].Date is not %v, but %v", i, c[i], v.Date)
		}
	}
}

func TestSortByTimeFieldDescending(t *testing.T) {
	is := items()
	DescByField(is, "Date")
	c := dates()
	l := len(is)
	for i, v := range is {
		if !v.Date.Equal(c[l-i-1]) {
			t.Errorf("is[%d].Date is not %v, but %v", i, c[l-i], v.Date)
		}
	}
}

type TestStruct struct {
	TimePtr    *time.Time
	Invalid    InvalidType
	unexported int
}

type InvalidType struct {
	Foo string
	Bar int
}

func testStructs() []TestStruct {
	return []TestStruct{
		{
			TimePtr:    &now,
			Invalid:    InvalidType{"foo", 123},
			unexported: 5,
		},
	}
}

func testPointers() []*TestStruct {
	return []*TestStruct{
		{
			TimePtr:    &now,
			Invalid:    InvalidType{"foo", 123},
			unexported: 5,
		},
	}
}

func TestSortInvalidType(t *testing.T) {
	// Sorting an invalid type should cause a panic
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Sorting an unrecognized type didn't cause a panic")
		}
	}()
	is := testStructs()
	AscByField(is, "Invalid")
}

func TestSortUnexportedType(t *testing.T) {
	// Sorting an unexported type should cause a panic
	// TODO: This should test on a field outside the package
	return // TEMP
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Sorting an unexported type didn't cause a panic")
		}
	}()
	is := testStructs()
	AscByField(is, "unexported")
}

func TestSortByPointer(t *testing.T) {
	// Sorting by a pointer type shouldn't cause a panic
	is := testStructs()
	AscByField(is, "TimePtr")
}

func TestSortPointerSlice(t *testing.T) {
	// Sorting a slice of pointers shouldn't cause a panic
	is := pointers()
	AscByField(is, "Id")
}

func TestSortPointerSliceByPointer(t *testing.T) {
	// Sorting a slice of pointers by a pointer type shouldn't cause a panic
	is := testPointers()
	AscByField(is, "TimePtr")
}

func TestReverse(t *testing.T) {
	ints := []int{4, 2, 6, 4, 8}
	correct := []int{8, 4, 6, 2, 4}
	Reverse(ints)
	if !reflect.DeepEqual(ints, correct) {
		t.Errorf("Reversed slice was not %v: %v", correct, ints)
	}
}

func BenchmarkSortByInt64(b *testing.B) {
	var is []Item
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		is = append(is, items()...)
	}
	b.StartTimer()
	sort.Sort(SortableItems(is))
}

func BenchmarkSortPointersByInt64(b *testing.B) {
	var is []*Item
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		is = append(is, pointers()...)
	}
	b.StartTimer()
	sort.Sort(SortablePointers(is))
}

func BenchmarkSortReverseByInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		is := items()
		SortReverseInterface(SortableItems(is))
	}
}

func benchmarkInts(n int) []int {
	ints := make([]int, n, n)
	v := 0
	for i := 0; i < n; i++ {
		ints[i] = v
		v++
		if v > 10 {
			v = 0
		}
	}
	return ints
}

func benchmarkInt64s(n int) []int64 {
	ints := make([]int64, n, n)
	v := int64(0)
	for i := 0; i < n; i++ {
		ints[i] = v
		v++
		if v > 10 {
			v = 0
		}
	}
	return ints
}

func BenchmarkSortInts(b *testing.B) {
	b.StopTimer()
	ints := benchmarkInts(b.N)
	b.StartTimer()
	sort.Sort(sort.IntSlice(ints))
}

func BenchmarkAscInts(b *testing.B) {
	b.StopTimer()
	ints := benchmarkInts(b.N)
	b.StartTimer()
	Asc(ints)
}

func BenchmarkDescInts(b *testing.B) {
	b.StopTimer()
	ints := benchmarkInts(b.N)
	b.StartTimer()
	Desc(ints)
}

func BenchmarkAscInt64s(b *testing.B) {
	b.StopTimer()
	ints := benchmarkInt64s(b.N)
	b.StartTimer()
	Asc(ints)
}

func BenchmarkDescInt64s(b *testing.B) {
	b.StopTimer()
	ints := benchmarkInt64s(b.N)
	b.StartTimer()
	Desc(ints)
}

func benchmarkItems(l int) []Item {
	is := make([]Item, l, l)
	n := names()
	d := dates()
	v := 0
	for i := range is {
		is[i] = Item{int64(v), n[v], d[v], true}
		v++
		if v > 5 {
			v = 0
		}
	}
	return is
}

func benchmarkPointers(l int) []*Item {
	is := make([]*Item, l, l)
	n := names()
	d := dates()
	v := 0
	for i := range is {
		is[i] = &Item{int64(v), n[v], d[v], true}
		v++
		if v > 5 {
			v = 0
		}
	}
	return is
}

func BenchmarkAscByInt64(b *testing.B) {
	b.StopTimer()
	is := benchmarkItems(b.N)
	b.StartTimer()
	AscByField(is, "Id")
}

func BenchmarkAscPointersByInt64(b *testing.B) {
	b.StopTimer()
	is := benchmarkPointers(b.N)
	b.StartTimer()
	AscByField(is, "Id")
}

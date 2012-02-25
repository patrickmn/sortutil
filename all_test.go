package sortutil

import (
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

func nestedIntSlice() [][]int {
	return [][]int{
		{4, 5, 1},
		{2, 1, 7},
		{9, 3, 3},
		{1, 6, 2},
	}
}

func TestSortReverse(t *testing.T) {
	is := items()
	SortReverse(SortableItems(is))
	for i, v := range is {
		if v.Id != int64(i+1) {
			t.Errorf("is[%d].Id was not %d, but %d", i, i+1, v.Id)
		}
	}
}

// func TestSortNoGetterReverse(t *testing.T) {

func TestSortByStringFieldAscending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Name"), Ascending)
	c := names()
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name was not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldDescending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Name"), Descending)
	c := names()
	Reverse(sort.StringSlice(c))
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name was not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldCaseInsensitiveAscending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Name"), CaseInsensitiveAscending)
	c := namesInsensitive()
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name was not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByStringFieldCaseInsensitiveDescending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Name"), CaseInsensitiveDescending)
	c := namesInsensitive()
	Reverse(sort.StringSlice(c))
	for i, v := range is {
		if v.Name != c[i] {
			t.Errorf("is[%d].Name was not %s, but %s", i, c[i], v.Name)
		}
	}
}

func TestSortByInt64FieldAscending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Id"), Ascending)
	for i, v := range is {
		if v.Id != int64(i+1) {
			t.Errorf("is[%d].Id was not %d, but %d", i, i+1, v.Id)
		}
	}
}

func TestSortByInt64FieldDescending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Id"), Descending)
	l := len(is)
	for i, v := range is {
		if v.Id != int64(l-i) {
			t.Errorf("is[%d].Id was not %d, but %d", i, l-i, v.Id)
		}
	}
}

func TestSortByIntIndexAscending(t *testing.T) {
	is := nestedIntSlice()
	Sort(is, IndexGetter(2), Ascending)
	if !sort.IntsAreSorted([]int{is[0][2], is[1][2], is[2][2], is[3][2]}) {
		t.Errorf("Nested int slice was not sorted by index 2 in child slices: %v", is)
	}
}

func TestSortByTimeFieldAscending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Date"), Ascending)
	c := dates()
	for i, v := range is {
		if !v.Date.Equal(c[i]) {
			t.Errorf("is[%d].Date was not %v, but %v", i, c[i], v.Date)
		}
	}
}

func TestSortByTimeFieldDescending(t *testing.T) {
	is := items()
	Sort(is, FieldGetter("Date"), Descending)
	c := dates()
	l := len(is)
	for i, v := range is {
		if !v.Date.Equal(c[l-i-1]) {
			t.Errorf("is[%d].Date was not %v, but %v", i, c[l-i], v.Date)
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

func TestSortInvalidType(t *testing.T) {
	// Sorting an invalid type should cause a panic
	defer func() {
		if x := recover(); x == nil {
			t.Fatal("Sorting an unrecognized type didn't cause a panic")
		}
	}()
	s := testStructs()
	Sort(s, FieldGetter("Invalid"), Ascending)
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
	s := testStructs()
	Sort(s, FieldGetter("unexported"), Ascending)
}

func TestSortPointerType(t *testing.T) {
	// Sorting a pointer type shouldn't cause a panic
	s := testStructs()
	Sort(s, FieldGetter("TimePtr"), Ascending)
}

func BenchmarkSortStructByInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		is := items()
		sort.Sort(SortableItems(is))
	}
}

func BenchmarkSortReverseStructByInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		is := items()
		SortReverse(SortableItems(is))
	}
}

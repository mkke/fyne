package binding

import "errors"

var errOutOfBounds = errors.New("index out of bounds")

// DataList is the base interface for all bindable data lists.
//
// Since: 2.0.0
type DataList interface {
	DataItem
	GetItem(int) (DataItem, error)
	Length() int
}

type listBase struct {
	base
	items []DataItem
}

// GetItem returns the DataItem at the specified index.
func (b *listBase) GetItem(i int) (DataItem, error) {
	if i < 0 || i >= len(b.items) {
		return nil, errOutOfBounds
	}

	return b.items[i], nil
}

// Length returns the number of items in this data list.
func (b *listBase) Length() int {
	return len(b.items)
}

func (b *listBase) appendItem(i DataItem) {
	b.items = append(b.items, i)

	b.trigger()
}

func (b *listBase) prependItem(i DataItem) {
	b.items = append([]DataItem{i}, b.items...)

	b.trigger()
}

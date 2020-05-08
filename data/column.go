package data

import "time"

type columnType int

const (
	columnTime columnType = iota
	columnString
	columnInt
	columnFloat
)

// Column data column
type Column struct {
	index      int
	name       string
	t          columnType
	timeParse  func(string) time.Time
	timeFormat func(time.Time) string
}

// NewStringColumn create string column
func NewStringColumn(name string, idx int) Column {
	return Column{index: idx, name: name, t: columnString}
}

// NewIntColumn create int column
func NewIntColumn(name string, idx int) Column {
	return Column{index: idx, name: name, t: columnInt}
}

// NewFloatColumn create int column
func NewFloatColumn(name string, idx int) Column {
	return Column{index: idx, name: name, t: columnFloat}
}

// NewTimeColumn create time column
func NewTimeColumn(name string, idx int, parse func(string) time.Time, format func(time.Time) string) Column {
	return Column{
		index:      idx,
		name:       name,
		t:          columnTime,
		timeParse:  parse,
		timeFormat: format,
	}
}

// GetName get name of column
func (c *Column) GetName() string {
	return c.name
}

// GetIndex get index of column
func (c *Column) GetIndex() int {
	return c.index
}

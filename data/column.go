package data

import "time"

type columnType int

const (
	columnTime columnType = iota
	columnString
	columnInt
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

package data

import (
	"fmt"
	"time"
)

// Cell data cell
type Cell struct {
	t          columnType
	s          string
	i          int
	f          float64
	ts         time.Time
	timeFormat func(time.Time) string
	empty      bool
}

func (c *Cell) String() string {
	if c.empty {
		return "<null>"
	}
	switch c.t {
	case columnString:
		return c.s
	case columnInt:
		return fmt.Sprintf("%d", c.i)
	case columnTime:
		return c.timeFormat(c.ts)
	default:
		return ""
	}
}

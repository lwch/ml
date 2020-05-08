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

func (c *Cell) div(target *Cell) {
	switch c.t {
	case columnInt:
		switch target.t {
		case columnInt:
			c.t = columnFloat
			c.f = float64(c.i) / float64(target.i)
		case columnFloat:
			c.t = columnFloat
			c.f = float64(c.i) / target.f
		}
	case columnFloat:
		switch target.t {
		case columnInt:
			c.t = columnFloat
			c.f = c.f / float64(target.i)
		case columnFloat:
			c.t = columnFloat
			c.f = c.f / target.f
		}
	}
}

// Float get float value
func (c *Cell) Float() float64 {
	if c.t != columnFloat {
		return 0
	}
	return c.f
}

package data

type numberFunc func(*Data, *Column) (*Cell, bool)
type hashFunc func(*Cell) int

// Mean number func get mean value
func Mean(d *Data, c *Column) (*Cell, bool) {
	switch c.t {
	case columnInt:
		var total int
		for _, row := range d.cellsByIndex {
			total += row[c.index].i
		}
		return &Cell{
			t: c.t,
			i: total / len(d.cellsByIndex),
		}, true
	case columnFloat:
		var total float64
		for _, row := range d.cellsByIndex {
			total += row[c.index].f
		}
		return &Cell{
			t: c.t,
			f: total / float64(len(d.cellsByIndex)),
		}, true
	default:
		return nil, false
	}
}

// Max number func get max value
func Max(d *Data, c *Column) (*Cell, bool) {
	switch c.t {
	case columnInt:
		cell := d.cellsByIndex[0][c.index]
		for _, row := range d.cellsByIndex {
			if row[c.index].i > cell.i {
				cell = row[c.index]
			}
		}
		return &Cell{
			t: columnInt,
			i: cell.i,
		}, true
	case columnFloat:
		cell := d.cellsByIndex[0][c.index]
		for _, row := range d.cellsByIndex {
			if row[c.index].f > cell.f {
				cell = row[c.index]
			}
		}
		return &Cell{
			t: columnFloat,
			f: cell.f,
		}, true
	default:
		return nil, false
	}
}

// Length hash func for length
func Length(c *Cell) int {
	if c.t != columnString {
		return 0
	}
	return len(c.s)
}

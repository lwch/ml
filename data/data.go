package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"ml/constant"
	"sort"
	"strconv"
	"time"
)

// Data data
type Data struct {
	columnsByIndex map[int]*Column
	columnsByName  map[string]*Column

	cellsByIndex []map[int]*Cell
	cellsByName  []map[string]*Cell

	loaded bool
}

// NewData create data
func NewData() *Data {
	return &Data{
		columnsByIndex: make(map[int]*Column),
		columnsByName:  make(map[string]*Column),
	}
}

// AddColumn add column definition
func (d *Data) AddColumn(col Column) {
	d.columnsByIndex[col.index] = &col
	d.columnsByName[col.name] = &col
}

// GetColumnByIndex get column by index
func (d *Data) GetColumnByIndex(idx int) *Column {
	return d.columnsByIndex[idx]
}

// GetColumnByName get column by name
func (d *Data) GetColumnByName(name string) *Column {
	return d.columnsByName[name]
}

// LoadFromCSV read data from csv
func (d *Data) LoadFromCSV(r io.Reader, skipHeader bool) error {
	if len(d.columnsByIndex) == 0 {
		return constant.ErrNoColumns
	}
	cr := csv.NewReader(r)
	var rowIndex int
	for {
		row, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				d.loaded = true
				return nil
			}
			return err
		}
		rowIndex++
		if rowIndex == 1 && skipHeader {
			continue
		}
		d.addRow(row)
	}
}

func (d *Data) addRow(row []string) {
	index := make(map[int]*Cell)
	name := make(map[string]*Cell)
	for idx, col := range d.columnsByIndex {
		str := row[idx]
		var cell Cell
		cell.t = col.t
		if len(str) == 0 {
			cell.empty = true
		} else {
			switch col.t {
			case columnString:
				cell.s = str
			case columnInt:
				n, _ := strconv.ParseInt(str, 10, 64)
				cell.i = int(n)
			case columnFloat:
				n, _ := strconv.ParseFloat(str, 10)
				cell.f = n
			case columnTime:
				cell.ts = col.timeParse(str)
				cell.timeFormat = col.timeFormat
			}
		}
		index[col.index] = &cell
		name[col.name] = &cell
	}
	d.cellsByIndex = append(d.cellsByIndex, index)
	d.cellsByName = append(d.cellsByName, name)
}

// CSV format data to csv
func (d *Data) CSV() (ret string) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	defer func() {
		w.Flush()
		ret = buf.String()
	}()
	index := make([]int, 0, len(d.columnsByIndex))
	for idx := range d.columnsByIndex {
		index = append(index, idx)
	}
	sort.Ints(index)
	header := make([]string, len(index))
	for idx, i := range index {
		header[idx] = d.columnsByIndex[i].name
	}
	if err := w.Write(header); err != nil {
		return
	}
	for _, row := range d.cellsByIndex {
		cols := make([]string, len(index))
		for idx, i := range index {
			cols[idx] = row[i].String()
		}
		if err := w.Write(cols); err != nil {
			return
		}
	}
	return
}

// Columns get columns of data
func (d *Data) Columns() []*Column {
	ret := make([]*Column, len(d.columnsByIndex))
	for i, col := range d.columnsByIndex {
		ret[i] = col
	}
	return ret
}

// Statistics statistics by column
func (d *Data) Statistics(c *Column) string {
	if !d.loaded {
		return ""
	}
	switch c.t {
	case columnTime:
		return d.statisticsTime(c)
	case columnInt:
		return d.statisticsInt(c)
	case columnFloat:
		return d.statisticsFloat(c)
	case columnString:
		return d.statisticsString(c)
	default:
		return ""
	}
}

func (d *Data) statisticsTime(c *Column) string {
	var min, max time.Time
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			continue
		}
		min = cell.ts
		max = cell.ts
	}
	var valid int
	var missing int
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			missing++
			continue
		}
		if cell.ts.Before(min) {
			min = cell.ts
		}
		if cell.ts.After(max) {
			max = cell.ts
		}
		valid++
	}
	return fmt.Sprintf("valid: %d\nmissing: %d\nmin: %s\nmax: %s\n",
		valid, missing, c.timeFormat(min), c.timeFormat(max))
}

func (d *Data) statisticsInt(c *Column) string {
	var valid int
	var missing int
	var total int
	values := make([]int, 0, len(d.cellsByIndex))
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			missing++
			continue
		}
		total += cell.i
		values = append(values, cell.i)
		valid++
	}
	sort.Ints(values)
	avg := float64(total) / float64(valid)
	var totalDiff float64
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			continue
		}
		n := float64(cell.i) - avg
		totalDiff += n * n
	}
	return fmt.Sprintf("valid: %d\nmissing: %d\nmean: %d\nstd dev: %d\nmin: %d; 25%%: %d; 50%%: %d; 75%%: %d; max: %d\n",
		valid, missing, int(avg), int(totalDiff)/valid,
		values[0], values[len(values)/4], values[len(values)/2], values[len(values)*3/4], values[len(values)-1])
}

func (d *Data) statisticsFloat(c *Column) string {
	var valid int
	var missing int
	var total float64
	values := make([]float64, 0, len(d.cellsByIndex))
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			missing++
			continue
		}
		total += cell.f
		values = append(values, cell.f)
		valid++
	}
	sort.Float64s(values)
	avg := total / float64(valid)
	var totalDiff float64
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			continue
		}
		n := cell.f - avg
		totalDiff += n * n
	}
	return fmt.Sprintf("valid: %d\nmissing: %d\nmean: %f\nstd dev: %f\nmin: %f; 25%%: %f; 50%%: %f; 75%%: %f; max: %f\n",
		valid, missing, avg, totalDiff/float64(valid),
		values[0], values[len(values)/4], values[len(values)/2], values[len(values)*3/4], values[len(values)-1])
}

func (d *Data) statisticsString(c *Column) string {
	var valid int
	var missing int
	values := make(map[string]int, len(d.cellsByIndex))
	for _, row := range d.cellsByIndex {
		cell := row[c.index]
		if cell.empty {
			missing++
			continue
		}
		values[cell.s]++
		valid++
	}
	var top string
	max := 0
	for k, v := range values {
		if v > max {
			top = k
			max = v
		}
	}
	return fmt.Sprintf("valid: %d\nmissing: %d\nuniq: %d\ntop: %s\n",
		valid, missing, len(values), top)
}

// Fill fill missing data
func (d *Data) Fill(c *Column, fn numberFunc) {
	cell, ok := fn(d, c)
	if !ok {
		return
	}
	for i, row := range d.cellsByIndex {
		if !row[c.index].empty {
			continue
		}
		row[c.index] = cell
		d.cellsByIndex[i] = row
	}
	for i, row := range d.cellsByName {
		if !row[c.name].empty {
			continue
		}
		row[c.name] = cell
		d.cellsByName[i] = row
	}
}

// Normalize normalize data
func (d *Data) Normalize(c *Column, max numberFunc) {
	cell, ok := max(d, c)
	if !ok {
		return
	}
	for i, row := range d.cellsByIndex {
		if row[c.index].empty {
			continue
		}
		row[c.index].div(cell)
		d.cellsByIndex[i] = row
	}
	c.t = columnFloat
}

// NormalizeString normalize string data
func (d *Data) NormalizeString(c *Column, hash hashFunc) {
	if c.t != columnString {
		return
	}
	for i, row := range d.cellsByIndex {
		if row[c.index].empty {
			continue
		}
		row[c.index].i = hash(row[c.index])
		row[c.index].t = columnInt
		d.cellsByIndex[i] = row
	}
	c.t = columnInt
}

// NormalizeStringEncode normalize string by encode
func (d *Data) NormalizeStringEncode(c *Column) {
	if c.t != columnString {
		return
	}
	encode := make(map[string]int)
	for _, row := range d.cellsByIndex {
		if row[c.index].empty {
			continue
		}
		encode[row[c.index].s] = 0
	}
	var i int
	for k := range encode {
		encode[k] = i
		i++
	}
	for i, row := range d.cellsByIndex {
		if row[c.index].empty {
			continue
		}
		row[c.index].t = columnInt
		row[c.index].i = encode[row[c.index].s]
		d.cellsByIndex[i] = row
	}
	c.t = columnInt
}

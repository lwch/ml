package data

import (
	"bytes"
	"encoding/csv"
	"io"
	"ml/constant"
	"sort"
	"strconv"
)

// Data data
type Data struct {
	columnsByIndex map[int]*Column
	columnsByName  map[string]*Column

	cellsByIndex []map[int]*Cell
	cellsByName  []map[string]*Cell
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

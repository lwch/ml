package ml

import (
	"fmt"
	"math/rand"
	"ml/data"
	"ml/model"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHouseInLondon(t *testing.T) {
	d := data.NewData()
	// def columns
	d.AddColumn(data.NewTimeColumn("date", 0, func(str string) time.Time {
		t, _ := time.Parse("2006-01-02", str)
		return t
	}, func(t time.Time) string {
		return t.Format("2006-01-02")
	}))
	d.AddColumn(data.NewStringColumn("area", 1))
	d.AddColumn(data.NewIntColumn("average_price", 2))
	d.AddColumn(data.NewStringColumn("code", 3))
	d.AddColumn(data.NewIntColumn("houses_sold", 4))
	d.AddColumn(data.NewFloatColumn("no_of_crimes", 5))
	d.AddColumn(data.NewIntColumn("borough_flag", 6))
	// load data
	f, err := os.Open("test_data/house_in_london.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = d.LoadFromCSV(f, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, col := range d.Columns() {
		fmt.Printf("column %s:\n", col.GetName())
		fmt.Println(d.Statistics(col))
	}
	fmt.Println("=============== fill data ===================")
	d.Fill(d.GetColumnByName("houses_sold"), data.Mean)
	d.Fill(d.GetColumnByName("no_of_crimes"), data.Mean)
	for _, col := range d.Columns() {
		fmt.Printf("column %s:\n", col.GetName())
		fmt.Println(d.Statistics(col))
	}
	fmt.Println("=============== normalize ===================")
	d.Normalize(d.GetColumnByName("average_price"), data.Max)
	d.Normalize(d.GetColumnByName("houses_sold"), data.Max)
	d.Normalize(d.GetColumnByName("no_of_crimes"), data.Max)
	d.Normalize(d.GetColumnByName("borough_flag"), data.Max)
	// d.NormalizeString(d.GetColumnByName("area"), data.Length)
	// d.NormalizeString(d.GetColumnByName("code"), data.Length)
	d.NormalizeStringEncode(d.GetColumnByName("area"))
	d.NormalizeStringEncode(d.GetColumnByName("code"))
	d.Normalize(d.GetColumnByName("area"), data.Max)
	d.Normalize(d.GetColumnByName("code"), data.Max)
	d.AddX0()
	for _, col := range d.Columns() {
		fmt.Printf("column %s:\n", col.GetName())
		fmt.Println(d.Statistics(col))
	}
	fmt.Println("=============== train ===================")
	var lr model.LinearRegression
	features := []int{
		d.GetColumnByName("x0").GetIndex(),
		d.GetColumnByName("area").GetIndex(),
		d.GetColumnByName("average_price").GetIndex(),
		d.GetColumnByName("code").GetIndex(),
		d.GetColumnByName("houses_sold").GetIndex(),
		d.GetColumnByName("no_of_crimes").GetIndex(),
	}
	label := d.GetColumnByName("borough_flag").GetIndex()
	lr.Begin(features)
	show := func() {
		params := lr.Params()
		arr := make([]string, len(params))
		for i, feature := range features {
			arr[i] = d.GetColumnByIndex(feature).GetName() + "=" + fmt.Sprintf("%f", params[i])
		}
		fmt.Println(strings.Join(arr, "; "))
	}
	for i := 0; i < 5000; i++ {
		lr.Train(0.1, d, features, label)
		fmt.Printf("%d: loss=%f\n", i, lr.Loss(d, features, label))
		if i%50 == 0 {
			show()
		}
	}
	show()
	fmt.Println("=============== predict ===================")
	for i := 0; i < 10; i++ {
		row := rand.Int() % d.Total()
		cells := d.GetColumns(row, features)
		values := make([]float64, len(cells))
		for i, cell := range cells {
			values[i] = cell.Float()
		}
		fmt.Println("predict=", lr.Predict(values), "accurate=", d.GetFloat(row, label))
	}
}

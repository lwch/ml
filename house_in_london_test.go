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
	cols := []string{"x0", "area", "average_price", "code", "houses_sold", "no_of_crimes"}
	index := make([]int, len(cols))
	for i, col := range cols {
		index[i] = d.GetColumnByName(col).GetIndex()
	}
	features := d.GetMatrix(index...)
	var lr model.LinearRegression
	labels := d.GetLables(d.GetColumnByName("borough_flag"))
	lr.Begin(len(features[0]))
	show := func() {
		params := lr.Params()
		arr := make([]string, len(params))
		for i, col := range cols {
			arr[i] = col + "=" + fmt.Sprintf("%f", params[i])
		}
		fmt.Println(strings.Join(arr, "; "))
	}
	for i := 0; i < 100000; i++ {
		lr.Train(0.01, features, labels)
		fmt.Printf("%d: loss=%f\n", i, lr.Loss(features, labels))
		if i%50 == 0 {
			show()
		}
	}
	show()
	fmt.Println("=============== predict ===================")
	for i := 0; i < 10; i++ {
		row := rand.Int() % len(features)
		fmt.Println("predict=", lr.Predict(features[row]), "accurate=", labels[row])
	}
}

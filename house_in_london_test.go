package ml

import (
	"fmt"
	"ml/data"
	"os"
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
	d.AddColumn(data.NewIntColumn("no_of_crimes", 5))
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
	fmt.Println(d.CSV())
}

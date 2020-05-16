package ml

import (
	"fmt"
	"math/rand"
	"ml/data"
	"ml/model"
	"os"
	"testing"
	"time"

	"github.com/olekukonko/tablewriter"
)

const batchCount = 100000
const learnRate = 0.1
const epoch = 500
const predictCount = 100

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
	// d.NormalizeStringEncode(d.GetColumnByName("area"))
	// d.NormalizeStringEncode(d.GetColumnByName("code"))
	// d.Normalize(d.GetColumnByName("area"), data.Max)
	// d.Normalize(d.GetColumnByName("code"), data.Max)
	d.NormalizeStringOneHot(d.GetColumnByName("area"))
	d.NormalizeStringOneHot(d.GetColumnByName("code"))
	d.AddX0()
	for _, col := range d.Columns() {
		fmt.Printf("column %s:\n", col.GetName())
		fmt.Println(d.Statistics(col))
	}
	fmt.Println("=============== train ===================")
	cols := []string{"x0", "average_price", "houses_sold", "no_of_crimes"}
	cols = append(cols, d.GetOneHotColumnNames("area")...)
	cols = append(cols, d.GetOneHotColumnNames("code")...)
	index := make([]int, len(cols))
	for i, col := range cols {
		index[i] = d.GetColumnByName(col).GetIndex()
	}
	features := d.GetMatrix(index...)
	labels := d.GetLables(d.GetColumnByName("borough_flag"))
	linearRegression(features, labels, cols)
	// logisticRegression(features, labels, cols)
}

func linearRegression(features [][]float64, labels []float64, cols []string) {
	var lr model.LinearRegression
	lr.Begin(len(features[0]))
	samples, test := selectSampleAndTest(features, labels, 0.6)
	show := func(i int) {
		w := tablewriter.NewWriter(os.Stdout)
		w.Append([]string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%f", lr.Loss(test.features, test.labels)),
		})
		params := lr.Params()
		for i, col := range cols {
			w.Append([]string{
				col,
				fmt.Sprintf("%f", params[i]),
			})
		}
		w.Render()
	}
	var offset int
	const count = 100
	for i := 0; i < batchCount; i++ {
		batch := selectBatch(samples, offset, count)
		lr.Train(learnRate, batch.features, batch.labels)
		if i%epoch == 0 {
			show(i)
		}
		offset += count
	}
	fmt.Println("=============== predict ===================")
	for i := 0; i < predictCount; i++ {
		row := rand.Int() % len(features)
		fmt.Println("predict=", lr.Predict(features[row]), "accurate=", labels[row])
	}
}

func logisticRegression(features [][]float64, labels []float64, cols []string) {
	var lr model.LogisticRegression
	lr.Begin(len(features[0]))
	samples, test := selectSampleAndTest(features, labels, 0.6)
	show := func(i int) {
		w := tablewriter.NewWriter(os.Stdout)
		w.Append([]string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%f", lr.Loss(test.features, test.labels)),
		})
		params := lr.Params()
		for i, col := range cols {
			w.Append([]string{
				col,
				fmt.Sprintf("%f", params[i]),
			})
		}
		w.Render()
	}
	var offset int
	const count = 100
	for i := 0; i < batchCount; i++ {
		batch := selectBatch(samples, offset, count)
		lr.Train(learnRate, batch.features, batch.labels)
		if i%epoch == 0 {
			show(i)
		}
		offset += count
	}
	fmt.Println("=============== predict ===================")
	for i := 0; i < predictCount; i++ {
		row := rand.Int() % len(features)
		fmt.Println("predict=", lr.Predict(features[row]), "accurate=", labels[row])
	}
}

type kv struct {
	features [][]float64
	labels   []float64
}

// return samples, test
func selectSampleAndTest(features [][]float64, labels []float64, sampleRate float64) (kv, kv) {
	if sampleRate < 0 || sampleRate > 1 {
		sampleRate = .6
	}
	var samples, test kv
	samples.features = make([][]float64, 0, int(float64(len(features))*sampleRate))
	samples.labels = make([]float64, 0, len(samples.features))
	test.features = make([][]float64, 0, len(features)-len(samples.features))
	test.labels = make([]float64, 0, len(features)-len(samples.features))
	for i, row := range features {
		if rand.Float64() <= sampleRate {
			samples.features = append(samples.features, row)
			samples.labels = append(samples.labels, labels[i])
		} else {
			test.features = append(test.features, row)
			test.labels = append(test.labels, labels[i])
		}
	}
	return samples, test
}

func selectBatch(matrix kv, offset, count int) kv {
	var ret kv
	ret.features = make([][]float64, count)
	ret.labels = make([]float64, count)
	j := 0
	n := len(matrix.features)
	for i := 0; i < count; i++ {
		ret.features[j] = matrix.features[(offset+i)%n]
		ret.labels[j] = matrix.labels[(offset+i)%n]
		j++
	}
	return ret
}

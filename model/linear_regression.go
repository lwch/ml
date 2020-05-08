package model

import (
	"ml/data"
	"sort"
)

// LinearRegression linear regression
type LinearRegression struct {
	theta []float64
}

func (lr *LinearRegression) hop(row []*data.Cell) float64 {
	var ret float64
	for i, cell := range row {
		ret += lr.theta[i] * cell.Float()
	}
	return ret
}

// Loss loss func
func (lr *LinearRegression) Loss(d *data.Data, features []int, label int) float64 {
	sort.Ints(features)
	var total float64
	for i := 0; i < d.Total(); i++ {
		n := lr.hop(d.GetColumns(i, features)) - d.GetFloat(i, label)
		total += n * n
	}
	return total / (2. * float64(d.Total()))
}

// Begin begin train
func (lr *LinearRegression) Begin(features []int) {
	lr.theta = make([]float64, len(features))
}

// Train train data one times
func (lr *LinearRegression) Train(rate float64, d *data.Data, features []int, label int) {
	sort.Ints(features)
	for j := 0; j < len(features); j++ {
		var total float64
		for i := 0; i < d.Total(); i++ {
			loss := lr.hop(d.GetColumns(i, features)) - d.GetFloat(i, label)
			total += loss * d.GetFloat(i, features[j])
		}
		lr.theta[j] -= rate * total / float64(d.Total())
	}
}

// Predict predict score
func (lr *LinearRegression) Predict(features []float64) float64 {
	var total float64
	for i := 0; i < len(features); i++ {
		total += lr.theta[i] * features[i]
	}
	return total
}

// Params get params
func (lr *LinearRegression) Params() []float64 {
	return lr.theta
}

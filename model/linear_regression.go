package model

// LinearRegression linear regression
type LinearRegression struct {
	theta []float64
}

func (lr *LinearRegression) hop(row []float64) float64 {
	var ret float64
	for i, cell := range row {
		ret += lr.theta[i] * cell
	}
	return ret
}

// Loss loss func
func (lr *LinearRegression) Loss(features [][]float64, labels []float64) float64 {
	var total float64
	for i, row := range features {
		n := lr.hop(row) - labels[i]
		total += n * n
	}
	return total / (2. * float64(len(features)))
}

// Begin begin train
func (lr *LinearRegression) Begin(featureCount int) {
	lr.theta = make([]float64, featureCount)
}

// Train train data one times
func (lr *LinearRegression) Train(rate float64, features [][]float64, labels []float64) {
	for j := 0; j < len(features[0]); j++ {
		var total float64
		for i, row := range features {
			loss := lr.hop(row) - labels[i]
			total += loss * row[j]
		}
		lr.theta[j] -= rate * total / float64(len(features))
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

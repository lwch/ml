package model

import "math"

// LogisticRegression logistic regression
type LogisticRegression struct {
	theta []float64
}

func (lr *LogisticRegression) sigmoid(z float64) float64 {
	return 1. / (1 + math.Exp(-z))
}

// Loss loss func
func (lr *LogisticRegression) Loss(features [][]float64, labels []float64) float64 {
	var total float64
	for i, row := range features {
		n := lr.Predict(row) - labels[i]
		total += n * n
	}
	return total / (2. * float64(len(features)))
}

// Begin begin train
func (lr *LogisticRegression) Begin(featureCount int) {
	lr.theta = make([]float64, featureCount)
}

// Train train data one times
func (lr *LogisticRegression) Train(rate float64, features [][]float64, labels []float64) {
	for j := 0; j < len(features[0]); j++ {
		var total float64
		for i, row := range features {
			loss := lr.Predict(row) - labels[i]
			total += loss * row[j]
		}
		lr.theta[j] -= rate * total / float64(len(features))
	}
}

// Predict predict score
func (lr *LogisticRegression) Predict(features []float64) float64 {
	var total float64
	for i := 0; i < len(features); i++ {
		total += lr.theta[i] * features[i]
	}
	return lr.sigmoid(total)
}

// Params get params
func (lr *LogisticRegression) Params() []float64 {
	return lr.theta
}

package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

const sampleCount = 10000
const lr = 0.01
const batch = 100000
const epoch = 500

type model struct {
	a float64
	b float64
}

func (m *model) rand(x float64) float64 {
	switch rand.Int() % 3 {
	case 1:
		return m.a*x + m.b + rand.Float64()*0.5
	case 2:
		return m.a*x + m.b - rand.Float64()*0.5
	default:
		return m.a*x + m.b
	}
}

func (m *model) loss(a, b float64, x, y []float64) float64 {
	var total float64
	for i := range x {
		d := a*x[i] + b - y[i]
		d *= d
		total += d
	}
	return total / float64(len(x)) / .5
}

func (m *model) optimize(a, b float64, x, y []float64) (float64, float64) {
	var da, db float64
	for i := range x {
		dy := a*x[i] + b - y[i]
		da += dy * x[i]
		db += dy
	}
	da /= float64(len(x))
	db /= float64(len(x))
	a -= lr * da
	b -= lr * db
	return a, b
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func train(x, y []float64, m model, count int) ([]float64, []float64) {
	var a, b float64
	ar := make([]float64, 0, count/epoch)
	br := make([]float64, 0, count/epoch)
	for i := 0; i < count; i++ {
		a, b = m.optimize(a, b, x, y)
		if i > 0 && i%epoch == 0 {
			testX := []float64{rand.Float64()}
			testY := []float64{m.a*testX[0] + m.b}
			fmt.Printf("round=%d, a=%f, b=%f, loss=%f\n", i, a, b, m.loss(a, b, testX, testY))
			ar = append(ar, a)
			br = append(br, b)
		}
	}
	return ar, br
}

func main() {
	m := model{5., .3}
	x := make([]float64, sampleCount)
	y := make([]float64, sampleCount)
	for i := 0; i < sampleCount; i++ {
		x[i] = rand.Float64()
		y[i] = m.rand(x[i])
	}
	a, b := train(x, y, m, batch)
	fmt.Println(a[len(a)-1], b[len(b)-1])
	save(x, y, a, b)
}

func save(x, y, a, b []float64) {
	p, err := plot.New()
	assert(err)

	fs := make([]plot.Plotter, len(a))
	for i := range a {
		func(i int) {
			fs[i] = plotter.NewFunction(func(x float64) float64 {
				return a[i]*x + b[i]
			})
			fs[i].(*plotter.Function).Color = color.RGBA{B: 255, A: 255}
		}(i)
	}
	p.Add(fs...)

	points := make(plotter.XYs, len(x))
	for i := range x {
		points[i] = plotter.XY{X: x[i], Y: y[i]}
	}
	ps, err := plotter.NewScatter(points)
	assert(err)
	ps.Radius = 1
	ps.Color = color.RGBA{R: 255, A: 255}
	p.Add(ps)

	result := plotter.NewFunction(func(x float64) float64 {
		return a[len(a)-1]*x + b[len(b)-1]
	})
	result.Color = color.Gray{128}
	result.Width = 3
	p.Add(result)

	assert(p.Save(5*vg.Inch, 5*vg.Inch, "lr.jpg"))
}

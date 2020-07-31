package draw

import (
	"bytes"
	"image/color"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"github.com/rwlist/autotrade-bot/pkg/formula"

	"gonum.org/v1/plot"

	"github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot/plotter"
)

type Plot struct {
	Plot *plot.Plot
}

func NewPlot() *Plot {
	p, _ := plot.New()
	return &Plot{
		Plot: p,
	}
}

func (p *Plot) AddEnv() {
	p.Plot.X.Tick.Marker = AllTimeTicks{}
	p.Plot.Y.Tick.Marker = AllPriceTicks{}
	p.Plot.Add(plotter.NewGrid())
}

func (p *Plot) AddHelpLines(lastPrice, minPrice, maxPrice float64, startTime int64) {
	p.Plot.Add(MakeHorLine(float64(startTime), lastPrice, 0, 0, 255))
	p.Plot.Add(MakeHorLine(float64(startTime), minPrice, 255, 0, 0))
	p.Plot.Add(MakeHorLine(float64(startTime), maxPrice, 0, 255, 0))
}

func (p *Plot) AddFunction(f formula.Formula, yMax, xMax float64) {
	if yMax == -1 {
		yMax = convert.Float64(f.Calc(f.Start().Add(secDay)))
	}
	p.Plot.X.Max = xMax
	p.Plot.Y.Max = yMax
	lambda := func(x float64) float64 {
		return convert.Float64(f.Calc(time.Unix(int64(x), 0)))
	}
	fu := plotter.NewFunction(lambda)
	fu.XMin = float64(f.Start().Unix())
	fu.XMax = xMax
	fu.Width = 2
	fu.Color = color.RGBA{R: 255, B: 0, G: 165, A: 255}
	p.Plot.Add(fu)
}

func (p *Plot) AddRateGraph(klines *Klines) {
	bars, _ := custplotter.NewCandlesticks(klines)
	bars.ColorUp = color.RGBA{
		R: 2,
		G: 192,
		B: 118,
		A: 255,
	}
	bars.ColorDown = color.RGBA{
		R: 217,
		G: 48,
		B: 78,
		A: 255,
	}
	bars.FixedLineColor = false
	bars.Width = 1
	p.Plot.Add(bars)
}

const DefaultWidth = 900
const DefaultHeight = 400

func (p *Plot) SaveToBuffer() (*bytes.Buffer, error) {
	w, err := p.Plot.WriterTo(DefaultWidth, DefaultHeight, "png")
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer([]byte{})
	_, err = w.WriteTo(buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

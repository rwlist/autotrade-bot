package draw

import (
	"bytes"
	"image/color"

	"gonum.org/v1/plot"

	"github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot/plotter"
)

type Plot struct {
	Plot *plot.Plot
}

func NewPlot() (Plot, error) {
	p, err := plot.New()
	return Plot{p}, err
}

func (p Plot) DrawEnv() {
	p.Plot.X.Tick.Marker = AllTimeTicks{}
	p.Plot.Y.Tick.Marker = AllPriceTicks{}
	p.Plot.Add(plotter.NewGrid())
}

func (p Plot) DrawHelpLines(lastPrice, minPrice, maxPrice, startTime float64) {
	p.Plot.Add(MakeHorLine(startTime, lastPrice, 0, 0, 255))
	p.Plot.Add(MakeHorLine(startTime, minPrice, 255, 0, 0))
	p.Plot.Add(MakeHorLine(startTime, maxPrice, 0, 255, 0))
}

func (p Plot) DrawMainGraph(klines Klines) error {
	bars, err := custplotter.NewCandlesticks(klines)
	if err != nil {
		return err
	}
	bars.ColorUp = color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	}
	bars.ColorDown = color.RGBA{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	}
	p.Plot.Add(bars)
	return nil
}

const DefaultWidth = 900
const DefaultHeight = 400

func (p Plot) SaveToBuffer() (*bytes.Buffer, error) {
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

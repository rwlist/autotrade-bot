package draw

import (
	"bytes"
	"image/color"

	"github.com/rwlist/autotrade-bot/binance"
	"gonum.org/v1/plot"

	"github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot/plotter"
)

type Plot struct {
	plot *plot.Plot
}

func NewPlot() (Plot, error) {
	p, err := plot.New()
	return Plot{p}, err
}

func (p Plot) DrawEnv() {
	p.plot.X.Tick.Marker = AllTimeTicks{}
	p.plot.Y.Tick.Marker = AllPriceTicks{}
	p.plot.Add(plotter.NewGrid())
}

func (p Plot) DrawHelpLines(lastPrice, minPrice, maxPrice, startTime float64) {
	p.plot.Add(MakeHorLine(startTime, lastPrice, 0, 0, 255))
	p.plot.Add(MakeHorLine(startTime, minPrice, 255, 0, 0))
	p.plot.Add(MakeHorLine(startTime, maxPrice, 0, 255, 0))
}

func (p Plot) DrawMainGraph(klines binance.TOHLCVs) error {
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
	p.plot.Add(bars)
	return nil
}

func (p Plot) SaveToBuffer() (*bytes.Buffer, error) {
	w, err := p.plot.WriterTo(900, 400, "png")
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

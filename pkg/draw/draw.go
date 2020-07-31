package draw

import (
	"image/color"
	"time"

	"github.com/rwlist/autotrade-bot/pkg/convert"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type AllTimeTicks struct{}

func (AllTimeTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		ut := time.Unix(int64(t.Value), 0)
		tks[i].Label = ut.Format("02.01\n15:04")
	}
	return tks
}

type AllPriceTicks struct{}

func (AllPriceTicks) Ticks(min, max float64) []plot.Tick {
	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" {
			tks[i].Label = convert.Str64(int64(t.Value))
		}
	}
	return tks
}

const secDay = 86400

func MakeHorLine(x, y float64, r, g, b uint8) *plotter.Line {
	pts := make(plotter.XYs, 2)
	pts[0].Y = y
	pts[1].Y = y
	pts[0].X = x - 100
	pts[1].X = x + 2*secDay + 100
	line, _ := plotter.NewLine(pts)
	line.Color = color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	}
	return line
}

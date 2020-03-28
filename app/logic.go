package app

import (
	"bytes"
	"fmt"
	"github.com/pplcc/plotext/custplotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"image/color"
	"time"

	"github.com/adshao/go-binance"
)

type Logic struct {
	b *MyBinance
}

func NewLogic(b *MyBinance) *Logic {
	return &Logic{
		b: b,
	}
}

type Balance struct {
	usd    string
	asset  string
	free   string
	locked string
}

type Status struct {
	total    string
	rate     string
	balances []*Balance
}

func (l *Logic) CommandStatus() (*Status, error) {
	rate, err := l.b.GetRate()
	if err != nil {
		return nil, err
	}
	allBalances, err := l.b.AccountBalance()
	if err != nil {
		return nil, err
	}

	var balances []*Balance
	var total float64
	for _, bal := range allBalances {
		if isEmptyBalance(bal.Free) && isEmptyBalance(bal.Locked) {
			continue
		}

		balUSD, err := l.b.BalanceToUSD(&bal)
		if err != nil {
			return &Status{}, err
		}
		total += balUSD
		resBal := &Balance{
			usd:    float64ToStr(balUSD, 2),
			asset:  bal.Asset,
			free:   bal.Free,
			locked: bal.Locked,
		}
		balances = append(balances, resBal)
	}

	res := &Status{
		total:    float64ToStr(total, 2),
		rate:     rate,
		balances: balances,
	}
	return res, err
}

const sleepDur = time.Duration(1) * time.Second

func (l *Logic) CommandBuy(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.BuyAll()
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeBuy)))
			return
		}
		s.Send(startMessage(orderNew))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID())
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeBuy)))
			return
		}
		s.Send(orderStatusMessage(order))
		err = l.b.CancelOrder(order.OrderID())
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeBuy)))
			return
		}
	}
}

func (l *Logic) CommandSell(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.SellAll()
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeSell)))
			return
		}
		s.Send(startMessage(orderNew))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID())
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeSell)))
			return
		}
		s.Send(orderStatusMessage(order))
		err = l.b.CancelOrder(order.OrderID())
		if err != nil {
			s.Send(errorMessage(err, string(binance.SideTypeSell)))
			return
		}
	}
}

func(l* Logic) CommandDraw(s *Sender) {
	klines, lastPrice, minPrice, maxPrice, startTime, err := l.b.GetKlines()
	if err != nil {
		s.Send(errorMessage(err, "Draw GetKlines"))
		return
	}
	p, err := plot.New()
	if err != nil {
		s.Send(errorMessage(err, "Draw in plot.New()"))
		return
	}

	p.Title.Text = "Candlesticks"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Price"
	p.X.Tick.Marker = AllTimeTicks{}
	p.Y.Tick.Marker = AllPriceTicks{}
	p.Add(plotter.NewGrid())
	p.Add(MakeInfHorLine(startTime, lastPrice, 0, 0, 255))
	p.Add(MakeInfHorLine(startTime, minPrice, 255, 0, 0))
	p.Add(MakeInfHorLine(startTime, maxPrice, 0, 255, 0))

	bars, err := custplotter.NewCandlesticks(klines)
	if err != nil {
		s.Send(errorMessage(err, "Draw in custplotter.NewCandlesticks(klines)"))
		return
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
	p.Add(bars)
	w, err := p.WriterTo(900, 400, "png")
	if err != nil {
		s.Send(errorMessage(err, "Draw in p.WriterTo(900, 400, \"png\")"))
		return
	}
	buffer := bytes.NewBuffer([]byte{})
	_, err = w.WriteTo(buffer)
	if err != nil {
		s.Send(errorMessage(err, "Draw in w.WriteTo(buffer)"))
		return
	}
	err = s.SendPhoto("graph.png", buffer.Bytes())
	if err != nil {
		s.Send(errorMessage(err, "Draw in SendPhoto"))
		return
	}
}

//--------------------------------------TEMPLATES FOR SENDER----------------------------------------------
func errorMessage(err error, command string) string {
	return fmt.Sprintf("Error while %v:\n\n%s", command, err)
}

func startMessage(order Order) string {
	return fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for %s", order.Side(), order.Price(), sleepDur)
}

func orderStatusMessage(order Order) string {
	return fmt.Sprintf("Side: %v\nDone %v / %v\nStatus: %v", order.Side(), order.ExecutedQuantity(), order.OrigQuantity(), order.Status())
}

//-------------------------------------------------------------------------------

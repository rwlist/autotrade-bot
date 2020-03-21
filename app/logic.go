package app

import (
	"fmt"
	"github.com/adshao/go-binance"
	"time"
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

const sleepSec = 1
const sleepDur = time.Duration(sleepSec) * time.Second

func (l *Logic) CommandBuy(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.BuyAll()
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeBuy))
			return
		}
		s.Send(startMessage(&OrderNew{orderNew}))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID)
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeBuy))
			return
		}
		s.Send(orderStatusMessage(&OrderExist{order}))
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeBuy))
			return
		}
	}
}

func (l *Logic) CommandSell(s *Sender) {
	for i := 0; i < 5; i++ {
		orderNew, err := l.b.SellAll()
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeSell))
			return
		}
		s.Send(startMessage(&OrderNew{orderNew}))
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderNew.OrderID)
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeSell))
			return
		}
		s.Send(orderStatusMessage(&OrderExist{order}))
		err = l.b.CancelOrder(order.OrderID)
		if err != nil {
			s.Send(errorMessage(err, binance.SideTypeSell))
			return
		}
	}
}

//--------------------------------------TEMPLATES FOR SENDER----------------------------------------------
func errorMessage(err error, side binance.SideType) string {
	return fmt.Sprintf("Error while %v:\n\n%s", side, err)
}

func startMessage(order Order) string {
	return fmt.Sprintf("A %v BTC/USDT order was placed with price = %v.\nWaiting for %v seconds..", order.Side(), order.Price(), sleepSec)
}

func orderStatusMessage(order Order) string {
	return fmt.Sprintf("Side: %v\nDone %v / %v\nStatus: %v", order.Side(), order.ExecutedQuantity(), order.OrigQuantity(), order.Status())
}
//-------------------------------------------------------------------------------

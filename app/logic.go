package app

import (
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
	total	 string
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
			   usd:    float64ToStr(balUSD),
			   asset:  bal.Asset,
			   free:   bal.Free,
			   locked: bal.Locked,
		}
		balances = append(balances, resBal)
	}

	res := &Status{
		total:	  float64ToStr(total),
		rate:     rate,
		balances: balances,
	}
	return res, err
}

type OrderInfo struct {
	InfoType int
	Status binance.OrderStatusType
	Side binance.SideType
	Price string
	OrigQuantity string
	ExecutedQuantity string
	Err error
}

const sleepDur = time.Duration(2) * time.Second

func toChan1(ch chan *OrderInfo, order *binance.CreateOrderResponse, err error, InfoType int) {
	orderInfo := &OrderInfo{
		InfoType: InfoType,
		Status: order.Status,
		Side: "",
		Price: "",
		OrigQuantity: "",
		ExecutedQuantity: "",
		Err: err,
	}
	if  err == nil {
		orderInfo.Side = order.Side
		orderInfo.Price = order.Price
	}
	ch <- orderInfo
}

func toChan2(ch chan *OrderInfo, order *binance.Order, err error, InfoType int) {
	orderInfo := &OrderInfo{
		InfoType: InfoType,
		Status: order.Status,
		Side: "",
		Price: "",
		OrigQuantity: "",
		ExecutedQuantity: "",
		Err: err,
	}
	if  err == nil {
		orderInfo.OrigQuantity = order.OrigQuantity
		orderInfo.ExecutedQuantity = order.ExecutedQuantity
	}
	ch <- orderInfo
}

func (l *Logic) CommandBuy(ch chan *OrderInfo) {
	for i := 0; i < 5; i++ {
		orderCreate, err := l.b.BuyAll()
		toChan1(ch, orderCreate, err, 1)
		if err != nil {
			close(ch)
			return
		}
		time.Sleep(sleepDur)
		order, err := l.b.GetOrder(orderCreate.OrderID)
		toChan2(ch, order, err, 2)
		if err != nil {
			close(ch)
			return
		}
		if order.Status == binance.OrderStatusTypeFilled {
			close(ch)
			return
		}
		err = l.b.CancelOrder(order.OrderID)
		toChan2(ch, order, err, 0)
		if err != nil {
			close(ch)
			return
		}
	}
	close(ch)
}

//TEST COMMANDS
func (l *Logic) CommandTestOrderAll() error {
	return l.b.TestBuyAll()
}

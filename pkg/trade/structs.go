package trade

type Order struct {
	Symbol           string
	OrderID          int64
	Price            string
	OrigQuantity     string
	ExecutedQuantity string
	Status           string
	Side             string
}

type Status struct {
	Order *Order
	Done  bool
	Err   error
}

type Balance struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

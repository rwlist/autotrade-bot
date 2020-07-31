package formula

import (
	"testing"

	"github.com/shopspring/decimal"
)

func ok(a, b []decimal.Decimal) bool {
	for i := range a {
		if !a[i].Equal(b[i]) {
			return false
		}
	}
	return true
}

func TestParseBasic(t *testing.T) {
	type test struct {
		in  string
		out []decimal.Decimal
	}
	tests := []test{
		{
			in: "/draw rate-1+1*(now-start)^1",
			out: []decimal.Decimal{
				decimal.NewFromFloat(-1),
				decimal.NewFromFloat(1),
				decimal.NewFromFloat(1),
			},
		},
		{
			in: "/draw rate+2.2+0.0002*(now-start)^2.2",
			out: []decimal.Decimal{
				decimal.NewFromFloat(2.2),
				decimal.NewFromFloat(0.0002),
				decimal.NewFromFloat(2.2),
			},
		},
		{
			in: "/draw rate-3-0.0003*(now-start)^1.3",
			out: []decimal.Decimal{
				decimal.NewFromFloat(-3),
				decimal.NewFromFloat(-0.0003),
				decimal.NewFromFloat(1.3),
			},
		},
		{
			in: "/draw rate+4-0.0004*(now-start)^1.4",
			out: []decimal.Decimal{
				decimal.NewFromFloat(4),
				decimal.NewFromFloat(-0.0004),
				decimal.NewFromFloat(1.4),
			},
		},
	}

	for _, val := range tests {
		res, err := parseBasic(val.in)
		if err != nil {
			t.Error(err)
		}
		if !ok(res, val.out) {
			t.Errorf("%v incorrect result for %v", res, val.in)
		}
	}
}

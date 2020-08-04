package history

type History struct {
	formulas []string
}

func New() *History {
	return &History{}
}

const maxLen = 10

func (h *History) AddFormula(s string) {
	h.formulas = append(h.formulas, s)
	if len(h.formulas) > maxLen {
		h.formulas = h.formulas[1:]
	}
}

func (h *History) GetFormulasList() []string {
	return h.formulas
}

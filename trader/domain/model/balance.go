package model

type Balance struct {
	currencyCode string
	amount       float64
	available    float64
}

func NewBalance(currencyCode string, amount, available float64) *Balance {
	if currencyCode == "" {
		return nil
	}

	if amount < 0 {
		return nil
	}

	if available < 0 {
		return nil
	}

	return &Balance{
		currencyCode: currencyCode,
		amount:       amount,
		available:    available,
	}
}

func (b *Balance) CurrencyCode() string {
	return b.currencyCode
}

func (b *Balance) Amount() float64 {
	return b.amount
}

func (b *Balance) Available() float64 {
	return b.available
}

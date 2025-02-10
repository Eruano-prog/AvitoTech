package service

type CoinService struct {
}

func (c CoinService) SendCoin(fromUser int, toUser string, amount int) error {
	panic("implement me")
}

func (c CoinService) BuyItem(id int, item string) error {
	panic("implement me")
}

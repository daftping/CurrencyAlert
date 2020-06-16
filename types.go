package main

const (
	apiCurrencyMono        = "https://api.monobank.ua/bank/currency"
	apiCurrencyP24Bussines = "https://otp24.privatbank.ua/v3/api/1/info/currency/get"
	apiInterBank           = "https://charts.finance.ua/ua/currency/data-daily?for=interbank&source=1&indicator=usd"
	// ISO 4217
	USD = 840
	UAH = 980
)

//Mono
type MonoPair struct {
	CurrencyCodeA int
	CurrencyCodeB int
	Date          int64
	RateBuy       float32
	RateSell      float32
}

type MonoCurrency struct {
	Pairs []MonoPair
}

//P24Bussines
type P24BussinesCurrency struct {
	USD P24BussinesBuyAndSell
	EUR P24BussinesBuyAndSell
}

type P24BussinesBuyAndSell struct {
	B P24BussinesRate
	S P24BussinesRate
}

type P24BussinesRate struct {
	Date      string
	Rate      string
	RateDelta string `json:"rate_delta"`
	NbuRate   string
}

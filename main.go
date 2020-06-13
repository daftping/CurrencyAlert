package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	apiCurrencyMono        = "https://api.monobank.ua/bank/currency"
	apiCurrencyP24Bussines = "https://otp24.privatbank.ua/v3/api/1/info/currency/get"
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

func (c MonoCurrency) getPair(A int, B int) (MonoPair, error) {
	for _, pair := range c.Pairs {
		if pair.CurrencyCodeA == A && pair.CurrencyCodeB == B {
			return pair, nil
		}
	}
	return MonoPair{}, fmt.Errorf("Currency pair has not been found")
}

func getHTTP(api string) []byte {
	res, errHTTP := http.Get(api)
	if errHTTP != nil {
		log.Fatal(errHTTP)
	}
	respBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return respBody
}

func getMono() float32 {
	respBody := getHTTP(apiCurrencyMono)

	pairs := []MonoPair{}
	errJSON := json.Unmarshal(respBody, &pairs)
	if errJSON != nil {
		fmt.Println("error:", errJSON)
	}

	exchangeRate := MonoCurrency{Pairs: pairs}

	pair, errPair := exchangeRate.getPair(USD, UAH)
	if errPair != nil {
		log.Fatal(errPair)
	}

	return pair.RateSell
}

func getP24Bussines() float32 {
	respBody := getHTTP(apiCurrencyP24Bussines)

	currency := P24BussinesCurrency{}
	errJSON := json.Unmarshal(respBody, &currency)
	if errJSON != nil {
		fmt.Println("error:", errJSON)
	}
	rate, _ := strconv.ParseFloat(currency.USD.B.Rate, 32)

	return float32(rate)
}

// Only for local development
func localRun() {
	mono := getMono()
	pb24 := getP24Bussines()
	diff := mono - pb24
	fmt.Printf("Mono: %v\n", mono)
	fmt.Printf("PB24: %v\n", pb24)
	fmt.Printf("Exchange rate difference is: %v\n", diff)
	publishToSNS(fmt.Sprint(diff))
}

func HandleRequest(ctx context.Context) (string, error) {
	diff := getMono() - getP24Bussines()
	publishToSNS(fmt.Sprint(diff))
	return fmt.Sprintf("%v", diff), nil
}

func main() {
	// For local development
	if runtime.GOOS == "darwin" {
		localRun()
	} else {
		lambda.Start(HandleRequest)
	}

}

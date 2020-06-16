package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hashicorp/go-retryablehttp"
)

func (c MonoCurrency) getPair(A int, B int) (MonoPair, error) {
	for _, pair := range c.Pairs {
		if pair.CurrencyCodeA == A && pair.CurrencyCodeB == B {
			return pair, nil
		}
	}
	return MonoPair{}, fmt.Errorf("Currency pair has not been found")
}

func getHTTP(api string) []byte {
	res, errHTTP := retryablehttp.Get(api)
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

func getInterBank() float32 {
	respBody := getHTTP(apiInterBank)

	var interBank [][]string
	errJSON := json.Unmarshal(respBody, &interBank)
	if errJSON != nil {
		fmt.Println("error:", errJSON)
	}
	rateNow, _ := strconv.ParseFloat(interBank[len(interBank)-1][1], 32)

	return float32(rateNow)
}

// Only for local development
func localRun() {
	output := formatOutput()
	fmt.Println(output)
	publishToSNS(output)
}

func formatOutput() string {
	mono := getMono()
	pb24 := getP24Bussines()
	interBank := getInterBank()
	return fmt.Sprintf("МоноБанк:\t%.2f\nПриват24:\t%.2f\nРізниця:\t%.2f\nМіжбанк:\t%.2f\n", mono, pb24, mono-pb24, interBank)
}

func HandleRequest(ctx context.Context) (string, error) {
	output := formatOutput()
	publishToSNS(output)
	return output, nil
}

func main() {
	// For local development
	if runtime.GOOS == "darwin" {
		localRun()
	} else {
		lambda.Start(HandleRequest)
	}

}

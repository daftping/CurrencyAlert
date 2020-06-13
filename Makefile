lambda:
	go get github.com/aws/aws-lambda-go/lambda
	GOOS=linux go build .
	zip CurrencyAlert.zip CurrencyAlert
	rm CurrencyAlert
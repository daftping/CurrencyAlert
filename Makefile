lambda: get
	GOOS=linux go build .
	zip CurrencyAlert.zip CurrencyAlert
	rm CurrencyAlert
	aws lambda update-function-code \
    --function-name  CurrencyAlert \
    --zip-file fileb://CurrencyAlert.zip

run:
	echo $SNS_TOPIC_ARN
	@go run .	

get:
	go get -u .
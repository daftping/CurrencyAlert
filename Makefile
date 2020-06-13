lambda: get
	GOOS=linux go build .
	zip CurrencyAlert.zip CurrencyAlert
	rm CurrencyAlert

run:
	@go run .	

get:
	go get -u .
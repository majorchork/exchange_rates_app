run:
	go run main.go

mock:
	mockgen -source=service/rates.go -destination=service/mocks/rates_mock.go -package=mocks

test:
	go test ./... -cover


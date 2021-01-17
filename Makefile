.PHONY:

build:
	go build -o booking cmd/app/main.go

tests:
	psql -c "\i scripts/test_init.sql;" -U postgres && go test ./... -v

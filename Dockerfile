FROM golang:1.15 AS build

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GOGC=off go build -a -installsuffix cgo -ldflags="-w -s" -v -o ./app ./cmd/app

EXPOSE 9000

CMD ./app
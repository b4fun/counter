FROM golang:1.10.3-alpine3.7

WORKDIR /go/src/github.com/b4fun/counter
COPY . .

CMD go run ./cmd/*.go

EXPOSE 8081

FROM golang:1.10.3-alpine3.7 as builder

WORKDIR /go/src/github.com/b4fun/counter
COPY . .

RUN go build -o counter ./cmd/*.go

FROM alpine:3.7

COPY --from=builder /go/src/github.com/b4fun/counter/counter /

CMD ["/counter"]

EXPOSE 8081

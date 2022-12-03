FROM golang:1.19.0 AS builder

ADD . /app
WORKDIR /app
# GOOS/GOARCH as you build not from go alpine
RUN GOOS=linux GOARCH=amd64 go build -o go-url-shortener-app ./cmd/go-url-shortener-example

FROM alpine:3.15 AS app
WORKDIR /app
COPY --from=builder /app/go-url-shortener-app /app
CMD ["/app/go-url-shortener-app"]
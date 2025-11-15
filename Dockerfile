FROM golang:1.25.1 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY go.mod .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app

EXPOSE 8080

CMD ["./app"]
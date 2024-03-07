FROM golang:1.20.4-alpine AS builder

WORKDIR /banking

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /banking

EXPOSE 8080

CMD ["./banking", "apiserver"]

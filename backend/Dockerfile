FROM golang:alpine3.19 AS builder

WORKDIR /app 

COPY . .

RUN go mod download

RUN go build -o /app/build/shortly /app/cmd/shortly/main.go

FROM alpine:latest

COPY --from=builder /app/build/shortly /app/build/shortly

ENTRYPOINT [ "/app/build/shortly" ]

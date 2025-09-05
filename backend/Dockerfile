FROM golang:alpine3.22 AS builder

WORKDIR /app 

COPY . .

RUN go mod download

RUN go build -o /app/build/shortly /app/cmd/shortly/main.go

FROM alpine:latest

WORKDIR /app

RUN addgroup --system --gid 1000 golang && \
    adduser --system --uid 1000 shortly && \
    mkdir -p /app/keys && \
    chown -R 1000:1000 /app

COPY --from=builder /app/build/shortly /app/build/shortly

USER shortly

ENTRYPOINT [ "./build/shortly" ]

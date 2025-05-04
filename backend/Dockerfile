FROM golang:alpine3.19 AS builder

WORKDIR /app 

COPY . .

RUN go mod download

RUN go build -o /app/build/shortly /app/cmd/shortly/main.go

FROM alpine:latest

RUN addgroup --system --gid 1000 golang
RUN adduser --system --uid 1000 shortly

COPY --from=builder /app/build/shortly /app/build/shortly

RUN chown -R 1000:1000 /app

USER shortly

ENTRYPOINT [ "/app/build/shortly" ]

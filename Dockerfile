FROM golang:1.24.1 AS builder

WORKDIR /app

COPY . .

RUN GOARCH=amd64 GOOS=linux go mod tidy && go build -o /app/bin/app ./cmd

FROM alpine:latest

RUN apk --no-cache add libc6-compat

WORKDIR /root/

COPY --from=builder /app/config /root/config
COPY --from=builder /app/bin/app /root/

CMD ["./app"]

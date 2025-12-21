FROM golang:latest as builder
WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build --ldflags="-w -s" -o service-b cmd/main.go

FROM alpine:latest
COPY --from=builder /app/service-b /app/service-b
CMD ["/app/service-b"] 
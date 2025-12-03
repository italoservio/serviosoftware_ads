FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin ./bin
EXPOSE 8080
CMD ["./bin"]

# Build stage
FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o dstgbot .

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/dstgbot .



EXPOSE 8080
CMD ["./dstgbot"]
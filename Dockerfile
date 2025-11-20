# Stage 1: Build the application
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /app/main ./app/main.go

# Stage 2: Create the final, small image
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/assets ./assets

EXPOSE 5000

CMD ["./main"]

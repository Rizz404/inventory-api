# Stage 1: Build the application
FROM golang:1.25-alpine AS builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /app/main ./app/main.go

# Buat seednya
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /app/seeder ./cmd/seed/main.go

# Stage 2: Create the final, small image
FROM alpine:latest
# Install ca-certificates & tzdata (PENTING untuk API calls & Timezone Indonesia)
# Install curl untuk healthcheck
RUN apk --no-cache add ca-certificates tzdata curl
WORKDIR /root/

# Set timezone ke WIB (Jakarta)
ENV TZ=Asia/Jakarta

# 2. Copy Binary Goose dari builder ke final image
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# 3. Copy Folder Migrasi kamu ke dalam image
COPY --from=builder /app/db/migrations ./db/migrations

COPY --from=builder /app/main .
# Buat seeder
COPY --from=builder /app/seeder .
COPY --from=builder /app/assets ./assets

EXPOSE 5000

CMD ["./main"]

# ========================
# Stage 1: Build
# ========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Dependency cache ভালো রাখার জন্য আগে go.mod ও go.sum কপি করি
COPY go.mod go.sum ./
RUN go mod download

# Source code কপি করি
COPY . .

# Static binary বানাই
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# ========================
# Stage 2: Final Image (ছোট)
# ========================
FROM alpine:3.20

WORKDIR /app

# Security এর জন্য non-root user
RUN adduser -D -g '' appuser
USER appuser

# শুধু binary কপি করি
COPY --from=builder /app/main .

# যদি .env বা config file থাকে তাহলে এখানে কপি করো
# COPY --from=builder /app/.env .
# COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./main"]
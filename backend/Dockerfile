# Stage 1: Build the binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Run the binary
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Render uses the PORT env var, but we expose it for documentation
EXPOSE 8080
CMD ["./main"]
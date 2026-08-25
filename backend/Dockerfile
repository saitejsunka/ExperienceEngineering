# Build stage
FROM golang:alpine AS builder

# Install CA certificates and tzdata for SSL/TLS and timezones
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory (which includes the service folder)
COPY . .

# Build the Go application
# -ldflags="-w -s" strips debug info for a smaller binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./service

# Runtime stage (Scratch for minimal attack surface and size)
FROM scratch

# Import certs and timezone data from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the compiled binary
COPY --from=builder /app/server /server

# Expose gRPC port (Assuming standard Cloud Run default mapping or custom 8080)
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/server"]

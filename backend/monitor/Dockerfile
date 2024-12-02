# Use the official Golang image to create a build artifact
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o monitor-service

# Use an Alpine image for the final container
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates netcat-openbsd

WORKDIR /root/
COPY --from=builder /app/monitor-service .
COPY .env .env
COPY scripts/startup.sh ./startup.sh

# Make the startup script executable
RUN chmod +x ./startup.sh

# Command to run the startup script
ENTRYPOINT ["./startup.sh"]

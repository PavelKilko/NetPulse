# Use the official Golang 1.22 image to create a build artifact
FROM golang:1.22-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache build-base

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY .. .

# Build the Go app
RUN go build -o netpulse-backend

# Use an Alpine image for the final container
FROM alpine:latest

# Install necessary packages for Alpine (ca-certificates, nc, curl, postgres client)
RUN apk --no-cache add ca-certificates netcat-openbsd curl postgresql-client

# Install migrate tool
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz && mv migrate /usr/local/bin/

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/netpulse-backend .

# Copy migration files from the previous stage
COPY --from=builder /app/db/migration ./db/migration

# Copy the unified startup script
COPY --from=builder /app/scripts/startup.sh ./startup.sh

# Make the startup script executable
RUN chmod +x ./startup.sh

# Command to run the unified startup script
ENTRYPOINT ["./startup.sh"]

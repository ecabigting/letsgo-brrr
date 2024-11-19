# Use a lightweight golang image as the base
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Install air
RUN go install github.com/air-verse/air@latest

# Copy the source code
COPY . .

# Build the application
RUN go mod download

# Use a smaller image for the final image
# FROM alpine:latest

# Copy the built binary from the builder stage
# COPY --from=builder /app/main /app/main

# Set the working directory
# WORKDIR /app

# Expose the port
EXPOSE 8080

ARG PORT=8080

# Start the application with air
CMD ["air"]

# Use Go 1.22.5 as the base image for building the application
FROM golang:1.22.5-alpine AS builder

# Set the working directory
WORKDIR /banking

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /banking

# Use a minimal Alpine image for the final runtime container
FROM alpine:latest

# Set the working directory
WORKDIR /banking

# Copy the built binary from the builder stage
COPY --from=builder /banking /banking

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./banking", "apiserver"]

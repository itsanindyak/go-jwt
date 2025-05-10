# Use the official Golang image with version 1.24.3 and Alpine 3.21 as the base image for building the application
FROM golang:1.24.3-alpine3.21 AS builder

# Set the working directory inside the container to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the Go module dependencies specified in go.mod and go.sum
RUN go mod download

# Copy the entire project directory into the container's working directory
COPY . .

# Build the Go application as a statically linked binary named "goapp"
RUN CGO_ENABLED=0 GOOS=linux go build -o goapp .

# Use the latest Alpine Linux image as the base image for the final container
FROM alpine:latest

# Set the working directory inside the container to /root
WORKDIR /root

# Copy the built Go binary from the builder stage to the final image
COPY --from=builder /app/goapp .

# Expose port 9050 to allow external access to the application
EXPOSE 9050

# Specify the command to run the Go application when the container starts
CMD [ "./goapp" ]
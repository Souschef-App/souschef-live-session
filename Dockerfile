# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application inside the container
RUN go build -o main .

# Expose the port the WebSocket server will listen on
EXPOSE 8080

# Command to run the WebSocket server
CMD ["./main"]

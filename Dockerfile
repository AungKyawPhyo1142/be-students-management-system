FROM golang:1.22-alpine

# Set the working directory to root
WORKDIR /

# Copy go.mod and go.sum to the root directory
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project to the container’s root
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port (adjust if your app listens on a different port)
EXPOSE 8080

# Run the application
CMD ["./main"]

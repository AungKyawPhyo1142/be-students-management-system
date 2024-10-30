FROM golang:1.21-alpine

# Add git and build tools
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Set environment variables
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
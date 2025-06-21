# Start from the official Go image
FROM golang:1.24-alpine

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o sensor-consumer

# Run the binary
CMD ["./sensor-consumer"]
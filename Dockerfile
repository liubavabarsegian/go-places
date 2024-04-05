ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}

# Set the working directory in the container
WORKDIR /app
# Copy the Go module files (go.mod and go.sum) to the working directory
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the working directory
COPY . .
# Build the Go app
RUN go build -o  main cmd/main.go
# Expose port 8080 for the application
EXPOSE 9200
# Command to run the executable
CMD ["./main"]

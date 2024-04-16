ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}
ENV CGO_ENABLED = 0
RUN apk update
RUN apk add --no-cache git gcc

RUN mkdir /app
ADD . /app/
# Set the working directory in the container
WORKDIR /app
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the working directory
COPY . .
# Build the Go app
RUN go build -o  main cmd/main.go
# Expose port 8888 for the application
EXPOSE 8888
# Command to run the executable
CMD ["./main"]

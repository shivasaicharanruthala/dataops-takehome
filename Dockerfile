# Use the official Golang image to create a build artifact.
# This stage is referred to as the builder stage.
FROM golang:1.21 AS builder

# Set the Current Working Directory inside the container
WORKDIR /dataops-takehome

# Copy go mod and go sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go get ./...

# Copy the source from the current directory to the Working Directory inside the container
COPY ./database ./etl ./log ./model init.sql main.go go.mod ./

RUN ls
# Build the Go app
RUN GOARCH=amd64 GOOS=linux go build -tags musl -o dataops-takehome

# Use a minimal Ubuntu image for the second stage
FROM ubuntu:22.04

# Set the Current Working Directory inside the container
WORKDIR /dataops-takehome

# Copy the Pre-built binary file from the builder stage
COPY --from=builder /dataops-takehome .

# Command to run the executable
CMD ["./dataops-takehome"]

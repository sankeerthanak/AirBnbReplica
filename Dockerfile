# Use an official Golang runtime as a parent image
FROM golang:1.22-alpine

# Set the current working directory inside the container
WORKDIR /AirbnbReplica

# Copy the current directory contents into the container at /app
COPY . .

# Install any needed dependencies specified in go.mod and go.sum
RUN go mod download

# Build the Go app
RUN go build -o main ./cmd

# Make port 8081 available to the world outside this container
EXPOSE 8081

# Run the executable
CMD ["./main"]

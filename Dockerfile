# Use an official Go runtime as a base image
FROM golang:1.18 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies (can be cached)
RUN go mod download

# Copy the src directory containing your Go files into the container
COPY src/ /app/src/

# Set the working directory to where the Go files are located
WORKDIR /app/src

# Build the Go application
RUN go build -o /app/out

# Final stage: Create a smaller image with only the necessary files
FROM golang:1.18-alpine

# Set the working directory in the final image
WORKDIR /app

# Copy the built binary from the build stage
COPY --from=build /app/out .

# Expose the port that your Go app runs on (change the port if necessary)
EXPOSE 3000

# Command to run the binary
CMD ["./out"]

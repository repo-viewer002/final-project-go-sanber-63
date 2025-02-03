# Build stage
FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src/ /app/src/
WORKDIR /app/src

# Set architecture and OS for building the Go binary
ENV GOARCH=amd64
ENV GOOS=linux

# Build the Go binary
RUN go build -o /app/out
RUN ls -l /app  # Debugging step to check if 'out' exists

# Final stage: smaller image
FROM golang:1.23-alpine

# Install dependencies for dynamic linking
RUN apk --no-cache add libc6-compat

WORKDIR /app

# Copy the built binary
COPY --from=build /app/out .

EXPOSE 3000

# Default command to run the app
CMD ["./out"]

FROM golang:1.24 AS build

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Use a smaller base image for the final stage
FROM alpine:3.21

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/api .

# Set the entrypoint
ENTRYPOINT ["/app/api"]


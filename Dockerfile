# Stage 1: Build the application
FROM golang:1.20 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Stage 2: Create the final image
FROM alpine:latest

# Install ca-certificates to enable HTTPS support
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/myapp .

# Ensure the binary is executable
RUN chmod +x /app/myapp

EXPOSE 1323

# Specify the entry point for the container
CMD ["./myapp"]

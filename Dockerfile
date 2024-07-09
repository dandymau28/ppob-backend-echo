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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ppob .

# Stage 2: Create the final image
FROM alpine:latest

# Install ca-certificates to enable HTTPS support
RUN apk --no-cache add ca-certificates

# create mirror dir
RUN mkdir /mirror

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/ppob /app/ppob
COPY --from=build /app/.env /app/.env

# Ensure the binary is executable
RUN chmod +x /app/ppob

RUN ls -la /app

EXPOSE 1323

# Specify the entry point for the container
CMD ["./ppob"]

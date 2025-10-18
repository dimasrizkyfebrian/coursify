# Stage 1: The "Builder" Stage
# Go application compilation
FROM golang:1.25-alpine AS builder

# Setting the working directory inside the container
WORKDIR /app

# Copying dependency files first to take advantage of Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copying the remaining application source code
COPY . .

# Building a Go application into a static binary.
# CGO_ENABLED=0 is important for creating a binary that does not depend on the C library
# GOOS=linux ensures that this binary can run in a Linux environment (Docker)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /coursify ./cmd/server

# ---

# Stage 2: The "Final" Stage
# Preparing the final image that will be executed
FROM alpine:latest

# Working directory inside the final container
WORKDIR /

# Copying ONLY the binary files that have been compiled from the 'builder' stage
COPY --from=builder /coursify /coursify

# Copy the uploads folder so that it can be accessed by the application
# Note: For persistent data, use a volume in docker-compose
COPY ./uploads ./uploads

# Inform Docker that this container will run on port 8080
EXPOSE 8080

# Run the application
CMD ["/coursify"]
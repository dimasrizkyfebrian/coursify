# Stage 1: The "Builder" Stage
# Go application compilation
FROM golang:1.25 AS builder

# Install C library needed by kolesa-team/go-webp
RUN apt-get update && apt-get install -y --no-install-recommends libwebp-dev gcc build-essential

# Setting the working directory inside the container
WORKDIR /app

# Copying dependency files first to take advantage of Docker layer caching
COPY go.mod go.sum ./

# Download dependencies based on go.mod/go.sum
RUN go mod download

# Copying the remaining application source code
COPY . .

# Building a Go application into a static binary.
RUN go build -a -installsuffix cgo -o /coursify ./cmd/server

# ---

# Stage 2: The "Final" Stage
# Preparing the final image that will be executed
FROM debian:bookworm-slim AS final

# Working directory inside the final container
WORKDIR /

# Install runtime library needed by kolesa-team/go-webp
RUN apt-get update && apt-get install -y --no-install-recommends libwebp7 && rm -rf /var/lib/apt/lists/*

# Copying ONLY the binary files that have been compiled from the 'builder' stage
COPY --from=builder /coursify /coursify

# Copy the uploads folder so that it can be accessed by the application
# Note: For persistent data, use a volume in docker-compose
COPY ./uploads ./uploads

# Inform Docker that this container will run on port 8080
EXPOSE 8080

# Run the application
CMD ["/coursify"]
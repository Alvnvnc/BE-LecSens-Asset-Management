FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git postgresql-client

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates postgresql-client tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy any necessary config files or static assets if needed
COPY --from=builder /app/data-layer/migration/seeder/ ./data-layer/migration/seeder/

EXPOSE 3122

# Create entrypoint script with proper database wait and environment variables
RUN echo "#!/bin/sh" > /entrypoint.sh && \
    echo "echo 'Waiting for database to be ready...'" >> /entrypoint.sh && \
    echo "until pg_isready -h db -U root -d asset_management; do" >> /entrypoint.sh && \
    echo "  echo 'Database not ready, waiting...'" >> /entrypoint.sh && \
    echo "  sleep 3" >> /entrypoint.sh && \
    echo "done" >> /entrypoint.sh && \
    echo "echo 'Database is ready! Starting application...'" >> /entrypoint.sh && \
    echo "./main" >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]

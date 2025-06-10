FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Build the command tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cmd ./helpers/cmd/cmd.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates postgresql-client tzdata bash

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/cmd .

# Copy .env file if it exists
COPY .env* ./

# Copy any necessary config files or static assets if needed
COPY --from=builder /app/data-layer/migration/seeder/ ./data-layer/migration/seeder/
COPY --from=builder /app/helpers/ ./helpers/

EXPOSE 3122

# Run the main application
CMD ["./main"]

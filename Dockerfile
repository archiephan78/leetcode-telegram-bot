# Build stage
FROM golang:1.21-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM alpine:latest

# Install ca-certificates and sqlite
RUN apk --no-cache add ca-certificates sqlite tzdata

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the problems file
COPY --from=builder /app/problem_deduplicated.yaml .

# Create data directory for database
RUN mkdir -p /data

# Expose port (if needed for health checks)
EXPOSE 8080

# Set environment variables
ENV DATABASE_PATH=/data/leetcode_bot.db
ENV PROBLEMS_FILE_PATH=./problem_deduplicated.yaml
ENV TIMEZONE=Asia/Ho_Chi_Minh

# Run the application
CMD ["./main"] 
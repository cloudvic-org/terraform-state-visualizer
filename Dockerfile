# Multi-stage build for terraform-state-visualizer
FROM golang:1.25.3-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod ./

# Download dependencies (if go.sum exists)
RUN if [ -f go.sum ]; then go mod download; fi

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o terraform-state-visualizer .

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/terraform-state-visualizer .

# Change ownership to non-root user
RUN chown appuser:appgroup terraform-state-visualizer

# Switch to non-root user
USER appuser

# Expose port (if needed for future web interface)
EXPOSE 8080

# Set default command
ENTRYPOINT ["./terraform-state-visualizer"]

# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install packr2 for static file embedding
RUN go install github.com/gobuffalo/packr/v2/packr2@latest 

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Embed static files and build
RUN packr2 build -o main .

# Runtime stage
FROM alpine:latest

# Install system dependencies
RUN apk add --no-cache \
    ca-certificates \
    curl \
    jq \
    && update-ca-certificates

# Install yq (YAML processor)
RUN wget https://github.com/mikefarah/yq/releases/download/v4.34.1/yq_linux_amd64 -O /usr/local/bin/yq \
    && chmod +x /usr/local/bin/yq

# Install Helm
ENV VERIFY_CHECKSUM=false
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 \
    && chmod 700 get_helm.sh \
    && sh ./get_helm.sh \
    && rm get_helm.sh

# Create app directory
WORKDIR /app

# Copy built binary
COPY --from=builder /app/main .
COPY dist /app/dist

# Expose port
EXPOSE 5358

# Start application
CMD ["./main"]

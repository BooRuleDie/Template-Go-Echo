# Main application Dockerfile using FROM scratch
FROM golang:1.25.1-alpine AS build

# Install ca-certificates in the build stage
RUN apk add --no-cache ca-certificates git tzdata

WORKDIR /build

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary with all optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/server

# Final stage - FROM scratch for minimal image (~5-10MB)
FROM scratch

# Copy ca-certificates from build stage (needed for HTTPS)
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data (if your app needs it)
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the static binary
COPY --from=build /build/main /main

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/main"]

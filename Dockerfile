# Stage 1: Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o personal-vault-backend .

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=builder /app/personal-vault-backend .

EXPOSE 8080

CMD ["./personal-vault-backend"]

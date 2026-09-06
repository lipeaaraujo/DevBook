FROM golang:1.21.0 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN go build -o main .

FROM debian:bookworm-slim

WORKDIR /root/

COPY --from=builder --chmod=755 /app/main .
COPY --from=builder /app/migrations ./migrations
RUN chmod +x main

EXPOSE 4000

# Run
CMD ["./main"]

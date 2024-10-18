FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY main.go main.go
COPY internal/ internal/

RUN go build -o main .

FROM gcr.io/distroless/base

# Copy the Go binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

CMD ["main"]
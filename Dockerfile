FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY main.go main.go
COPY internal/ internal/

RUN go build -o git-delta .

FROM gcr.io/distroless/static-debian12

# Copy the Go binary from the builder stage
COPY --from=builder /app/git-delta /usr/local/bin/git-delta

CMD ["git-delta"]
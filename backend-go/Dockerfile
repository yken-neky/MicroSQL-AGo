### Build stage
FROM golang:1.24.4-alpine AS builder
WORKDIR /src

# install git/ca-certificates for go modules over https
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /app/bin/server ./cmd/server

# Build wait-for-mysql binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /app/bin/wait_for_mysql ./cmd/wait
 
# Build migrator binary (with CGO enabled for SQLite support)
RUN apk add --no-cache gcc musl-dev sqlite-dev
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /app/bin/migrate ./cmd/migrate

### Final image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates sqlite-libs
WORKDIR /app

COPY --from=builder /app/bin/server ./server
COPY --from=builder /app/bin/wait_for_mysql ./wait_for_mysql
COPY --from=builder /app/bin/migrate ./migrate

# Copy migrations or other static files if needed
EXPOSE 8000
ENV GIN_MODE=release
ENTRYPOINT ["/bin/sh","-c","/app/wait_for_mysql && exec /app/server"]

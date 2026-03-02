FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.857

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install Tailwind CSS standalone
RUN wget -qO /usr/local/bin/tailwindcss https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 && \
    chmod +x /usr/local/bin/tailwindcss

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate code
RUN templ generate
RUN sqlc generate
RUN tailwindcss -i internal/static/css/input.css -o internal/static/css/output.css --minify

# Build
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/bin/server ./cmd/server
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /app/bin/reset-password ./cmd/reset-password

# --- Runtime ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/server /usr/local/bin/server
COPY --from=builder /app/bin/reset-password /usr/local/bin/reset-password

VOLUME /data
ENV DATA_DIR=/data
EXPOSE 8090

ENTRYPOINT ["server"]

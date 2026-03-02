FROM golang:1.26.0-alpine AS builder

RUN apk add --no-cache gcc musl-dev

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

# Install Tailwind CSS standalone
RUN apk add --no-cache nodejs npm && \
    npm install -g tailwindcss@4.2.1 @tailwindcss/cli@4.2.1

ENV NODE_PATH=/usr/local/lib/node_modules

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
FROM alpine:3.23

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/server /usr/local/bin/server
COPY --from=builder /app/bin/reset-password /usr/local/bin/reset-password

VOLUME /data
ENV DATA_DIR=/data
EXPOSE 8090

ENTRYPOINT ["server"]

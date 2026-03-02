FROM golang:1.26.0-alpine

RUN apk add --no-cache gcc musl-dev nodejs npm

# Install Go tools
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001 && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 && \
    go install github.com/air-verse/air@v1.64.5

# Install Tailwind CSS
RUN npm install -g tailwindcss@4.2.1 @tailwindcss/cli@4.2.1

ENV NODE_PATH=/usr/local/lib/node_modules

WORKDIR /app

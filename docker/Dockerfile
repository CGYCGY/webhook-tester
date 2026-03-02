FROM golang:1.23-alpine

RUN apk add --no-cache gcc musl-dev nodejs npm

# Install Go tools
RUN go install github.com/a-h/templ/cmd/templ@v0.3.857 && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 && \
    go install github.com/air-verse/air@v1.61.7

# Install Tailwind CSS
RUN npm install -g tailwindcss@3

WORKDIR /app

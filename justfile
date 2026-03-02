# Default recipe
default:
    @just --list

# Run development server with Air hot-reload
dev:
    air

# Generate templ templates
templ:
    templ generate

# Generate sqlc queries
sqlc:
    sqlc generate

# Build Tailwind CSS
css:
    tailwindcss -i internal/static/css/input.css -o internal/static/css/output.css --minify

# Watch Tailwind CSS
css-watch:
    tailwindcss -i internal/static/css/input.css -o internal/static/css/output.css --watch

# Generate all (templ + sqlc + css)
generate: templ sqlc css

# Build the server binary
build: generate
    go build -o bin/server ./cmd/server

# Build the reset-password CLI
build-cli:
    go build -o bin/reset-password ./cmd/reset-password

# Run the server
run: build
    ./bin/server

# Run tests
test:
    go test ./...

# Docker build
docker-build:
    docker build -t webhook-tester .

# Docker run
docker-run:
    docker run -p 8090:8090 \
        -v webhook-tester-data:/data \
        -e ADMIN_EMAIL=admin@example.com \
        -e ADMIN_PASSWORD=changeme123 \
        -e JWT_SECRET=change-this-secret \
        webhook-tester

# Clean build artifacts
clean:
    rm -rf bin/ tmp/ internal/static/css/output.css

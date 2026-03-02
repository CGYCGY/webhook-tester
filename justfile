compose := "docker compose -f docker/docker-compose.yml --env-file .env"
exec := compose + " exec dev"

# Default recipe
default:
    @just --list

# Start dev container
up:
    {{compose}} up -d --build

# Stop dev container
down:
    {{compose}} down

# Generate templ templates
templ:
    {{exec}} templ generate

# Generate sqlc queries
sqlc:
    {{exec}} sqlc generate

# Build Tailwind CSS
css:
    {{exec}} tailwindcss -i internal/static/css/input.css -o internal/static/css/output.css --minify

# Generate all (templ + sqlc + css)
generate: templ sqlc css

# Run Air hot-reload server inside container
dev:
    {{compose}} exec -e PORT=8090 dev air

# Run Go tests
test:
    {{exec}} go test ./...

# Run Go vet
vet:
    {{exec}} go vet ./...

# Build production Docker image
build:
    docker build -f deploy/Dockerfile -t webhook-tester .

# Build, push, and deploy via Coolify
deploy:
    bash deploy/deploy.sh

# Remove build artifacts
clean:
    rm -rf bin/ tmp/ internal/static/css/output.css

# Remove everything including Docker volumes
clean-all: clean
    {{compose}} down -v

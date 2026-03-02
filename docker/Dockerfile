FROM golang:1.26.0-alpine

RUN apk add --no-cache gcc musl-dev nodejs npm

# Install dev tools via prebuilt binaries
RUN wget -qO- https://github.com/a-h/templ/releases/download/v0.3.1001/templ_Linux_x86_64.tar.gz | tar -xz -C /usr/local/bin templ && \
    wget -qO- https://github.com/sqlc-dev/sqlc/releases/download/v1.30.0/sqlc_1.30.0_linux_amd64.tar.gz | tar -xz -C /usr/local/bin sqlc && \
    wget -qO /usr/local/bin/air https://github.com/air-verse/air/releases/download/v1.64.5/air_1.64.5_linux_amd64 && \
    chmod +x /usr/local/bin/air

# Install Tailwind CSS
RUN npm install -g tailwindcss@4.2.1 @tailwindcss/cli@4.2.1

ENV NODE_PATH=/usr/local/lib/node_modules

WORKDIR /app

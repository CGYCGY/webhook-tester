#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Load config
if [[ -f "$SCRIPT_DIR/.env.deploy" ]]; then
    set -a
    source "$SCRIPT_DIR/.env.deploy"
    set +a
fi

# Required vars
: "${GITHUB_ORG:?GITHUB_ORG is required in deploy/.env.deploy}"
: "${REPO_NAME:?REPO_NAME is required in deploy/.env.deploy}"

TAG="${1:-latest}"
IMAGE="ghcr.io/${GITHUB_ORG}/${REPO_NAME}"
GIT_SHA="$(git -C "$PROJECT_ROOT" rev-parse --short HEAD)"

echo "Building ${IMAGE}:${TAG} ..."
docker build -f "$SCRIPT_DIR/Dockerfile" -t "${IMAGE}:${TAG}" -t "${IMAGE}:${GIT_SHA}" "$PROJECT_ROOT"

echo "Pushing ${IMAGE}:${TAG} ..."
docker push "${IMAGE}:${TAG}"

echo "Pushing ${IMAGE}:${GIT_SHA} ..."
docker push "${IMAGE}:${GIT_SHA}"

if [[ -n "${COOLIFY_WEBHOOK_URL:-}" && -n "${COOLIFY_API_TOKEN:-}" ]]; then
    echo "Triggering Coolify redeploy ..."
    curl -fsSL -H "Authorization: Bearer ${COOLIFY_API_TOKEN}" "$COOLIFY_WEBHOOK_URL" || echo "Warning: Coolify webhook failed"
fi

echo "Done! Deployed ${IMAGE}:${TAG} (${GIT_SHA})"

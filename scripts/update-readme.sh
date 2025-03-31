#!/bin/bash

# Get versions from project files
NEXT_VERSION=$(cd frontend && grep -m 1 '"next":' package.json | sed 's/.*"next": "\([^"]*\)",/\1/')
TS_VERSION=$(cd frontend && grep -m 1 '"typescript":' package.json | sed 's/.*"typescript": "\([^"]*\)",\{0,1\}/\1/')
TAILWIND_VERSION=$(cd frontend && grep -m 1 '"tailwindcss":' package.json | sed 's/.*"tailwindcss": "\([^"]*\)",\{0,1\}/\1/')
GO_VERSION=$(cd backend && grep -m 1 'go [0-9]' go.mod | sed 's/go \([0-9.]*\)/\1/')
POSTGRES_VERSION=$(cd backend && grep -m 1 'bitnami/postgresql' docker-compose.yaml | sed 's/.*bitnami\/postgresql:\([^"]*\)/\1/')
DOCKER_VERSION=$(cd backend && grep -m 1 'version:' docker-compose.yaml | sed 's/version: "\([^"]*\)"/\1/')

# Default to 'latest' if version is not found
POSTGRES_VERSION=${POSTGRES_VERSION:-latest}
DOCKER_VERSION=${DOCKER_VERSION:-latest}

# Update frontend README
sed -i.bak \
  -e "s/badge\/Next.js-[^-]*-black/badge\/Next.js-${NEXT_VERSION}-black/" \
  -e "s/badge\/TypeScript-[^-]*-blue/badge\/TypeScript-${TS_VERSION}-blue/" \
  -e "s/badge\/Tailwind-[^-]*-38B2AC/badge\/Tailwind-${TAILWIND_VERSION}-38B2AC/" \
  frontend/README.md

# Update backend README
sed -i.bak \
  -e "s/badge\/Go-[^-]*-00ADD8/badge\/Go-${GO_VERSION}-00ADD8/" \
  -e "s/badge\/PostgreSQL-[^-]*-336791/badge\/PostgreSQL-${POSTGRES_VERSION}-336791/" \
  -e "s/badge\/Docker-[^-]*-2496ED/badge\/Docker-${DOCKER_VERSION}-2496ED/" \
  backend/README.md

# Update root README
sed -i.bak \
  -e "s/badge\/Next.js-[^-]*-black/badge\/Next.js-${NEXT_VERSION}-black/" \
  -e "s/badge\/TypeScript-[^-]*-blue/badge\/TypeScript-${TS_VERSION}-blue/" \
  -e "s/badge\/Go-[^-]*-00ADD8/badge\/Go-${GO_VERSION}-00ADD8/" \
  -e "s/badge\/PostgreSQL-[^-]*-336791/badge\/PostgreSQL-${POSTGRES_VERSION}-336791/" \
  README.md

# Remove backup files
rm -f frontend/README.md.bak backend/README.md.bak README.md.bak

echo "READMEs updated with current versions:"
echo "Next.js: ${NEXT_VERSION}"
echo "TypeScript: ${TS_VERSION}"
echo "Tailwind CSS: ${TAILWIND_VERSION}"
echo "Go: ${GO_VERSION}"
echo "PostgreSQL: ${POSTGRES_VERSION}"
echo "Docker: ${DOCKER_VERSION}" 
# Builds the frontend and the backend, and runs both in one container on port 8080.
#   docker build -t sezzle-calculator .
#   docker run --rm -p 8080:8080 sezzle-calculator

# 1. Frontend: static files in /app/dist
FROM node:24-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# 2. Backend: a static Go binary
FROM golang:1.27-alpine AS backend
WORKDIR /src
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# 3. Runtime: small image, non-root user
FROM alpine:3.22
RUN adduser -D -H app
WORKDIR /app
COPY --from=backend /out/server ./server
COPY --from=frontend /app/dist ./static
ENV PORT=8080 STATIC_DIR=/app/static
EXPOSE 8080
USER app
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["./server"]

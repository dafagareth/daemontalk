# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26-bookworm AS builder
WORKDIR /app

# templ CLI (pinned to the version in go.mod) for component code generation.
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Regenerate *_templ.go from .templ sources, then build a static binary.
# modernc.org/sqlite is pure Go, so CGO can stay disabled (smaller, portable).
RUN templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o daemontalk .

# ---- Runtime stage ----
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --home /app app

# Runtime needs: the binary, blog content (read at startup), static assets
# (served + chroma.css regenerated here), and a writable data dir for SQLite.
COPY --from=builder /app/daemontalk .
COPY --from=builder /app/web/static/ web/static/
COPY --from=builder /app/content/ content/
RUN mkdir -p data && chown -R app:app data web/static

USER app
EXPOSE 8080
EXPOSE 2222
ENV PORT=8080
ENV SSH_PORT=2222
# Persist the comments/views database across container recreations.
VOLUME ["/app/data"]
CMD ["./daemontalk"]

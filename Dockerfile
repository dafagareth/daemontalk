# syntax=docker/dockerfile:1
# Production multi-stage build: compiles static Go binary with Templ and packs into a minimal Debian runtime.

FROM golang:1.26-bookworm AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /app

RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN templ generate && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o daemontalk .

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates tzdata git \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --uid 10001 --home /app app

COPY --from=builder /app/daemontalk .
COPY --from=builder /app/google*.html ./
COPY --from=builder /app/web/static/ web/static/
COPY --from=builder /app/content/ content/
RUN mkdir -p data content/posts web/static/images/posts && chown -R app:app data content web/static
RUN git config --system --add safe.directory /app

USER app
EXPOSE 8080
EXPOSE 2222
ENV PORT=8080
ENV SSH_PORT=2222
VOLUME ["/app/data"]
CMD ["./daemontalk"]

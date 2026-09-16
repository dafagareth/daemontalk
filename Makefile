.PHONY: all full build build-prod run dev preview draft assets generate templ css css-dev bundle-js test test-v test-race test-cover test-html vet fmt lint tidy clean deploy backup restore health setup-vps new-post new-uid list-posts stats-posts validate-posts archive-post restore-post delete-post docker-build docker-up docker-down docker-logs docker-restart docker-ps docker-dev docker-dev-down

all: full

full: fmt assets vet test build
	@echo "[ok] Full pipeline completed: assets generated, tests passed, and binary built!"

TAILWIND := $(shell which tailwindcss 2>/dev/null || echo $(CURDIR)/bin/tailwindcss)

$(CURDIR)/bin/tailwindcss:
	@mkdir -p $(CURDIR)/bin
	@echo "[info] Downloading Tailwind CSS standalone CLI..."
	@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	ARCH=$$(uname -m); \
	if [ "$$ARCH" = "x86_64" ]; then ARCH="x64"; fi; \
	if [ "$$ARCH" = "aarch64" ]; then ARCH="arm64"; fi; \
	URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$${OS}-$${ARCH}"; \
	curl -sL "$$URL" -o $(CURDIR)/bin/tailwindcss && chmod +x $(CURDIR)/bin/tailwindcss

run: build
	@./daemontalk

dev: generate $(TAILWIND)
	@trap 'kill 0' EXIT INT TERM; \
	templ generate --watch < /dev/null & \
	$(TAILWIND) -i web/static/css/input.css -o web/static/css/main.css --watch < /dev/null & \
	air

preview: build
	@./daemontalk

draft: preview

generate: templ

templ:
	@templ generate

css: $(TAILWIND)
	@$(TAILWIND) -i web/static/css/input.css \
		-o web/static/css/main.css --minify

css-dev: $(TAILWIND)
	@$(TAILWIND) -i web/static/css/input.css \
		-o web/static/css/main.css

bundle-js:
	@cat web/static/js/post/utils.js \
		web/static/js/post/bookmarks.js \
		web/static/js/post/reading-modes.js \
		web/static/js/post/code-blocks.js \
		web/static/js/post/share.js \
		web/static/js/post/toc.js \
		web/static/js/post/task-lists.js \
		web/static/js/post/footnotes.js \
		web/static/js/post/code-tabs.js \
		web/static/js/post/read-status.js \
		web/static/js/post/comments.js > web/static/js/post.bundle.js

assets: generate css bundle-js

build: assets
	@go build -trimpath -ldflags="-s -w" -o daemontalk .

build-prod: assets
	@CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o daemontalk .

test:
	@go test ./...

test-v:
	@go test -v ./...

test-race:
	@go test -race ./...

test-cover:
	@go test -cover ./...

test-html:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@rm -f coverage.out

fmt:
	@go fmt ./...
	@templ fmt web/templates/ 2>/dev/null || true

vet:
	@go vet ./...

lint: fmt vet

tidy:
	@go mod tidy
	@go mod verify

new-post:
	@./scripts/post.sh new

new-uid:
	@./scripts/post.sh new --uid

list-posts:
	@./scripts/post.sh list

stats-posts:
	@./scripts/post.sh stats

validate-posts:
	@./scripts/post.sh validate

archive-post:
	@./scripts/post.sh archive $(SLUG)

restore-post:
	@./scripts/post.sh restore $(SLUG)

delete-post:
	@./scripts/post.sh delete $(SLUG)

docker-build:
	@docker build -t daemontalk .

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

docker-logs:
	@docker compose logs -f

docker-restart:
	@docker compose restart

docker-ps:
	@docker compose ps

docker-dev:
	@docker compose -f docker-compose.dev.yml up --build

docker-dev-down:
	@docker compose -f docker-compose.dev.yml down

backup:
	@./scripts/backup.sh

restore:
	@./scripts/restore.sh

health:
	@./scripts/healthcheck.sh

deploy:
	@./scripts/deploy.sh

setup-vps:
	@./scripts/setup-vps.sh

clean:
	@rm -f daemontalk
	@rm -f web/static/css/main.css
	@rm -f coverage.out

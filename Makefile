.PHONY: build dev css generate test vet clean deploy

generate:
	templ generate

css:
	npx @tailwindcss/cli -i web/static/css/input.css \
		-o web/static/css/main.css --minify

build: generate css
	go build -trimpath -ldflags="-s -w" -o daemontalk .

dev: generate
	npx @tailwindcss/cli -i web/static/css/input.css -o web/static/css/main.css
	templ generate --watch < /dev/null &
	npx @tailwindcss/cli -i web/static/css/input.css \
		-o web/static/css/main.css --watch < /dev/null &
	air

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -f daemontalk
	rm -f web/static/css/main.css

deploy:
	@./scripts/deploy.sh

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


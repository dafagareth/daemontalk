.PHONY: build dev css generate test vet clean deploy

generate:
	templ generate

css:
	npx @tailwindcss/cli -i web/static/css/input.css \
		-o web/static/css/main.css --minify

build: generate css
	go build -o daemontalk .

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

deploy: build
	fly deploy

new-post:
	@SLUG=$$(openssl rand -hex 4); \
	FILE="content/posts/$$SLUG.md"; \
	IMG_DIR="web/static/images/posts/$$SLUG"; \
	mkdir -p $$IMG_DIR; \
	echo "---" > $$FILE; \
	echo "title: \"Judul Postingan Baru\"" >> $$FILE; \
	echo "slug: $$SLUG" >> $$FILE; \
	echo "aliases: []" >> $$FILE; \
	echo "date: $$(date +%Y-%m-%d)" >> $$FILE; \
	echo "tags: []" >> $$FILE; \
	echo "lang: id" >> $$FILE; \
	echo "draft: true" >> $$FILE; \
	echo "---" >> $$FILE; \
	echo "Postingan baru (draft) berhasil dibuat di: $$FILE"; \
	echo "Folder aset gambar berhasil dibuat di: $$IMG_DIR/"

archive-post:
	@./scripts/manage-post.sh archive $(SLUG)

restore-post:
	@./scripts/manage-post.sh restore $(SLUG)

delete-post:
	@./scripts/manage-post.sh delete $(SLUG)


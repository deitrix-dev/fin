styles:
	npx tailwindcss -i input.css -o web/assets/style.css

templ:
	templ generate

build: templ
	CGO_ENABLED=0 go build -o ./bin/finserve ./cmd/finserve

finserve: styles build


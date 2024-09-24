styles:
	npx tailwindcss -i input.css -o web/assets/style.css

build:
	CGO_ENABLED=0 go build -o ./bin/finserve ./cmd/finserve

finserve: styles build


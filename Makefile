styles:
	npx tailwindcss -i input.css -o web/assets/style.css

templ:
	TEMPL_EXPERIMENT=rawgo go run github.com/a-h/templ/cmd/templ generate

build: templ
	CGO_ENABLED=0 go build -o ./bin/finserve ./cmd/finserve

finserve: styles build


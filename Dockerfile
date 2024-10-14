FROM node:20 as assets

WORKDIR /opt/fin

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npx rollup -c
RUN npx tailwindcss -i input.css -o web/assets/style.css

FROM golang:1.23 as build

RUN apt update && apt install -y ca-certificates

RUN useradd -u 10001 fin

WORKDIR /opt/fin

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

COPY . .
COPY --from=assets /opt/fin/web/assets/style.css web/assets/style.css
COPY --from=assets /opt/fin/web/assets/index.js web/assets/index.js

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o ./bin/finserve ./cmd/finserve

VOLUME /var/lib/fin

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /etc/passwd /etc/passwd
USER fin

WORKDIR /opt/fin

COPY --from=build --chown=fin /opt/fin/bin/finserve /opt/fin/bin/finserve
COPY --from=build --chown=fin /opt/fin/store/mysql/schema.sql /opt/fin/store/mysql/schema.sql

ENTRYPOINT ["/opt/fin/bin/finserve"]

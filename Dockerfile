FROM node:20 as assets

WORKDIR /opt/fin

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npx tailwindcss -i input.css -o web/assets/style.css

FROM golang:1.23 as build

WORKDIR /opt/fin

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY --from=assets /opt/fin/web/assets/style.css web/assets/style.css

RUN go generate ./...
RUN CGO_ENABLED=0 go build -o ./bin/finserve ./cmd/finserve

VOLUME /var/lib/fin

FROM scratch

WORKDIR /opt/fin

COPY --from=build /opt/fin/bin/finserve /opt/fin/bin/finserve

ENTRYPOINT ["/opt/fin/bin/finserve"]

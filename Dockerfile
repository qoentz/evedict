# === Assets Stage ===
FROM node:20 AS assets

WORKDIR /evedict

RUN apt-get update && apt-get install -y git

COPY package.json ./
RUN npm install

COPY tailwind.config.js ./
COPY internal/view/css/input.css ./internal/view/css/input.css

RUN npx tailwindcss -i ./internal/view/css/input.css -o ./dist/styles.css --minify

# === Go binary ===
FROM golang:1.23-alpine AS builder

WORKDIR /evedict

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

COPY internal/promptgen/prompts.yaml ./internal/promptgen/prompts.yaml

RUN go build -o /evedict/main ./cmd/app/main.go

# === Runtime image ===
FROM alpine:latest AS final

WORKDIR /evedict

COPY --from=builder /evedict/main ./main
COPY --from=builder /evedict/internal/promptgen/prompts.yaml ./internal/promptgen/prompts.yaml
COPY --from=assets /evedict/dist/styles.css ./static/css/styles.css
COPY static ./static

EXPOSE 8080

CMD ["./main"]




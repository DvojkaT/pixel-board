FROM golang:1.26-alpine AS development

RUN go install github.com/air-verse/air@latest

WORKDIR /app

ENTRYPOINT ["air", "-c", ".air.toml"]

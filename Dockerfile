FROM ghcr.io/a-h/templ:latest AS generate
WORKDIR /app
COPY --chown=65532:65532 . .
RUN ["templ", "generate"]

FROM node:24-alpine AS frontend-build
WORKDIR /app
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install
COPY --from=generate /app /app
RUN cd frontend && npm run minify:css && npm run bundle:js

FROM golang:1.26.2-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY --from=generate /app .
COPY --from=frontend-build /app/frontend/public ./frontend/public
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/api/main.go

FROM golang:1.26.2-alpine AS watch
WORKDIR /app
RUN apk add --no-cache nodejs npm wget make
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest && \
    go clean -cache
COPY frontend/package*.json ./frontend/
RUN cd frontend && npm install && npm cache clean --force
CMD ["air", "-c", ".air.docker.toml"]

FROM alpine:3.23 AS prod
RUN apk add --no-cache wget
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/frontend /app/frontend
ARG PORT=8080
ARG TLS_PORT=8443
EXPOSE ${PORT}
EXPOSE ${TLS_PORT}
CMD ["/app/main"]
# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS api-build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.22 AS api

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=api-build /out/server ./server

EXPOSE 8080

USER nobody:nobody

ENTRYPOINT ["/app/server"]

FROM node:24-alpine AS frontend-build

WORKDIR /src/frontend

RUN corepack enable && corepack prepare pnpm@11.25.0 --activate

COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend/index.html frontend/vite.config.js ./
COPY frontend/src ./src

RUN pnpm build

FROM nginx:alpine AS frontend

COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=frontend-build /src/frontend/dist /usr/share/nginx/html

EXPOSE 80

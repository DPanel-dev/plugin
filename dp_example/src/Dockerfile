FROM node:22-alpine AS frontend

WORKDIR /build/app
COPY app/package.json ./
RUN npm install --no-audit --no-fund
COPY app/ ./
RUN npm run build

FROM golang:1.23-alpine AS backend

WORKDIR /build/server
COPY server/go.mod server/main.go ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/example .

FROM alpine:3.21
WORKDIR /app
COPY --from=backend /out/example ./example
COPY --from=frontend /build/app/dist ./asset
EXPOSE 18080
CMD ["./example"]

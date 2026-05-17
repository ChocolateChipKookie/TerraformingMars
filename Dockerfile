# Build fronend
FROM node:24-alpine AS frontend
WORKDIR /frontend
COPY package*.json .
COPY ./backend/shared/game-data.json /backend/shared/game-data.json
COPY ./frontend .
RUN npm ci
RUN npm run build

# Build backend
FROM golang:1.25-alpine AS backend
WORKDIR /backend

# Install build dependencies for libwebp
RUN apk add --no-cache gcc musl-dev libwebp-dev

COPY ./backend/ .
RUN go build -o tm-server

# Create the service
FROM alpine:latest
RUN apk add --no-cache libwebp ca-certificates
COPY --from=backend /backend/tm-server /tm-server
COPY --from=frontend /frontend/dist /dist
USER 65534:65534
EXPOSE 8080
ENTRYPOINT ["/tm-server"]

# syntax=docker/dockerfile:1.7
FROM golang:1.25-alpine AS build
WORKDIR /src

# deps first for caching
COPY go.mod go.sum ./
RUN go mod download

# app code
COPY . .

# build the api under cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /bin/app ./cmd/api

# --- runtime
FROM alpine:3.20
# install certs + wget for healthcheck
RUN apk --no-cache add ca-certificates wget && adduser -D -g '' app

# copy as root, then drop privileges
COPY --from=build /bin/app /app
USER app

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --retries=5 CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/app"]
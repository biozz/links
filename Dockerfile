# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24 AS builder
WORKDIR /src

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY main.go ./
COPY migrations ./migrations
COPY web ./web

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/links .

FROM alpine:3.24.3
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /app
COPY --from=builder /out/links /app/links
COPY VERSION /app/VERSION

EXPOSE 8080
ENTRYPOINT ["/app/links"]
CMD ["serve", "--http=0.0.0.0:8080"]

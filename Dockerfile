FROM golang:1.24.3-alpine as builder
WORKDIR /app/
COPY . .
RUN CGO_ENABLED=0 go build -o bin/links main.go

FROM alpine:3.20
WORKDIR /app/
COPY --from=builder /app/bin/links .
EXPOSE 8080
ENTRYPOINT ["./links"]
CMD ["serve", "--http=0.0.0.0:8080"]

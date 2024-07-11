FROM golang:1.20.1 as builder
WORKDIR /app
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -v -o stores-server ./cmd/stores-server

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/stores-server /usr/local/bin/stores-server
COPY --from=builder /app/data /data
CMD ["/usr/local/bin/stores-server"]

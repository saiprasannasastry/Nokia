# build the server binary
FROM golang:1.14 AS builder
LABEL stage=server-intermediate
WORKDIR /go/src/github.com/Nokia
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server

# copy the server binary from builder stage; run the server binary
FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /go/src/github.com/Nokia/www/ www/
COPY --from=builder /go/src/github.com/Nokia/bin/server .
Expose 8080
ENTRYPOINT ["server"]

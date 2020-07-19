FROM golang:1.14 AS builder
LABEL stage=server-intermediate
WORKDIR /go/src/github.com/Nokia
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/consumer ./cmd/consumer

# copy the consumer from binary 
FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /go/src/github.com/Nokia/bin/consumer .
Expose 8080
ENTRYPOINT ["consumer"]
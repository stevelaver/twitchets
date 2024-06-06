# Builder Image
FROM golang:1.21 as builder

WORKDIR /twitchets
COPY . .
RUN go mod download
RUN go build -v -o bin/twitchets

# Ditribution Image
FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=builder /twitchets/bin/twitchets /twitchets

EXPOSE 5656

ENTRYPOINT ["/twitchets"]

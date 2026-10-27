FROM golang:1.21-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=${VERSION}" -o /out/service ./cmd/service

FROM alpine:3.21
RUN adduser -D -H -u 65532 app
COPY --from=builder /out/service /usr/local/bin/service
USER 65532:65532
EXPOSE 8080
ENV APP_ADDR=:8080
ENTRYPOINT ["/usr/local/bin/service"]

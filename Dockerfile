FROM golang:1.20 AS builder
RUN apt update
WORKDIR /usr/local/bin
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make deps build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/.build/dalil .
COPY --from=builder /usr/local/bin/configs ./configs
RUN addgroup -g 1000 app
RUN adduser -S app -u 1000 -G app
RUN chown -R app:app ./
USER app
EXPOSE 8080
CMD ["./dalil"]

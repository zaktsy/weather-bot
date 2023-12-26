FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod  go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /weather-bot

FROM alpine:3
COPY --from=builder weather-bot /bin/main
ENTRYPOINT ["/bin/main"]
FROM golang:1.14-alpine AS builder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /workdir

COPY ./go.mod go.mod
COPY ./go.sum go.sum

RUN go mod download

COPY ./cmd ./cmd

RUN go build -ldflags="-w -s" -o /go/bin/eventing-upgrade-check  ./cmd/

FROM alpine:3.11

COPY --from=builder /go/bin/eventing-upgrade-check /eventing-upgrade-check

ENTRYPOINT ["/eventing-upgrade-check"]

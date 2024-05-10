FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o zenexport main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/zenexport /build/zenexport

ENTRYPOINT ["/build/zenexport", "-d"]
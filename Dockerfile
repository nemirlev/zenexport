## Стадия сборки
#FROM golang:1.21 AS builder
#
#WORKDIR /app
#
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify
#
#COPY . .
#
#RUN go build -o /app/build/ ./...
#
#FROM alpine:latest
#
#WORKDIR /app
#
#COPY --from=builder --chmod=755 /app/build/ /app/
#
#RUN chmod 755 /app/build/
#
#CMD /app zenexport -d -interval 180

FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o zenexport main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/zenexport /build/zenexport

ENTRYPOINT ["/build/zenexport", "-d"]
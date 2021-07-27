FROM golang:1.16.6-alpine as builder

RUN mkdir /build
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o botArena .

FROM alpine:3.14

COPY --from=builder /build/botArena /usr/bin
ENTRYPOINT [ "/usr/bin/botArena" ]
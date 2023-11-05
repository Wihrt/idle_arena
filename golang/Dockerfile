FROM golang:1.16.6-alpine as builder

ARG GOBINARY
ARG GOFILE

RUN mkdir /build
COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ${GOBINARY} ${GOFILE}

FROM alpine:3.15

ARG GOBINARY

COPY --from=builder /build/$GOBINARY /usr/bin
RUN chmod a+x /usr/bin/$GOBINARY
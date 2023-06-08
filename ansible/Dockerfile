FROM python:3-alpine3.18

ARG ANSIBLE_VERSION

RUN apk add --no-cache --update --virtual ansible_build rust cargo libffi-dev openssl-dev build-base && \
    pip install ansible==${ANSIBLE_VERSION} && \
    apk del --purge ansible_build && \
    apk add --no-cache --update openssh
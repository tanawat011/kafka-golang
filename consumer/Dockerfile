FROM golang:1.13.6-alpine as build-stage

ENV GO111MODULE=on

WORKDIR /go/src/iOMXG-consumer

COPY ./go.mod .

RUN mkdir /fresh_conf
COPY ./fresh.conf /fresh_conf

RUN apk update
RUN apk add --no-cache tzdata librdkafka librdkafka-dev build-base krb5 libsasl
ENV TZ Asia/Bangkok
# RUN apk update

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

#gen version
RUN echo "{\"version\":\"1.0.0\", \"build\":\"`date +%Y/%m/%d\ %H\:%M\:%S`\"}" >> /go/src/iOMXG-consumer/api.version

# hot-reloaded with FRESH
CMD go get github.com/pilu/fresh \
    && fresh -c /fresh_conf/fresh.conf
# CMD go run main.go

EXPOSE 8080

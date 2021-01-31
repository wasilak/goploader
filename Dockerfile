FROM  --platform=$BUILDPLATFORM quay.io/wasilak/golang:1.15-alpine as builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM

RUN apk add --update --no-cache git
RUN go get github.com/GeertJohan/go.rice
RUN go get github.com/GeertJohan/go.rice/rice
RUN mkdir -p /goploader/
ADD ./ /goploader/
WORKDIR /goploader/server

COPY --from=tonistiigi/xx:golang / /

RUN rice embed-go
RUN go build -o server .

FROM --platform=$BUILDPLATFORM quay.io/wasilak/alpine:3
COPY --from=builder /goploader/server/server /goploader

RUN mkdir /up/
CMD ["/goploader", "--conf=./"]

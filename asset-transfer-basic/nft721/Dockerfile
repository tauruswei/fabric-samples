# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0

ARG GO_VER=1.14.12
ARG ALPINE_VER=3.12

FROM golang:${GO_VER}-alpine${ALPINE_VER}

WORKDIR /go/src/github.com/hyperledger/fabric-samples/asset-transfer-basic/nft721
COPY . .
ENV GOPROXY https://goproxy.cn,direct
RUN apk add gcc && apk add g++

#RUN go get -d -v ./...
#RUN go install -v ./...
RUN go install -mod=vendor
#COPY benft /go/bin/

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
EXPOSE 9998
CMD ["nft721"]
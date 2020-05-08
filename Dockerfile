# Build Geth in a stock Go builder container
FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers bash cmake g++

ADD . /PlatON-Go
RUN cd /PlatON-Go && make platon

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates libstdc++ gcc
COPY --from=builder /PlatON-Go/build/bin/platon /usr/local/bin/

EXPOSE 8545 8546 30303 30303/udp
ENTRYPOINT ["platon"]

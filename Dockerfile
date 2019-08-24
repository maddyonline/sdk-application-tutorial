FROM golang:1.12

RUN apt-get update && \
    apt install -y tmux

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY . .
RUN make install


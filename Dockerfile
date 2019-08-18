FROM golang:1.12

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY . .
RUN make install


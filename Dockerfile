FROM golang:1.14-alpine

WORKDIR /go/src/Goff
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT /go/bin/goff
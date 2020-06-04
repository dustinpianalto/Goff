FROM golang:1.14-alpine

WORKDIR /go/src/Goff
COPY ./go.mod .

RUN apk add --no-cache git

RUN go get -d -v ./...

COPY . .
RUN go install -v ./...

ENTRYPOINT /go/bin/goff

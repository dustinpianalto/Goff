FROM golang:1.14-alpine as dev

WORKDIR /go/src/Goff
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .
RUN go install github.com/dustinpianalto/goff

CMD [ "go", "run", "goff.go"]

from alpine

WORKDIR /bin

COPY --from=dev /go/bin/goff ./goff

CMD [ "goff" ]

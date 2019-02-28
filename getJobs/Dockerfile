FROM golang:1.11-alpine

RUN apk update && apk upgrade && apk add --no-cache bash git
RUN go get google.golang.org/grpc
RUN go get github.com/gocql/gocql
RUN go get github.com/thedevsaddam/renderer

ENV SOURCES /app/
COPY . ${SOURCES}

RUN cd ${SOURCES} && CGO_ENABLED=0 go build

WORKDIR ${SOURCES}
CMD ${SOURCES}app

EXPOSE 8080

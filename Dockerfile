FROM golang:1.11

WORKDIR /simplestorage
ENV PROD=DOCKER

ADD . .
ARG GO111MODULE=on
RUN go get -v
RUN go install simplestorage

EXPOSE 8080

ENTRYPOINT /go/bin/simplestorage


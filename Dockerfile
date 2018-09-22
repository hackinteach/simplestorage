FROM golang:1.11

WORKDIR /simplestorage
ENV PROD=DOCKER

ADD . .
#ARG GO111MODULE=on
#RUN go get -v
#RUN go install simplestorage

RUN go get github.com/codegangsta/gin
RUN go get
RUN go install

EXPOSE 3000
CMD gin run
#ENTRYPOINT /go/bin/simplestorage
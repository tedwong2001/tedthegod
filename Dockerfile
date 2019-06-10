FROM golang:latest

RUN mkdir -p /go/src/tedthegod

WORKDIR /go/src/tedthegod

COPY . /go/src/tedthegod

RUN go install tedthegod

CMD /go/bin/tedthegod

EXPOSE 8080
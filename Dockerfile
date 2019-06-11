FROM golang:latest

RUN mkdir -p /go/src/tedthegod

WORKDIR /go/src/tedthegod

COPY . /go/src/tedthegod

RUN go install tedthegod

RUN go get -u github.com/lib/pq

RUN go get -u github.com/gin-gonic/gin

CMD /go/bin/tedthegod

EXPOSE 8080

FROM golang

ADD . /code

WORKDIR /code

RUN go get github.com/tools/godep

RUN godep restore

RUN go get github.com/codegangsta/gin


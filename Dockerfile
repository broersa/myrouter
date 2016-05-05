# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ADD . /go/src/github.com/broersa/myrouter

RUN go get github.com/broersa/myrouter
RUN go install github.com/broersa/myrouter

ENTRYPOINT /go/bin/myrouter

EXPOSE 1700

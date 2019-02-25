FROM golang:1.9.7

MAINTAINER TripleZ "me@triplez.cn"

WORKDIR $GOPATH/src/a2os/safeu-backend
ADD . $GOPATH/src/a2os/safeu-backend/

# Add dependencies
RUN go get -u github.com/kardianos/govendor && go get github.com/pilu/fresh && govendor sync

# Build package
RUN go build .

EXPOSE 8080

ENTRYPOINT ["./safeu-backend"]

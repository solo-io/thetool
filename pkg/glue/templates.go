package glue

const (
	buildScript = `#!/bin/bash

set -ex

cd $GOPATH
mkdir -p -v src/github.com/solo-io
cd src/github.com/solo-io
ln -s /glue/external/glue .
cd glue
pwd

go get -u github.com/golang/dep/cmd/dep
dep ensure
GOOS=linux CGO_ENABLED=0 go build -o glue cmd/glue/main.go
cp glue /glue/glue
`

	dockerFile = `FROM alpine:3.5
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY glue .
CMD ["./glue"]
`
)

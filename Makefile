SOURCES := $(shell find . -name *.go)
BINARY:=thetool

build: $(BINARY)

$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-X main.Version=`git describe --tags`" -v -o $@ *.go

test:
	ginkgo -r -v .

clean:
	rm -f $(BINARY)

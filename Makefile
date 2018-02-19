SRC := vnc-screenshot.go
TARGET := vnc-screenshot
VERSION := $(shell git describe --tags)

all: $(TARGET)

$(TARGET): $(SRC)
	go build -ldflags "-X main.version=$(VERSION) -X main.builddate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -o $@

clean:
	$(RM) $(TARGET)

release:
	rm -rf dist/
	mkdir dist
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr
	CGO_ENABLED=0 gox -ldflags "-X main.version=$(VERSION) -X main.builddate=`date -u +%Y-%m-%dT%H:%M:%SZ`" -output "dist/vnc-screenshot_{{.OS}}_{{.Arch}}" -os="linux"
	ghr -u erasche -replace $(VERSION) dist/

.PHONY: all gofmt clean release

VERSION = 0.0.1
GO      = go
NAME    = curltab

# TODO: Add deps management, glide.sh
# TODO: Add golint
# TODO: Add tests
# IDEA: possibly set VERSION = $(shell git describe --abbrev=1 --tags --always)

build = GOOS=$(1) GOARCH=$(2) $(GO) build $(BUILDFLAGS) -ldflags="-X main.VERSION=$(VERSION)" -o build/$(NAME)-$(VERSION)-$(1)-$(2) main.go
tar   = cd build && tar -zcvf $(NAME)-$(VERSION)-$(1)-$(2).tar.gz $(NAME)-$(VERSION)-$(1)-$(2) && rm $(NAME)-$(VERSION)-$(1)-$(2)
zip   = cd build && zip $(NAME)-$(VERSION)-$(1)-$(2).zip $(NAME)-$(VERSION)-$(1)-$(2) && rm $(NAME)-$(VERSION)-$(1)-$(2)

.PHONY: all linux darwin clean dev-build
all: linux darwin

clean:
	rm -rf build/

dev-build:
	$(GO) build $(BUILDFLAGS) -o bin/$(NAME) main.go

linux: build/$(NAME)-$(VERSION)-linux-amd64.tar.gz

build/$(NAME)-$(VERSION)-linux-amd64.tar.gz:
	$(call build,linux,amd64)
	$(call tar,linux,amd64)

darwin: build/$(NAME)-$(VERSION)-darwin-amd64.zip

build/$(NAME)-$(VERSION)-darwin-amd64.zip:
	$(call build,darwin,amd64)
	$(call zip,darwin,amd64)

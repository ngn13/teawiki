RELEASE = 0
VERSION = $(shell git rev-parse --short HEAD)
GOSRCS  = $(wildcard *.go) $(wildcard */*.go)
HTML    = $(wildcard views/*.html views/*/*.html)
SCSS    = $(wildcard static/css/*.scss)
CSS     = $(patsubst %.scss,%.css,$(SCSS))

all: teawiki.elf

teawiki.elf: $(GOSRCS) $(HTML) $(CSS)
ifeq ($(RELEASE),1)
	CGO_ENABLED=0 go build \
		-ldflags "-X github.com/ngn13/teawiki/consts.VERSION=$(VERSION)" \
		-o $@
else
	go build -o $@
endif

%.css: %.scss
	sassc $^ $@

run: teawiki.elf
	TW_URL=http://127.0.0.1:8080 TW_REPO_PATH=. ./teawiki.elf

format:
	gofmt -s -w .

clean:
	rm -f static/css/*.css
	rm *.elf

test: teawiki.elf
	tests/run.sh

.PHONY: run format clean test

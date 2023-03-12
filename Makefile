APP_NAME=dalil

GO=go

all: build

build:
	@mkdir -p .build
	$(GO) build -o .build/dalil cmd/dalil/dalil.go

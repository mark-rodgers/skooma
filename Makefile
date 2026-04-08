.PHONY: build clean

OS := $(shell go env GOOS)
ifeq ($(OS),windows)
	OUTPUT := bin/skooma.exe
else
	OUTPUT := bin/skooma
endif

build:
	go build -o $(OUTPUT)

clean:
ifeq ($(OS),windows)
	@if exist bin rmdir /s /q bin
else
	rm -rf bin
endif

.DEFAULT_GOAL := build

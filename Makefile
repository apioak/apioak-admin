.PHONY: build clean

build:
	@echo "Building APIOAK"
	@mkdir -p bin
	@go build -o bin/apioak .

clean:
	@echo "Remove APIOAK"
	@rm -f bin/*

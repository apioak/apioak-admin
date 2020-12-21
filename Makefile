.PHONY: build clean

build:
	@echo "Building APIOAK"
	@mkdir -p bin etc
	@cp conf/apioak.yaml etc/apioak.yaml
	@go build -o bin/apioak .

clean:
	@echo "Remove APIOAK"
	@rm -rf bin etc

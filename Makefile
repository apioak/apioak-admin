.PHONY: build clean

build:
	@echo "Building APIOAK"
	@mkdir -p bin etc
	@cp conf/apioak-admin.yaml etc/apioak-admin.yaml
	@go build -o bin/apioak-admin .

clean:
	@echo "Remove APIOAK"
	@rm -rf bin etc

.PHONY: build
build:
	go build -o apioak-admin main.go

.PHONY: build-all
build-all:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o apioak-admin_linux_amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o apioak-admin_linux_arm64 main.go

.PHONY: run
run:
	@go run ./main.go

.PHONY: help
help:
	@echo "make build : 仅根据当前平台编辑"
	@echo "make build-all : 编辑 linux/amd64、linux/amd64"
	@echo "make run : 直接运行 Go 代码"

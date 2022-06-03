OS:=$(shell go env GOOS)

BINARY_DIR:=bin/PrometheusAlertFire
BINARY_NAME:=PrometheusAlertFire
BUILD_TIME:=$(shell date "+%Y%m%d%H%M") 


#mac
build-darwin:
	CGO_ENABLED=0 GOOS=darwin go build -o $(BINARY_DIR)/$(BINARY_NAME)-darwin
# windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-win.exe
# linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-linux

build:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)

# 全平台
build-all:
	make build-darwin
	make build-win
	make build-linux



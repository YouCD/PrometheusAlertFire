OS:=$(shell go env GOOS)

BINARY_DIR:=bin
BINARY_NAME:=PrometheusAlertFire
BUILD_TIME:=$(shell date "+%Y%m%d%H%M") 


#mac
build-darwin:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -o $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64/$(BINARY_NAME)
# windows
build-win:
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64/$(BINARY_NAME).exe
# linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64   go build -o $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64/$(BINARY_NAME)

build:
	@CGO_ENABLED=0 GOOS=$(OS) GOARCH=amd64   go build -o $(BINARY_DIR)/$(BINARY_NAME)

# 全平台
build-all:
	@make build-darwin
	@make build-win
	@make build-linux
	@cp config_exapmle.yaml $(BINARY_DIR)/$(BINARY_NAME)-darwin-amd64/config.yaml
	@cp config_exapmle.yaml $(BINARY_DIR)/$(BINARY_NAME)-windows-amd64/config.yaml
	@cp config_exapmle.yaml $(BINARY_DIR)/$(BINARY_NAME)-linux-amd64/config.yaml
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-darwin-amd64.txz $(BINARY_NAME)-darwin-amd64&&rm -rf $(BINARY_NAME)-darwin-amd64
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-windows-amd64.txz $(BINARY_NAME)-windows-amd64&&rm -rf $(BINARY_NAME)-windows-amd64
	@cd $(BINARY_DIR)&& tar Jcf $(BINARY_NAME)-linux-amd64.txz $(BINARY_NAME)-linux-amd64&&rm -rf $(BINARY_NAME)-linux-amd64



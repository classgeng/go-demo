.PHONY: build
build:
	@echo build local
	@go build

	@echo build linux_amd64_x86
	@GOOS=linux  GOARCH=amd64 CGO_ENABLED=0 go build -o storagesecurity_x86

	@echo build linux_arm64_x86
	@GOOS=linux  GOARCH=arm64 CGO_ENABLED=0 go build -o storagesecurity_arm
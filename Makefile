.PHONY: all build run gotool clean help

APPNAME=api

all: gotool build run

build:
	go build -o bin/${APPNAME} cmd/${APPNAME}/*.go

run:
	bin/${APPNAME}

gotool:
	go fmt ./...
	go vet ./...

clean:
	@if [ -f bin/${APPNAME} ] ; then rm bin/${APPNAME} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
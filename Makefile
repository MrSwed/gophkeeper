version?=1.1

SERVER_NAME=gophkeeper-server
CLIENT_NAME=gophkeeper
buildVersion=$(shell git log --pretty=format:"%h" -1)
buildCommit=$(shell git log --pretty=format:"%s (%ad)" --date=rfc2822 -1)
buildDate=$(shell date +'%Y-%m-%d %H:%M:%S')

.PHONY: dir clean_server clean_client clean build_server build_client_linux build_client_windows build_all run_server run_client test race install_go_cover_treemap cover proto gosec

dir:
	mkdir -p ./bin

clean_server: dir
	rm -f ./bin/${SERVER_NAME}

clean_client: dir
	rm -f ./bin/${CLIENT_NAME}

clean: clean_client clean_server

build_server: clean_server
	CGO_ENABLED=1 \
	go build -ldflags "-w -s -X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${SERVER_NAME}" ./cmd/server/*.go

build_client_linux: clean_client
	go build -ldflags "-w -s -X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${CLIENT_NAME}" ./cmd/client/*.go

# sudo apt-get install gcc-mingw-w64-i686 and sudo apt-get install gcc-mingw-w64-x86-64
build_client_windows:
	GOOS=windows \
	CGO_ENABLED=1 \
	CC="i686-w64-mingw32-gcc" \
	GOARCH=386 \
	go build -ldflags "-w -s -X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${CLIENT_NAME}.exe" ./cmd/client/*.go
	upx -1 "./bin/${CLIENT_NAME}.exe"

build_all: build_server build_client_linux build_client_windows

run_server: build_server
	./bin/${SERVER_NAME}

run_client: build_client_linux
	./bin/${CLIENT_NAME}

test:
	go test -v -count=1 ./...

race:
	go test -v -race -count=1 ./...

install_go_cover_treemap:
	command -v go-cover-treemap | grep go-cover-treemap > /dev/null || go install github.com/nikolaydubina/go-cover-treemap@latest

cover:
	go test -v -count=1 -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out
	@command -v go-cover-treemap | grep go-cover-treemap > /dev/null && \
	( echo "go-cover-treemap -coverprofile coverage.out > coverage.out.svg";\
		go-cover-treemap -coverprofile coverage.out > coverage.out.svg ) || \
		( echo "coverage.out.svg is not created, please install go-cover-treemap by command" ;\
		echo "  make install_go_cover_treemap" )

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      internal/grpc/proto/service.proto

gosec:
	@command -v go-cover-treemap | grep go-cover-treemap > /dev/null || \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	gosec ./...


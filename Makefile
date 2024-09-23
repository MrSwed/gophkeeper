version?=1.1

SERVER_NAME=server
CLIENT_NAME=client
buildVersion=$(shell git log --pretty=format:"%h" -1)
buildCommit=$(shell git log --pretty=format:"%s (%ad)" --date=rfc2822 -1)
buildDate=$(shell date +'%Y-%m-%d %H:%M:%S')

dir:
	mkdir -p ./bin

clean_server: dir
	rm -f ./bin/${SERVER_NAME}
clean_client: dir
	rm -f ./bin/${CLIENT_NAME}

clean: clean_client clean_server

build_server: clean_server
	CGO_ENABLED=1 \
    go build -ldflags "-X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${SERVER_NAME}" ./cmd/${SERVER_NAME}/*.go

build_client: clean_client
	go build -ldflags "-X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${CLIENT_NAME}" ./cmd/${CLIENT_NAME}/*.go

build_all: build_server build_client

run_app: build_server
	./bin/${SERVER_NAME}

run_client: build_client
	./bin/${CLIENT_NAME}

test:
	go test -v -count=1 ./...

race:
	go test -v -race -count=1 ./...

install_go_cover_treemap:
	go install github.com/nikolaydubina/go-cover-treemap@latest

run_go_cover_treemap:
	go-cover-treemap -coverprofile coverage.out > coverage.out.svg


.PHONY: cover
cover: install_go_cover_treemap
	go test -v -count=1 -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out
	go-cover-treemap -coverprofile coverage.out > coverage.out.svg

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
      --go-grpc_out=. --go-grpc_opt=paths=source_relative \
      internal/grpc/proto/service.proto


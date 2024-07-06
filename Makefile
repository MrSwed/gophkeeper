version?=1.0

APP_NAME=server
AGENT_NAME=agent
buildVersion=$(shell git log --pretty=format:"%h" -1)
buildCommit=$(shell git log --pretty=format:"%s (%ad)" --date=rfc2822 -1)
buildDate=$(shell date +'%Y-%m-%d %H:%M:%S')

dir:
	mkdir -p ./bin

clean_app: dir
	rm -f ./bin/${APP_NAME}
clean_agent: dir
	rm -f ./bin/${AGENT_NAME}

clean: clean_agent clean_app

build_app: clean_app
	CGO_ENABLED=1 \
    go build -ldflags "-X 'main.buildVersion=${version} (${buildVersion})' -X 'main.buildDate=${buildDate}' -X 'main.buildCommit=${buildCommit}'" -o "./bin/${APP_NAME}" ./cmd/${APP_NAME}/*.go

build_agent: clean_agent
	go build -ldflags "-X 'app.buildVersion=${version} (${buildVersion})' -X 'app.buildDate=${buildDate}' -X 'app.buildCommit=${buildCommit}'" -o "./bin/${AGENT_NAME}" ./cmd/${AGENT_NAME}/*.go

build_all: build_app build_agent

run_app: build_app
	./bin/${APP_NAME}

run_agent: build_agent
	./bin/${AGENT_NAME}

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


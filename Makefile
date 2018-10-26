

.PHONY: kpay-grpc kauth-grpc kauth-rest kauth-docs install-server
all: kpay-grpc kauth-grpc kauth-rest kauth-docs install-server

kpay-grpc: ## generate kpay grpc output from proto
	@protoc -I kpay \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:kpay \
	kpay/kpay.proto

kauth-grpc: ## generate kauth grpc output
	@protoc -I kauth \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:kauth \
	kauth/kauth.proto

kauth-rest: ## generate kauth rest output
	@protoc  -I kauth \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:kauth \
	kauth/kauth.proto

kauth-docs: ## generate kauth doc files
    @protoc -I kauth \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--swagger_out=logtostderr=true:kauth \
	kauth/kauth.proto

install-server: ## install server
	@go install ./...

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

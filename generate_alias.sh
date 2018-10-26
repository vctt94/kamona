#!/bin/bash


# generate alias for kauth
alias kauth-grpc="protoc -I kauth -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:kauth kauth/kauth.proto"
alias kauth-rest="protoc  -I kauth  -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --grpc-gateway_out=logtostderr=true:kauth kauth/kauth.proto"
alias kauth-rest-docs="protoc  -I kauth  -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis  --swagger_out=logtostderr=true:katuh kauth/kauth.proto"

# # generate alias for kpay
alias kpay-grpc="protoc -I kpay -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:kpay kpay/kpay.proto"

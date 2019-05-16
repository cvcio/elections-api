REGISTRY=reg.plagiari.sm
PROJECT=elections-api
ORG=cvcio
TAG=`cat $(GOPATH)/src/github.com/$(ORG)/$(PROJECT)/VERSION`

keys:
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048

tools:
	go get github.com/oxequa/realize
	go get github.com/golangci/golangci-lint

init:
	cp .realize.yaml.template .realize.yaml
	mkdir pkg/proto && protocols

protocols:
	mkdir -p pkg/proto && \
		protoc -I/usr/local/include -I. -I${GOPATH}/src --go_out=plugins=grpc:./pkg proto/twitter.proto && \
		protoc -I/usr/local/include -I. -I${GOPATH}/src --go_out=plugins=grpc:./pkg proto/classification.proto

dev:
	realize start

lint:
	golangci-lint run -e vendor

# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: release
release: custom 

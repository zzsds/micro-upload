GOPATH:=$(shell go env GOPATH)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --tags --always --match "v*")
# IMAGE_TAG=$(GIT_TAG)-$(GIT_COMMIT)
IMAGE_TAG = 0.0.1
NAME = upload-api-http
IMAGE_NAME = micro-welfare/${NAME}
NETWORK = micro-welfare
PROTO = upload
CFG_CLUSTER=prod
all: build

proto:
	@echo execute ${PROTO} proto file generate
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/${PROTO}/${PROTO}.proto
	
vendor:
	go mod tidy
	go mod vendor

build:
	go build -o ${NAME} *.go

test:
	go test -v ./... -cover

docker: vendor docker-build docker-run

docker-build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest
	# docker push $(IMAGE_NAME):$(IMAGE_TAG)
	# docker push $(IMAGE_NAME):latest

docker-run:
	docker run --rm --name ${NAME} -d -p :50051 --network ${NETWORK} -e CFG_CLUSTER=${CFG_CLUSTER} -e MICRO_ADDRESS=:50051 -e MICRO_REGISTRY=mdns ${IMAGE_NAME}:${IMAGE_TAG}

.PHONY: vendor build proto clean vet test docker-build docker-run
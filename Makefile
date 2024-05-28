SERVICE ?= headnexus
REGISTRY ?= hub.hitosea.com/cossim
TAG ?= latest
IMG ?= ${REGISTRY}/${SERVICE}:${TAG}
DOCKER_BUILD_PATH ?= "cmd/main.go"
PLATFORMS ?= linux/amd64

GOPROXY ?= https://goproxy.cn
# 防止命令行参数被误认为是目标
%:
	@:

.PHONY: dep
dep: ## Get the dependencies
	@go mod tidy

.PHONY: lint
lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: test
test: fmt vet## Run unittests
	@go test -short ./...

.PHONY: gen
gen:
	@#oapi-codegen -generate gin -o ports/openapi.gen.go -package ports api/openapi/*.yaml
	@oapi-codegen -generate spec,gin -o api/http/v1/openapi.gen.go -package v1 api/openapi/*.yaml
	@oapi-codegen -generate types -o api/http/v1//openapi.types.go -package v1 api/openapi/*.yaml

# If you wish built the manager image targeting other platforms you can use the --platform flag.
# (i.e. docker build --platform linux/arm64 ). However, you must enable docker buildKit for it.
# More info: https://docs.docker.com/develop/develop-images/build_enhancements/
docker-build: dep test## Build docker image with the manager.
	docker build --platform $(PLATFORMS)  \
			--build-arg BUILD_PATH=${DOCKER_BUILD_PATH} \
			--build-arg GOPROXY="${GOPROXY}" \
			-t "${IMG}" .

docker-push:
	docker push ${IMG}
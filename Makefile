OUTPUT=aws-iam-access-key-disabler
GO_VERSION=golang:1.13-stretch
GO_FLAGS=-a -ldflags='-s -w' -installsuffix cgo
ARTEFACTS_DIR=artefacts

.DEFAULT_GOAL=build

.PHONY: staticcheck
staticcheck:
	docker run --rm -v "$(PWD)":/code -it $(GO_VERSION) sh -c "\
		go get honnef.co/go/tools/cmd/staticcheck > /dev/null && \
		cd /code && \
		echo \"\n** Running static analyze\n\" && \
		staticcheck *.go \
	"

.PHONY: vet
vet:
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go -it $(GO_VERSION) sh -c "\
		cd /code && \
		echo \"\n** Running vet\n\" && \
		go vet \
	"

.PHONY: build
build: vet staticcheck
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go -it $(GO_VERSION) sh -c "\
		apt-get update && apt-get install -y upx && \
		cd /code && \
		echo \"\n** Building ${OUTPUT}\n\" && \
		make all \
	"

.PHONY: install
install: vet staticcheck
	go install $(GO_FLAGS)

.PHONY: clean
clean:
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go -it $(GO_VERSION) sh -c "\
		rm -rf /go/* /code/artefacts/* \
	"

all: linux64 darwin64 upx
linux64: $(ARTEFACTS_DIR)-$(OUTPUT)-linux-amd64
darwin64: $(ARTEFACTS_DIR)-$(OUTPUT)-darwin-amd64

$(ARTEFACTS_DIR)-$(OUTPUT)-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o $(ARTEFACTS_DIR)/$(OUTPUT)-linux-amd64

$(ARTEFACTS_DIR)-$(OUTPUT)-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) -o $(ARTEFACTS_DIR)/$(OUTPUT)-darwin-amd64

.PHONY: upx
upx:
	upx -q ${ARTEFACTS_DIR}/$(OUTPUT)-linux-amd64

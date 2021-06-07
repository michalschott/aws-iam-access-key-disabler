OUTPUT=aws-iam-access-key-disabler
OUTPUT_ZIP=aws-iam-access-key-disabler.zip
GO_VERSION=golang:1.13-stretch
GO_FLAGS=-a -ldflags='-s -w' -installsuffix cgo
ARTEFACTS_DIR=.

.DEFAULT_GOAL=build

.PHONY: staticcheck
staticcheck:
	docker run --rm -v "$(PWD)":/code $(GO_VERSION) sh -c "\
		go get honnef.co/go/tools/cmd/staticcheck > /dev/null && \
		cd /code && \
		echo \"\n** Running static analyze\n\" && \
		staticcheck *.go \
	"

.PHONY: vet
vet:
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go $(GO_VERSION) sh -c "\
		cd /code && \
		echo \"\n** Running vet\n\" && \
		go vet \
	"

.PHONY: build
build: vet staticcheck
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go $(GO_VERSION) sh -c "\
		apt-get update && apt-get upgrade -y && apt-get install -y zip && \
		cd /code && \
		echo \"\n** Building ${OUTPUT}\n\" && \
		make all \
	"

.PHONY: install
install: vet staticcheck
	go install $(GO_FLAGS)

.PHONY: clean
clean:
	docker run --rm -v "$(PWD)":/code -v "$(PWD)"/.cache:/go $(GO_VERSION) sh -c "\
		rm -rf /go/* \
	"

all: linux64
linux64: $(ARTEFACTS_DIR)-$(OUTPUT) $(ARTEFACTS_DIR)-$(OUTPUT_ZIP)

$(ARTEFACTS_DIR)-$(OUTPUT):
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o $(ARTEFACTS_DIR)/$(OUTPUT)

$(ARTEFACTS_DIR)-$(OUTPUT_ZIP):
	zip $(ARTEFACTS_DIR)/$(OUTPUT_ZIP) $(ARTEFACTS_DIR)/$(OUTPUT)

# Jupiter Golang Application Standard Makefile
SHELL:=/bin/bash
BASE_PATH:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
PATH_CONV:=$(shell command -v cygpath 2> /dev/null)
GO_ROOT:=$(shell go env GOROOT)
ifdef PATH_CONV
    BASE_PATH:=$(shell $(PATH_CONV) -u "$(BASE_PATH)")
    GO_ROOT:=$(shell $(PATH_CONV) -u "$(GO_ROOT)")
endif

DIST_PATH:=$(BASE_PATH)/build/dist
CFGS_PATH:=$(BASE_PATH)/build/config
SCRIPT_PATH:=$(BASE_PATH)/build/script
BUILD_TIME:=$(shell date +%Y-%m-%d--%T)

PACKAGE:=$(shell go list -m)
GO_FILES:=`find . -name "*.go" -type f -not -path "./vendor/*"`
CUE_FILES:=`find "$(CFGS_PATH)/config/" -name "*.cue"`
CUE_FILES+="${BASE_PATH}/common/config/template.cue"
TESTCOVER_FILES=`go list ./... | grep -Ev "${PACKAGE}/(tests|tool|package)"`

LOCALPKGS:=github.com/kentalee
export APP_NAME=ddns

all:clean print fmt lint build

clean:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@rm $(BASE_PATH)/Dockerfile ||:
	@rm -r $(BASE_PATH)/build/dist/* ||:
	@rm -r $(BASE_PATH)/xls_output ||:

print:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making print<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@echo SHELL:$(SHELL)
	@echo BASE_PATH:$(BASE_PATH)
	@echo SCRIPT_PATH:$(SCRIPT_PATH)
	@echo BUILD_TIME:$(BUILD_TIME)
	@echo -e "\n"

fmt:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making fmt<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	cue fmt $(CUE_FILES)
	gofmt -e -w -s ${GO_FILES}
	goimports -local ${LOCALPKGS} -l -e -w ${GO_FILES}
	@echo -e "\n"

lint:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making lint<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	cue vet -E -c --strict -p config --list $(CUE_FILES) "$(DIST_PATH)/.defines.cue"
	golangci-lint run -v
	@echo -e "\n"

testCover:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making testCover<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	$(eval coverageFile := $(shell mktemp))
	go test ${TESTCOVER_FILES} -v -coverprofile $(coverageFile) && go tool cover --func=$(coverageFile)
	$(shell rm $(coverageFile))

build:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making build<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	bash $(SCRIPT_PATH)/build.sh $(DIST_PATH)/$(APP_NAME) package/version $(BASE_PATH)/cmd/$(APP_NAME)
	@echo -e "\n"

exportDefines:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>export defines<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	SYS_ENV="testing" SYS_NAME="hkg" SYS_NODE="local01" go run $(BASE_PATH)/cmd/$(APP_NAME) -cDef "$(DIST_PATH)/.defines.cue"

.PHONY: build
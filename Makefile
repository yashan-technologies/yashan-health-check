# env defines
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
ARCH_AMD=x86_64
ARCH_ARM=aarch64
OS=$(shell if [ $(GOOS)a != ""a ]; then echo $(GOOS); else echo "linux"; fi)
ARCH=$(shell if [ $(GOARCH)a == "arm64"a ]; then echo $(ARCH_ARM); else echo $(ARCH_AMD); fi)
VERSION=$(shell cat ./VERSION)
GO_VERSION=$(shell go env GOVERSION)
GIT_COMMIT_ID=$(shell git rev-parse HEAD)
GIT_DESCRIBE=$(shell git describe --always)

# go command defines
GO_BUILD=go build
GO_MOD_TIDY=$(go mod tidy -compat 1.19)
GO_BUILD_WITH_INFO=$(GO_BUILD) -ldflags "\
	-X 'yhc/defs/compiledef._appVersion=$(VERSION)' \
	-X 'yhc/defs/compiledef._goVersion=$(GO_VERSION)'\
	-X 'yhc/defs/compiledef._gitCommitID=$(GIT_COMMIT_ID)'\
	-X 'yhc/defs/compiledef._gitDescribe=$(GIT_DESCRIBE)'"

# package defines
PKG_PERFIX=yashan-health-check
PKG=$(PKG_PERFIX)-$(VERSION)-$(OS)-$(ARCH).tar.gz

BUILD_PATH=./build
PKG_PATH=$(BUILD_PATH)/$(PKG_PERFIX)
BIN_PATH=$(PKG_PATH)/bin
LOG_PATH=$(PKG_PATH)/log
DOCS_PATH=$(PKG_PATH)/docs
RESULTS_PATH=$(PKG_PATH)/results

# build defines
BIN_YHCCTL=$(BUILD_PATH)/yhcctl
BIN_FILES=$(BIN_YHCCTL)

DIR_TO_MAKE=$(BIN_PATH) $(LOG_PATH) $(RESULTS_PATH) $(DOCS_PATH)
FILE_TO_COPY=./config ./scripts ./static

# functions
clean:
	rm -rf $(BUILD_PATH)

define build_yhcctl
	$(GO_BUILD_WITH_INFO) -o $(BIN_YHCCTL) ./cmd/yhcctl/*.go
endef

go_build: 
	$(GO_MOD_TIDY)
	$(call build_yhcctl)

build: go_build
	@mkdir -p $(DIR_TO_MAKE) 
	@cp -r $(FILE_TO_COPY) $(PKG_PATH)
	# @cp -r ./yhc-doc $(DOCS_PATH)/markdown
	# @cp ./yhc.pdf $(DOCS_PATH)
	@mv $(BIN_FILES) $(BIN_PATH)
	@> $(LOG_PATH)/yhcctl.log
	@cd $(PKG_PATH);ln -s ./bin/yhcctl ./yhcctl
	@cd $(BUILD_PATH);tar -cvzf $(PKG) $(PKG_PERFIX)/

force: clean build
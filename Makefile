# env defines
GOOS=$(shell go env GOOS)
ARCH=$(shell arch)
VERSION=$(shell cat ./VERSION)
GO_VERSION=$(shell go env GOVERSION)
GIT_COMMIT_ID=$(shell git rev-parse HEAD)
GIT_DESCRIBE=$(shell git describe --always)
OS=$(if $(GOOS),$(GOOS),linux)

# go command defines
GO_BUILD=go build
YARN_REPLACE_SOURCE=yarn config set registry https://registry.npmmirror.com
YARN_INSTALL=yarn install
YARN_BUILD=yarn build
GO_MOD_TIDY=$(go mod tidy -compat 1.19)
GO_BUILD_WITH_INFO=$(GO_BUILD) -ldflags "\
	-X 'yhc/defs/compiledef._appVersion=$(VERSION)' \
	-X 'yhc/defs/compiledef._goVersion=$(GO_VERSION)'\
	-X 'yhc/defs/compiledef._gitCommitID=$(GIT_COMMIT_ID)'\
	-X 'yhc/defs/compiledef._gitDescribe=$(GIT_DESCRIBE)'"

# package defines
PKG_PERFIX=yashan-health-check-$(VERSION)
PKG=$(PKG_PERFIX)-$(OS)-$(ARCH).tar.gz

BUILD_PATH=./build
PKG_PATH=$(BUILD_PATH)/$(PKG_PERFIX)
BIN_PATH=$(PKG_PATH)/bin
LOG_PATH=$(PKG_PATH)/log
DOCS_PATH=$(PKG_PATH)/docs
RESULTS_PATH=$(PKG_PATH)/results
HTML_PATH=$(PKG_PATH)/html-template


TEMPLATE_PATH=./html-template
TEMPLATE_BUILD_PATH=$(TEMPLATE_PATH)/dist

# build defines
BIN_YHCCTL=$(BUILD_PATH)/yhcctl
BIN_FILES=$(BIN_YHCCTL)

SCRIPTS_PATH=$(PKG_PATH)/scripts
SCRIPTS_YASDB_GO=$(BUILD_PATH)/yasdb-go
SCRIPTS_FILES=$(SCRIPTS_YASDB_GO)

DIR_TO_MAKE=$(BIN_PATH) $(LOG_PATH) $(RESULTS_PATH) $(DOCS_PATH) $(HTML_PATH)
FILE_TO_COPY=./config ./scripts

WORD_GENNER_PATH=./wordgenner
WORD_GENNER_DIST=$(WORD_GENNER_PATH)/dist/wordgenner

.PHONY: clean force go_build

build: pre_build go_build
	@cp ./template.html $(HTML_PATH)/
	@cp -r ./yhc-doc $(DOCS_PATH)/markdown
	@cp ./yhc.pdf $(DOCS_PATH)
	@mv $(BIN_FILES) $(BIN_PATH)
	@mv $(SCRIPTS_FILES) $(SCRIPTS_PATH)
	@> $(LOG_PATH)/yhcctl.log
	@> $(LOG_PATH)/console.out
	@cd $(PKG_PATH);ln -s ./bin/yhcctl ./yhcctl
	@cd $(BUILD_PATH);tar -cvzf $(PKG) $(PKG_PERFIX)/

clean:
	rm -rf $(BUILD_PATH)
	@cd $(WORD_GENNER_PATH);make clean


go_build: 
	$(GO_MOD_TIDY)
	$(GO_BUILD_WITH_INFO) -o $(BIN_YHCCTL) ./cmd/yhcctl/*.go
	$(GO_BUILD_WITH_INFO) -o $(SCRIPTS_YASDB_GO) ./cmd/yasdb-go/*.go

build_template:
	@cd $(TEMPLATE_PATH);$(YARN_REPLACE_SOURCE);$(YARN_INSTALL);$(YARN_BUILD)
	@cp $(TEMPLATE_BUILD_PATH)/index.html ./template.html

build_wordgenner:
	@cd $(WORD_GENNER_PATH);make build
	@cp -r $(WORD_GENNER_DIST) scripts/

pre_build:
	@mkdir -p $(DIR_TO_MAKE) 
	@cp -r $(FILE_TO_COPY) $(PKG_PATH)

force: clean build
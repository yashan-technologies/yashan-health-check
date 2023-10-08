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
PKG_PERFIX=yashan-health-check
PKG=$(PKG_PERFIX)-$(VERSION)-$(OS)-$(ARCH).tar.gz

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
FILE_TO_COPY=./config ./scripts ./static


PIP_INSTALL=pip3 install -i https://mirrors.aliyun.com/pypi/simple/
PYINSTALLER=pyinstaller
WORD_GENNER_PATH=./wordgenner
WORD_GENNER_DIST=$(WORD_GENNER_PATH)/dist/wordgenner
WORD_GENNER_INSTALL=$(PIP_INSTALL) -r requirements.txt
PYTHON_DOCX_PATH=$(WORD_GENNER_PATH)/python-docx
PYTHON_DOCX_INSTALL=$(PIP_INSTALL) .
WORD_GENNER_BUILD=$(PYINSTALLER) main.py --name wordgenner

.PHONY: clean force go_build

build: pre_build go_build wordgenner_build
	@cp ./template.html $(HTML_PATH)/
	# @cp -r ./yhc-doc $(DOCS_PATH)/markdown
	# @cp ./yhc.pdf $(DOCS_PATH)
	@mv $(BIN_FILES) $(BIN_PATH)
	@mv $(SCRIPTS_FILES) $(SCRIPTS_PATH)
	@> $(LOG_PATH)/yhcctl.log
	@> $(LOG_PATH)/console.out
	@cd $(PKG_PATH);ln -s ./bin/yhcctl ./yhcctl
	@cd $(BUILD_PATH);tar -cvzf $(PKG) $(PKG_PERFIX)/

clean:
	rm -rf $(BUILD_PATH)
	rm -rf $(WORD_GENNER_DIST)

go_build: 
	$(GO_MOD_TIDY)
	$(GO_BUILD_WITH_INFO) -o $(BIN_YHCCTL) ./cmd/yhcctl/*.go
	$(GO_BUILD_WITH_INFO) -o $(SCRIPTS_YASDB_GO) ./cmd/yasdb-go/*.go

build_template:
	@cd $(TEMPLATE_PATH);$(YARN_REPLACE_SOURCE);$(YARN_INSTALL);$(YARN_BUILD)
	@cp $(TEMPLATE_BUILD_PATH)/index.html ./template.html

wordgenner_build:
	@cd $(PYTHON_DOCX_PATH);$(PYTHON_DOCX_INSTALL)
	@cd $(WORD_GENNER_PATH);$(WORD_GENNER_INSTALL)
	@cd $(WORD_GENNER_PATH);$(WORD_GENNER_BUILD)
	@cd $(WORD_GENNER_DIST);mkdir -p docx/parts/
	@cp -r $(WORD_GENNER_DIST) $(SCRIPTS_PATH)

pre_build:
	@mkdir -p $(DIR_TO_MAKE) 
	@cp -r $(FILE_TO_COPY) $(PKG_PATH)

force: clean build
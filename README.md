# Yashan Health Check

## Summary

*Yashan Health Check* is a lightweight database health check tool, which can perform in-depth checks on Yashan database and the server where it is located, and give professional suggestions based on the actual situation.

## Features

1. Database health checks
2. The actual alerting policy
3. Professional handling advice
4. Custom checks
5. Reports in multiple formats

## Uasge

// TODO: yhc-pdf

## For Development

### Dependencies

#### Development And Compilation Tools

| TOOL   | VERSION |
| ------ | ------- |
| go     | go 1.19 |
| python | python3 |
| make   | v3+     |
| gcc    | v7+     |

##### go

```bash
# example for x86_64
cd ~
# x86_64
wget https://dl.google.com/go/go1.19.7.linux-amd64.tar.gz

# aarch64 
# wget https://dl.google.com/go/go1.19.7.linux-arm64.tar.gz

tar -xzvf go1.19.7.linux-amd64.tar.gz
mkdir golang
mv go golang/

vim ~/.golang.env
<insert>
export GOROOT=${HOME}/golang/go
export GOPATH=${HOME}/golang/src
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
export PATH=$PATH:${GOROOT}/bin:${GOPATH}/bin
<wq>

vim ~/.bashrc
<insert>
source ~/.golang.env
<wq>

source ~/.bashrc
```

##### python

```bash
# example for centos 7
yum install -y python3
```

##### make

```bash
# example for centos 7 x86_64
yum install -y make.x86_64
```

##### gcc

```bash
# example for centos 7 x86_64
yum install -y centos-release-scl
yum install -y devtoolset-7-gcc.x86_64
yum install -y devtoolset-7-gcc-c++.x86_64

# set env
vim ~/.bashrc
<insert>
source /opt/rh/devtoolset-7/enable
<wq>
```

#### Code Checking Tools(Optional)

| TOOL              | VERSION |
| ----------------- | ------- |
| golangci-lint     | v1.53.2 |
| goimports-reviser | v3.4.1  |
| yapf              | 0.32.0  |
| mypy              | 0.950   |
| shellcheck        | 0.3.8   |

##### golangci-lint

```bash

# use script for installation
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh |\
sh -s -- -b $(go env GOPATH)/bin v1.53.2

# or use package for installation
# Releases: https://github.com/golangci/golangci-lint/releases/tag/v1.53.2
# Available Resource: wget http://192.168.8.236:8888/golang/golangci-lint-1.53.2-linux-amd64.tar.gz

# NOTE: such go install/go get installation aren't guaranteed to work. We recommend using binary installation.
```

##### goimports-reviser

```bash
go install -v github.com/incu6us/goimports-reviser/v3@v3.4.1
```

##### yapf

```bash
pip3 install yapf==0.32.0
```

##### mypy

```bash
pip3 install mypy==0.950
```

##### shellcheck

```bash
# example for centos 7
yum install -y epel-release
yum install -y ShellCheck-0.3.8-1.el7
```

#### Unit Test Tools(Optional)

| TOOL                | VERSION |
| ------------------- | ------- |
| gotestsum           | v1.10.1 |
| go-test-html-report | v1.1.0  |
| junit2html          | 30.1.3  |

##### gotestsum

```bash
go install gotest.tools/gotestsum@v1.10.1
```

##### go-test-html-report

```bash
go install github.com/Thatooine/go-test-html-report@v1.1.0
```

##### junit2html

```bash
pip3 install junit2html==30.1.3
```

### Configurations

#### Private Repository Configuration

#### go env

```bash
go env -w GOPRIVATE=git.yasdb.com
```

#### git config

```bash
vim ~/.gitconfig
[url "git@git.yasdb.com:"]
    insteadOf = https://git.yasdb.com/
```

### Compilation Script

There's a `build.py` file located at the root directory of the project.

```bash
# clean compilation data
python3 build.py clean

# start compilation
python3 build.py build
# --skip-check  build without checking code (default: False)
# --skip-test   build without running unit test (default: False)
# -c, --clean   clean before building (default: False)
# -f, --force   clean before building, then build without checking code and running unit test (default: False)

# check code
python3 build.py check

# run unit test
python3 build.py test
```

### Quick Start

```bash
cd yashan-trace-collector
export YHC_HOME=$(pwd)
export YHC_DEBUG_MODE=true

go run cmd/yhcctl/*.go

go run cmd/yhcd/*.go
```

PKG = github.com/luopengift/version
MAIN = *.go

APP = $(shell basename `pwd`)
APPVERSION = 1.0.0
GOVERSION = $(shell go version)
TIME = $(shell date "+%F %T")
GIT = $(shell git rev-parse HEAD)

FLAG = "-X '${PKG}.APP=${APP}' -X '${PKG}.APPVERSION=${APPVERSION}' -X '${PKG}.GOVERSION=${GOVERSION}' -X '${PKG}.GIT=${GIT}' -X '${PKG}.TIME=${TIME}'"

DEPLOY_HOST = ""

build: 
	go build -ldflags $(FLAG) -o ${APP} ${MAIN}
deploy: build
	scp -i ~/.ssh/id_rsa -P 22 ${APP} root@${DEPLOY_HOST}:/tmp
debug:
	go run *.go -debug
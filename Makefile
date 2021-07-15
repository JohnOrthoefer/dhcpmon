SRC=check.go httpServer.go lookupEnv.go main.go parseLeases.go startcmd.go config.go macdb.go
GOLANG=/usr/local/go/bin/go
CURL=/usr/bin/curl
GIT=/usr/bin/git
REPONAME=$(shell basename `git rev-parse --show-toplevel`)
SHA1=$(shell git rev-parse --short HEAD)
NOW=$(shell date +%Y-%m-%d_%T)
export HTTP_PROXY=http://llproxy.llan.ll.mit.edu:8080
export HTTPS_PROXY=${HTTP_PROXY}

dhcpmon: ${SRC}
	echo ${REPONAME}
	${GOLANG} build -ldflags "-X main.sha1ver=${SHA1} -X main.buildTime=${NOW} -X main.repoName=${REPONAME}"

update: update-go update-json

update-go:
	${GOLANG} get github.com/fsnotify/fsnotify
	${GOLANG} get gopkg.in/ini.v1

update-json:
	${CURL} -O https://macaddress.io/database/macaddress.io-db.json

clean:
	rm dhcpmon dhcpmon.tar.gz

archive:
	${GIT} archive --format=tar.gz  --prefix=dhcpmon/ HEAD > dhcpmon.tar.gz

BUILD_ORG   := ca17
BUILD_VERSION   := latest
BUILD_TIME      := $(shell date "+%F %T")
BUILD_NAME      := teamsedge
RELEASE_VERSION := v1.0.1
SOURCE          := main.go
RELEASE_DIR     := ./release
COMMIT_SHA1     := $(shell git show -s --format=%H )
COMMIT_DATE     := $(shell git show -s --format=%cD )
COMMIT_USER     := $(shell git show -s --format=%ce )
COMMIT_SUBJECT     := $(shell git show -s --format=%s )

clean:
	rm -f teamsedge

gen:
	go generate

build:
	go generate
	CGO_ENABLED=0 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static"\
	' \
    -o ${BUILD_NAME} ${SOURCE}

build-linux:
	go generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \
	'\
	-X "main.BuildVersion=${BUILD_VERSION}"\
	-X "main.ReleaseVersion=${RELEASE_VERSION}"\
	-X "main.BuildTime=${BUILD_TIME}"\
	-X "main.BuildName=${BUILD_NAME}"\
	-X "main.CommitID=${COMMIT_SHA1}"\
	-X "main.CommitDate=${COMMIT_DATE}"\
	-X "main.CommitUser=${COMMIT_USER}"\
	-X "main.CommitSubject=${COMMIT_SUBJECT}"\
	-s -w -extldflags "-static"\
	' \
    -o ${RELEASE_DIR}/${BUILD_NAME} ${SOURCE}

pubbuild-pre:
	make build-linux
	make upx
	echo 'FROM python:3.9.6-alpine3.14' > .build
	echo 'RUN apk add --no-cache tzdata' >> .build
	echo 'RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime' >> .build
	echo 'RUN apk add --no-cache curl' >> .build
	echo 'RUN pip install pysocks requests' >> .build
	echo 'ARG CACHEBUST="$(shell date "+%F %T")"' >> .build
	echo 'COPY ./teamsedge /teamsedge' >> .build
	echo 'RUN chmod +x /teamsedge' >> .build
	echo 'ENTRYPOINT ["/teamsedge"]' >> .build


pubdev:
	make pubbuild-pre
	docker build -t teamsedge . -f .teamsedgebuild \
	&& sudo docker tag teamsedge alab.189csp.cn:5000/teamsedge-hy:dev \
	&& sudo docker push alab.189csp.cn:5000/teamsedge-hy:dev \
	&& rm -f /tmp/teamsedge \
	&& rm -f /tmp/.teamsedgebuild "
	rm -f .build

fastpub:
	make pubbuild-pre
	ssh DockerServer "cd /tmp \
	&& sudo docker build -t teamsedge . -f .teamsedgebuild \
	&& sudo docker tag teamsedge alab.189csp.cn:5000/teamsedge-hy:latest \
	&& sudo docker push alab.189csp.cn:5000/teamsedge-hy:latest \
	&& rm -f /tmp/teamsedge \
	&& rm -f /tmp/.teamsedgebuild "
	rm -f .build


upx:
	upx ${RELEASE_DIR}/${BUILD_NAME}

ci:
	@read -p "type commit message: " cimsg; \
	git ci -am "$(shell date "+%F %T") $${cimsg}"

syncwjt:
	@read -p "提示:同步操作尽量在完成一个完整功能特性后进行，请输入提交描述 (wjt):  " cimsg; \
	git commit -am "$(shell date "+%F %T") : $${cimsg}" || echo "no commit"
	# 切换主分支并更新
	git checkout develop
	git pull origin develop
	# 切换开发分支变基合并提交
	git checkout wjt
	git rebase -i develop
	# 切换回主分支并合并开发者分支，推送主分支到远程，方便其他开发者合并
	git checkout develop
	git merge --no-ff wjt
	git push origin develop
	# 切换回自己的开发分支继续工作
	git checkout wjt

gitlog:
	git log --oneline

.PHONY: clean build



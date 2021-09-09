FROM golang:1.16.3-buster AS build

WORKDIR $GOPATH/src

RUN mkdir -p /release && mkdir -p $WORKDIR/teamsedge

ENV GO111MODULE=on

COPY ./go.mod $WORKDIR/teamsedge/go.mod
COPY ./go.sum $WORKDIR/teamsedge/go.sum

RUN cd $WORKDIR/teamsedge && go mod download

ARG BTIME
ENV RELEASE_VERSION=v1.0.1
ENV BUILD_TIME=${BTIME:-latest}

COPY ./assets $WORKDIR/teamsedge/assets
COPY ./common $WORKDIR/teamsedge/common
COPY ./config $WORKDIR/teamsedge/config
COPY ./jobs $WORKDIR/teamsedge/jobs
COPY ./service $WORKDIR/teamsedge/service
COPY ./main.go $WORKDIR/teamsedge/main.go

RUN cd $WORKDIR/teamsedge && \
  CGO_ENABLED=0 go build -a -ldflags  '-s -w -extldflags "-static"'  -o /release/teamsedge main.go

#FROM python:3.9.6-alpine3.14
FROM alpine
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
#RUN apk add --no-cache curl
#RUN pip install pysocks requests
COPY --from=build /release/teamsedge /usr/bin/teamsedge

RUN chmod +x /usr/bin/teamsedge

CMD ["/usr/bin/teamsedge"]
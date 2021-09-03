FROM golang:1.16.3-buster AS build

WORKDIR $GOPATH/src

RUN mkdir -p /release && mkdir -p $WORKDIR/mediaedge

ENV GO111MODULE=on

COPY ./go.mod $WORKDIR/mediaedge/go.mod
COPY ./go.sum $WORKDIR/mediaedge/go.sum

RUN cd $WORKDIR/mediaedge && go mod download

ARG BTIME
ENV RELEASE_VERSION=v1.0.1
ENV BUILD_TIME=${BTIME:-latest}

COPY ./assets $WORKDIR/mediaedge/assets
COPY ./common $WORKDIR/mediaedge/common
COPY ./config $WORKDIR/mediaedge/config
COPY ./grpcservice $WORKDIR/mediaedge/grpcservice
COPY ./models $WORKDIR/mediaedge/models
COPY ./web $WORKDIR/mediaedge/web
COPY ./main.go $WORKDIR/mediaedge/main.go

COPY npc/npc_arm64 /release/npc_arm64
COPY assets/mediaserver.ini /release/mediaserver.ini

RUN cd $WORKDIR/mediaedge && \
  CGO_ENABLED=0 go build -a -ldflags  '-s -w -extldflags "-static"'  -o /release/mediaedge main.go

FROM superiot/mediaserver
RUN mv /opt/media/bin/MediaServer /usr/bin/MediaServer
COPY --from=build /release/mediaserver.ini /etc/config/mediaserver.ini
COPY --from=build /release/mediaedge /usr/bin/mediaedge
COPY --from=build /release/npc_arm64 /usr/bin/npc

RUN chmod +x /usr/bin/mediaedge && \
    chmod +x /usr/bin/npc && \
    chmod +x /usr/bin/MediaServer


CMD ["/usr/bin/mediaedge"]
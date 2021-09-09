package app

import (
	"fmt"
	"log"

	"github.com/ca17/teamsedge/common"
	"github.com/ca17/teamsedge/common/sysid"
	"github.com/ca17/teamsedge/config"
	"github.com/ca17/teamsedge/service/models"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

var (
	Config           *config.AppConfig
	pubsock mangos.Socket
	subsock mangos.Socket
)

// Init 全局初始化调用
func Init(cfg *config.AppConfig) {
	Config = cfg
	StartPublish()

	// 发布启动消息
	Publish(models.TeamsEdgeBootstrap, models.EdgeBootstrapMessage{Eid: sysid.GetSystemSid()})
}

// StartSubscribe 启动边缘节点消息订阅
func StartSubscribe() error {
	var err error
	if subsock, err = sub.NewSocket(); err != nil {
		return fmt.Errorf("subscribe client create error %s", err.Error())
	}

	err = subsock.SetOption(mangos.OptionDialAsynch, true)
	if err != nil {
		return err
	}

	// 订阅客户端连接服务端的发布端口
	if err = subsock.Dial(Config.Teamsacs.PubAddr); err != nil {
		return fmt.Errorf("subscribe client connect to server error %s", err.Error())
	}

	log.Println(fmt.Sprintf("subscribe client connect to server pubaddr %s", Config.Teamsacs.PubAddr))

	// 轮训接收消息
	for {
		msg, err := subsock.Recv()
		if err != nil {
			log.Printf("Subscriber recv Message error %s", err.Error())
			continue
		}
		go onMessage(msg)
	}
}

// StartPublish 启动发布客户端
func StartPublish() {
	var err error
	if pubsock, err = pub.NewSocket(); err != nil {
		common.Must(fmt.Errorf("publish client create error %s", err.Error()))
	}

	err = pubsock.SetOption(mangos.OptionDialAsynch, true)
	common.Must(err)
	// 发布客户端连接服务端的订阅端口
	if err = pubsock.Dial(Config.Teamsacs.SubAddr); err != nil {
		common.Must(fmt.Errorf("publish client connect to server error %s", err.Error()))
	}
	log.Println(fmt.Sprintf("publish client connect to server subaddr %s", Config.Teamsacs.SubAddr))
}



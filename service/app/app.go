package app

import (
	"encoding/json"
	"fmt"
	"github.com/teamsweHere/teamsedge-hy/common"
	"github.com/teamsweHere/teamsedge-hy/config"
	"github.com/teamsweHere/teamsedge-hy/service/models"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
	"log"
)

var (
	Config           *config.AppConfig
	SubscribeChannel chan models.EdgeMessage // 消息订阅接收通道

	pubsock mangos.Socket
	subsock mangos.Socket
)

// Init 全局初始化调用
func Init(cfg *config.AppConfig) {
	Config = cfg
	SubscribeChannel = make(chan models.EdgeMessage, 10)
	StartPublish()
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

	if err = subsock.Dial(Config.Teamsacs.PubAddr); err != nil {
		return fmt.Errorf("subscribe client connect to server error %s", err.Error())
	}

	log.Println(fmt.Sprintf("subscribe client connect to server %s", Config.Teamsacs.PubAddr))

	// 轮训接收消息
	for {
		msg, err := subsock.Recv()
		if err != nil {
			log.Printf("recv Message error %s", err.Error())
			continue
		}
		var qmsg models.EdgeMessage
		err = json.Unmarshal(msg, &qmsg)
		if err != nil {
			log.Printf("Unmarshal Message error %s", err.Error())
			continue
		}
		SubscribeChannel <- qmsg
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
	if err = pubsock.Dial(Config.Teamsacs.PubAddr); err != nil {
		common.Must(fmt.Errorf("publish client connect to server error %s", err.Error()))
	}

	log.Println(fmt.Sprintf("publish client connect to server %s", Config.Teamsacs.PubAddr))

}

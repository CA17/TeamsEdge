package app

import (
	"bytes"
	"log"

	"github.com/ca17/teamsedge/common"
	"github.com/ca17/teamsedge/common/msgutil"
)

// onMessage
func onMessage(src []byte) {
	log.Printf("onMessage %s", string(src))
	switch {
	default:
	}
}

// Publish 向 TeamsACS 发布消息
func Publish(topic string, msg interface{}) {
	// 加密编码消息
	_msg, err := msgutil.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	var buff = bytes.NewBuffer([]byte(topic))
	buff.Write(_msg)
	log.Printf("publinsh teamsedge message %s", common.ToJson(msg))
	err = pubsock.Send(buff.Bytes())
	if err != nil {
		log.Println(err)
	}
}

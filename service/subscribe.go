package service

import (
	"github.com/gorilla/websocket"
	"log"
)

const LIVE_DANMU_SUBSCRIBE_PATTERN string  = "live::danmu::*"

var (
	liveDanmuSubscribeHandler = LiveDanmuSubscribeHandler{
		make (chan *SubscribeMessage, 1000),
	}
)

type SubscribeHandler interface {
	ProcessLoop ()
	Subscribe(channel string, message string)
}

type SubscribeMessage struct {
	 channel string
	 message string
}

type LiveDanmuSubscribeHandler struct {
	InChan chan *SubscribeMessage
}

func (this LiveDanmuSubscribeHandler) Subscribe(channel string, message string)  {
	this.InChan <- &SubscribeMessage{channel, message}
}

func (this LiveDanmuSubscribeHandler) ProcessLoop () {

	channelPrefix := LIVE_DANMU_SUBSCRIBE_PATTERN[:len(LIVE_DANMU_SUBSCRIBE_PATTERN)-1]

	bucketManager:= GetBucketManager()

	for {
		select {
			case msg:= <-this.InChan:
				go func() {
					channel := msg.channel
					message := msg.message
					roomId := channel[len(channelPrefix):]
					log.Printf("收到redis message -> %s, roomId -> %s", message, roomId)
					bucket := bucketManager.getBucket(roomId)
					data := []byte(message)
					for i:=range bucket.Conn.Iter() {
						wsConnection := i.(*WsConnection)
						go wsConnection.wsWrite(websocket.TextMessage, data)
					}
				}()
		}
	}
}

func InitSubscribes () {
	go liveDanmuSubscribeHandler.ProcessLoop()
	go PSub(LIVE_DANMU_SUBSCRIBE_PATTERN, liveDanmuSubscribeHandler) //监听直播房间弹幕消息队列
}

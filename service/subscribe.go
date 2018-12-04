package service

const LIVE_DANMU_SUBSCRIBE_PATTERN string  = "live::danmu::*"
const VOD_DANMU_SUBSCRIBE_PATTERN string = "vod::danmu::*"

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
	//for {
	//	select {
	//		case msg:= <-this.InChan:
	//			go func() {
	//				//danmuCache.WriteLiveDanmu(msg)
	//			}()
	//	}
	//}
}

func InitSubscribes () {

	go liveDanmuSubscribeHandler.ProcessLoop()
	go PSub(LIVE_DANMU_SUBSCRIBE_PATTERN, liveDanmuSubscribeHandler) //监听直播房间弹幕消息队列

	//go PSub(VOD_DANMU_SUBSCRIBE_PATTERN) //监听回放弹幕消息队列


}

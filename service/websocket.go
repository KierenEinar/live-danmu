package service

import (
	"errors"
	"github.com/gorilla/websocket"
	"sync"
)

type WsConnection struct {
	WsSocketConn *websocket.Conn // 底层websocket
	InChan chan *WsMsgType // 读队列
	OutChan chan *WsMsgType // 写队列
	mutex sync.Mutex	// 避免重复关闭管道
	IsClosed bool
	CloseChan chan byte  // 关闭通知
}

type WsMsgType struct {
	MessageType int
	data []byte
}

func (this *WsConnection) WsHeartBeat() error {



}

func (this *WsConnection) wsWrite (messageType int, data []byte) error {
	select {
		case this.OutChan <- &WsMsgType{messageType, data}:
		case <- this.CloseChan:
			return errors.New("websocket closed")
	}
	return nil
}

func (this *WsConnection) wsRead () (*WsMsgType, error) {
	select {
	case msg := <- this.InChan:
		return msg, nil
	case <- this.CloseChan:
		return nil, errors.New("websocket closed")
	}
}

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
	Mutex sync.Mutex	// 避免重复关闭管道
	IsClosed bool
	CloseChan chan byte  // 关闭通知
	RoomId string
}

type WsMsgType struct {
	MessageType int
	data []byte
}

func (this *WsConnection) WsHeartBeat() error {

	for {
		msg, err := this.wsRead()
		if err != nil {
			this.WsClose()
			return errors.New("websocket closed")
		}
		this.wsWrite(msg.MessageType, msg.data)
	}

	return nil

}

func (this * WsConnection) WsReadLoop () error {

	for {

		msgType, data, err := this.WsSocketConn.ReadMessage()
		if err != nil {
			this.WsClose()
		}

		select {
			case this.InChan <- &WsMsgType{msgType, data}:
			case <-this.CloseChan:
				this.WsClose()
				return errors.New("websocket closed")
		}
	}
}

func (this * WsConnection) WsWriteLoop () error {

	for {
		select {
			case msg := <-this.OutChan:
				err := this.WsSocketConn.WriteMessage(msg.MessageType, msg.data)
				if err != nil {
					this.WsClose()
				}
			case <- this.CloseChan:
				this.WsClose()
				return errors.New("websocket closed")
		}
	}
}

func (this *WsConnection) WsClose () {
	bucketManager := GetBucketManager()
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if !this.IsClosed {
		bucketManager.DelConn4Buckets(this)
		this.WsSocketConn.Close()
		this.IsClosed = true
		close(this.CloseChan)
	}
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

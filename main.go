package main

import (
	"github.com/gorilla/websocket"
	"live-danmu/service"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main () {

	appConfig := service.LoadConfig()
	redisConfig := appConfig.GetRedisConfig()
	redisConfig.Connect()
	go service.Sub("hello")

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		var (
			conn * websocket.Conn
			err error
		)

		conn, err = upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}

		wsConnection := service.WsConnection{
			conn,
			make (chan *service.WsMsgType, 1000),
			make (chan * service.WsMsgType, 1000),
			sync.Mutex{},
			false,
			make (chan byte, 1),
		}
		go wsConnection.WsHeartBeat()

		go wsConnection.WsReadLoop()

		go wsConnection.WsWriteLoop()
	})

	http.ListenAndServe("0.0.0.0:7777",nil)

}

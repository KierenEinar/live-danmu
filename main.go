package main

import (
	"github.com/gorilla/websocket"
	"live-danmu/service"
	"net/http"
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
			data []byte
		)

		conn, err = upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}

		for {
			_, data ,err = conn.ReadMessage()
			if err != nil {
				goto ERR
			}
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				goto ERR
			}
		}
	ERR:
		conn.Close()
	})

	http.ListenAndServe("0.0.0.0:7777",nil)

}

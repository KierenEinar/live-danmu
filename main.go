package main

import (
	"github.com/gorilla/websocket"
	"live-danmu/service"
	"log"
	"net/http"
	"strings"
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
	service.InitSubscribes()

	bucketManager := service.GetBucketManager()

	http.HandleFunc("/ws/", func(writer http.ResponseWriter, request *http.Request) {
		var (
			conn * websocket.Conn
			err error
			roomId string
		)

		requestUri := request.RequestURI

		if !strings.Contains(requestUri, "/live/") && !strings.Contains(requestUri, "/vod/") {
			log.Printf("非法请求")
			return
		}


		if strings.Contains(requestUri, "/live/") {
			roomId = requestUri[strings.Index(requestUri, "/live/")+1:]
			roomId = strings.Replace(roomId, "/", ":", -1)
		} else {
			roomId = requestUri[strings.Index(requestUri, "/vod/")+1:]
			roomId = strings.Replace(roomId, "/", ":", -1)
		}


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
			roomId,
		}

		bucketManager.AddConn2Buckets(&wsConnection)

		go wsConnection.WsHeartBeat()

		go wsConnection.WsReadLoop()

		go wsConnection.WsWriteLoop()
	})

	http.ListenAndServe(appConfig.ServerAddr, nil)

}

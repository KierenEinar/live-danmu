package service

import (
	"bytes"
	"github.com/deckarep/golang-set"
	"github.com/stathat/consistent"
	"strconv"
	"sync"
)

const BUCKET_PARTITION = 1000

/**
一个房间对应一个buckets
*/

type Bucket struct {
	OutChannel chan *SubscribeMessage //队列
	RWMutex sync.RWMutex	//读写锁
	index int	//buckets 第几个
	Conn mapset.Set
}

type BucketManager struct {
	Buckets map[string]*Bucket //一致性hash, roomId -> 打散到bucket
	consistent *consistent.Consistent
}

var (
	bucketManager *BucketManager
	mutex sync.Mutex
)

func GetBucketManager () *BucketManager{
	if bucketManager != nil {
		return bucketManager
	}
	mutex.Lock()
	defer mutex.Unlock()

	c := consistent.New()

	bucketMap := make(map[string]*Bucket,BUCKET_PARTITION)

	for  i := 0; i< BUCKET_PARTITION ; i++ {
		var buffer bytes.Buffer
		buffer.WriteString(strconv.Itoa(i))
		partition := buffer.String()
		bucketMap[partition] = &Bucket{
			make (chan *SubscribeMessage, 1000),
			sync.RWMutex{},
			i,
			mapset.NewThreadUnsafeSet(),
		}
		c.Add(partition)
	}
	bucketManager = &BucketManager{bucketMap, c}
	return bucketManager
}

func (this *BucketManager) AddConn2Buckets (connection *WsConnection, roomId string) {
	partition,_ := this.consistent.Get(roomId)
	bucket := this.Buckets[partition]
	bucket.RWMutex.Lock()
	defer bucket.RWMutex.Unlock()
	bucket.Conn.Add(connection)
}
package service

import (
	"bytes"
	"github.com/deckarep/golang-set"
	"github.com/stathat/consistent"
	"log"
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
	Conn mapset.Set
}


type Segment struct {
	Buckets map[string]*Bucket
	RWMutex sync.RWMutex
}

type BucketManager struct {
	Segments map[string]*Segment //一致性hash, roomId -> 打散到bucket
	consistent *consistent.Consistent
}

var (
	bucketManager *BucketManager
	once = sync.Once{}
)

func GetBucketManager () *BucketManager{

	once.Do(func() {

		c := consistent.New()

		segment := make(map[string]*Segment,BUCKET_PARTITION)

		for  i := 0; i< BUCKET_PARTITION ; i++ {
			var buffer bytes.Buffer
			buffer.WriteString(strconv.Itoa(i))
			partition := buffer.String()
			segment[partition] = &Segment{make(map[string]*Bucket), sync.RWMutex{}}
			c.Add(partition)
		}
		bucketManager = &BucketManager{segment, c}
	})

	return bucketManager
}

func (this *BucketManager) AddConn2Buckets (connection *WsConnection) {
	roomId := connection.RoomId
	bucket := this.getBucket(roomId)
	bucket.RWMutex.Lock()
	defer bucket.RWMutex.Unlock()
	bucket.Conn.Add(connection)
}

func (this * BucketManager) getBucket (roomId string) *Bucket {
	log.Printf("%p", this)
	segment := this.getSegment(roomId)
	buckets := segment.Buckets
	bucket := buckets[roomId]

	if bucket == nil {
		segment.RWMutex.Lock()
		bucket = &Bucket{
			make (chan *SubscribeMessage, 1000),
			sync.RWMutex{},
			mapset.NewSet(),
		}
		buckets[roomId] = bucket
		segment.RWMutex.Unlock()
	}
	return bucket
}

func (this * BucketManager) getSegment (roomId string) *Segment {
	partition,_ := this.consistent.Get(roomId)
	segment := this.Segments[partition]
	return segment
}

func (this *BucketManager) DelConn4Buckets (connection *WsConnection) {
	roomId := connection.RoomId
	bucket := this.getBucket(roomId)
	bucket.RWMutex.Lock()
	defer bucket.RWMutex.Unlock()
	bucket.Conn.Remove(connection)
	if bucket.Conn.Cardinality() == 0 {
		segment := this.getSegment(roomId)
		delete (segment.Buckets, roomId)
	}

}

func (this *Bucket) push () {

}
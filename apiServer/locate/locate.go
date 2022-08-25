package locate

import (
	"encoding/json"
	"git_test/src/lib/rabbitmq"
	"git_test/src/lib/rs"
	"git_test/src/lib/types"
	"os"
	"time"
)

//1s的时间定位存放对象的数据节点
func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}

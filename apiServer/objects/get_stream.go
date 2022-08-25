package objects

import (
	"fmt"
	"git_test/apiServer/heartbeat"
	"git_test/apiServer/locate"
	"git_test/src/lib/rs"
)

//定位数据节点，根据数据节点返回一个objectstream.NewGetStream，该结构体实现io.reader
func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail,result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}

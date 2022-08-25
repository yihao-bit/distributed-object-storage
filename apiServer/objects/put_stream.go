package objects

import (
	"fmt"
	"git_test/apiServer/heartbeat"
	"git_test/src/lib/rs"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return rs.NewRSPutStream(servers, hash, size)
}

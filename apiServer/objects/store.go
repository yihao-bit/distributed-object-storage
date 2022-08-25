package objects

import (
	"fmt"
	"git_test/apiServer/locate"
	"git_test/src/lib/utils"
	"io"
	"net/http"
	"net/url"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,calculated=%d,requested=%ss", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}
